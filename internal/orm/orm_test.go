package orm

/*
	Please run create_tables.sql before running the tests
	in order to make sure they are all going to pass
*/

import (
	"example/challenges/internal"
	"testing"
)

func TestOrm(t *testing.T) {

	table, err := NewTaskTable()

	if err != nil {
		t.Errorf("failed to create new task table, %v", err)
	}

	internal.TestCRUD(t, &table)
}

func TestDBConnection(t *testing.T) {

	t.Run("ping db", func(t *testing.T) {
		db, err := NewConnection()

		if err != nil {
			t.Errorf("connection failed, %v", err)
		}

		sqlDB, _ := db.DB()
		pingErr := sqlDB.Ping()
		if pingErr != nil {
			t.Errorf("db ping failed, %v", pingErr)
		}

	})
}
