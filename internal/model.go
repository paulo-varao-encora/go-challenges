package internal

type Task struct {
	ID        int64
	Name      string
	Completed bool
}

// TaskTable defines the database methods. It
// sets the (CRUD) functions that Repository
// and ORM must implement
type TaskTable interface {
	RetrieveAll() ([]Task, error)
	FindByID(id int64) (Task, error)
	Create(t Task) (int64, error)
	Delete(id int64) (int64, error)
	Update(task Task) (int64, error)
	Filter(completed bool) ([]Task, error)
}
