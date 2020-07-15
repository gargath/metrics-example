package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/gargath/metrics-example/pkg/backend"
)

//ListUsers handles GET requests for /user
func (a *API) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := a.b.ListUsers()
	if err != nil {
		writeError(w, fmt.Sprintf("failed to list users: %v", err), http.StatusInternalServerError)
		return
	}

	entityJSON, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		writeError(w, fmt.Sprintf("failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\n\tusers:\n"))
	w.Write(entityJSON)
	w.Write([]byte("\n}\n"))
}

// AddUser handles POST requests to /user
func (a *API) AddUser(w http.ResponseWriter, r *http.Request) {
	var u backend.User
	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		writeError(w, fmt.Sprintf("Content-Type %s not supported", r.Header.Get("Content-Type")), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var errors []error
	err := decoder.Decode(&u)
	if err != nil {
		errors = append(errors, err)
	}
	if u.ID != "" {
		writeError(w, "id must be empty", http.StatusBadRequest)
		return
	}

	id := uuid.NewV4()
	if err != nil {
		writeError(w, fmt.Sprintf("failed to allocate id for user: %v", err), http.StatusInternalServerError)
		return
	}
	u.ID = id.String()

	if len(errors) > 0 {
		log.Printf("error parsing %v as JSON:", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\n\t\"error\": \"malformed request document\"\n\t\"errors\": [\n"))
		errstrs := []string{}
		for _, er := range errors {
			errstrs = append(errstrs, er.Error())
		}
		w.Write([]byte(strings.Join(errstrs, ", ")))
		w.Write([]byte("\t]\n}"))
		return
	}

	err = a.b.AddUser(u)
	if err != nil {
		writeError(w, fmt.Sprintf("failed to persist user: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s/user/%s", a.Prefix, u.ID))
	w.WriteHeader(http.StatusCreated)
}

// DeleteUser handles DELETE requests to /user/{:id}
func (a *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	err := a.b.DeleteUser(id)
	if err != nil {
		if err == backend.ErrDoesNotExist {
			writeError(w, fmt.Sprintf("user with id %s does not exist", id), http.StatusNotFound)
		} else {
			writeError(w, fmt.Sprintf("failed to delete user: %v", err), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//GetUser handles GET requests for /user/{id}
func (a *API) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := a.b.GetUser(id)
	if err != nil {
		if err == backend.ErrDoesNotExist {
			writeError(w, fmt.Sprintf("user with id %s does not exist", id), http.StatusNotFound)
		} else {
			writeError(w, fmt.Sprintf("failed to get user: %v", err), http.StatusInternalServerError)
		}
		return
	}

	entityJSON, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		writeError(w, fmt.Sprintf("failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(entityJSON)
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	errorResponse := &ErrorResponse{
		Error: message,
	}
	errJSON, _ := json.MarshalIndent(errorResponse, "", "\t")
	w.Write(errJSON)
	log.Printf("ERROR: %s", message)
}
