package api

/*
	Please run create_tables.sql before running the tests
	in order to make sure they are all going to pass
*/

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

var request *http.Request
var response *httptest.ResponseRecorder

func TestCrudServer(t *testing.T) {
	server, err := NewTaskServer()

	if err != nil {
		t.Errorf("failed to create new task server, %v", err)
	}

	t.Run("list all tasks", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks", nil)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		got := getTasksFromResponse(t, response.Body)
		assertTasks(t, got, defaultTasks)
	})

	t.Run("create a new task and returns its ID", func(t *testing.T) {
		body := newTask(t, "Sample task", true)

		updateRequestAndResponse(t, server, http.MethodPost, "/tasks", body)

		assertStatus(t, response.Code, http.StatusOK)
		assertNewID(t, response.Body.String())
	})

	t.Run("return bad request when creating an empty task", func(t *testing.T) {
		body := newTask(t, "", false) // nil

		updateRequestAndResponse(t, server, http.MethodPost, "/tasks", body)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("get single task by its ID", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks/1", nil)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		task := repository.Task{}
		err := json.Unmarshal(response.Body.Bytes(), &task)

		if err != nil {
			t.Errorf("converting response body to task failed, %v", err)
		}

		assertSingleTask(t, task, defaultTasks[0])
	})

	t.Run("redirect to /tasks case id is empty", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks/", nil)

		taskRedirect(t)
	})

	t.Run("send bad request response case find invalid id", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks/100", nil)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("update single task by its id", func(t *testing.T) {
		body := newTask(t, "Sample task updated", false)

		updateRequestAndResponse(t, server, http.MethodPut, "/tasks/5", body)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("send bad request response case update invalid id", func(t *testing.T) {
		body := newTask(t, "Sample task updated", false)

		updateRequestAndResponse(t, server, http.MethodPut, "/tasks/100", body)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("send bad request response case update having no changes", func(t *testing.T) {
		body := newTask(t, "Sample task updated", false)

		updateRequestAndResponse(t, server, http.MethodPut, "/tasks/5", body)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("send bad request response case method is invalid for /tasks", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodPut, "/tasks", nil)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("send bad request response case method is invalid for /tasks/id", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodPost, "/tasks/1", nil)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("list all completed tasks", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks?completed=true", nil)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		got := getTasksFromResponse(t, response.Body)
		gotSize := len(got)

		if gotSize != 2 {
			t.Errorf("got %v, expected 2", gotSize)
		}
	})

	t.Run("redirect to /tasks case completed is empty", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks?completed=", nil)

		taskRedirect(t)
	})

	t.Run("send bad request response case completed param is invalid", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks?completed=x", nil)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})
}

func taskRedirect(t testing.TB) {
	t.Helper()
	assertStatus(t, response.Code, http.StatusOK)
	assertContentType(t, response, jsonContentType)

	got := getTasksFromResponse(t, response.Body)

	if len(got) == 0 {
		t.Error("no task was found")
	}
}

func newTask(t testing.TB, name string, completed bool) *bytes.Buffer {
	t.Helper()
	newTask := repository.Task{Name: name, Completed: completed}
	body, err := json.Marshal(newTask)

	if err != nil {
		t.Errorf("converting task to json failed, %v", err)
	}

	return bytes.NewBuffer(body)
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

func assertNewID(t testing.TB, got string) {
	t.Helper()
	id, err := strconv.Atoi(got)

	if err != nil {
		t.Errorf("failed to convert string to int, %v", err)
	}

	if id < 5 {
		t.Errorf("got %v expeted greater than 4", got)
	}
}

func assertSingleTask(t testing.TB, got, want repository.Task) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func updateRequestAndResponse(t testing.TB, server *TaskServer, method, url string, body io.Reader) {
	t.Helper()

	var err error
	request, err = http.NewRequest(method, url, body)

	if err != nil {
		t.Errorf("request failed, %v", err)
	}

	response = httptest.NewRecorder()

	server.ServeHTTP(response, request)
}
