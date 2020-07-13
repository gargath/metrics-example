package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "myapp_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	responseCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_respose_total",
			Help: "Number of HTTP responses sent.",
		},
		[]string{"code", "method", "path"},
	)
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		crw := &CapturingResponseWriter{w: w}
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(crw, r)
		timer.ObserveDuration()
		responseCounter.With(prometheus.Labels{"code": strconv.Itoa(crw.Status), "method": r.Method, "path": path}).Inc()
	})
}
