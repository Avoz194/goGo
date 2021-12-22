package APIHandler

import entities "github.com/Avoz194/goGo/Entities"

//	An intermediate struct between the input from the API to Person (PATCH, POST).
//	in an input from API will fill the Name, Email, ProgLang fields.
//	in an output to API all the fields will be fill.
type PersonHolder struct {
	Name 		string	`json:"name"`
	Email		string	`json:"emails"`
	ProgLang	string	`json:"favoriteProgrammingLanguage"`
	ActiveTasks	int		`json:"activeTaskCount"`
	Id			string	`json:"id"`
}

//	An intermediate struct between the input from the API to Task (PATCH, POST).
//	in an input from API will fill the Title, Details, DueDate, Status fields.
//	in an output to API all the fields will be fill.
type TaskHolder struct {
	Title   string	`json:"title"`
	Details string	`json:"details"`
	DueDate string	`json:"dueDate"`
	Status 	string	`json:"status"`
	OwnerId string	`json:"ownerId"`
	Id		string	`json:"id"`
}

//	Transfer from Task to TaskHolder.
//	fill all the fields and change the DueDate and Status fields to String.
func taskToHolder(task entities.Task) TaskHolder{
	var holder TaskHolder
	holder.Id = task.GetTaskId()
	holder.Title = task.Title
	holder.OwnerId = task.OwnerId
	holder.Details = task.Details
	holder.DueDate = task.DueDate.Format("2006-01-02")
	holder.Status = task.Status.String()
	return holder
}

func tasksToHolders(tasks []entities.Task) []TaskHolder{
	holders:= []TaskHolder{}
	for _,task := range tasks {
		holders = append(holders, taskToHolder(task))
	}
	return holders
}

//	Transfer from Task to PersonHolder.
//	fill all the fields and change the ActiveTasks field to String.
func personToHolder(person entities.Person) PersonHolder{
	var holder PersonHolder
	holder.Id = person.GetPersonId()
	holder.Name = person.Name
	holder.Email = person.Email
	holder.ActiveTasks = person.GetActiveTasks()
	holder.ProgLang = person.ProgLang
	return holder
}

func personsToHolders(persons []entities.Person) []PersonHolder{
	holders := []PersonHolder{}
	for _,person := range persons {
		holders = append(holders, personToHolder(person))
	}
	return holders
}