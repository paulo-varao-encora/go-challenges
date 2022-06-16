package api

import (
	"bytes"
	"encoding/json"
	"example/challenges/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

var defaultTasks = []repository.Task{
	{ID: 1, Name: "Pay bills", Completed: true},
	{ID: 2, Name: "Walk the dog", Completed: false},
	{ID: 3, Name: "Buy groceries", Completed: false},
	{ID: 4, Name: "Exercise", Completed: true},
}

func TestCrudServer(t *testing.T) {
	server, err := NewTaskServer()

	if err != nil {
		t.Errorf("failed to create new task server, %v", err)
	}

	t.Run("list all tasks", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/tasks", nil)

		if err != nil {
			t.Errorf("request failed, %v", err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		got := getTasksFromResponse(t, response.Body)
		assertTasks(t, got, defaultTasks)
	})

	t.Run("create a new task and returns its ID", func(t *testing.T) {
		newTask := repository.Task{Name: "Sample task", Completed: true}
		body, err := json.Marshal(newTask)

		if err != nil {
			t.Errorf("converting task to json failed, %v", err)
		}

		request, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))

		if err != nil {
			t.Errorf("request failed, %v", err)
		}

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertNewId(t, response.Body.String())
	})
}

func getTasksFromResponse(t testing.TB, body io.Reader) (tasks []repository.Task) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&tasks)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Task, '%v'", body, err)
	}

	return
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of application/json, got %v",
			response.Result().Header)
	}
}

func assertTasks(t testing.TB, got, want []repository.Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertNewId(t testing.TB, got string) {
	t.Helper()
	id, err := strconv.Atoi(got)

	if err != nil {
		t.Errorf("failed to convert string to int, %v", err)
	}

	if id < 5 {
		t.Errorf("got %v expeted greater than 4", got)
	}
}
