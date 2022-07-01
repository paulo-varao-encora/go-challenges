package repository

import (
	"database/sql"
	"example/challenges/internal"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

// Implements internal.TaskTable interface
type TaskTable struct {
	db *sql.DB
}

// Build new Repository TaskTable so services may
// access database
func NewTaskTable() (TaskTable, error) {
	db, err := NewConnection()

	if err != nil {
		return TaskTable{}, fmt.Errorf("connecting to database failed, %v", err)
	}

	return TaskTable{db}, nil
}

func (c *TaskTable) RetrieveAll() ([]internal.Task, error) {

	rows, err := c.db.Query("SELECT * FROM tasks")

	if err != nil {
		return nil, fmt.Errorf("failed to list all tasks, %v", err)
	}
	defer rows.Close()

	return getTasksFromRows(rows)
}

func (c *TaskTable) FindByID(id int64) (internal.Task, error) {

	row := c.db.QueryRow("SELECT * FROM tasks WHERE ID = ?", id)

	var task internal.Task

	if err := row.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
		return task, fmt.Errorf("failed to convert row into task, %v", err)
	}

	return task, nil
}

func (c *TaskTable) Create(t internal.Task) (int64, error) {

	result, err := c.db.Exec("INSERT INTO tasks (Name, Completed) VALUES (?, ?)",
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

func (c *TaskTable) Delete(id int64) (int64, error) {
	result, err := c.db.Exec("DELETE FROM tasks WHERE ID = ?", id)

	if err != nil {
		return 0, fmt.Errorf("failed to execute delete statement, %v", err)
	}

	return result.RowsAffected()
}

func (c *TaskTable) Update(task internal.Task) (int64, error) {
	result, err := c.db.Exec("UPDATE tasks SET Name=?, Completed=? WHERE ID = ?",
		task.Name, task.Completed, task.ID)

	if err != nil {
		return 0, fmt.Errorf("failed to execute update statement, %v", err)
	}

	return result.RowsAffected()
}

func (c *TaskTable) Filter(completed bool) ([]internal.Task, error) {
	rows, err := c.db.Query("SELECT * FROM tasks WHERE Completed = ?", completed)

	if err != nil {
		return nil, fmt.Errorf("failed to filter tasks, %v", err)
	}
	defer rows.Close()

	return getTasksFromRows(rows)
}

// Run 'source env.sh' in bash to create/update env variables
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

// Convert *sql.Rows to []internal.Task
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
