package api_test

import (
	"fmt"

	"github.com/gargath/metrics-example/pkg/backend"
)

type BrokenBackend struct{}

func (b *BrokenBackend) AddUser(u backend.User) error {
	return fmt.Errorf("broken backend")
}

func (b *BrokenBackend) GetUser(s string) (*backend.User, error) {
	return nil, fmt.Errorf("broken backend")
}

func (b *BrokenBackend) ListUsers() ([]backend.User, error) {
	return []backend.User{}, fmt.Errorf("broken backend")
}

func (b *BrokenBackend) DeleteUser(s string) error {
	return fmt.Errorf("broken backend")
}

func NewBrokenBackend() backend.Backend {
	return &BrokenBackend{}
}
