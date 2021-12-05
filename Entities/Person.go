package entities

import (
	"github.com/google/uuid"
)

type Person struct {
	Name		string	`json:"name"`
	Email 		string	`json:"email"`
	ProgLang	string	`json:"favoriteProgrammingLanguage"`
	activeTasks	int		`json:"activeTaskCount"` //activeTask as private member
	id			string	`json:"id"`
}

func CreatePerson(name, email, progLang string) Person {
	id := uuid.New()
	return Person{Name: name, Email: email,ProgLang: progLang,activeTasks: 0, id: id.String()}
}

func (p Person) GetActiveTasks() int {
	return p.activeTasks
}
func (p Person) GetPersonId() string {
	return p.id
}
func (p *Person) SetPersonId(id string) {
	p.id = id
}
func (p *Person) SetActiveTasks(numOfTasks int)  {
	p.activeTasks = numOfTasks
}
