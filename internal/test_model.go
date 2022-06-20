package internal

import (
	"reflect"
	"testing"
)

var DefaultTasks = []Task{
	{1, "Pay bills", true},
	{2, "Walk the dog", false},
	{3, "Buy groceries", false},
	{4, "Exercise", true},
}

func TestCRUD(t *testing.T, table TaskTable) {

	t.Run("list all tasks", func(t *testing.T) {
		currentTasks, err := table.RetrieveAll()

		if err != nil {
			t.Errorf("failed to list tasks, %v", err)
		}

		assertListedRows(t, currentTasks, DefaultTasks)
	})

	t.Run("find task by id", func(t *testing.T) {
		task, err := table.FindByID(1)

		if err != nil {
			t.Errorf("failed to retrieve task, %v", err)
		}

		if task != DefaultTasks[0] {
			t.Errorf("got %v want %v", task, DefaultTasks[0])
		}
	})

	t.Run("create a new task", func(t *testing.T) {
		task := Task{Name: "Clean dishes", Completed: false}

		id, err := table.Create(task)

		if err != nil {
			t.Errorf("failed to create task, %v", err)
		}

		if id < 5 {
			t.Errorf("got %v expected greater than 4", id)
		}
	})

	t.Run("delete task", func(t *testing.T) {
		rowsAffected, err := table.Delete(5)

		if err != nil {
			t.Errorf("failed to delete task, %v", err)
		}

		if rowsAffected != 1 {
			t.Errorf("got %v want 1", rowsAffected)
		}
	})

	t.Run("update task", func(t *testing.T) {
		task := DefaultTasks[2]

		task.Name = "Buy medicines"
		task.Completed = true

		rowsAffected, err := table.Update(task)

		if err != nil {
			t.Errorf("failed to update task, %v", err)
		}

		if rowsAffected != 1 {
			t.Errorf("got %v want 1", rowsAffected)
		}
	})

	t.Run("filter completed tasks", func(t *testing.T) {
		tasks, err := table.Filter(true)

		if err != nil {
			t.Errorf("failed to filter tasks, %v", err)
		}

		got := len(tasks)

		if got != 3 {
			t.Errorf("got %v want 3", got)
		}
	})
}

func assertListedRows(t testing.TB, got, want []Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
