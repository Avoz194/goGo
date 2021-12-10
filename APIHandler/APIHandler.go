package APIHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	entities "github.com/Avoz194/goGo/Entities"
	mod "github.com/Avoz194/goGo/Model"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type PersonHolder struct {
	Name 		string	`json:"name"`
	Email		string	`json:"emails"`
	ProgLang	string	`json:"favoriteProgrammingLanguage"`
	ActiveTasks	int		`json:"activeTaskCount"`
	Id			string	`json:"id"`
}

type TaskHolder struct {
	Title   string	`json:"title"`
	Details string	`json:"details"`
	DueDate string	`json:"dueDate"`
	Status 	string	`json:"status"`
	OwnerId string	`json:"ownerId"`
	Id		string	`json:"id"`
}


func CreateServer(){
	server := mux.NewRouter()
	server.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusOK)
		})
	server.HandleFunc("/api/people/", functionHandler).Methods("POST", "GET")
	server.HandleFunc("/api/people/{id}", functionHandler).Methods("GET","PATCH", "DELETE")
	server.HandleFunc("/api/tasks/{id}", functionHandler).Methods("GET", "PATCH", "DELETE")
	server.HandleFunc("/api/tasks/{id}/status", functionHandler).Methods("GET", "PUT")
	server.HandleFunc("/api/tasks/{id}/owner", getOwnerId).Methods("GET", "PUT")

	//Different format for the optional query
	server.Path("/api/people/{id}/tasks/").Queries("status", "{status}").HandlerFunc(functionHandler).Methods("GET")
	server.Path("/api/people/{id}/tasks/").HandlerFunc(functionHandler).Methods("GET", "POST")

	http.Handle("/", server)

	c :=  cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST","OPTIONS","GET","PATCH","DELETE","PUT", "FETCH"},
		AllowedHeaders: []string{"*"},
	})
	println("\nserver on...")
	log.Fatal(http.ListenAndServe(":8080",c.Handler(server)))
}
func functionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uri := r.RequestURI
	method := r.Method
	params := mux.Vars(r)
	switch uri {
		case "/api/people/":
			{
				if method == "POST" {
					addPerson(w, r)
				} else if method == "GET" {
					getPeople(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		case fmt.Sprintf("/api/people/%s", params["id"]):
			{
				if method == "PATCH" {
					updatePerson(w, r)
				} else if method == "DELETE" {
					deletePerson(w, r)
				} else if method == "GET" {
					getPerson(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		case fmt.Sprintf("/api/people/%s/tasks/", params["id"]):
			{
				if method == "GET" {
					getPersonTasks(w, r)
				} else if method == "POST" {
					addNewTask(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		case fmt.Sprintf("/api/people/%s/tasks/?status=%s", params["id"],params["status"]):
			{
				if method == "GET" {
					getPersonTasks(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
	case fmt.Sprintf("/api/tasks/%s", params["id"]):
			{
				if method == "GET" {
					getTask(w, r)
				} else if method == "PATCH" {
					updateTask(w, r)
				} else if method == "DELETE" {
					removeTask(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
	case fmt.Sprintf("/api/tasks/%s/status", params["id"]):
			{
				if method == "GET" {
					getTaskStatus(w, r)
				} else if method == "PUT" {
					setTaskStatus(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
	case fmt.Sprintf("/api/tasks/%s/owner", params["id"]):
			{
				if method == "GET" {
					getOwnerId(w, r)
				} else if method == "PUT" {
					setOwner(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
	default:
		{
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	var holder PersonHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err,p := mod.AddPerson(holder.Name, holder.Email, holder.ProgLang)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte (err.Error()))
	} else {

		w.Header().Set("Location",fmt.Sprintf("/api/people/%s", p.GetPersonId()))
		w.Header().Set("x-Created-Id", p.GetPersonId())

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(personToHolder(p))
	}
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	err, people := mod.GetAllPersons()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(personsToHolders(people))
	}
}

//need to add case of not exist
func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err,p := mod.GetPerson(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(personToHolder(p))
	}
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder PersonHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err, p := mod.SetPersonDetails(params["id"], holder.Name, holder.Email, holder.ProgLang)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(personToHolder(p))
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := mod.RemovePerson(params["id"])
	// return err in case of failure
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	}else{
		w.WriteHeader(http.StatusOK)
	}
}

func getPersonTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	err, tasks := mod.GetPersonTasks(params["id"], params["status"])
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasksToHolders(tasks))
	}
}

func addNewTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder TaskHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err, t := mod.AddNewTask(params["id"], holder.Title, holder.Details, holder.Status, holder.DueDate)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		errDesc := err.Error()
		if strings.Contains(errDesc, "already") {
			w.WriteHeader(http.StatusBadRequest) //If the error is due to already exists, use this status
		} else if strings.Contains(errDesc, "does not") {
			w.WriteHeader(http.StatusNotFound) //If the error is due to doesn't exists, use this status
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte (err.Error()))
	} else {
		w.Header().Set("Location",fmt.Sprintf("/api/tasks/%s", t.GetTaskId()))
		w.Header().Set("x-Created-Id", t.GetTaskId())
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(taskToHolder(t))
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, t := mod.GetTaskDetails(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(taskToHolder(t))
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder TaskHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err, t := mod.SetTaskDetails(params["id"], holder.Title, holder.Details, holder.Status, holder.DueDate, holder.OwnerId)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(taskToHolder(t))
	}
}

func removeTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := mod.RemoveTask(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	} else{
		w.WriteHeader(http.StatusOK)
	}
}

func getTaskStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, s := mod.GetStatusForTask(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.String())
	}
}

func setTaskStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder string
	json.NewDecoder(r.Body).Decode(&holder)
	err := mod.SetTaskStatus(params["id"], holder)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		errDesc := err.Error()
		if strings.Contains(errDesc, "does not") {
			w.WriteHeader(http.StatusNotFound) //If the error is due to doesn't exists, use this status
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte (err.Error()))
	} else{
		w.WriteHeader(http.StatusNoContent)
	}
}

func getOwnerId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, id := mod.GetOwnerForTask(params["id"])

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(id)
	}
}

func setOwner(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ownerID string
	json.NewDecoder(r.Body).Decode(&ownerID)
	err := mod.SetTaskOwner(params["id"], ownerID)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		errDesc := err.Error()
		if strings.Contains(errDesc, "does not") {
			w.WriteHeader(http.StatusNotFound) //If the error is due to doesn't exists, use this status
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte (err.Error()))
	}else{
		w.WriteHeader(http.StatusNoContent)
	}
}

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
	var holders []TaskHolder
	for _,task := range tasks {
		holders = append(holders, taskToHolder(task))
	}
	return holders
}

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
	var holders []PersonHolder
	for _,person := range persons {
		holders = append(holders, personToHolder(person))
	}
	return holders
}