package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Run starts the server
func (s *Server) Run() error {
	s.srv = &http.Server{
		Addr:         s.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	r, err := buildRouter()
	if err != nil {
		return fmt.Errorf("error building routers: %v", err)
	}

	s.srv.Handler = r

	log.Printf("Server listening on %s", s.Addr)

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %v", err)
	}
	return nil
}

// Shutdown handles the graceful shutdown of the HTTP server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.GracefulShutdownPeriod)
	defer cancel()
	log.Println("Beginning Shutdown...")
	s.srv.SetKeepAlivesEnabled(false)
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shut down server: %v", err)
	}
	log.Println("Shutdown complete")
	return nil
}
