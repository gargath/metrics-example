package backend

import (
	"database/sql"
	"fmt"
	"time"

	// import SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
)

func ListUsers() ([]User, error) {
	timer := prometheus.NewTimer(dbDuration.WithLabelValues("user"))
	u, e := listUsers()
	timer.ObserveDuration()
	return u, e
}

func listUsers() ([]User, error) {
	var users []User

	db, err := sql.Open("sqlite3", "./anaximander.db")
	if err != nil {
		return users, fmt.Errorf("Failed to open database: %v", err)
	}
	rows, err := db.Query("SELECT id, name, dob, address FROM users")
	if err != nil {
		return users, fmt.Errorf("Failed to retrieve users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		var dob string
		var address string

		err = rows.Scan(&id, &name, &dob, &address)
		if err != nil {
			return users, fmt.Errorf("Failed to read user record from database: %v", err)
		}
		dobtime, err := time.Parse("2006-01-02 15:04:05-07:00", dob)
		if err != nil {
			return users, fmt.Errorf("Encountered invalid date %s: %v", dob, err)
		}
		u := User{
			ID:      id,
			Name:    name,
			DoB:     dobtime,
			Address: address,
		}
		users = append(users, u)
	}
	return users, nil
}

// AddUser will add a given user entity to the database
func AddUser(u User) error {
	timer := prometheus.NewTimer(dbDuration.WithLabelValues("user"))
	e := addUser(u)
	timer.ObserveDuration()
	return e
}

func addUser(u User) error {
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
