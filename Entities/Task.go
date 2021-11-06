package entities

import (
	"time"
)

type Task struct {
	id      string
	title   string
	ownerId string
	details string
	status  Status
	dueDate time.Time
}

func (t *Task) setOwner(newOwnerid string) {
	t.ownerId = newOwnerid
}

func (t *Task) setStatus(newStatus Status) {
	t.status = newStatus
}

func (t *Task) isDone() bool {
	return t.status.isDone()
}

func (t *Task) getStatus() Status {
	return t.status
}

func (t *Task) getOwnerID() string {
	return t.ownerId
}
