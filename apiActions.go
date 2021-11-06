package main

import (
	ent "github.com/Avoz194/goGo/Entities"
	"math/rand"
	"strconv"
)

func createPerson(name, email string) ent.Person {
	return ent.Person{Id: strconv.Itoa(rand.Intn(10000)), Name: name, Email: email, Tasks: nil}
}

func addPerson(person ent.Person) {
	persons = append(persons, person)
}