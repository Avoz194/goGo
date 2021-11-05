package main

import (
	ent "./Entities"
	// Task	"goGo/Entities"
	"fmt"
	"math/rand"
	"strconv"
)

var persons = []ent.Person{}

func main() {
	p1 := createPerson("Gil", "ss")
	addPerson(p1)
	fmt.Println(persons)

}

func createPerson(name, email string) ent.Person {
	return ent.Person{Id: strconv.Itoa(rand.Intn(10000)), Name: name, Email: email, Tasks: nil}
}

func addPerson(person ent.Person) {
	persons = append(persons, person)
}
