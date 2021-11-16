package main

import (
	"fmt"
	api "github.com/Avoz194/goGo/APIHandler"
	db "github.com/Avoz194/goGo/DBHandler"
	ent "github.com/Avoz194/goGo/entities"

	"os"
)

var persons = []ent.Person{}
var tasks = []ent.Task{}

func main() {

	os.Setenv("MYSQL_DBUSER", "Aviv")
	os.Setenv("MYSQL_DBPASS", "123456goGO")
	db.CreateDatabase()

	api.CreateServer()


	fmt.Println(persons)
}
