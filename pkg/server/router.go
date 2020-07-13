package server

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gargath/metrics-example/pkg/api"
)

func buildRouter() (*mux.Router, error) {
	router := mux.NewRouter()

	router.Use(prometheusMiddleware)

	api := &api.API{Prefix: "/api", ResponseCount: responseCounter}
	api.AddRoutes(router)

	prometheus.MustRegister(responseCounter)
	router.Path("/metrics").Handler(promhttp.Handler())

	return router, nil
}
