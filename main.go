package main

import (
	"fmt"
	ent "github.com/Avoz194/goGo/entities"
)

var persons = []ent.Person{}
var tasks = []ent.Task{}

func main() {

	addPerson("Gilsss", "ss")
	addPerson("aviv", "dd")

	fmt.Println(persons)
}
