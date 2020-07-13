package backend

import (
	"database/sql"
	"fmt"

	// import SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// AddUser will add a given user entity to the database
func AddUser(u User) error {
	db, err := sql.Open("sqlite3", "./anaximander.db")
	if err != nil {
		return fmt.Errorf("Failed to open database: %v", err)
	}
	rows, err := db.Query("SELECT id FROM users WHERE id=?", u.ID)
	if err != nil {
		return fmt.Errorf("Failed to check for existing user: %v", err)
	}
	defer rows.Close()
	if rows.Next() {
		if e := rows.Err(); e != nil {
			return fmt.Errorf("Failed to check for existing user: %v", err)
		}
		return ErrAlreadyExists
	}
	result, err := db.Exec("INSERT INTO users(id, name, dob, address) VALUES (?,?,?,?)", u.ID, u.Name, u.DoB, u.Address)
	if err != nil {
		return fmt.Errorf("Failed to insert user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unexpected error checking rows affected: %v", err)
	} else if rowsAffected != 1 {
		return fmt.Errorf("Unexpected rows affected by insert. Expected 1, got %d", rowsAffected)
	}
	return nil
}
