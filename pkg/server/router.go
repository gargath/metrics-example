package server

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gargath/metrics-example/pkg/api"
	"github.com/gargath/metrics-example/pkg/backend"
)

func buildRouter() (*mux.Router, error) {
	router := mux.NewRouter()

	router.Use(prometheusMiddleware)

	backend, err := backend.NewSqliteBackend("./anaximander.db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize backend: %v", err)
	}
	api := api.NewAPI("/api", backend)
	api.AddRoutes(router)

	prometheus.MustRegister(responseCounter)
	router.Path("/metrics").Handler(promhttp.Handler())

	return router, nil
}
