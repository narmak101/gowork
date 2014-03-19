package gowork

// CustomTask actually handles doing the work that needs to be done by the Worker
type CustomTask func(...interface{}) (interface{}, error)

// Task Interface - All workers should implement this
type Task interface {
    Work(CustomTask, ...interface{}) error
}

type Worker struct {
    Name     string
    WorkData interface{}
}

func (w *Worker) Work(b CustomTask, args ...interface{}) (interface{}, error) {
    return b(args...)
}
