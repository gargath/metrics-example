package api_test

import (
	"time"

	"github.com/gargath/metrics-example/pkg/backend"
)

var newUser = &backend.User{
	ID:      "1234-5678-90123",
	Name:    "New User",
	DoB:     time.Now(),
	Address: "25 Fobar Road, 12345 Alphabet Town",
}
