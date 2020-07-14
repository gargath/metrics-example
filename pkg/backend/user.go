package backend

import (
	"fmt"
	"time"

	// import SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
)

// ListUsers returns all users from the database
func (b *SqliteBackend) ListUsers() ([]User, error) {
	timer := prometheus.NewTimer(dbDuration.WithLabelValues("user"))
	u, e := b.listUsers()
	timer.ObserveDuration()
	return u, e
}

func (b *SqliteBackend) listUsers() ([]User, error) {
	var users []User

	rows, err := b.db.Query("SELECT id, name, dob, address FROM users")
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
func (b *SqliteBackend) AddUser(u User) error {
	timer := prometheus.NewTimer(dbDuration.WithLabelValues("user"))
	e := b.addUser(u)
	timer.ObserveDuration()
	return e
}

func (b *SqliteBackend) addUser(u User) error {
	rows, err := b.db.Query("SELECT id FROM users WHERE id=?", u.ID)
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
	result, err := b.db.Exec("INSERT INTO users(id, name, dob, address) VALUES (?,?,?,?)", u.ID, u.Name, u.DoB, u.Address)
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

// DeleteUser removes the user with the given ID from the database
// If a user with that ID does not exist, it returns ErrDoesNotExist
func (b *SqliteBackend) DeleteUser(id string) error {
	timer := prometheus.NewTimer(dbDuration.WithLabelValues("user"))
	e := b.deleteUser(id)
	timer.ObserveDuration()
	return e
}

func (b *SqliteBackend) deleteUser(id string) error {
	rows, err := b.db.Query("SELECT id FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("Failed to check for existing user: %v", err)
	}
	if !rows.Next() {
		if e := rows.Err(); e != nil {
			return fmt.Errorf("Failed to check for existing user: %v", err)
		}
		return ErrDoesNotExist
	}
	rows.Close()

	result, err := b.db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("Failed to delete user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unexpected error checking rows affected: %v", err)
	} else if rowsAffected != 1 {
		return fmt.Errorf("Unexpected rows affected by insert. Expected 1, got %d", rowsAffected)
	}
	return nil
}

// GetUser returns the user with the given ID.
// If no such user exists, it returns ErrDoesNotExist
func (b *SqliteBackend) GetUser(id string) (*User, error) {
	timer := prometheus.NewTimer(dbDuration.WithLabelValues("user"))
	u, e := b.getUser(id)
	timer.ObserveDuration()
	return u, e
}

func (b *SqliteBackend) getUser(id string) (*User, error) {
	rows, err := b.db.Query("SELECT id, name, dob, address FROM users WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("Failed to check for existing user: %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		if e := rows.Err(); e != nil {
			return nil, fmt.Errorf("Failed to check for existing user: %v", err)
		}
		return nil, ErrDoesNotExist
	}

	var uid string
	var name string
	var dob string
	var address string

	err = rows.Scan(&uid, &name, &dob, &address)
	if err != nil {
		return nil, fmt.Errorf("Failed to read user record from database: %v", err)
	}
	dobtime, err := time.Parse("2006-01-02 15:04:05-07:00", dob)
	if err != nil {
		return nil, fmt.Errorf("Encountered invalid date %s: %v", dob, err)
	}
	u := User{
		ID:      uid,
		Name:    name,
		DoB:     dobtime,
		Address: address,
	}

	return &u, nil
}
