package handlers

import (
	"github.com/vinchauhan/task-scheduler/internal/services"
	"net/http"
)

func (h *handler) getAgents(writer http.ResponseWriter, request *http.Request) {
	var agentsOut []services.Agent

	//Fetch all Agents in the system.
	agentsOut, err := h.GetAgents(request.Context())
	if err == services.ErrFailedToFetchAgents {
		http.Error(writer, services.ErrFailedToFetchAgents.Error(), http.StatusInternalServerError)
		return
	}

	respond(writer, agentsOut, http.StatusOK)
}
