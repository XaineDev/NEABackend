package database

import (
	"NEABackend/src/util"
	"database/sql"
	"log"
	"strconv"
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
		`insert into users(username) values(?)`,
		username,
		pass,
		salt,
	)
	if err != nil {
		return err
	}
	_, err = c.Database.Exec(
		`insert into security(user_id, password, password_salt, token) values((select id from users where username = ?), ?, ?, ?)`,
		username,
		pass,
		salt,
		"",
	)
	return nil
}

func (c *Connection) GetUserByUsername(username string) (*User, error) {
	var user User
	err := c.Database.QueryRow(
		`SELECT u.id,
					   u.username,
					   u.admin,
					   u.created_at,
					   s.password,
					   s.password_salt,
					   s.token
				FROM   users u
					   INNER JOIN security s
							   ON u.id = s.user_id
				WHERE  u.username = ? `,
		username,
	).Scan(&user.ID, &user.Username, &user.IsAdmin, &user.CreatedAt, &user.Password, &user.PasswordSalt, &user.Token)
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

func (c *Connection) GetUserByToken(token string) (*User, error) {
	var userId int
	err := c.Database.QueryRow("select user_id from security where token = ?", token).Scan(&userId)
	if err != nil {
		return nil, err
	}
	return c.GetUserById(userId)
}

func (c *Connection) GetUserById(id int) (*User, error) {
	var user User
	err := c.Database.QueryRow(
		`SELECT u.id,
					   u.username,
					   u.admin,
					   u.created_at,
					   s.password,
					   s.password_salt,
					   s.token
				FROM   users u
					   INNER JOIN security s
							   ON u.id = s.user_id
				WHERE  u.id = ? `,
		id,
	).Scan(&user.ID, &user.Username, &user.IsAdmin, &user.CreatedAt, &user.Password, &user.PasswordSalt, &user.Token)
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

func (c *Connection) UpdateUserToken(user *User, token string) error {
	_, err := c.Database.Exec("update security set token = ? where user_id = ?", token, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) IsAdminToken(token string) (bool, error) {
	var isAdmin bool
	var userId int
	err := c.Database.QueryRow("select user_id from security where token = ?", token).Scan(&userId)
	if err != nil {
		return false, err
	}
	err = c.Database.QueryRow("select admin from users where id = ?", userId).Scan(&isAdmin)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func (c *Connection) GetBooksFromUser(user User) ([]Book, error) {
	rows, err := c.Database.Query(`select id, title, author from books where currentOwner = ?`, user.ID)
	if err != nil {
		return nil, err
	}
	var books []Book
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			return nil, err
		}
		book.CurrentOwner = user.ID
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
	book.ID = int(id)
	return nil
}

func (c *Connection) ValidateRequest(token string, userID string, username string) *User {
	user, err := c.GetUserByUsername(username)
	if user == nil || err != nil {
		if err != nil {
			log.Printf("Failed to validate user, %s\n", err.Error())
		}
		return nil
	}
	if user.Token != token {
		return nil
	}

	userIDInt, err := strconv.Atoi(userID)

	if err != nil || user.ID != userIDInt {
		if err != nil {
			log.Printf("Failed to validate user, %s\n", err.Error())
		}
		return nil
	}
	return user
}

func (c *Connection) GetBookFromId(id int) (*Book, error) {
	bookOwner := sql.NullInt16{}
	row := c.Database.QueryRow(`select id, title, author, currentOwner from books where id = ?`, id)
	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &bookOwner)
	if err != nil {
		return nil, err
	}
	if bookOwner.Valid {
		book.CurrentOwner = int(bookOwner.Int16)
	} else {
		book.CurrentOwner = 0
	}
	return &book, nil
}

func (c *Connection) UpdateBook(book *Book) error {
	_, err := c.Database.Exec(`update books set title = ?, author = ?, currentOwner = ? where id = ?`, book.Title, book.Author, book.CurrentOwner, book.ID)
	return err
}

func (c *Connection) UnclaimBook(book *Book) error {
	_, err := c.Database.Exec(`update books set currentOwner = null where id = ?`, book.ID)
	return err
}

func (c *Connection) GetBooks() ([]*Book, error) {
	rows, err := c.Database.Query(`select * from books`)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows while getting books: ", err)
		}
	}(rows)
	if err != nil {
		return nil, err
	}

	var books []*Book
	for rows.Next() {
		var book Book
		ownerHolder := sql.NullInt64{}
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &ownerHolder)
		if err != nil {
			return nil, err
		}
		if ownerHolder.Valid {
			book.CurrentOwner = int(ownerHolder.Int64)
		} else {
			book.CurrentOwner = 0
		}
		books = append(books, &book)
	}

	return books, nil
}

func (c *Connection) GetUsers(amount int, page int) ([]User, error) {
	var users []User
	rows, err := c.Database.Query(`select * from users limit ? offset ?`, amount, amount*page)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows while getting users: ", err)
		}
	}(rows)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
