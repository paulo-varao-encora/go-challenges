package orm

import "testing"

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
