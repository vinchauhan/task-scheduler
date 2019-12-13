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
	return baseRouter
}
