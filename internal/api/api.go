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
	router.Handle("/tasks", http.HandlerFunc(server.retrieveHandler))

	server.Handler = router

	return server, nil
}

func (t *TaskServer) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.crud.RetrieveAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	} else {
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(tasks)
	}
}
