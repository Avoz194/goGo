package DBHandler

import (
	"database/sql"
	"fmt"
	ent "github.com/Avoz194/goGo/Entities"
	erro "github.com/Avoz194/goGo/GoGoError"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

const IP = "127.0.0.1:3306"
const DATABASE_NAME = "goGODB"
const CREATE_PERSON_TABLE = "CREATE TABLE IF NOT EXISTS Persons(id varchar(50) NOT NULL, name varchar(50), email varchar(50) UNIQUE , progLang varchar(50), PRIMARY KEY (id));"
const CREATE_TASK_STATUS_TABLE = "CREATE TABLE IF NOT EXISTS TaskStatus(id integer NOT NULL, title varchar(50), PRIMARY KEY (id)); "
const CREATE_TASKS_TABLE = "CREATE TABLE IF NOT EXISTS Tasks(id varchar(50) NOT NULL, title varchar(50), ownerId varchar(50) NOT NULL, details varchar(50), statusID integer NOT NULL, dueDate date, PRIMARY KEY (id), CONSTRAINT FK_TaskToOwner FOREIGN KEY (ownerId) REFERENCES Persons(id), CONSTRAINT FK_TaskToStatus FOREIGN KEY (statusID) REFERENCES TaskStatus(id));"

func openConnection() (erro.GoGoError,*sql.DB) {
	cfg := mysql.Config{
		User:   os.Getenv("GOGODBUSER"),
		Passwd: os.Getenv("GOGODBPASS"),
		Net:    "tcp",
		Addr:   IP,
		DBName: DATABASE_NAME,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Failed To connect to MySQL")
		goErr := erro.GoGoError{ErrorNum: erro.TechnicalFailrue, EntityType: erro.GoGoError{}, ErrorOnValue: "", ErrorOnKey: "", AdditionalMsg:"Failed To connect to MySQL.", Err: err}
		return goErr,nil
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		goErr := erro.GoGoError{ErrorNum: erro.TechnicalFailrue, EntityType: erro.GoGoError{}, ErrorOnValue: "", ErrorOnKey: "", AdditionalMsg:"Failed To connect to MySQL.", Err: pingErr}
		return goErr, nil
	}
	return erro.GoGoError{},db
}

func CreateDatabase(){
	cfg := mysql.Config{
		User:   os.Getenv("GOGODBUSER"),
		Passwd: os.Getenv("GOGODBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
	}
	db, err:= sql.Open("mysql", cfg.FormatDSN())
	//db.Exec("DROP DATABASE " + DATABASE_NAME)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + DATABASE_NAME)
	if err != nil {
		panic(err)
	}
	_,err = db.Exec("USE "+ DATABASE_NAME)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(CREATE_PERSON_TABLE)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(CREATE_TASK_STATUS_TABLE)
	if err != nil {
		panic(err)
	}
	goErr := insertStatuses() //On Each Exec of the DB, go over current statuses and try to insert them (PK will block if exist)
	if goErr.GetError() != nil {
		panic(err)
	}
	_, err = db.Exec(CREATE_TASKS_TABLE)
	if err != nil {
		panic(err)
	}
}

func insertStatuses() erro.GoGoError{
	goErr,db := openConnection()
	if goErr.GetError()!=nil {
		return goErr
	}
	defer db.Close()

	for i, statID := range ent.AllStatusIDs{
		q := "INSERT INTO TaskStatus VALUES ( ?, ? ) "
		_, _ = db.Query(q, statID, ent.AllStatuses[i])
	}
	return erro.GoGoError{}
}

func DeletePerson(person ent.Person) erro.GoGoError{
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

func DeleteTask(task ent.Task) erro.GoGoError {
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

func GetPerson(id string) (erro.GoGoError,ent.Person) {
	goErr, db := openConnection()
	if goErr.GetError() != nil {
		return goErr,ent.Person{}
	}
	defer db.Close()

	var p ent.Person
	activeTasks := 0
	var personID = ""
	err := db.QueryRow("SELECT DISTINCT Persons.*, count(Tasks.id) over (partition by Persons.id) as numOfActiveTasks FROM Persons left join Tasks on Persons.id = Tasks.ownerId AND Tasks.statusID = 1 where Persons.id = ? ",id).Scan(&personID, &p.Name, &p.Email, &p.ProgLang, &activeTasks)
	p.SetActiveTasks(activeTasks)
	p.SetPersonId(personID)

	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.NoSuchEntityError, EntityType: p,ErrorOnValue: id, ErrorOnKey: "id", Err: err, AdditionalMsg: ""}
		return goErr, ent.Person{}
	}
	return erro.GoGoError{},p
}

func GetTask(id string) (erro.GoGoError,ent.Task) {
	goErr, db := openConnection()
	if goErr.GetError()!= nil {
		return goErr,ent.Task{}
	}
	defer db.Close()

	var t ent.Task
	var date string
	var taskID string
	err := db.QueryRow("SELECT id, title, ownerID, details, statusID, dueDate FROM Tasks where id = ?",id).Scan(&taskID, &t.Title, &t.OwnerId, &t.Details, &t.Status, &date)
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.NoSuchEntityError, EntityType: t,ErrorOnValue: id, ErrorOnKey: "id", Err: err, AdditionalMsg: ""}
		return goErr, ent.Task{}
	}
	t.DueDate = getTime(date)
	t.SetTaskId(taskID)
	return erro.GoGoError{},t
}

