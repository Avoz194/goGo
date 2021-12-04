package entities

import (
	"github.com/google/uuid"
)

type Person struct {
	Name		string	`json:"name"`
	Email 		string	`json:"email"`
	ProgLang	string	`json:"favoriteProgrammingLanguage"`
	ActiveTasks	int		`json:"activeTaskCount"`
	Id			string	`json:"id"`
}

func CreatePerson(name, email, progLang string) Person {
	id := uuid.New()
	return Person{Name: name, Email: email,ProgLang: progLang,ActiveTasks: 0, Id: id.String()}
}
