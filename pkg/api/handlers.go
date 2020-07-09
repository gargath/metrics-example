package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ListFoos handles GET request for /foo without a specific entity ID
func (a *API) ListFoos(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "List of Foo:\n - <empty>\n")
}

// GetFoo handles GET request for /foo/:id
func (a *API) GetFoo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Requested Foo ID: %v\n", vars["id"])
}
