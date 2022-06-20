package repository

/*
	Please run create_tables.sql before running the tests
	in order to make sure they are all going to pass
*/

import (
	"example/challenges/internal"
	"testing"
)

func TestRepository(t *testing.T) {

	crud, err := NewTaskCrud()

	if err != nil {
		t.Errorf("failed to create new task crud, %v", err)
	}

	internal.TestCRUD(t, crud)
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
