package database

import (
	"NEABackend/src/util"
	"database/sql"
)

type Connection struct {
	Database *sql.DB
}

func (c *Connection) Close() error {
	return c.Database.Close()
}

func (c *Connection) CreateUser(username string, password string) error {
	pass, salt, err := util.HashPassowrd(password)
	if err != nil {
		return err
	}
	_, err = c.Database.Exec(
		`insert into users(username, password, passwordSalt) values(?, ?, ?)`,
		username,
		pass,
		salt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) GetUserByUsername(username string) (*User, error) {
	var user User
	err := c.Database.QueryRow(
		`select id, username, admin, password, passwordSalt, token, created_at from users where username = ?`,
		username,
	).Scan(&user.Id, &user.Username, &user.IsAdmin, &user.Password, &user.PasswordSalt, &user.Token, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Connection) UpdateUserToken(username string, token string) error {
	_, err := c.Database.Exec("update users set token = ? where username = ?", token, username)
	if err != nil {
		return err
	}
	return nil
}
