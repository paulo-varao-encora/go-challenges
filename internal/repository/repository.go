package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Task struct {
	ID        int64
	Name      string
	Completed bool
}

type TaskCrud struct {
	DbConn *sql.DB
}

func (c *TaskCrud) RetrieveAll() ([]Task, error) {

	rows, err := c.DbConn.Query("SELECT * FROM tasks")

	if err != nil {
		return nil, fmt.Errorf("failed to list all tasks, %v", err)
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task

		if err := rows.Scan(&t.ID, &t.Name, &t.Completed); err != nil {
			return nil, fmt.Errorf("failed to convert row into task, %v", err)
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (c *TaskCrud) FindById(id int64) (Task, error) {

	row := c.DbConn.QueryRow("SELECT * FROM tasks WHERE ID = ?", id)

	var task Task

	if err := row.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
		return task, fmt.Errorf("failed to convert row into task, %v", err)
	}

	return task, nil
}

func (c *TaskCrud) Create(t Task) (int64, error) {

	result, err := c.DbConn.Exec("INSERT INTO tasks (Name, Completed) VALUES (?, ?)",
		t.Name, t.Completed)

	if err != nil {
		return 0, fmt.Errorf("failed to execute insert statement, %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted id, %v", err)
	}

	return id, nil
}

func (c *TaskCrud) Delete(id int64) (int64, error) {
	result, err := c.DbConn.Exec("DELETE FROM tasks WHERE ID = ?", id)

	if err != nil {
		return 0, fmt.Errorf("failed to execute delete statement, %v", err)
	}

	return result.RowsAffected()
}

func (c *TaskCrud) Update(task Task) (int64, error) {
	result, err := c.DbConn.Exec("UPDATE tasks SET Name=?, Completed=? WHERE ID = ?",
		task.Name, task.Completed, task.ID)

	if err != nil {
		return 0, fmt.Errorf("failed to execute update statement, %v", err)
	}

	return result.RowsAffected()
}

func NewConnection() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "zq1o0m",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "go_challenges",
	}

	return sql.Open("mysql", cfg.FormatDSN())
}
