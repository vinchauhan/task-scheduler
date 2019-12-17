package services_test

import (
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/vinchauhan/task-scheduler/internal/mock"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"testing"
)

func TestNewClient(t *testing.T) {

	Convey("Setup", t, func() {
		mockService := &ServiceMock{}
		mockService.GetServiceCall.Returns.Service = &services.Service{}

		Convey("Call GetService method will return a service", func() {

		})
	})
}
