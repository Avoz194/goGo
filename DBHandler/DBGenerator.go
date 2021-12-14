package DBHandler

import (
	"database/sql"
	ent "github.com/Avoz194/goGo/Entities"
	erro "github.com/Avoz194/goGo/GoGoError"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const CREATE_PERSON_TABLE = "CREATE TABLE IF NOT EXISTS Persons(id varchar(50) NOT NULL, name varchar(50), email varchar(50) UNIQUE , progLang varchar(50), PRIMARY KEY (id));"
const CREATE_TASK_STATUS_TABLE = "CREATE TABLE IF NOT EXISTS TaskStatus(id integer NOT NULL, title varchar(50), PRIMARY KEY (id)); "
const CREATE_TASKS_TABLE = "CREATE TABLE IF NOT EXISTS Tasks(id varchar(50) NOT NULL, title varchar(50), ownerId varchar(50) NOT NULL, details varchar(50), statusID integer NOT NULL, dueDate date, PRIMARY KEY (id), CONSTRAINT FK_TaskToOwner FOREIGN KEY (ownerId) REFERENCES Persons(id), CONSTRAINT FK_TaskToStatus FOREIGN KEY (statusID) REFERENCES TaskStatus(id));"


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