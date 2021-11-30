package APIHandler

import (
	"encoding/json"
	ent "github.com/Avoz194/goGo/Entities"
	mod "github.com/Avoz194/goGo/Model"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

type PersonHolder struct {
	Name 	string	`json:"name"`
	Email	string	`json:"emails"`
}

type TaskHolder struct {
	Title   string		`json:"title"`
	Details string		`json:"details"`
	DueDate time.Time	`json:"dueDate"`
	Status 	ent.Status	`json:"status"`
}


func CreateServer(){
	server := mux.NewRouter()
	server.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusOK)
		})
	server.HandleFunc("/api/people/", functionHandler).Methods("POST", "GET")
	server.HandleFunc("/api/people/{id}", functionHandler).Methods("GET","PATCH", "DELETE")
	server.HandleFunc("/api/people/{id}/tasks/", functionHandler).Methods("GET", "POST")
	server.HandleFunc("/api/tasks/{id}", functionHandler).Methods("GET", "PATCH", "DELETE")
	server.HandleFunc("/api/tasks/{id}/status", functionHandler).Methods("GET", "PUT")
	server.HandleFunc("/api/tasks/{id}/owner", getOwnerId).Methods("GET", "PUT")
	http.Handle("/", server)

	c :=  cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST","OPTIONS","GET","PATCH","DELETE","PUT", "FETCH"},
		AllowedHeaders: []string{"*"},
	})
	print("\nserver on...")
	log.Fatal(http.ListenAndServe(":8080",c.Handler(server)))
}
func functionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uri := r.RequestURI
	method := r.Method
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
		case "/api/people/{id}":
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
		case "/people/{id}/tasks/":
			{
				if method == "GET" {
					getPersonTasks(w, r)
				} else if method == "POST" {
					addNewTask(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		case "/api/tasks/{id}":
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
		case "/api/tasks/{id}/status":
			{
				if method == "GET" {
					getTaskStatus(w, r)
				} else if method == "PUT" {
					setTaskStatus(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		case "/api/tasks/{id}/owner":
			{
				if method == "GET" {
					getOwnerId(w, r)
				} else if method == "PUT" {
					setOwner(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
	}
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	var holder PersonHolder
	json.NewDecoder(r.Body).Decode(&holder)
	p := mod.AddPerson(holder.Name, holder.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	people := mod.GetAllPersons()
	json.NewEncoder(w).Encode(people)
}

//need to add case of not exist
func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	p := mod.GetPerson(params["id"])
	json.NewEncoder(w).Encode(p)

}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var holder PersonHolder
	json.NewDecoder(r.Body).Decode(&holder)
	p := mod.SetPersonDetails(params["id"], holder.Name, holder.Email)
	json.NewEncoder(w).Encode(p)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	mod.RemovePerson(params["id"])
	// return err in case of failure
}

func getPersonTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	tasks := mod.GetPersonTasks(params["id"])
	json.NewEncoder(w).Encode(tasks)
}

func addNewTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var holder TaskHolder
	json.NewDecoder(r.Body).Decode(&holder)
	t := mod.AddNewTask(params["id"], holder.Title, holder.Details, holder.Status, holder.DueDate)

	json.NewEncoder(w).Encode(t)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	t := mod.GetTaskDetails(params["id"])
	json.NewEncoder(w).Encode(t)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var holder TaskHolder
	json.NewDecoder(r.Body).Decode(&holder)
	t := mod.SetTaskDetails(params["id"], holder.Title, holder.Details, holder.Status, holder.DueDate)
	json.NewEncoder(w).Encode(t)
}

func removeTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	mod.RemoveTask(params["id"])
}

func getTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	s := mod.GetStatusForTask(params["id"])
	json.NewEncoder(w).Encode(s)
}

func setTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var holder string
	json.NewDecoder(r.Body).Decode(&holder)
	mod.SetTaskStatus(params["id"], holder)
}

func getOwnerId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := mod.GetOwnerForTask(params["id"])
	json.NewEncoder(w).Encode(id)
}

func setOwner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var ownerID string
	json.NewDecoder(r.Body).Decode(&ownerID)
	mod.SetTaskOwner(params["id"], ownerID)
}
