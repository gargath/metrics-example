package api

import (
	"github.com/gargath/metrics-example/pkg/backend"
)

// API provides the API and handler functions
type API struct {
	Prefix string
	b      backend.Backend
}
