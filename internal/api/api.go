package api

import (
	"encoding/json"
	"example/challenges/internal/repository"
	"fmt"
	"net/http"
)

const jsonContentType = "application/json"

type TaskServer struct {
	crud *repository.TaskCrud
	http.Handler
}

func NewTaskServer() (*TaskServer, error) {
	server := new(TaskServer)
	crud, err := repository.NewTaskCrud()

	if err != nil {
		return nil, err
	}

	server.crud = crud

	router := http.NewServeMux()
	router.Handle("/tasks", http.HandlerFunc(server.tasksHandler))

	server.Handler = router

	return server, nil
}

func (t *TaskServer) tasksHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		createTask(t, w, r)
	case http.MethodGet:
		retrieveTasks(t, w)
	}
}

func retrieveTasks(t *TaskServer, w http.ResponseWriter) {
	tasks, err := t.crud.RetrieveAll()

	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err.Error())
	} else {
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(tasks)
	}
}

func createTask(t *TaskServer, w http.ResponseWriter, r *http.Request) {
	var task repository.Task

	if r.Body == nil {
		errorHandler(w, http.StatusBadRequest, "can't create an empty task, body is nil")
	} else {
		json.NewDecoder(r.Body).Decode(&task)

		if task.Name == "" {
			errorHandler(w, http.StatusBadRequest, "can't create a nameless task")
		} else {
			id, err := t.crud.Create(task)

			if err != nil {
				errorHandler(w, http.StatusInternalServerError, err.Error())
			} else {
				fmt.Fprint(w, id)
			}
		}
	}
}

func errorHandler(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, msg)
}
