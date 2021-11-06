package entities

type Task struct {
	Name string
} // just for no errors

func (t *Task) isDone() bool {
	return true
}