func AddPerson(p ent.Person) (erro.GoGoError){
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

func AddTask(t ent.Task) erro.GoGoError{
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr
	}
	defer db.Close()
	q := "INSERT INTO Tasks VALUES ( ?, ? ,?, ?, ?, ? )"
	insertResult, err := db.Query(q, t.GetTaskId(), t.Title, t.OwnerId, t.Details, t.Status,t.DueDate.Format("2006-01-02"))
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.NoSuchEntityError, EntityType: ent.Person{},ErrorOnValue: t.OwnerId, ErrorOnKey: "id", Err: err, AdditionalMsg: ""}
		return goErr
	}
	defer insertResult.Close()

	return erro.GoGoError{}
}

func UpdateTask(t ent.Task) erro.GoGoError{
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr
	}
	defer db.Close()

	q := "UPDATE Tasks SET title = ? ,ownerID = ?, details = ?, statusID = ?, dueDate = ?  where id = ?"
	updateResult, err := db.Query(q, t.Title, t.OwnerId, t.Details, t.Status,t.DueDate.Format("2006-01-02"), t.GetTaskId())

	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: ent.Task{},ErrorOnValue: t.GetTaskId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Update Task with id '%s'", t.GetTaskId())}
		return goErr
	}
	defer updateResult.Close()
	return erro.GoGoError{}
}

func UpdatePerson(p ent.Person) erro.GoGoError{
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

func GetPersonTasks(p ent.Person, status ent.Status) (erro.GoGoError,[]ent.Task) {
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr, []ent.Task{}
	}
	var results *sql.Rows
	var err error
	defer db.Close()
	if status == ent.UnkownStatus{
		results, err = db.Query("SELECT * FROM Tasks where ownerid = ?",p.GetPersonId())
	}	else{
		results, err = db.Query("SELECT * FROM Tasks where ownerid = ? AND statusID = ?",p.GetPersonId(), status)
	}
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: ent.Person{},ErrorOnValue: p.GetPersonId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Get Tasks for Person with id '%s'", p.GetPersonId())}
		return goErr, []ent.Task{}
	}

	tasksList := []ent.Task{}
	for results.Next() {
		var task ent.Task
		var date string
		var id string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&id, &task.Title, &task.OwnerId, &task.Details, &task.Status, &date)
		if err != nil {
			goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: ent.Person{},ErrorOnValue: p.GetPersonId(), ErrorOnKey: "id", Err: err, AdditionalMsg: fmt.Sprintf("Get Tasks for Person with id '%s'", p.GetPersonId())}
			return goErr, []ent.Task{}
		}
		task.DueDate = getTime(date)
		task.SetTaskId(id)
		tasksList = append(tasksList, task)
	}
	return erro.GoGoError{},tasksList
}

func GetAllPersons() (erro.GoGoError,[]ent.Person) {
	goErr, db := openConnection()
	if goErr.GetError()!=nil {
		return goErr, []ent.Person{}
	}

	defer db.Close()

	results, err := db.Query("SELECT DISTINCT Persons.*, count(Tasks.id) over (partition by Persons.id) as numOfActiveTasks FROM Persons left join Tasks on Persons.id = Tasks.ownerId AND Tasks.statusID = 1")
	if err != nil {
		goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: ent.Person{},ErrorOnValue: "", ErrorOnKey: "", Err: err, AdditionalMsg: "Get All Persons"}
		return goErr, []ent.Person{}
	}

	personList := []ent.Person{}
	for results.Next() {
		var person ent.Person
		activeTasks := 0
		var personID string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&personID, &person.Name, &person.Email, &person.ProgLang, &activeTasks)
		person.SetActiveTasks(activeTasks)
		person.SetPersonId(personID)
		if err != nil {
			goErr = erro.GoGoError{ErrorNum: erro.FailedCommitingRequest, EntityType: ent.Person{},ErrorOnValue: "", ErrorOnKey: "", Err: err, AdditionalMsg: "Get All Persons"}
			return goErr, []ent.Person{}
		}
		personList = append(personList, person)
	}
	return erro.GoGoError{},personList
}

func getTime(date string) time.Time{
	dueDateT, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return dueDateT
}
