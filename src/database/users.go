package database

type User struct {
	ID           int    `json:"id"`                     // id of the user
	Username     string `json:"username"`               // username for logging in
	Password     string `json:"password,omitempty"`     // hashed password for logging in, hashed using argon2id
	PasswordSalt string `json:"passwordSalt,omitempty"` // salt used for hashing the password
	Token        string `json:"token"`                  // token used for fetching / updating information needing an account
	CreatedAt    int    `json:"created_at"`             // unix timestamp of account creation
	IsAdmin      bool   `json:"isAdmin"`                // if the user is an admin
	CurrentBooks []Book `json:"currentBooks"`           // books the user is currently owner of
}

func CreateUsersTable() error {
	// create the users table if it doesn't exist
	_, err := DatabaseConnection.Database.Exec(
		`create table if not exists users( 
    				id integer primary key autoincrement, 
    				username text not null unique, 
    				admin integer not null default 0,
    				created_at integer not null default (strftime('%s','now'))
            	)`,
	)
	if err != nil {
		return err
	}

	return nil
}
