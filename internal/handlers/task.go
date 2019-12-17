package handlers

import (
	"encoding/json"
	"github.com/matryer/way"
	"github.com/sanity-io/litter"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"net/http"
)

type TaskInput struct {
	Priority string
	Skills   []string
}

func (h *handler) createNewTask(w http.ResponseWriter, r *http.Request) {
	var in TaskInput
	defer r.Body.Close()

	//decode the incoming request to json
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := h.CreateTask(r.Context(), in.Priority, in.Skills)

	//check for Application define errors
	if err == services.ErrTaskDoesNotHaveSkills {
		http.Error(w, services.ErrTaskDoesNotHaveSkills.Error(), http.StatusBadRequest)
		return
	}

	//check if return task is null which means Task couldnt be created
	if task.Id == "" {
		http.Error(w, "No Agents found at this time", http.StatusInternalServerError)
		return
	}
	if err == services.ErrFailedToFindAnyAgentWithSkill || err == services.ErrSystemFindingAgent {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respond(w, task, http.StatusOK)
}


func (h *handler) completeTask(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	taskId := way.Param(ctx, "taskId")

	u, err := h.CompleteTask(ctx, taskId)
	litter.Dump(u)
	if err != nil {
		respondError(writer, err)
		return
	}
	respond(writer, u, http.StatusOK)
}


func (h *handler) beginTask(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	taskId := way.Param(ctx, "taskId")

	u, err := h.BeginTask(ctx, taskId)
	litter.Dump(u)
	if err != nil {
		respondError(writer, err)
		return
	}
	respond(writer, u, http.StatusOK)
}
