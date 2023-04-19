package database

import (
	"database/sql"
	"log"
)

type UserAction struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Action   string `json:"action"`
	Response string `json:"response"`
	Time     int    `json:"time"`
}

func CreateActionTable() error {
	// create the actions table if it doesn't exist
	_, err := DatabaseConnection.Database.Exec(
		`create table if not exists actions(
					id integer primary key autoincrement,
					user_id integer not null,
					action text not null,
					response text not null default '',
					time integer not null default (strftime('%s','now')),
					foreign key(user_id) references users(id)
				)`,
	)
	if err != nil {
		return err
	}

	return nil
}

func LogAction(userID int, action string, response string) error {
	// log an action
	_, err := DatabaseConnection.Database.Exec(
		`insert into actions(user_id, action, response) values(?, ?, ?)`,
		userID,
		action,
		response,
	)
	if err != nil {
		return err
	}

	return nil
}

func FetchActions(limit int, offset int) ([]*UserAction, error) {
	// fetch all actions
	rows, err := DatabaseConnection.Database.Query(
		`select id, user_id, action, response, time from actions order by id desc limit ? offset ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	actions := make([]*UserAction, 0)

	for rows.Next() {
		var action UserAction
		err := rows.Scan(&action.ID, &action.UserID, &action.Action, &action.Response, &action.Time)
		if err != nil {
			return nil, err
		}
		actions = append(actions, &action)
	}

	return actions, nil
}

func FetchActionsFromUser(userID int, limit int, offset int) ([]*UserAction, error) {
	// fetch all actions from a user
	rows, err := DatabaseConnection.Database.Query(
		`select id, user_id, action, response, time from actions where user_id = ? order by id desc limit ? offset ?`,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	actions := make([]*UserAction, 0)

	for rows.Next() {
		var action UserAction
		err := rows.Scan(&action.ID, &action.UserID, &action.Action, &action.Response, &action.Time)
		if err != nil {
			return nil, err
		}
		actions = append(actions, &action)
	}

	return actions, nil
}
