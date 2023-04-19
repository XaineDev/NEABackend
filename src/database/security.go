package database

func CreateSecurityTable() error {
	// create the security table if it doesn't exist
	_, err := DatabaseConnection.Database.Exec(
		`create table if not exists security(
					id integer primary key autoincrement,
					user_id integer not null,
					password text not null,
					password_salt text not null,
					token text not null,
					foreign key(user_id) references users(id)
				)`,
	)
	if err != nil {
		return err
	}

	return nil
}
