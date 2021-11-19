package entities

import (
	"github.com/google/uuid"
)

type Person struct {
	Id, Name, Email string
}

func CreatePerson(name, email string) Person {
	id := uuid.New()
	return Person{Id: id.String(), Name: name, Email: email}
}
