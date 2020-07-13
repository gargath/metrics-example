package backend

import (
	"database/sql"

	// import SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	db, err := sql.Open("sqlite3", "./anaximander.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users ( id TEXT PRIMARY KEY, name TEXT NOT NULL, dob TEXT NOT NULL, address, TEXT );")
	if err != nil {
		panic(err)
	}
}
