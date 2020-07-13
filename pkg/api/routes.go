package api

import (
	"github.com/gorilla/mux"
)

// AddRoutes adds the API routes to the provided router.
func (a *API) AddRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/foo", a.ListFoos).Methods("GET")
	apiRouter.HandleFunc("/foo/{id}", a.GetFoo).Methods("GET")

	apiRouter.HandleFunc("/user", a.AddUser).Methods("POST")
}
