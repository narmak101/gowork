package gowork

type BusyWork func(...interface{}) (interface{}, error)

type Work interface {
    Process(BusyWork, ...interface{}) error
}

type Worker struct {
    Name     string
    WorkData interface{}
}

func (w *Worker) Process(b BusyWork, args ...interface{}) (interface{}, error) {
    return b(args...)
}
