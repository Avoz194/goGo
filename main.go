package main

import (
	"fmt"
	ent "github.com/Avoz194/goGo/entities"
)

var persons = []ent.Person{}
var tasks = []ent.Task{}

func main() {
	p1 := createPerson("Gilsss", "ss")
	addPerson(p1)
	p2 := createPerson("aviv", "dd")
	addPerson(p2)

	fmt.Println(persons)
}


// func removePerson(){
// 	persons=
// }