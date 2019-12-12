package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Suite")
}

var s *httptest.Server
var serverURL string

var _ = BeforeSuite(func() {
	s := httptest.NewServer(nil)
	Expect(len(s.URL)).To(BeNumerically(">", 0))
	serverURL = s.URL
})

var _ = Describe("the server tests", func() {
	It("should respond to GET /task/{taskid} endpoint and return a task", func() {
		url := fmt.Sprintf("%s/task/1", serverURL)
		resp, err := http.Get(url)
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
	})
})
