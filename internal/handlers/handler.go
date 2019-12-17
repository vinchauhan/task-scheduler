package handlers

import (
	"github.com/matryer/way"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"net/http"
)

type handler struct {
	*services.Service
}

func SetUpRoutes(s *services.Service) http.Handler  {
	h := handler{s}
	baseRouter := way.NewRouter()
	baseRouter.HandleFunc("POST", "/task",h.createNewTask)
	baseRouter.HandleFunc("PUT", "/task/complete/:taskId",h.completeTask)
	baseRouter.HandleFunc("PUT", "/task/begin/:taskId",h.beginTask)
	baseRouter.HandleFunc("GET", "/agents",h.getAgents)
	return baseRouter
}
