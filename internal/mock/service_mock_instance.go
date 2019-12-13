package mock

import (
	"github.com/vinchauhan/task-scheduler/internal/services"
	"net/http"
)

type Handler struct {
	SetUpRoutesCall struct{
		Receives struct{
			ServiceInstance *services.Service
		}
		Returns struct{
			Routes http.Handler
			Error error
		}
	}
}

func (m *Handler) SetUpRoutes(serviceInstance *services.Service) (http.Handler, error)  {
	m.SetUpRoutesCall.Receives.ServiceInstance = serviceInstance
	return m.SetUpRoutesCall.Returns.Routes, m.SetUpRoutesCall.Returns.Error
}
