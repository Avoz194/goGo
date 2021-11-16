package model

import (
	db "github.com/Avoz194/goGo/DBHandler"
	ent "github.com/Avoz194/goGo/Entities"
	"time"
)

func AddPerson(name, email string) ent.Person {
	person := ent.CreatePerson(name, email)
	db.AddPerson(person)
	//the db fun also return a person, which one should return?
	return person
}

// returning list of person, should return a list of person in json probably
// need to protect from corruption (only admins)
func GetAllPersons() []ent.Person{
	return db.GetAllPersons()
}

func GetPerson(id string) ent.Person{
	return db.GetPerson(id)
}

func getPersonDetails(id string) ent.Person{
	return GetPerson(id)
	// may get an error*****************************************
}

//need to check how do we get the details (json?)
func SetPersonDetails(id string) ent.Person{


}
//should return error if id not exist?
func RemovePerson(id string){
	db.DeletePerson()
}

func GetPersonTasks(id string) []ent.Task{
	tasksList := []ent.Task{}
	// if person exist*************************************************
	for _, taskid := range getPerson(id).TasksId{
		tasksList = append(tasksList, getTaskDetails(taskid))
	}
	return tasksList
}

func AddNewTask(personId, title , details string, status ent.Status, dueDate time.Time  ) ent.Task{
	//create task with details
	task := CreateTask()
	// if person exist*************************************************
	getPerson(id).AddTask(task)
}

func addTask(task ent.Task) {
	main.tasks = append(main.tasks, task)
}
//RaiseError if no TaskID
func GetTaskDetails(taskId string) ent.Task {
	for _,task:= range main.tasks {
		if task.Id==taskId {
			return task
		}
	}
	return ent.Task{}
}
//added all the needed params
func SetTaskDetails(taskID , title , details string, status ent.Status, dueDate time.Time) ent.Task {
	var task = getTaskDetails(taskID)
}

func RemoveTask(id string) {
	indexToRemove := -1
	for index,task := range main.tasks {
		if task.Id == id
		{
			indexToRemove = index
			break
		}
	}
	if indexToRemove >-1{
		main.tasks[indexToRemove] = main.tasks[len(main.tasks)-1]
	}
}

func GetStatusForTask(taskId string) ent.Status{
	var task = getTaskDetails(taskId)
	return task.Status
}

func GetOwnerForTask(taskId string) string{
	var task = getTaskDetails(taskId)
	return task.OwnerId
}

//Validate Owner ID
func SetTaskOwner(taskId string, ownerID string){
	var task = getTaskDetails(taskId)
	if (getOwner(ownerID)!= -1) {
		task.OwnerId = ownerID
	}
}

func SetTaskStatus(taskId string, status string){
	var task = getTaskDetails(taskId)
	var stat = ent.CreateStatus(status)
	task.Status = stat
}
