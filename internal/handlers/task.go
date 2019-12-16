package handlers

import (
	"encoding/json"
	"github.com/matryer/way"
	"github.com/sanity-io/litter"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Task struct {
}

type TaskInput struct {
	//Id       string
	Priority string
	Skills   []string
	AgentId  string
}

//func (t *Task) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
//	_, _ = resp.Write([]byte(fmt.Sprintf(`{"id": "1","priority":"high", "skills":"skill1, skill2","agentId":"3"}`)))
//}

func (h *handler) createNewTask(w http.ResponseWriter, r *http.Request) {
	var in TaskInput
	var task services.TaskOutput
	defer r.Body.Close()

	//decode the incoming request to json
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	litter.Dump(in)
	litter.Dump(in.Priority)
	litter.Dump(in.Skills)
	task, err := h.CreateTask(r.Context(), in.Priority, in.Skills)

	//check for Application define errors
	if err == services.ErrTaskCannotBeAssigned {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		respondError(w, err)
		return
	}

	respond(w, task, http.StatusOK)

}


func (h *handler) createNewTaskUsingMongo(w http.ResponseWriter, r *http.Request) {
	var in TaskInput
	defer r.Body.Close()

	//decode the incoming request to json
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: Respond with created task in the response for working use case
	task, err := h.CreateTask(r.Context(), in.Priority, in.Skills)

	//check for Application define errors
	if err == services.ErrTaskCannotBeAssigned {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
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
	var out string
	ctx := request.Context()
	taskId := way.Param(ctx, "taskId")

	u, err := h.CompleteTask(ctx, taskId)
	litter.Dump(u)
	if err != nil {
		respondError(writer, err)
		return
	}
	if objID , ok := u.(primitive.ObjectID); ok {
		out = objID.String()
	}
	respond(writer, out, http.StatusOK)
}


func (h *handler) beginTask(writer http.ResponseWriter, request *http.Request) {
	//ctx := request.Context()
	//taskId := way.Param(ctx, "taskId")

	//u, err := h.UpdateTask(ctx, taskId)
	//
	//if err != nil {
	//	respondError(writer, err)
	//	return
	//}

	//respond(writer, u, http.StatusOK)
}
