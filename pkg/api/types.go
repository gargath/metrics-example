package api

import (
	"github.com/prometheus/client_golang/prometheus"
)

// API provides the API and handler functions
type API struct {
	Prefix        string
	ResponseCount *prometheus.CounterVec
}
