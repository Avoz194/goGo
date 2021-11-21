package entities

import (
	"github.com/google/uuid"
)

type Task struct {
	Id      string
	Title   string
	OwnerId string
	Details string
	Status  Status
	DueDate string
}

func (t *Task) IsDone() bool {
	return t.Status.isDone()
}

func CreateTask(title string, ownerID string, details string, status Status, dueDate string) Task {
	id := uuid.New()
	return Task{Id: id.String(), Title: title, OwnerId: ownerID,Details: details,Status: status,DueDate: dueDate}
}
