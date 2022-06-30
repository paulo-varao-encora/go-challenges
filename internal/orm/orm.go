package orm

import (
	"example/challenges/internal"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TaskTable struct {
	db *gorm.DB
}

func NewTaskTable() (TaskTable, error) {
	db, err := NewConnection()

	if err != nil {
		return TaskTable{}, fmt.Errorf("connecting to database failed, %v", err)
	}

	return TaskTable{db}, nil
}

func (o *TaskTable) RetrieveAll() ([]internal.Task, error) {

	var tasks []internal.Task
	result := o.db.Find(&tasks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list all tasks, %v", result.Error)
	}

	return tasks, nil
}

func (o *TaskTable) FindByID(id int64) (internal.Task, error) {
	var task internal.Task
	result := o.db.First(&task, id)

	if result.Error != nil {
		return task, fmt.Errorf("failed to find task by its id, %v", result.Error)
	}

	return task, nil
}

func (o *TaskTable) Create(task internal.Task) (int64, error) {
	result := o.db.Create(&task)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to create task, %v", result.Error)
	}

	return task.ID, nil
}

func (o *TaskTable) Delete(id int64) (int64, error) {
	result := o.db.Delete(&internal.Task{}, id)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete task, %v", result.Error)
	}

	return result.RowsAffected, nil
}

func (o *TaskTable) Update(task internal.Task) (int64, error) {
	_, err := o.FindByID(task.ID)

	if err != nil {
		return 0, err
	}

	result := o.db.Save(&task)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (o *TaskTable) Filter(completed bool) ([]internal.Task, error) {
	var tasks []internal.Task

	result := o.db.Where("Completed = ?", completed).Find(&tasks)

	if result.Error != nil {
		return tasks, fmt.Errorf("failed to filter tasks, %v", result.Error)
	}

	return tasks, nil
}

func NewConnection() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/go_challenges?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
