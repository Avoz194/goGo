package APIHandler

import (
	"encoding/json"
	//ent "github.com/Avoz194/goGo/Entities"
	mod "github.com/Avoz194/goGo/Model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type PersonHolder struct {
	Name 	string	`json:"name"`
	Email	string	`json:"emails"`
}

type TaskHolder struct {
	Title   string		`json:"title"`
	Details string		`json:"details"`
	DueDate string	`json:"dueDate"`
	Status 	string	`json:"status"`
}


func CreateServer() *mux.Router{
	server := mux.NewRouter()
	server.HandleFunc("/api/people/", addPerson).Methods("POST")
	server.HandleFunc("/api/people/", getPeople).Methods("GET")
	server.HandleFunc("/api/people/{id}", getPerson).Methods("GET")
	server.HandleFunc("/api/people/{id}", updatePerson).Methods("PATCH")
	server.HandleFunc("/api/people/{id}", deletePerson).Methods("DELETE")
	server.HandleFunc("/api/people/{id}/tasks/", getPersonTasks).Methods("GET")
	server.HandleFunc("/api/people/{id}/tasks/", addNewTask).Methods("POST")
	server.HandleFunc("/api/tasks/{id}", getTask).Methods("GET")
	server.HandleFunc("/api/tasks/{id}", updateTask).Methods("PATCH")
	server.HandleFunc("/api/tasks/{id}", removeTask).Methods("DELETE")
	server.HandleFunc("/api/tasks/{id}/status", getTaskStatus).Methods("GET")
	server.HandleFunc("/api/tasks/{id}/status", setTaskStatus).Methods("PUT")
	server.HandleFunc("/api/tasks/{id}/owner", getOwnerId).Methods("GET")
	server.HandleFunc("/api/tasks/{id}/owner", setOwner).Methods("PUT")
	print("\nserver on...")


	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers",
		"Accept", "X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "PATCH"})
	http.Handle("/", server)

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(server)))
	print("\nserver on...")

	return server
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	println("in add")
	w.Header().Set("Content-Type", "application/json")
	var holder PersonHolder
	json.NewDecoder(r.Body).Decode(&holder)
	p := mod.AddPerson(holder.Name, holder.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
