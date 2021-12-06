package DBHandler

import (
	"database/sql"
	"fmt"
	ent "github.com/Avoz194/goGo/Entities"
	mod "github.com/Avoz194/goGo/model"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const IP = "127.0.0.1:3306"
const DATABASE_NAME = "goGODB"
const CREATE_PERSON_TABLE = "CREATE TABLE IF NOT EXISTS Persons(id varchar(50) NOT NULL, name varchar(50), email varchar(50) UNIQUE , progLang varchar(50), PRIMARY KEY (id));"
const CREATE_TASK_STATUS_TABLE = "CREATE TABLE IF NOT EXISTS TaskStatus(id integer NOT NULL, title varchar(50), PRIMARY KEY (id)); "
const CREATE_TASKS_TABLE = "CREATE TABLE IF NOT EXISTS Tasks(id varchar(50) NOT NULL, title varchar(50), ownerId varchar(50) NOT NULL, details varchar(50), statusID integer NOT NULL, dueDate date, PRIMARY KEY (id), CONSTRAINT FK_TaskToOwner FOREIGN KEY (ownerId) REFERENCES Persons(id), CONSTRAINT FK_TaskToStatus FOREIGN KEY (statusID) REFERENCES TaskStatus(id));"

func openConnection() (error,*sql.DB) {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "noMoreTests123!",
		Net:    "tcp",
		Addr:   IP,
		DBName: DATABASE_NAME,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Failed To connect to MySQL")
		return mod.TechnicalFailrue("Failed To connect to MySQL.",err), nil
	}
	return nil,db
}

func CreateDatabase(){
	cfg := mysql.Config{
		User:   "root",
		Passwd: "noMoreTests123!",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
	}
	db, err:= sql.Open("mysql", cfg.FormatDSN())
	//db.Exec("DROP DATABASE " + DATABASE_NAME)
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
	err = insertStatuses() //On Each Exec of the DB, go over current statuses and try to insert them (PK will block if exist)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(CREATE_TASKS_TABLE)
	if err != nil {
		panic(err)
	}
}

func insertStatuses() error{
	err,db := openConnection()
	if err!=nil {
		return err
	}
	defer db.Close()

	for i, statID := range ent.AllStatusIDs{
		q := "INSERT INTO TaskStatus VALUES ( ?, ? ) "
		_, _ = db.Query(q, statID, ent.AllStatuses[i])
	}
	return nil
}

func DeletePerson(person ent.Person) error{
	err,db := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "DELETE FROM Persons WHERE id =?"
	_, err = db.Exec(query, person.GetPersonId())
	if err != nil {
		msg := fmt.Sprintf("Delete Person with id %s", person.GetPersonId())
		return mod.FailedCommitingRequest(msg,err)
	}
	return nil
}

func DeleteTask(task ent.Task) {
	db := openConnection()
	if db==nil {
		return
	}
	defer db.Close()

	query := "DELETE FROM Tasks WHERE id =?"
	_, err := db.Exec(query,task.GetTaskId())
	if err != nil {
		panic(err)
	}
}

func GetPerson(id string) (error,ent.Person) {
	err, db := openConnection()
	if err!= nil {
		return err,ent.Person{}
	}
	defer db.Close()

	var p ent.Person
	activeTasks := 0
	var personID = ""
	err = db.QueryRow("SELECT DISTINCT Persons.*, count(Tasks.id) over (partition by Persons.id) as numOfActiveTasks FROM Persons left join Tasks on Persons.id = Tasks.ownerId AND Tasks.statusID = 1 where Persons.id = ? ",id).Scan(&personID, &p.Name, &p.Email, &p.ProgLang, &activeTasks)
	p.SetActiveTasks(activeTasks)
	p.SetPersonId(personID)

	if err != nil {
		extraDetails :=	fmt.Sprintf("id %s", id)
		return mod.NoSuchEntityError(p, extraDetails,err), ent.Person{}
	}
	return nil,p
}

func GetTask(id string) ent.Task {
	db := openConnection()
	if db==nil {
		return ent.Task{}
	}
	defer db.Close()

	var t ent.Task
	var date string
	var taskID string
	err := db.QueryRow("SELECT id, title, ownerID, details, statusID, dueDate FROM Tasks where id = ?",id).Scan(&taskID, &t.Title, &t.OwnerId, &t.Details, &t.Status, &date)
	if err != nil {
		panic(err)
	}
	t.DueDate = getTime(date)
	t.SetTaskId(taskID)
	return t
}

