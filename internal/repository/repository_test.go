package repository

/*
	Please run create_tables.sql before running the tests
	in order to make sure they are all going to pass
*/

import (
	"reflect"
	"testing"
)

var defaultTasks = []Task{
	{1, "Pay bills", true},
	{2, "Walk the dog", false},
	{3, "Buy groceries", false},
	{4, "Exercise", true},
}

func TestCRUD(t *testing.T) {

	crud, err := NewTaskCrud()

	if err != nil {
		t.Errorf("failed to create new task crud, %v", err)
	}

	t.Run("list all tasks", func(t *testing.T) {
		currentTasks, err := crud.RetrieveAll()

		if err != nil {
			t.Errorf("failed to list tasks, %v", err)
		}

		assertListedRows(t, currentTasks, defaultTasks)
	})

	t.Run("find task by id", func(t *testing.T) {
		task, err := crud.FindByID(1)

		if err != nil {
			t.Errorf("failed to retrieve task, %v", err)
		}

		if task != defaultTasks[0] {
			t.Errorf("got %v want %v", task, defaultTasks[0])
		}
	})

	t.Run("create a new task", func(t *testing.T) {
		task := Task{Name: "Clean dishes", Completed: false}

		id, err := crud.Create(task)

		if err != nil {
			t.Errorf("failed to create task, %v", err)
		}

		if id < 5 {
			t.Errorf("got %v expected greater than 4", id)
		}
	})

	t.Run("delete task", func(t *testing.T) {
		rowsAffected, err := crud.Delete(5)

		if err != nil {
			t.Errorf("failed to delete task, %v", err)
		}

		if rowsAffected != 1 {
			t.Errorf("got %v want 1", rowsAffected)
		}
	})

	t.Run("update task", func(t *testing.T) {
		task := defaultTasks[2]

		task.Name = "Buy medicines"
		task.Completed = true

		rowsAffected, err := crud.Update(task)

		if err != nil {
			t.Errorf("failed to update task, %v", err)
		}

		if rowsAffected != 1 {
			t.Errorf("got %v want 1", rowsAffected)
		}
	})

	t.Run("filter completed tasks", func(t *testing.T) {
		tasks, err := crud.Filter(true)

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

func TestDBConnection(t *testing.T) {

	t.Run("ping db", func(t *testing.T) {
		db, err := NewConnection()

		if err != nil {
			t.Errorf("connection failed, %v", err)
		}

		pingErr := db.Ping()
		if pingErr != nil {
			t.Errorf("db ping failed, %v", pingErr)
		}
	})
}
