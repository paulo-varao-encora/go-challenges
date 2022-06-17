package api

import (
	"encoding/json"
	"example/challenges/internal/repository"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	router.Handle("/tasks/", http.HandlerFunc(server.singleTaskHandler))

	server.Handler = router

	return server, nil
}

func (t *TaskServer) tasksHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		processRequestBodyTask(t, w, r, -1, createTask)
	case http.MethodGet:
		retrieveTasks(t, w)
	}
}

func (t *TaskServer) singleTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	if idStr == "" {
		retrieveTasks(t, w)
	} else {
		id, err := strconv.Atoi(idStr)

		if err != nil {
			errorHandler(w, http.StatusBadRequest, fmt.Sprintf("failed to get id, %v", err))
		} else {

			switch r.Method {
			case http.MethodPut:
				processRequestBodyTask(t, w, r, int64(id), updateTask)
			case http.MethodGet:
				retrieveTaskByID(t, w, int64(id))
			}

		}
	}
}

func retrieveTasks(t *TaskServer, w http.ResponseWriter) {
	tasks, err := t.crud.RetrieveAll()

	sendTasks(w, tasks, err)
}

func createTask(t *TaskServer, w http.ResponseWriter, task repository.Task) {
	id, err := t.crud.Create(task)

	if err != nil {
		errorHandler(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to create task in database, %v", err))
	} else {
		fmt.Fprint(w, id)
	}
}

func retrieveTaskByID(t *TaskServer, w http.ResponseWriter, id int64) {
	task, err := t.crud.FindById(id)

	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		errorHandler(w, http.StatusBadRequest, "invalid task id")
	} else {
		sendTasks(w, task, err)
	}
}

func updateTask(t *TaskServer, w http.ResponseWriter, task repository.Task) {
	rows, err := t.crud.Update(task)

	if err != nil {
		errorHandler(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to update task %v in database, %v", task.ID, err))
	} else if rows < 1 {
		errorHandler(w, http.StatusBadRequest,
			fmt.Sprintf("no task was affected: invalid id %v or no changes were detected", task.ID))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func processRequestBodyTask(t *TaskServer, w http.ResponseWriter, r *http.Request, id int64,
	processTask func(t *TaskServer, w http.ResponseWriter, task repository.Task)) {
	var task repository.Task

	if r.Body == nil {
		errorHandler(w, http.StatusBadRequest, "can't process an empty task, body is nil")
	} else {
		err := json.NewDecoder(r.Body).Decode(&task)

		if err != nil {
			errorHandler(w, http.StatusInternalServerError,
				fmt.Sprintf("failed to decode request body to json, %v", err))
		} else if task.Name == "" {
			errorHandler(w, http.StatusBadRequest, "can't process a nameless task")
		} else {
			if id > 0 {
				task.ID = id
			}
			processTask(t, w, task)
		}
	}
}

func sendTasks(w http.ResponseWriter, respBody interface{}, err error) {
	if err != nil {
		errorHandler(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to query database, %v", err))
	} else {
		w.Header().Set("content-type", jsonContentType)
		err := json.NewEncoder(w).Encode(respBody)

		if err != nil {
			errorHandler(w, http.StatusInternalServerError,
				fmt.Sprintf("failed to encode response body to json, %v", err))
		}
	}
}

func errorHandler(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	fmt.Fprint(w, msg)
}
