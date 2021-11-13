package main

import (
	"fmt"
	ent "github.com/Avoz194/goGo/entities"
	db "github.com/Avoz194/goGo/DBHandler"
"os"
)

var persons = []ent.Person{}
var tasks = []ent.Task{}

func main() {

	os.Setenv("MYSQL_DBUSER", "Aviv")
	os.Setenv("MYSQL_DBPASS", "123456goGO")
	db.CreateDatabase()

	addPerson("Gilsss", "ss")
	addPerson("aviv", "dd")

	fmt.Println(persons)
}
