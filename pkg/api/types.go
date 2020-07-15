package api

import (
	"github.com/gargath/metrics-example/pkg/backend"
)

// API provides the API and handler functions
type API struct {
	Prefix string
	b      backend.Backend
}

// ErrorResponse holds errors to be returned as JSON to client requests
type ErrorResponse struct {
	Error string `json:"error"`
}
