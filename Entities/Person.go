package entities

import (
	"github.com/google/uuid"
)

// func init() {

// }

type Person struct {
	Id, Name, Email string
	TasksId           []string
}

func CreatePerson(name, email string) Person {
	id := uuid.New()
	return Person{Id: id.String(), Name: name, Email: email, TasksId: []string{}}
}

func (p *Person) AddTask(task string) {
	p.TasksId = append(p.TasksId, task)
}

func (p *Person) tasksList() []string {
	return p.TasksId
}

func (p *Person) isAllDone() bool {
	for _, task := range p.TasksId {
		if getTask(task).isDone() == false {
			return false
		}
	}
	return true
}

func (p *Person) getActiveTasks() []string {
	activeTasks := []string{}

	for _, task := range p.TasksId {
		if getTask(task).isDone() == false {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks
}