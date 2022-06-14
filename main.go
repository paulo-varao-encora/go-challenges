package main

import (
	"example/challenges/internal/repository"
	"fmt"
)

func main() {
	dbConn, _ := repository.NewConnection()
	crud := repository.TaskCrud{DbConn: dbConn}

	task := repository.Task{Name: "Sample task", Completed: false}

	// Create
	id, _ := crud.Create(task)
	fmt.Printf("task id: %v\n", id)

	// Retrieve
	tasks, _ := crud.RetrieveAll()
	total := len(tasks)
	fmt.Printf("total tasks: %v\n", total)

	// Update
	task = tasks[total-1]
	task.Name = "Simple task"
	task.Completed = true

	rowsAffected, _ := crud.Update(task)
	fmt.Printf("updated tasks: %v\n", rowsAffected)

	// Delete
	rowsAffected, _ = crud.Delete(id)
	fmt.Printf("deleted tasks: %v\n", rowsAffected)
}
