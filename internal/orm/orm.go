package orm

import (
	"example/challenges/internal"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TaskOrm struct {
	DBConn *gorm.DB
}

func NewTaskOrm() (TaskOrm, error) {
	DBConn, err := NewConnection()

	if err != nil {
		return TaskOrm{}, fmt.Errorf("connecting to database failed, %v", err)
	}

	return TaskOrm{DBConn}, nil
}

func (o *TaskOrm) RetrieveAll() ([]internal.Task, error) {

	var tasks []internal.Task
	result := o.DBConn.Find(&tasks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list all tasks, %v", result.Error)
	}

	return tasks, nil
}

func (o *TaskOrm) FindByID(id int64) (internal.Task, error) {
	var task internal.Task
	result := o.DBConn.First(&task, id)

	if result.Error != nil {
		return task, fmt.Errorf("failed to find task by its id, %v", result.Error)
	}

	return task, nil
}

func (o *TaskOrm) Create(task internal.Task) (int64, error) {
	result := o.DBConn.Create(&task)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to create task, %v", result.Error)
	}

	return task.ID, nil
}

func (o *TaskOrm) Delete(id int64) (int64, error) {
	result := o.DBConn.Delete(&internal.Task{}, id)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete task, %v", result.Error)
	}

	return result.RowsAffected, nil
}

func (o *TaskOrm) Update(task internal.Task) (int64, error) {
	_, err := o.FindByID(task.ID)

	if err != nil {
		return 0, err
	}

	result := o.DBConn.Save(&task)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (o *TaskOrm) Filter(completed bool) ([]internal.Task, error) {
	var tasks []internal.Task

	result := o.DBConn.Where("Completed = ?", completed).Find(&tasks)

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
