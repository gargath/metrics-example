package backend

import (
	"database/sql"
	"errors"
	"time"
)

// Backend is the presistence backend used by the CRUD API
type Backend interface {
	AddUser(User) error
	GetUser(string) (*User, error)
	ListUsers() ([]User, error)
	DeleteUser(string) error
}

// SqliteBackend is a backend implementation using sqlite
type SqliteBackend struct {
	db *sql.DB
}

// User represents a user entity
type User struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	DoB     time.Time `json:"dob"`
	Address string    `json:"address"`
}

// ErrAlreadyExists is a named error used when an entity to be added already exists
// in the backend
var ErrAlreadyExists = errors.New("entity already exists")

// ErrDoesNotExist is a named error used to indicate that the entity an operation was requested
// on does not exist
var ErrDoesNotExist = errors.New("entity does not exist")

// Equals returns true if the receiver user is field-equal to the argument user
// TODO: Generate this
func (u *User) Equals(comp *User) bool {
	if u.ID != comp.ID {
		return false
	}
	if u.Name != comp.Name {
		return false
	}
	if u.Address != comp.Address {
		return false
	}
	if u.DoB.Sub(comp.DoB) > 0 {
		return false
	}

	return true
}
