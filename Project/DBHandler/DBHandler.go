package DBHandler

import (
	"database/sql"
	"fmt"
	"github.com/Avoz194/goGo/Project/Entities"
	erro "github.com/Avoz194/goGo/Project/GoGoError"
	"time"
)

const IP = "127.0.0.1:3306"
const DATABASE_NAME = "goGODB"

//	Returning a FailedCommitingRequest error if the Person's id does not exist.
func DeletePerson(person entities.Person) erro.GoGoError {
	goErr,db := openConnection()
	if goErr.GetError() != nil {
		return goErr
	}
	defer db.Close()

	query := "DELETE FROM Persons WHERE id =?"
	_, err := db.Exec(query, person.GetPersonId())
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, ErrorOnValue: person.GetPersonId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Delete Person with id '%s'", person.GetPersonId())}
		return goErr
	}
	return erro.GoGoError{}
}

//	Returning a FailedCommitingRequest error if the Task's id does not exist.
func DeleteTask(task entities.Task) erro.GoGoError {
	goErr, db := openConnection()
	if goErr.GetError() != nil {
		return goErr
	}
	defer db.Close()

	query := "DELETE FROM Tasks WHERE id =?"
	_, err := db.Exec(query,task.GetTaskId())
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: task,ErrorOnValue: task.GetTaskId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Delete Task with id '%s'", task.GetTaskId())}
		return goErr
	}
	return erro.GoGoError{}
}

//	Returning a NoSuchEntityError error if the Person's id does not exist.
func GetPerson(id string) (erro.GoGoError, entities.Person) {
	goErr, db := openConnection()
	if goErr.GetError() != nil {
		return goErr, entities.Person{}
	}
	defer db.Close()

	var p entities.Person
	activeTasks := 0
	var personID = ""
	err := db.QueryRow("SELECT DISTINCT Persons.*, count(Tasks.id) over (partition by Persons.id) as numOfActiveTasks FROM Persons left join Tasks on Persons.id = Tasks.ownerId AND Tasks.statusID = 1 where Persons.id = ? ",id).Scan(&personID, &p.Name, &p.Email, &p.ProgLang, &activeTasks)
	p.SetActiveTasks(activeTasks)
	p.SetPersonId(personID)

	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.NoSuchEntityError, EntityType: p,ErrorOnValue: id, ErrorOnKey: "id", Err: err, AdditionalMsg: ""}
		return goErr, entities.Person{}
	}
	return erro.GoGoError{},p
}

//	Returning a NoSuchEntityError error if the Task's id does not exist.
func GetTask(id string) (erro.GoGoError, entities.Task) {
	goErr, db := openConnection()
	if goErr.GetError()!= nil {
		return goErr, entities.Task{}
	}
	defer db.Close()

	var t entities.Task
	var date string
	var taskID string
	err := db.QueryRow("SELECT id, title, ownerID, details, statusID, dueDate FROM Tasks where id = ?",id).Scan(&taskID, &t.Title, &t.OwnerId, &t.Details, &t.Status, &date)
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.NoSuchEntityError, EntityType: t,ErrorOnValue: id, ErrorOnKey: "id", Err: err, AdditionalMsg: ""}
		return goErr, entities.Task{}
	}
	_,t.DueDate = getTime(date)
	t.SetTaskId(taskID)
	return erro.GoGoError{},t
}

//	Returning a EntityAlreadyExists error if the Person's id already exist.
func AddPerson(p entities.Person) (erro.GoGoError){
	goErr,db := openConnection()
	if goErr.GetError()!= nil {
		return goErr
	}
	defer db.Close()
	q := "INSERT INTO Persons VALUES ( ?, ? ,?, ?) "
	insertResult, err := db.Query(q, p.GetPersonId(), p.Name, p.Email, p.ProgLang)
 	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.EntityAlreadyExists, EntityType: p,ErrorOnValue: p.Email, ErrorOnKey: "email", Err: err, AdditionalMsg: ""}
		return goErr
	}
	defer insertResult.Close()
	return erro.GoGoError{}
}

//	Returning a NoSuchEntityError error if the ownerId does not exist.
func AddTask(t entities.Task) erro.GoGoError {
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr
	}
	defer db.Close()
	q := "INSERT INTO Tasks VALUES ( ?, ? ,?, ?, ?, ? )"
	insertResult, err := db.Query(q, t.GetTaskId(), t.Title, t.OwnerId, t.Details, t.Status,t.DueDate.Format("2006-01-02"))
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.NoSuchEntityError, EntityType: entities.Person{},ErrorOnValue: t.OwnerId, ErrorOnKey: "id", Err: err, AdditionalMsg: ""}
		return goErr
	}
	defer insertResult.Close()

	return erro.GoGoError{}
}

