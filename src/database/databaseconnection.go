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

	books, err := c.GetBooksFromUser(user)
	if err != nil {
		return nil, err
	}
	user.CurrentBooks = books

	return &user, nil
}

func (c *Connection) UpdateUserToken(username string, token string) error {
	_, err := c.Database.Exec("update users set token = ? where username = ?", token, username)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) IsAdminToken(token string) (bool, error) {
	var isAdmin bool
	err := c.Database.QueryRow("select admin from users where token = ?", token).Scan(&isAdmin)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func (c *Connection) GetBooksFromUser(user User) ([]Book, error) {
	rows, err := c.Database.Query(`select id, title, author from books where currentOwner = ?`, user.Id)
	if err != nil {
		return nil, err
	}
	var books []Book
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.Id, &book.Title, &book.Author)
		if err != nil {
			return nil, err
		}
		book.CurrentOwner = user.Username
		books = append(books, book)
	}
	return books, nil
}

func (c *Connection) CreateBook(book *Book) error {
	result, err := c.Database.Exec(`insert into books(title, author) values(?, ?)`, book.Title, book.Author)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	book.Id = int(id)
	return nil
}
