package repository

import (
	"database/sql"
	"example/challenges/internal"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

type TaskCrud struct {
	DBConn *sql.DB
}

func NewTaskCrud() (TaskCrud, error) {
	DBConn, err := NewConnection()

	if err != nil {
		return TaskCrud{}, fmt.Errorf("connecting to database failed, %v", err)
	}

	return TaskCrud{DBConn}, nil
}

func (c *TaskCrud) RetrieveAll() ([]internal.Task, error) {

	rows, err := c.DBConn.Query("SELECT * FROM tasks")

	if err != nil {
		return nil, fmt.Errorf("failed to list all tasks, %v", err)
	}
	defer rows.Close()

	return getTasksFromRows(rows)
}

func (c *TaskCrud) FindByID(id int64) (internal.Task, error) {

	row := c.DBConn.QueryRow("SELECT * FROM tasks WHERE ID = ?", id)

	var task internal.Task

	if err := row.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
		return task, fmt.Errorf("failed to convert row into task, %v", err)
	}

	return task, nil
}

func (c *TaskCrud) Create(t internal.Task) (int64, error) {

	result, err := c.DBConn.Exec("INSERT INTO tasks (Name, Completed) VALUES (?, ?)",
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
	result, err := c.DBConn.Exec("DELETE FROM tasks WHERE ID = ?", id)

	if err != nil {
		return 0, fmt.Errorf("failed to execute delete statement, %v", err)
	}

	return result.RowsAffected()
}

func (c *TaskCrud) Update(task internal.Task) (int64, error) {
	result, err := c.DBConn.Exec("UPDATE tasks SET Name=?, Completed=? WHERE ID = ?",
		task.Name, task.Completed, task.ID)

	if err != nil {
		return 0, fmt.Errorf("failed to execute update statement, %v", err)
	}

	return result.RowsAffected()
}

func (c *TaskCrud) Filter(completed bool) ([]internal.Task, error) {
	rows, err := c.DBConn.Query("SELECT * FROM tasks WHERE Completed = ?", completed)

	if err != nil {
		return nil, fmt.Errorf("failed to filter tasks, %v", err)
	}
	defer rows.Close()

	return getTasksFromRows(rows)
}

// to create env variables run: source env.sh
func NewConnection() (*sql.DB, error) {
	address := os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   address,
		DBName: "go_challenges",
	}

	return sql.Open("mysql", cfg.FormatDSN())
}

func getTasksFromRows(rows *sql.Rows) ([]internal.Task, error) {
	var tasks []internal.Task

	for rows.Next() {
		var t internal.Task

		if err := rows.Scan(&t.ID, &t.Name, &t.Completed); err != nil {
			return nil, fmt.Errorf("failed to convert row into task, %v", err)
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}
