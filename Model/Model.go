package model

import (
	db "github.com/Avoz194/goGo/DBHandler"
	ent "github.com/Avoz194/goGo/Entities"
	goErr "github.com/Avoz194/goGo/GoGoError"
	"time"
)

func AddPerson(name, email, progLang string) (goErr.GoGoError, ent.Person) {
	p := ent.CreatePerson(name, email, progLang)
	err := db.AddPerson(p)
	if err.GetError() != nil{
		return err,ent.Person{}
	}
	return GetPerson(p.GetPersonId())
}

func GetAllPersons() (goErr.GoGoError,[]ent.Person){
	return db.GetAllPersons()
}

func GetPerson(id string) (goErr.GoGoError,ent.Person){
	return db.GetPerson(id)
}

//	If one of the parameters is empty, the original value will remain.
func SetPersonDetails(id, name, email, progLang string) (goErr.GoGoError ,ent.Person){
	goErr,p := GetPerson(id)
	if goErr.GetError() != nil{
		return goErr,ent.Person{}
	}
	//update if not empty
	if email!= ""{
		p.Email = email
	}
	if name!= ""{
		p.Name = name
	}
	if progLang!= ""{
		p.ProgLang = progLang
	}
	goErr = db.UpdatePerson(p)
	if goErr.GetError() != nil{
		return goErr, ent.Person{}
	}
	return GetPerson(p.GetPersonId())
}

//	Before deleting a person, delete all it's tasks.
func RemovePerson(id string) goErr.GoGoError{
	goErr, p := GetPerson(id)
	if goErr.GetError() != nil{
		return goErr
	}
	//delete all tasks.
	goErr, tasks := GetPersonTasks(id, "")
	for _,task:= range tasks{
		RemoveTask(task.GetTaskId())
	}
	return db.DeletePerson(p)
}

//	If status is empty returns 'active' and 'done' tasks.
//	else return all the 'active' or 'done' task according to its value.
func GetPersonTasks(id string, status string) (goErr.GoGoError, []ent.Task){
	err, p := GetPerson(id)
	if err.GetError() != nil{
		return err, []ent.Task{}
	}
	stat := ent.UnkownStatus
	if status != ""{
		stat = ent.CreateStatus(status)
		if stat == ent.UnkownStatus {
			err = goErr.GoGoError{ErrorNum: goErr.InvalidInput, EntityType: ent.Task{}, ErrorOnKey: "task status", ErrorOnValue: status}
			return err, []ent.Task{}
		}
	}
	return db.GetPersonTasks(p,stat)
}

func AddNewTask(personId, title , details string, status string, dueDate string) (goErr.GoGoError,ent.Task){
	err, dueDateT := getTime(dueDate)
	if err.GetError() != nil {
		return err, ent.Task{}
	}
	task := ent.CreateTask(title, personId, details, ent.CreateStatus(status) , dueDateT)
	if task.Status == ent.UnkownStatus {
		err := goErr.GoGoError{ErrorNum: goErr.InvalidInput, EntityType: ent.Task{}, ErrorOnKey: "task status", ErrorOnValue: status}
		return err, ent.Task{}
	}
	err = db.AddTask(task)
	if err.GetError() != nil {
		return err, ent.Task{}
	}
	return GetTaskDetails(task.GetTaskId())
}

//	return Time value, by the format: YYYY-MM-DD.
func getTime(date string) (goErr.GoGoError,time.Time){
	dueDateT, err := time.Parse("2006-01-02", date)
	if err != nil {
		err := goErr.GoGoError{ErrorNum: goErr.InvalidInput, EntityType: ent.Task{}, ErrorOnKey: "task dueDate", ErrorOnValue: date}
		return err, time.Time{}
	}
	return goErr.GoGoError{},dueDateT
}

func GetTaskDetails(taskId string) (goErr.GoGoError, ent.Task) {
	return db.GetTask(taskId)
}

//	If one of the parameters is empty, the original value will remain.
func SetTaskDetails(taskID , title , details string, status string, dueDate string, ownerid string) (goErr.GoGoError, ent.Task) {
	err, t := GetTaskDetails(taskID)
	if err.GetError() != nil{
		return err, ent.Task{}
	}

	if title!= "" {
		t.Title = title
	}
	if details!= "" {
		t.Details = details
	}
	if dueDate!= "" {
		err, t.DueDate = getTime(dueDate)
		if err.GetError() != nil {
			return err, ent.Task{}
		}
	}
	if status!=""{
		var stat = ent.CreateStatus(status)
		if stat == ent.UnkownStatus {
			err = goErr.GoGoError{ErrorNum: goErr.InvalidInput, EntityType: ent.Task{}, ErrorOnKey: "task status", ErrorOnValue: status}
			return err, ent.Task{}
		}
		t.Status = stat
	}
	if ownerid!=""{
		err, _ = GetPerson(ownerid)
		if err.GetError() != nil{
			return err, ent.Task{}
		}
		t.OwnerId = ownerid
	}

	err = db.UpdateTask(t)
	if err.GetError() != nil{
		return err, ent.Task{}
	}
	return GetTaskDetails(taskID)
}

func RemoveTask(id string) goErr.GoGoError {
	err, t := GetTaskDetails(id)
	if err.GetError() != nil{
		return err
	}
	return db.DeleteTask(t)
}

func GetStatusForTask(taskId string) (goErr.GoGoError, ent.Status){
	err, task := GetTaskDetails(taskId)
	if err.GetError() != nil{
		return err, -1
	}
	return goErr.GoGoError{},task.Status
}

func GetOwnerForTask(taskId string) (goErr.GoGoError, string){
	err, task := GetTaskDetails(taskId)
	if err.GetError() != nil{
		return err, ""
	}
	return goErr.GoGoError{}, task.OwnerId
}

//	only if the new ownerId is already exist, update the task.
func SetTaskOwner(taskId string, ownerID string) goErr.GoGoError{
	err, _ := GetPerson(ownerID)
	if err.GetError() != nil{
		return err
	}
	err, task := GetTaskDetails(taskId)
	if err.GetError() != nil{
		return err
	}
	task.OwnerId = ownerID
	return db.UpdateTask(task)
}

func SetTaskStatus(taskId string, status string) goErr.GoGoError{
	err, task := GetTaskDetails(taskId)
	if err.GetError() != nil{
		return err
	}

	var stat = ent.CreateStatus(status)
	if stat == ent.UnkownStatus {
		err = goErr.GoGoError{ErrorNum: goErr.InvalidInput, EntityType: ent.Task{}, ErrorOnKey: "task status", ErrorOnValue: status}
		return err
	}
	task.Status = stat
	return db.UpdateTask(task)
}
