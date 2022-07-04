package api

/*
	Please run create_tables.sql before running the tests
	in order to make sure they are all going to pass
*/

import (
	"bytes"
	"encoding/json"
	"example/challenges/internal"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
)

var bearerToken = os.Getenv("API_TOKEN")

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
		assertTasks(t, got, internal.DefaultTasks)
	})

	t.Run("create a new task and returns its ID", func(t *testing.T) {
		body := newTask(t, "Sample task", true)

		updateRequestAndResponse(t, server, http.MethodPost, "/tasks", body)

		assertStatus(t, response.Code, http.StatusOK)
		assertNewID(t, response.Body.String())
	})

	t.Run("get single task by its ID", func(t *testing.T) {
		updateRequestAndResponse(t, server, http.MethodGet, "/tasks/1", nil)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		task := internal.Task{}
		err := json.Unmarshal(response.Body.Bytes(), &task)

		if err != nil {
			t.Errorf("converting response body to task failed, %v", err)
		}

		assertSingleTask(t, task, internal.DefaultTasks[0])
	})

	t.Run("update single task by its id", func(t *testing.T) {
		body := newTask(t, "Sample task updated", false)

		updateRequestAndResponse(t, server, http.MethodPut, "/tasks/5", body)

		assertStatus(t, response.Code, http.StatusOK)
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

	t.Run("get unauthorized status when sending bad token", func(t *testing.T) {
		request, _ = http.NewRequest(http.MethodGet, "/tasks", nil)

		request.Header.Add("authorization", "wrongToken")
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	badRequestTests := map[string]struct {
		method string
		url    string
		body   io.Reader
	}{
		"bad request: create empty task": {http.MethodPost, "/tasks", newTask(t, "", false)},
		"bad request: find invalid id":   {http.MethodGet, "/tasks/100", nil},
		"bad request: update invalid id": {http.MethodPut, "/tasks/100",
			newTask(t, "Sample task updated", false)},
		"bad request: update having no changes": {http.MethodPut, "/tasks/5",
			newTask(t, "Sample task updated", false)},
		"bad request: invalid method for /tasks":    {http.MethodPut, "/tasks", nil},
		"bad request: invalid method for /tasks/id": {http.MethodPost, "/tasks/1", nil},
		"bad request: invalid completed param":      {http.MethodGet, "/tasks?completed=x", nil},
	}

	for name, tc := range badRequestTests {
		t.Run(name, func(t *testing.T) {
			updateRequestAndResponse(t, server, tc.method, tc.url, tc.body)
			assertStatus(t, response.Code, http.StatusBadRequest)
		})
	}

	redirectTests := map[string]string{
		"redirect: id is empty":        "/tasks/",
		"redirect: completed is empty": "/tasks?completed=",
	}

	for name, tc := range redirectTests {
		t.Run(name, func(t *testing.T) {
			updateRequestAndResponse(t, server, http.MethodGet, tc, nil)
			taskRedirect(t)
		})
	}
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
	newTask := internal.Task{Name: name, Completed: completed}
	body, err := json.Marshal(newTask)

	if err != nil {
		t.Errorf("converting task to json failed, %v", err)
	}

	return bytes.NewBuffer(body)
}

func getTasksFromResponse(t testing.TB, body io.Reader) (tasks []internal.Task) {
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

func assertTasks(t testing.TB, got, want []internal.Task) {
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

func assertSingleTask(t testing.TB, got, want internal.Task) {
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

	request.Header.Add("authorization", bearerToken)
	response = httptest.NewRecorder()

	server.ServeHTTP(response, request)
}
