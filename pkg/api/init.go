package api

import "github.com/gargath/metrics-example/pkg/backend"

// NewAPI returns an initialized API
func NewAPI(prefix string, b backend.Backend) *API {
	return &API{
		Prefix: prefix,
		b:      b,
	}
}
