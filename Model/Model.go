package model

import (
	"fmt"
	db "github.com/Avoz194/goGo/DBHandler"
	ent "github.com/Avoz194/goGo/Entities"
	"time"
)

func AddPerson(name, email, progLang string) ent.Person {
	person := ent.CreatePerson(name, email, progLang)
	return db.AddPerson(person)
}

// returning list of person, should return a list of person in json probably
// need to protect from corruption (only admins)
func GetAllPersons() []ent.Person{
	return db.GetAllPersons()
}

func GetPerson(id string) ent.Person{
	return db.GetPerson(id)
}

//need to check how do we get the details (json?)
func SetPersonDetails(id, name, email, progLang string) ent.Person{
	p := GetPerson(id)
	// we should decide how we want to make no change (maybe null), and how to delete (maybe "").
	p.Email = email
	p.Name = name
	p.ProgLang = progLang
	return db.UpdatePerson(p)
}
//should return error if id not exist?
func RemovePerson(id string){
	db.DeletePerson(GetPerson(id))
}

func GetPersonTasks(id string) []ent.Task{
	return db.GetPersonTasks(GetPerson(id))
}

func AddNewTask(personId, title , details string, status string, dueDate string) ent.Task{
	dueDateT := getTime(dueDate)
	task := ent.CreateTask(title, personId, details, ent.CreateStatus(status) , dueDateT)
	return db.AddTask(task)
}

func getTime(date string) time.Time{
	dueDateT, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println(err)
	}
	return dueDateT
}

//RaiseError if no TaskID
func GetTaskDetails(taskId string) ent.Task {
	return db.GetTask(taskId)
}
func SetTaskDetails(taskID , title , details string, status string, dueDate string) ent.Task {
	t := GetTaskDetails(taskID)
	t.Title = title
	t.Details = details
	t.Status = ent.CreateStatus(status)
	t.DueDate = getTime(dueDate)
	return db.UpdateTask(t)
}

func RemoveTask(id string) {
	db.DeleteTask(GetTaskDetails(id))
}

func GetStatusForTask(taskId string) ent.Status{
	var task = GetTaskDetails(taskId)
	return task.Status
}

func GetOwnerForTask(taskId string) string{
	var task = GetTaskDetails(taskId)
	return task.OwnerId
}

//Validate Owner ID
func SetTaskOwner(taskId string, ownerID string){
	var task = GetTaskDetails(taskId)
	task.OwnerId = ownerID
	db.UpdateTask(task)
}

func SetTaskStatus(taskId string, status string){
	var task = GetTaskDetails(taskId)
	var stat = ent.CreateStatus(status)
	task.Status = stat
	db.UpdateTask(task)
}
