package server

import (
	"net/http"
)

// Header calls the underlying writer's Header() function
func (c *CapturingResponseWriter) Header() http.Header {
	return c.w.Header()
}

// Write calls the underlying writer's Write() function
func (c *CapturingResponseWriter) Write(data []byte) (int, error) {
	return c.w.Write(data)
}

// WriteHeader captures the status code written, then calls the underlying writer's WriteHeader() function
func (c *CapturingResponseWriter) WriteHeader(statusCode int) {
	c.Status = statusCode
	c.w.WriteHeader(statusCode)
}

// WrapHandleFunc wraps a given http.Handler, wrapping its ResponseWriter in a CapturingResponseWriter
func WrapHandleFunc(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		crw := &CapturingResponseWriter{w: w}
		wrappedHandler.ServeHTTP(crw, req)
	})
}
