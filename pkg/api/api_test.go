package api_test

import (
	"fmt"

	"github.com/gargath/metrics-example/pkg/api"
	"github.com/gargath/metrics-example/pkg/backend"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Initializing the API", func() {
	It("adds required routes", func() {

		requiredRoutes := map[string][]string{
			"/api/user":      {"GET", "POST"},
			"/api/user/{id}": {"GET", "DELETE"},
		}

		routes := make(map[string][]string)

		router := mux.NewRouter()
		backend, err := backend.NewSqliteBackend(DBNAME)
		Expect(err).NotTo(HaveOccurred())
		api := api.NewAPI("/api", backend)
		api.AddRoutes(router)

		_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			tmpl, err := route.GetPathTemplate()
			Expect(err).NotTo(HaveOccurred())
			if tmpl == "/api" {
				return nil
			}
			ms, err := route.GetMethods()
			Expect(err).NotTo(HaveOccurred())
			routes[tmpl] = append(routes[tmpl], ms...)
			return nil
		})

		for k, v := range routes {
			if _, ok := requiredRoutes[k]; ok {
				for _, m := range v {
					if contains(requiredRoutes[k], m) {
						requiredRoutes[k] = remove(requiredRoutes[k], m)
					} else {
						fmt.Printf("Warning: API route %s contains spurious method %s\n", k, m)
					}
				}
			} else {
				fmt.Printf("Warning: API contains spurious route %s\n", k)
			}
		}

		for k, v := range requiredRoutes {
			Expect(len(v)).To(Equal(0), "API route %s is missing methods %v", k, v)
		}
	})
	It("uses the correct prefix", func() {

		var found bool

		router := mux.NewRouter()
		backend, err := backend.NewSqliteBackend(DBNAME)
		Expect(err).NotTo(HaveOccurred())
		api := api.NewAPI("/api", backend)
		api.AddRoutes(router)

		_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			tmpl, err := route.GetPathTemplate()
			Expect(err).NotTo(HaveOccurred())
			if tmpl == "/api" {
				found = true
			}
			return nil
		})
		Expect(found).To(BeTrue())
	})
})

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func remove(s []string, e string) []string {
	s2 := []string{}
	for _, v := range s {
		if v != e {
			s2 = append(s2, v)
		}
	}
	return s2
}
