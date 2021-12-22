package entities

import (
	"github.com/google/uuid"
)

//	Email: Must be unique.
//	activeTasks: private. number of all the person's active tasks. Gets the number from the DB.
//	id: private. Unique id.
type Person struct {
	Name		string	`json:"name"`
	Email 		string	`json:"email"`
	ProgLang	string	`json:"favoriteProgrammingLanguage"`
	activeTasks	int		`json:"activeTaskCount"` //activeTask as private member
	id			string	`json:"id"`
}

//	Been used only when creating a new Person in the system.
//	Generate unique id, and activeTasks set by default to 0.
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
