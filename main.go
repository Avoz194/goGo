package main

import (
	"fmt"
	"math/rand"
	"strconv"

	ent "github.com/yagi1/goGo/entities"
)

var persons = []ent.Person{}

func main() {
	p1 := createPerson("Gilsss", "ss")
	addPerson(p1)
	p2 := createPerson("aviv", "dd")
	addPerson(p2)

	fmt.Println(persons)

	// doest recoganize new field
	t1 := ent.Task{Name: "task1"}

	p1.addTask(t1)

}

func createPerson(name, email string) ent.Person {
	return ent.Person{Id: strconv.Itoa(rand.Intn(10000)), Name: name, Email: email, Tasks: nil}
}

func addPerson(person ent.Person) {
	persons = append(persons, person)
}

// func removePerson(){
// 	persons=
// }
