package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/gargath/metrics-example/pkg/backend"
)

//ListUsers handles GET requests for /user
func (a *API) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := backend.ListUsers()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\n\terror: \"failed to list users: %v \"\n}\n", err)))
		log.Printf("failed to list users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entityJson, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\n\terror: \"failed to marshal json: %v \"\n}\n", err)))
		log.Printf("failed to marshal JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\n\tusers:\n"))
	w.Write(entityJson)
	w.Write([]byte("\n}\n"))
}

// AddUser handles POST requests to /user
func (a *API) AddUser(w http.ResponseWriter, r *http.Request) {
	var u backend.User
	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\n\terror: \"Content Type %s not supported \"\n}\n", r.Header.Get("Content-Type"))))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var errors []error
	err := decoder.Decode(&u)
	if err != nil {
		errors = append(errors, err)
	}
	if u.ID != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("{\n\terror: \"User id must be empty \"\n}\n")))
		return
	}

	id := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\n\terror: \"Failed to allocate id for user: %v \"\n}\n", err)))
		return
	}
	u.ID = id.String()

	if len(errors) > 0 {
		log.Printf("error parsing %v as JSON:", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\n\terror: \"malformed request document\"\n\terrors: [\n"))
		errstrs := []string{}
		for _, er := range errors {
			errstrs = append(errstrs, er.Error())
		}
		w.Write([]byte(strings.Join(errstrs, ", ")))
		w.Write([]byte("\t]\n}"))
		return
	}

	err = backend.AddUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\n\terror: \"Failed to persist user: %v \"\n}\n", err)))
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s/user/%s", a.Prefix, u.ID))
	w.WriteHeader(http.StatusCreated)
}
