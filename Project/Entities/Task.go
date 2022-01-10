package entities

import (
	"github.com/google/uuid"
	"time"
)

//Task
//	id: private. unique Task id.
//	OwnerId: The Unique id of the Person who own the Task.
//	Status: 'active' or 'done'.
//	DueDate: Date by the format YYYY-MM-DD.
type Task struct {
	id      string
	Title   string
	OwnerId string
	Details string
	Status  Status
	DueDate time.Time
}

//	Generate unique id, and using valid values
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
