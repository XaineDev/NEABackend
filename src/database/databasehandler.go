package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DatabaseConnection *Connection
)

/*
Startup ~ Used to start up the database and create the default tables if they dont exist
*/
func Startup(dataFolder string) error {

	// open connection to the database
	conn, err := sql.Open("sqlite3", dataFolder+"/database.sqlite")
	if err != nil {
		return err
	}

	// set the database connection
	DatabaseConnection = &Connection{
		Database: conn,
	}

	// setup default tables
	err = CreateUsersTable()
	if err != nil {
		log.Println("Failed to create users table: ", err)
	}
	err = CreateBooksTable()
	if err != nil {
		log.Println("Failed to create books table: ", err)
	}
	err = CreateActionTable()
	if err != nil {
		log.Println("Failed to create actions table: ", err)
	}
	err = CreateSecurityTable()
	if err != nil {
		log.Println("Failed to create security table: ", err)
	}
	return nil
}
