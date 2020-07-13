package server

import (
	"net/http"
	"time"
)

// Server is the HTTP server struct used to serve the API
type Server struct {
	GracefulShutdownPeriod time.Duration
	Addr                   string `default:"0.0.0.0:8080"`
	srv                    *http.Server
}

// CapturingResponseWriter wraps an http.ResponseWriter, capturing the HTTP status code written to it
// for metrics collection
type CapturingResponseWriter struct {
	w      http.ResponseWriter
	Status int
}
