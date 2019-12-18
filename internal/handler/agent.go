package handler

import (
	"github.com/vinchauhan/task-scheduler/internal/service"
	"net/http"
)

func (h *handler) getAgents(writer http.ResponseWriter, request *http.Request) {
	var agentsOut []service.Agent

	//Fetch all Agents in the system.
	agentsOut, err := h.GetAgents(request.Context())
	if err == service.ErrFailedToFetchAgents {
		http.Error(writer, service.ErrFailedToFetchAgents.Error(), http.StatusInternalServerError)
		return
	}

	respond(writer, agentsOut, http.StatusOK)
}