func AddPerson(p ent.Person) (error){
	err,db := openConnection()
	if err!=nil {
		return err
	}
	defer db.Close()
	q := "INSERT INTO Persons VALUES ( ?, ? ,?, ?) "
	insertResult, err := db.Query(q, p.GetPersonId(), p.Name, p.Email, p.ProgLang)
 	if err != nil {
		extraDetails :=	fmt.Sprintf("email %s",p.Email)
		return mod.EntityAlreadyExists(p,extraDetails,err)
	}
	defer insertResult.Close()

	return nil
}

func AddTask(t ent.Task) ent.Task{
	db := openConnection()
	if db==nil {
		return ent.Task{}
	}
	defer db.Close()
	q := "INSERT INTO Tasks VALUES ( ?, ? ,?, ?, ?, ? )"
	insertResult, err := db.Query(q, t.GetTaskId(), t.Title, t.OwnerId, t.Details, t.Status,t.DueDate.Format("2006-01-02"))
	if err != nil {
		panic(err.Error())
	}
	defer insertResult.Close()

	return GetTask(t.GetTaskId())
}

func UpdateTask(t ent.Task) ent.Task{
	db := openConnection()
	if db==nil {
		return ent.Task{}
	}
	defer db.Close()

	q := "UPDATE Tasks SET title = ? ,ownerID = ?, details = ?, statusID = ?, dueDate = ?  where id = ?"
	updateResult, err := db.Query(q, t.Title, t.OwnerId, t.Details, t.Status,t.DueDate.Format("2006-01-02"), t.GetTaskId())

	if err != nil {
		panic(err)
	}

	defer updateResult.Close()

	var task ent.Task
	var id string
	var date string
 	err = updateResult.Scan(&id, &task.Title, &task.OwnerId, &t.Details, &task.Status, &date)
	if err != nil {
		panic(err)
	}
	task.DueDate = getTime(date)
	task.SetTaskId(id)
	return task
}

func UpdatePerson(p ent.Person) (error){
	err,db := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	q := "UPDATE Persons SET name = ? ,email = ?, progLang = ? where id = ?"
	updateResult, err := db.Query(q, p.Name, p.Email, p.ProgLang,  p.GetPersonId())

	if err != nil {
		extraDetails :=	fmt.Sprintf("email %s",p.Email)
		return mod.EntityAlreadyExists(p,extraDetails,err)
	}

	defer updateResult.Close()
	return nil
}

func GetPersonTasks(p ent.Person) (error,[]ent.Task) {

	err, db := openConnection()
	if err != nil {
		return err, []ent.Task{}
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM Tasks where ownerid = ?",p.GetPersonId())
	if err != nil {
		msg := fmt.Sprintf("Get Tasks for Person with id %s", p.GetPersonId())
		return mod.FailedCommitingRequest(msg,err), []ent.Task{}
	}

	tasksList := []ent.Task{}
	for results.Next() {
		var task ent.Task
		var date string
		var id string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&id, &task.Title, &task.OwnerId, &task.Details, &task.Status, &date)
		if err != nil {
			msg := fmt.Sprintf("Get Tasks for Person with id %s", p.GetPersonId())
			return mod.FailedCommitingRequest(msg,err), []ent.Task{}
		}
		task.DueDate = getTime(date)
		task.SetTaskId(id)
		tasksList = append(tasksList, task)
	}
	return nil,tasksList
}

func GetAllPersons() (error,[]ent.Person) {
	err,db := openConnection()
	if err != nil {
		return err,[]ent.Person{}
	}

	defer db.Close()

	results, err := db.Query("SELECT DISTINCT Persons.*, count(Tasks.id) over (partition by Persons.id) as numOfActiveTasks FROM Persons left join Tasks on Persons.id = Tasks.ownerId AND Tasks.statusID = 1")
	if err != nil {
		msg := "Get All Persons"
		return mod.FailedCommitingRequest(msg,err), []ent.Person{}
	}

	personList := []ent.Person{}
	for results.Next() {
		var person ent.Person
		activeTasks := 0
		var personID string
		// for each row, scan the result into our tag composite object
		person.SetActiveTasks(activeTasks)
		person.SetPersonId(personID)
		err = results.Scan(&personID, &person.Name, &person.Email, &person.ProgLang, &activeTasks)
		if err != nil {
			msg := "Get All Persons"
			return mod.FailedCommitingRequest(msg,err), []ent.Person{}
		}
		personList = append(personList, person)
	}
	return nil,personList
}

func getTime(date string) time.Time{
	dueDateT, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return dueDateT
}
