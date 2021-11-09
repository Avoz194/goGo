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
