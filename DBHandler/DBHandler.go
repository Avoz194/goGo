package DBHandler

import (
	"database/sql"
	ent "github.com/Avoz194/goGo/Entities"
	model "github.com/Avoz194/goGo/Model"
	"github.com/go-sql-driver/mysql"
)

const IP = "127.0.0.1:3306"
const DATABASE_NAME = "goGODB"
const CREATE_PERSON_TABLE = "CREATE TABLE IF NOT EXISTS Persons(id varchar(50) NOT NULL, name varchar(50), email varchar(50), PRIMARY KEY (id,email));"
const CREATE_TASKS_TABLE = "CREATE TABLE IF NOT EXISTS Tasks(id varchar(50) NOT NULL, title varchar(50), ownerId varchar(50) NOT NULL, details varchar(50), statusID int NOT NULL, dueDate date, PRIMARY KEY (id), CONSTRAINT FK_ownerId FOREIGN KEY (ownerId) REFERENCES Persons(id));"
const CREATE_STATUS_TABLE = "CREATE TABLE IF NOT EXISTS Status(id varchar(10) NOT NULL, title varchar(50), PRIMARY KEY (id));"

func openConnection() (*sql.DB,error) {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "noMoreTests123!",
		Net:    "tcp",
		Addr:   IP,
		DBName: DATABASE_NAME,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, model.TechnicalFailrue("",err)
	}
	return db, nil
}

func CreateDatabase(){
	cfg := mysql.Config{
		User:   "root",
		Passwd: "noMoreTests123!",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
	}
	db, err:= sql.Open("mysql", cfg.FormatDSN())

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
	_, err = db.Exec(CREATE_TASKS_TABLE)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(CREATE_STATUS_TABLE)
	if err != nil {
		panic(err)
	}
}

func DeletePerson(person ent.Person) {
	db := openConnection()
	if db==nil {
		return
	}
	defer db.Close()

	query := "DELETE FROM Persons WHERE id =?"
	_, err := db.Exec(query,person.Id)
	if err != nil {
		panic(err)
	}
}

func DeleteTask(task ent.Task) {
	db := openConnection()
	if db==nil {
		return
	}
	defer db.Close()

	query := "DELETE FROM Tasks WHERE id =?"
	_, err := db.Exec(query,task.Id)
	if err != nil {
		panic(err)
	}
}

func GetPerson(id string) ent.Person {
	db := openConnection()
	if db==nil {
		return ent.Person{}
	}
	defer db.Close()

	var p ent.Person

	err := db.QueryRow("SELECT * FROM Persons where id = ?",id).Scan(&p.Id, &p.Name, &p.Email)
	if err != nil {
		panic(err)
	}

	return p
}

func GetTask(id string) ent.Task {
	db := openConnection()
	if db==nil {
		return ent.Task{}
	}
	defer db.Close()

	var t ent.Task

	err := db.QueryRow("SELECT id, title, ownerID, details, statusID, dueDate FROM Tasks where id = ?",id).Scan(&t.Id, &t.Title, &t.OwnerId, &t.Details, &t.Status, &t.DueDate)
	if err != nil {
		panic(err)
	}
	return t
}

func AddPerson(p ent.Person) ent.Person{
	db := openConnection()
	if db==nil {
		return ent.Person{}
	}
	defer db.Close()
	q := "INSERT INTO Persons VALUES ( ?, ? ,? ) "
	insertResult, err := db.Query(q, p.Id, p.Name, p.Email)
 	if err != nil {
		panic(err.Error())
	}
	defer insertResult.Close()


	return GetPerson(p.Id)
}

func AddTask(t ent.Task) ent.Task{
	db := openConnection()
	if db==nil {
		return ent.Task{}
	}
	defer db.Close()
	q := "INSERT INTO Persons VALUES ( ?, ? ,?, ?, ?, ? )"
	insertResult, err := db.Query(q, t.Id, t.Title, t.OwnerId, t.Details, t.Status,t.DueDate)
	if err != nil {
		panic(err.Error())
	}
	defer insertResult.Close()

	return GetTask(t.Id)
}

func UpdateTask(t ent.Task) ent.Task{
	db := openConnection()
	if db==nil {
		return ent.Task{}
	}
	defer db.Close()

	q := "UPDATE Tasks SET title = ? ,ownerID = ?, details = ?, statusID = ?, dueDate = ?  where id = ?"
	updateResult, err := db.Query(q, t.Title, t.OwnerId, t.Details, t.Status,t.DueDate, t.Id)

	if err != nil {
		panic(err)
	}

	defer updateResult.Close()

	var task ent.Task
 	err = updateResult.Scan(&task.Id, &task.Title, &task.OwnerId, &task.Status, &task.DueDate)
	if err != nil {
		panic(err)
	}
	return task
}

func UpdatePerson(p ent.Person) ent.Person{
	db := openConnection()
	if db==nil {
		return ent.Person{}
	}
	defer db.Close()

	q := "UPDATE Persons SET name = ? ,email = ? where id = ?"
	updateResult, err := db.Query(q, p.Name, p.Email, p.Id)

	if err != nil {
		panic(err)
	}

	defer updateResult.Close()

	var person ent.Person
	err = updateResult.Scan(&person.Id, &person.Name, &person.Email)
	if err != nil {
		panic(err)
	}
	return person
}

func GetPersonTasks(p ent.Person) []ent.Task {

	db := openConnection()
	if db==nil {
		return []ent.Task{}
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM Tasks where ownerid = ?",p.Id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	tasksList := []ent.Task{}
	for results.Next() {
		var task ent.Task
		// for each row, scan the result into our tag composite object
		err = results.Scan(&task.Id, &task.Title, &task.OwnerId, &task.Status, &task.DueDate)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		tasksList = append(tasksList, task)
	}
	return tasksList
}

func GetAllPersons() []ent.Person {
	db := openConnection()
	if db==nil {
		return []ent.Person{}
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM Persons")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	personList := []ent.Person{}
	for results.Next() {
		var person ent.Person
		// for each row, scan the result into our tag composite object
		err = results.Scan(&person.Id, &person.Name, &person.Email)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		personList = append(personList, person)
	}
	return personList
}