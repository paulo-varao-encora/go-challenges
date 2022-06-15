package api

import (
	"encoding/json"
	"example/challenges/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	t.Run("retrieve all tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getTasksFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertTasks(t, got, defaultTasks)
		assertContentType(t, response, jsonContentType)
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

func assertTasks(t testing.TB, got, want []repository.Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
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
