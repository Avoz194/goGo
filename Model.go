package main

import (
	ent "github.com/Avoz194/goGo/Entities"
)

func addPerson(name, email string) {
	person := ent.CreatePerson(name, email)
	persons = append(persons, person)
}

// returning list of person, should return a list of person in json probably
// need to protect from corruption (only admins)
func getAllPersons() []ent.Person{
	return persons
}

func getPerson(id string) ent.Person{
	for _, person := range persons {
		if person.Id == id {
			return person
		}
	}
	//persons id not found, should return error **************************
}

func getPersonDetails(id string) ent.Person{
	return getPerson(id)
	// may get an error*****************************************
}

//need to check how do we get the details (json?)
func setPersonDetails(id string){
	for i, person := range persons {
		if person.Id == id {
			persons[i] =
			break
		}
	}
}
//should return error if id not exist?
func removePerson(id string){
	for i, person := range persons {
		if person.Id == id {
			persons[i] = persons[len(persons)-1]
			persons = persons[:len(persons) -1]
		}
	}
}

func getPersonTasks(id string) []ent.Task{
	tasksList := []ent.Task{}
	// if person exist*************************************************
	for _, taskid := range getPerson(id).TasksId{
		tasksList = append(tasksList, getTaskDetails(taskid))
	}
	return tasksList
}

func addNewTask(id string){
	//create task with details
	task := CreateTask()
	// if person exist*************************************************
	getPerson(id).AddTask(task)
}

func addTask(task ent.Task) {
	tasks = append(tasks, task)
}
//RaiseError if no TaskID
func getTaskDetails(taskId string) ent.Task {
	for _,task:= range tasks{
		if task.Id==taskId {
			return task
		}
	}
	return ent.Task{}
}

func setTaskDetails(taskID string) {
	var task = getTaskDetails(taskID)
}

func removeTask(id string) {
	indexToRemove := -1
	for index,task := range tasks {
		if task.Id == id
		{
			indexToRemove = index
			break
		}
	}
	if indexToRemove >0{
		tasks[indexToRemove] = tasks[len(tasks)-1]
	}
}

func getStatusForTask(taskId string) ent.Status{
	var task = getTaskDetails(taskId)
	return task.Status
}

func getOwnerForTask(taskId string) string{
	var task = getTaskDetails(taskId)
	return task.OwnerId
}

//Validate Owner ID
func setTaskOwner(taskId string, ownerID string){
	var task = getTaskDetails(taskId)
	if (getOwner(ownerID)!= -1) {
		task.OwnerId = ownerID
	}
}

func setTaskStatus(taskId string, status string){
	var task = getTaskDetails(taskId)
	var stat = ent.CreateStatus(status)
	task.Status = stat
}
