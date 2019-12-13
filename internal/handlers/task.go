package handlers

import (
	"encoding/json"
	"github.com/sanity-io/litter"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"net/http"
)

type Task struct {
}

type TaskInput struct {
	Id       string
	Priority string
	Skills   []string
	AgentId  string
}

//func (t *Task) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
//	_, _ = resp.Write([]byte(fmt.Sprintf(`{"id": "1","priority":"high", "skills":"skill1, skill2","agentId":"3"}`)))
//}

func (h *handler) createNewTask(w http.ResponseWriter, r *http.Request) {
	var in TaskInput
	defer r.Body.Close()

	//decode the incoming request to json
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	litter.Dump(in)
	litter.Dump(in.Id)
	litter.Dump(in.Priority)
	litter.Dump(in.Skills)
	err := h.CreateTask(r.Context(), in.Id, in.Priority, in.Skills)

	//check for Application define errors
	if err == services.ErrTaskCannotBeAssigned {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

}
