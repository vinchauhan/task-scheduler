package handler

import (
	"github.com/matryer/way"
	"github.com/vinchauhan/task-scheduler/internal/service"
	"net/http"
)

type handler struct {
	*service.Service
}

func SetUpRoutes(s *service.Service) http.Handler  {
	h := handler{s}
	baseRouter := way.NewRouter()
	baseRouter.HandleFunc("POST", "/task",h.createNewTask)
	baseRouter.HandleFunc("PUT", "/task/complete/:taskId",h.completeTask)
	baseRouter.HandleFunc("PUT", "/task/begin/:taskId",h.beginTask)
	baseRouter.HandleFunc("GET", "/agents",h.getAgents)
	return baseRouter
}

func NewRouter(s *service.Service) http.Handler {
	h := handler{s}
	router := way.NewRouter()
	router.HandleFunc("POST", "/task",h.createNewTask)
	router.HandleFunc("PUT", "/task/complete/:taskId",h.completeTask)
	router.HandleFunc("PUT", "/task/begin/:taskId",h.beginTask)
	router.HandleFunc("GET", "/agents",h.getAgents)
	return router
}
