package backend

import (
	"errors"
	"time"
)

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