//	Returning a FailedCommitingRequest error if the Task could not been update.
func UpdateTask(t entities.Task) erro.GoGoError {
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr
	}
	defer db.Close()

	q := "UPDATE Tasks SET title = ? ,ownerID = ?, details = ?, statusID = ?, dueDate = ?  where id = ?"
	updateResult, err := db.Query(q, t.Title, t.OwnerId, t.Details, t.Status,t.DueDate.Format("2006-01-02"), t.GetTaskId())

	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: entities.Task{},ErrorOnValue: t.GetTaskId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Update Task with id '%s'", t.GetTaskId())}
		return goErr
	}
	defer updateResult.Close()
	return erro.GoGoError{}
}

//	Returning a EntityAlreadyExists error if the Person's Email already exist.
func UpdatePerson(p entities.Person) erro.GoGoError {
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr
	}
	defer db.Close()

	q := "UPDATE Persons SET name = ? ,email = ?, progLang = ? where id = ?"
	updateResult, err := db.Query(q, p.Name, p.Email, p.ProgLang,  p.GetPersonId())

	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.EntityAlreadyExists, EntityType: p,ErrorOnValue: p.Email, ErrorOnKey: "email", Err: err, AdditionalMsg: ""}
		return goErr
	}

	defer updateResult.Close()
	return erro.GoGoError{}
}

//	Get all of the 'active' or 'done' tasks of a Person according to the status value.
//	return both 'active' and 'done' tasks of a Person if the status value is 'unknown'
//	first gets the tasks values in 'results' then put each of them in a Task array.
//	return a FailedCommitingRequest error if could not get the Tasks.
func GetPersonTasks(p entities.Person, status entities.Status) (erro.GoGoError,[]entities.Task) {
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr, []entities.Task{}
	}
	var results *sql.Rows
	var err error
	defer db.Close()
	//get the tasks values
	if status == entities.UnknownStatus {
		results, err = db.Query("SELECT * FROM Tasks where ownerid = ?",p.GetPersonId())
	}	else{
		results, err = db.Query("SELECT * FROM Tasks where ownerid = ? AND statusID = ?",p.GetPersonId(), status)
	}
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: entities.Person{},ErrorOnValue: p.GetPersonId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Get Tasks for Person with id '%s'", p.GetPersonId())}
		return goErr, []entities.Task{}
	}
	//put the values in Task array.
	tasksList := []entities.Task{}
	for results.Next() {
		var task entities.Task
		var date string
		var id string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&id, &task.Title, &task.OwnerId, &task.Details, &task.Status, &date)
		if err != nil {
			goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: entities.Person{},ErrorOnValue: p.GetPersonId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Get Tasks for Person with id '%s'", p.GetPersonId())}
			return goErr, []entities.Task{}
		}
		_,task.DueDate = getTime(date)
		task.SetTaskId(id)
		tasksList = append(tasksList, task)
	}
	return erro.GoGoError{},tasksList
}

//	first gets the Persons values in 'results' then put each of them in a Person array.
//	return a FailedCommitingRequest error if could not get the Persons.
func GetAllPersons() (erro.GoGoError,[]entities.Person) {
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr, []entities.Person{}
	}

	defer db.Close()
	//get the persons values
	results, err := db.Query("SELECT DISTINCT Persons.*, count(Tasks.id) over (partition by Persons.id) as numOfActiveTasks FROM Persons left join Tasks on Persons.id = Tasks.ownerId AND Tasks.statusID = 1")
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: entities.Person{},ErrorOnValue: "", ErrorOnKey: "", Err: err, AdditionalMsg: "Get All Persons"}
		return goErr, []entities.Person{}
	}
	// put the values in a Person array.
	personList := []entities.Person{}
	for results.Next() {
		var person entities.Person
		activeTasks := 0
		var personID string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&personID, &person.Name, &person.Email, &person.ProgLang, &activeTasks)
		person.SetActiveTasks(activeTasks)
		person.SetPersonId(personID)
		if err != nil {
			goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: entities.Person{},ErrorOnValue: "", ErrorOnKey: "", Err: err, AdditionalMsg: "Get All Persons"}
			return goErr, []entities.Person{}
		}
		personList = append(personList, person)
	}
	return erro.GoGoError{},personList
}

//	return Time value by the format: YYYY-MM-DD.
//	if the Time format is invalid returning InvalidInput error.
func getTime(date string) (erro.GoGoError,time.Time){
	dueDateT, err := time.Parse("2006-01-02", date)
	if err != nil {
		err := erro.GoGoError{ErrorNum: erro.InvalidInput, EntityType: entities.Task{}, ErrorOnKey: "task dueDate", ErrorOnValue: date}
		return err, time.Time{}
	}
	return erro.GoGoError{},dueDateT
}
