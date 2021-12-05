package entities

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	id      string
	Title   string
	OwnerId string
	Details string
	Status  Status
	DueDate time.Time
}

func CreateTask(title string, ownerID string, details string, status Status, dueDate time.Time) Task {
	id := uuid.New()
	return Task{id: id.String(), Title: title, OwnerId: ownerID,Details: details,Status: status,DueDate: dueDate}
}

func (t *Task) IsDone() bool {
	return t.Status.isDone()
}

func (t *Task) GetTaskId() string {
	return t.id
}

func (t *Task) SetTaskId(taskID string) {
	t.id = taskID
}
