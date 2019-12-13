package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vinchauhan/task-scheduler/internal/handlers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Task", func() {
	It("should return a JSON for created task", func() {
		//taskHandler := handlers.Task{}

		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/xyz", nil)
		h := handlers.New()
		taskHandler := http.HandlerFunc(handlers.CreateNewTask)
		taskHandler.ServeHTTP(resp, req)

		Expect(resp.Code).To(Equal(200))

		respBytes, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(respBytes).To(MatchJSON(`{"id": "1","priority":"high", "skills":"skill1, skill2","agentId":"3"}`))
	})
})
