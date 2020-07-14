package api_test

import (
	"net/http"
	"net/http/httptest"

	uuid "github.com/satori/go.uuid"

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
			dbfile := uuid.NewV4().String()
			b, err := backend.NewSqliteBackend(dbfile + "_" + DBNAME)
			Expect(err).NotTo(HaveOccurred())

			err = b.AddUser(*newUser)
			Expect(err).NotTo(HaveOccurred())
			api := api.NewAPI("/api", b)
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
			Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))
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
			Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))
			By("returning an error in the body")
		})
	})

})
