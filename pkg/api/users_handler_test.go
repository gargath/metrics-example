package api_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gargath/metrics-example/pkg/api"
	"github.com/gargath/metrics-example/pkg/backend"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The GET user handler", func() {
	Context("with a healthy backend", func() {

		var server *httptest.Server
		var client *http.Client

		BeforeEach(func() {
			router := mux.NewRouter()
			backend, err := backend.NewSqliteBackend(DBNAME)
			Expect(err).NotTo(HaveOccurred())
			_ = backend.DeleteUser("1234-5678-90123")
			err = backend.AddUser(*newUser)
			Expect(err).NotTo(HaveOccurred())
			api := api.NewAPI("/api", backend)
			api.AddRoutes(router)

			client = &http.Client{}
			server = httptest.NewServer(router)
		})

		AfterEach(func() {
			server.Close()
		})

		It("returns valid users", func() {
			req, err := http.NewRequest("GET", server.URL+"/api/user/1234-5678-90123", nil)
			Expect(err).NotTo(HaveOccurred())
			By("allowing GET request for a specific ID")
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			By("returning HTTP status 200")
			Expect(res.StatusCode).To(Equal(200))
			By("setting the correct content type")

			By("returning the user as JSON")

		})
		It("return 404 for invalid user IDs", func() {
			req, err := http.NewRequest("GET", server.URL+"/api/user/1234-5678-90124", nil)
			Expect(err).NotTo(HaveOccurred())
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())

			By("returning HTTP status 404")
			Expect(res.StatusCode).To(Equal(404))
			By("setting the correct content type")
			By("returning an error in the body")
		})
	})

})
