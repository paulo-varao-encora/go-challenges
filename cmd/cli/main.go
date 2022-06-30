package main

import (
	"example/challenges/internal"
	"example/challenges/internal/repository"
	"fmt"
)

func main() {
	table, _ := repository.NewTaskTable()

	task := internal.Task{Name: "Sample task", Completed: false}

	// Create
	id, _ := table.Create(task)
	fmt.Printf("task id: %v\n", id)

	// Retrieve
	tasks, _ := table.RetrieveAll()
	total := len(tasks)
	fmt.Printf("total tasks: %v\n", total)

	// Update
	task = tasks[total-1]
	task.Name = "Simple task"
	task.Completed = true

	rowsAffected, _ := table.Update(task)
	fmt.Printf("updated tasks: %v\n", rowsAffected)

	// Delete
	rowsAffected, _ = table.Delete(id)
	fmt.Printf("deleted tasks: %v\n", rowsAffected)
}
