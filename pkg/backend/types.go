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
