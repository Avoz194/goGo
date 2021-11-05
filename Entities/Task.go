package entities

import (
	"time"
)

type Task struct {
	id      string
	ownerId string
	details string
	status  Status
	dueDate time.Time
}

func (t *Task) setOwner(id string) {
	t.ownerId = id
}

func (t *Task) setStatus(newStatus Status) {
	t.status = newStatus
}

func (t *Task) isDone() {
	return t.status.isDone()
}
