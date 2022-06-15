package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrudServer(t *testing.T) {
	server, err := NewTaskServer()

	if err != nil {
		t.Errorf("failed to create new task server, %v", err)
	}

	t.Run("retrieve all tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
