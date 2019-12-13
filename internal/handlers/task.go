package handlers

import "net/http"

type Task struct {

}

type TaskInput struct {
	taskId string
	priority string
	skills []string
}

//func (t *Task) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
//	_, _ = resp.Write([]byte(fmt.Sprintf(`{"id": "1","priority":"high", "skills":"skill1, skill2","agentId":"3"}`)))
//}

func (h *handler) createNewTask(resp http.ResponseWriter, req *http.Request)  {

}
