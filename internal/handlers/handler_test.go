package handlers_test

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
	"github.com/vinchauhan/task-scheduler/internal/handlers"
	"github.com/vinchauhan/task-scheduler/internal/mock"
	"github.com/vinchauhan/task-scheduler/internal/services"
)

type Service interface {
	GetService() *services.Service
}

var _ = Describe("Test for Handler", func() {
	It("should return a httpHandler", func() {
		handler := handlers.SetUpRoutes()
	})
})
