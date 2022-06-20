package api

import (
	"encoding/json"
	"example/challenges/internal"
	"example/challenges/internal/orm"
	"example/challenges/internal/repository"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const jsonContentType = "application/json"

type TaskServer struct {
	table internal.TaskTable
	http.Handler
}

func NewTaskServer() (*TaskServer, error) {
	server := new(TaskServer)

	dbImpl := os.Getenv("DB_IMPL")
	var table internal.TaskTable
	var err error

	if dbImpl == "orm" {
		rep, repErr := orm.NewTaskOrm()
		table = &rep
		err = repErr
	} else {
		rep, repErr := repository.NewTaskCrud()
		table = &rep
		err = repErr
	}

	if err != nil {
		return nil, err
	}

	server.table = table

	router := http.NewServeMux()
	router.Handle("/tasks", http.HandlerFunc(server.tasksHandler))
	router.Handle("/tasks/", http.HandlerFunc(server.singleTaskHandler))

	server.Handler = router

	return server, nil
}

func (t *TaskServer) tasksHandler(w http.ResponseWriter, r *http.Request) {

	if auth := checkAuthorization(w, r); !auth {
		return
	}

	switch r.Method {
	case http.MethodPost:
		processRequestBodyTask(t, w, r, -1, createTask)
	case http.MethodGet:
		retrieveTasks(t, w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (t *TaskServer) singleTaskHandler(w http.ResponseWriter, r *http.Request) {

	if auth := checkAuthorization(w, r); !auth {
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	if idStr == "" {
		retrieveTasks(t, w, r)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Sprintf("failed to get id, %v", err))
		return
	}

	switch r.Method {
	case http.MethodPut:
		processRequestBodyTask(t, w, r, int64(id), updateTask)
	case http.MethodGet:
		retrieveTaskByID(t, w, int64(id))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func retrieveTasks(t *TaskServer, w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("completed")

	if filter == "" {
		tasks, err := t.table.RetrieveAll()
		sendTasks(w, tasks, err)
		return
	}

	completed, err := strconv.ParseBool(filter)

	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Sprintf("failed to filter param, %v", err))
	} else {
		tasks, err := t.table.Filter(completed)
		sendTasks(w, tasks, err)
	}

}

func createTask(t *TaskServer, w http.ResponseWriter, task internal.Task) {
	id, err := t.table.Create(task)

	if err != nil {
		errorHandler(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to create task in database, %v", err))
	} else {
		fmt.Fprint(w, id)
	}
}

func retrieveTaskByID(t *TaskServer, w http.ResponseWriter, id int64) {
	task, err := t.table.FindByID(id)

	if err != nil && (strings.Contains(err.Error(), "sql: no rows in result set") ||
		strings.Contains(err.Error(), "record not found")) {
		errorHandler(w, http.StatusBadRequest, "invalid task id")
	} else {
		sendTasks(w, task, err)
	}
}

func updateTask(t *TaskServer, w http.ResponseWriter, task internal.Task) {
	rows, err := t.table.Update(task)

	if err != nil && !strings.Contains(err.Error(), "record not found") {
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
	processTask func(t *TaskServer, w http.ResponseWriter, task internal.Task)) {
	var task internal.Task

	if r.Body == nil {
		errorHandler(w, http.StatusBadRequest, "can't process an empty task, body is nil")
		return
	}

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

func sendTasks(w http.ResponseWriter, respBody interface{}, err error) {
	if err != nil {
		errorHandler(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to query database, %v", err))
		return
	}

	w.Header().Set("content-type", jsonContentType)
	encodeErr := json.NewEncoder(w).Encode(respBody)

	if encodeErr != nil {
		errorHandler(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to encode response body to json, %v", encodeErr))
	}

}

func checkAuthorization(w http.ResponseWriter, r *http.Request) bool {

	token := os.Getenv("API_TOKEN")
	requestToken := r.Header.Get("authorization")

	if token != requestToken {
		errorHandler(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	return true
}

func errorHandler(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	fmt.Fprint(w, msg)
}
