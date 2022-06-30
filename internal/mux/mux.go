package mux

import (
	"example/challenges/internal"
	"example/challenges/internal/orm"
	"example/challenges/internal/repository"
	"os"
)

func GetTable() (internal.TaskTable, error) {
	dbImpl := os.Getenv("DB_IMPL")

	if dbImpl == "orm" {
		table, err := orm.NewTaskTable()
		return &table, err
	}

	table, err := repository.NewTaskTable()
	return &table, err
}
