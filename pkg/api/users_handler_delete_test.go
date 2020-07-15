package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/gargath/metrics-example/pkg/api"
	"github.com/gargath/metrics-example/pkg/backend"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The DELETE user handler", func() {
	Context("with a healthy backend", func() {

		var server *httptest.Server
		var client *http.Client
		var dbFileName string

		BeforeEach(func() {
			router := mux.NewRouter()
			dbfile := uuid.NewV4().String()
			dbFileName = dbfile + "_" + DBNAME
			b, err := backend.NewSqliteBackend(dbFileName)
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
			os.Remove(dbFileName)
		})

		It("deletes existing users", func() {
			req, err := http.NewRequest("DELETE", server.URL+"/api/user/1234-5678-90123", nil)
			Expect(err).NotTo(HaveOccurred())
			By("allowing GET request for a specific ID")
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			By("returning HTTP status 200")
			Expect(res.StatusCode).To(Equal(204))
			By("returning an empty body")
			body, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(body)).To(Equal(0), "Expected empty body, got %s", string(body))
			By("actuall having removed the user")
			req2, err := http.NewRequest("GET", server.URL+"/api/user/1234-5678-90124", nil)
			Expect(err).NotTo(HaveOccurred())
			res2, err := client.Do(req2)
			Expect(err).NotTo(HaveOccurred())
			Expect(res2.StatusCode).To(Equal(404))
		})
		It("return 404 for invalid user IDs", func() {
			req, err := http.NewRequest("DELETE", server.URL+"/api/user/1234-5678-90126", nil)
			Expect(err).NotTo(HaveOccurred())
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			By("returning HTTP status 404")
			Expect(res.StatusCode).To(Equal(404))
			By("setting the correct content type")
			Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))
			By("returning an error in the body")
			var e api.ErrorResponse
			err = json.NewDecoder(res.Body).Decode(&e)
			Expect(err).NotTo(HaveOccurred())
			Expect(e.Error).To(Equal("user with id 1234-5678-90126 does not exist"))
		})
		It("returns 405 for DELETE requests on the users root", func() {
			req, err := http.NewRequest("DELETE", server.URL+"/api/user", nil)
			Expect(err).NotTo(HaveOccurred())
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.StatusCode).To(Equal(405))
		})
	})
	Context("with a broken backend", func() {

		var server *httptest.Server
		var client *http.Client

		BeforeEach(func() {
			router := mux.NewRouter()
			b := NewBrokenBackend()

			api := api.NewAPI("/api", b)
			api.AddRoutes(router)

			client = &http.Client{}
			server = httptest.NewServer(router)
		})

		AfterEach(func() {
			server.Close()
		})

		It("returns an error to the client", func() {
			req, err := http.NewRequest("DELETE", server.URL+"/api/user/1234-5678-90123", nil)
			Expect(err).NotTo(HaveOccurred())
			By("allowing GET request for a specific ID")
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			By("returning HTTP status 500")
			Expect(res.StatusCode).To(Equal(500))
			By("setting the correct content type")
			Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))
			By("returning the error as JSON")
			var errResp api.ErrorResponse
			err = json.NewDecoder(res.Body).Decode(&errResp)
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Contains(errResp.Error, "failed to delete user")).To(BeTrue(), "Expected error to be like 'failed to get user', but got '%s'", errResp.Error)
		})
	})
})
