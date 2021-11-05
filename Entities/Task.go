package entities

type Task struct{} // just for no errors

func (t *Task) isDone() bool {
	return true
}
