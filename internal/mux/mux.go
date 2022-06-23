package mux

import (
	"example/challenges/internal"
	"example/challenges/internal/orm"
	"example/challenges/internal/repository"
	"os"
)

func SelectDBImpl() (internal.TaskTable, error) {
	dbImpl := os.Getenv("DB_IMPL")
	var table internal.TaskTable
	var err error

	if dbImpl == "orm" {
		rep, repErr := orm.NewTaskOrm()
		table = &rep
		err = repErr
	} else {
		rep, repErr := repository.NewTaskCrud()
		table = &rep
		err = repErr
	}

	if err != nil {
		return nil, err
	}

	return table, nil
}
