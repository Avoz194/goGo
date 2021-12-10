package model

import (
	"fmt"
	db "github.com/Avoz194/goGo/DBHandler"
	ent "github.com/Avoz194/goGo/Entities"
	"time"
)

func AddPerson(name, email, progLang string) (error, ent.Person) {
	p := ent.CreatePerson(name, email, progLang)
	err := db.AddPerson(p)
	if err != nil{
		return err,ent.Person{}
	}
	return GetPerson(p.GetPersonId())
}

// returning list of person, should return a list of person in json probably
// need to protect from corruption (only admins)
func GetAllPersons() (error,[]ent.Person){
	return db.GetAllPersons()
}

func GetPerson(id string) (error,ent.Person){
	return db.GetPerson(id)
}

//need to check how do we get the details (json?)
func SetPersonDetails(id, name, email, progLang string) (error ,ent.Person){
	err,p := GetPerson(id)
	if err != nil{
		return err,ent.Person{}
	}
	// we should decide how we want to make no change (maybe null), and how to delete (maybe "").
	p.Email = email
	p.Name = name
	p.ProgLang = progLang

	err = db.UpdatePerson(p)
	if err != nil{
		return err, ent.Person{}
	}
	return GetPerson(p.GetPersonId())
}
func RemovePerson(id string) error{
	err, p := GetPerson(id)
	if err != nil{
		return err
	}
	return db.DeletePerson(p)
}

func GetPersonTasks(id string) (error, []ent.Task){
	err, p := GetPerson(id)
	if err != nil{
		return err, []ent.Task{}
	}
	return db.GetPersonTasks(p)
}

func AddNewTask(personId, title , details string, status string, dueDate string) (error,ent.Task){
	dueDateT := getTime(dueDate)
	task := ent.CreateTask(title, personId, details, ent.CreateStatus(status) , dueDateT)
	err := db.AddTask(task)
	if err != nil {
		return err, ent.Task{}
	}
	return GetTaskDetails(task.GetTaskId())
}

func getTime(date string) time.Time{
	dueDateT, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println(err)
	}
	return dueDateT
}

//RaiseError if no TaskID
func GetTaskDetails(taskId string) (error, ent.Task) {
	return db.GetTask(taskId)
}
func SetTaskDetails(taskID , title , details string, status string, dueDate string) (error, ent.Task) {
	err, t := GetTaskDetails(taskID)
	if err != nil{
		return err, ent.Task{}
	}
	t.Title = title
	t.Details = details
	t.Status = ent.CreateStatus(status)
	t.DueDate = getTime(dueDate)

	err = db.UpdateTask(t)
	if err != nil{
		return err, ent.Task{}
	}
	return GetTaskDetails(taskID)
}

func RemoveTask(id string) error {
	err, t := GetTaskDetails(id)
	if err != nil{
		return err
	}
	return db.DeleteTask(t)
}

func GetStatusForTask(taskId string) (error, ent.Status){
	err, task := GetTaskDetails(taskId)
	if err != nil{
		return err, -1
	}
	return nil,task.Status
}

func GetOwnerForTask(taskId string) (error, string){
	err, task := GetTaskDetails(taskId)
	if err != nil{
		return err, ""
	}
	return nil, task.OwnerId
}

//Validate Owner ID
func SetTaskOwner(taskId string, ownerID string) error{
	err, task := GetTaskDetails(taskId)
	if err != nil{
		return err
	}
	task.OwnerId = ownerID
	return db.UpdateTask(task)
}

func SetTaskStatus(taskId string, status string) error{
	err, task := GetTaskDetails(taskId)
	if err != nil{
		return err
	}
	var stat = ent.CreateStatus(status)
	task.Status = stat
	return db.UpdateTask(task)
}
