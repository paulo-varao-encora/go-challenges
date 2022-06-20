package orm

import (
	"example/challenges/internal"
	"testing"
)

func TestOrm(t *testing.T) {

	orm, err := NewTaskOrm()

	if err != nil {
		t.Errorf("failed to create new task orm, %v", err)
	}

	internal.TestCRUD(t, orm)
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
