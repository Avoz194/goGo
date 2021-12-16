package APIHandler

import (
	"encoding/json"
	"fmt"
	"github.com/Avoz194/goGo/GoGoError"
	mod "github.com/Avoz194/goGo/Model"
	"github.com/gorilla/mux"
	"net/http"
)

//	get the path, method and the params of the API and according to them use the appropriate func.
//	If the data was invalid, return StatusNotFound.
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

//	Decode the information of the Person from the Json input, and put in inside the PersonHolder.
//	any other field in the Json input except 'Name', 'Email' and 'ProgLang' will be ignored.
//	on success, return StatusCreated-201 code with Location and 'x-Created-Id' headers.
//	on failure, return 'text/plain' with the information of the error.
func addPerson(w http.ResponseWriter, r *http.Request) {
	var holder PersonHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err,p := mod.AddPerson(holder.Name, holder.Email, holder.ProgLang)
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	} else {
		w.Header().Set("Location",fmt.Sprintf("/api/people/%s", p.GetPersonId()))
		w.Header().Set("x-Created-Id", p.GetPersonId())
		w.WriteHeader(http.StatusCreated)
	}
}

//	Get a Person array of all the people in the DB.
//	on success, return StatusOK-200 code, transform the Person array to PersonHolder array and encode it.
//	on failure, return 'text/plain' with the information of the error.
func getPeople(w http.ResponseWriter, r *http.Request) {
	err, people := mod.GetAllPersons()
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(personsToHolders(people))
	}
}

//	Get the Person with the id given in the path.
//	on success, return StatusOK-200 code, transform the Person to PersonHolder and encode it.
//	on failure, return 'text/plain' with the information of the error.
func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err,p := mod.GetPerson(params["id"])
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(personToHolder(p))
	}
}

//	Update a Person with the id given in the path, with the details which decoded from the Json input.
//	any other field in the Json input except 'Name', 'Email' and 'ProgLang' will be ignored.
//	on success, return StatusOK-200 code, transform the Person to PersonHolder and encode it.
//	on failure, return 'text/plain' with the information of the error.
func updatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder PersonHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err, p := mod.SetPersonDetails(params["id"], holder.Name, holder.Email, holder.ProgLang)
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(personToHolder(p))
	}
}

//	Delete the Person with the id given in the path.
//	on success, return StatusOK-200 code.
//	on failure, return 'text/plain' with the information of the error.
func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := mod.RemovePerson(params["id"])
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else{
		w.WriteHeader(http.StatusOK)
	}
}

//	Get the Task array of a Person with id given in the path.
//	if a 'status' parameter was given, return all the Tasks according to the parameter value.
//	if a 'status' parameter was not given, return all Tasks.
//	on success, return StatusOK-200 code, transform the Task array to TaskHolder array and encode it.
//	on failure, return 'text/plain' with the information of the error.
func getPersonTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	err, tasks := mod.GetPersonTasks(params["id"], params["status"])
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasksToHolders(tasks))
	}
}

//	Add a new Task to the Person with the id given in the path.
//	Decode the information of the Task from the Json input, and put in inside the TaskHolder.
//	any other field in the Json input except 'Title', 'Details', 'Status' and 'DueDate' will be ignored.
//	on success, return StatusCreated-201 code with Location and 'x-Created-Id' headers.
//	on failure, return 'text/plain' with the information of the error.
func addNewTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder TaskHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err, t := mod.AddNewTask(params["id"], holder.Title, holder.Details, holder.Status, holder.DueDate)

	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	} else {
		w.Header().Set("Location",fmt.Sprintf("/api/tasks/%s", t.GetTaskId()))
		w.Header().Set("x-Created-Id", t.GetTaskId())
		w.WriteHeader(http.StatusCreated)
	}
}

//	Get the Task with the id given in the path.
//	on success, return StatusOK-200 code, transform the Task to TaskHolder and encode it.
//	on failure, return 'text/plain' with the information of the error.
func getTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, t := mod.GetTaskDetails(params["id"])
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(taskToHolder(t))
	}
}

//	Update a Task with the id given in the path, with the details which decoded from the Json input.
//	any other field in the Json input except 'Title', 'Details', 'Status', 'DueDate' and 'OwnerId' will be ignored.
//	on success, return StatusOK-200 code, transform the Task to TaskHolder and encode it.
//	on failure, return 'text/plain' with the information of the error.
func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder TaskHolder
	json.NewDecoder(r.Body).Decode(&holder)
	err, t := mod.SetTaskDetails(params["id"], holder.Title, holder.Details, holder.Status, holder.DueDate, holder.OwnerId)
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(taskToHolder(t))
	}
}

//	Delete the Task with the id given in the path.
//	on success, return StatusOK-200 code.
//	on failure, return 'text/plain' with the information of the error.
func removeTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := mod.RemoveTask(params["id"])
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	} else{
		w.WriteHeader(http.StatusOK)
	}
}

//	Get the Status of Task with the id given in the path.
//	on success, return StatusOK-200 code, and encode the task's status.
//	on failure, return 'text/plain' with the information of the error.
func getTaskStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, s := mod.GetStatusForTask(params["id"])
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.String())
	}
}

//	Update the Status of the Task with the id given in the path.
//	set the status to the value which decoded from the input.
//	on success, return StatusNoContent-204 code.
//	on failure, return 'text/plain' with the information of the error.
func setTaskStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var holder string
	json.NewDecoder(r.Body).Decode(&holder)
	err := mod.SetTaskStatus(params["id"], holder)
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	} else{
		w.WriteHeader(http.StatusNoContent)
	}
}

//	Get the Person ID who own the Task with the id given in the path.
//	on success, return StatusOK-200 code, and encode the task's OwnerId.
//	on failure, return 'text/plain' with the information of the error.
func getOwnerId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, id := mod.GetOwnerForTask(params["id"])

	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(id)
	}
}

//	Update the OwnerId of the Task with the id given in the path.
//	set the OwnerId to the value which decoded from the input.
//	on success, return StatusNoContent-204 code.
//	on failure, return 'text/plain' with the information of the error.
func setOwner(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ownerID string
	json.NewDecoder(r.Body).Decode(&ownerID)
	err := mod.SetTaskOwner(params["id"], ownerID)
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else{
		w.WriteHeader(http.StatusNoContent)
	}
}

//	For each type of goGoError, return the appropriate http status code.
func getAPIStatusForError(goErr GoGoError.GoGoError) int {
	switch goErr.ErrorNum {
	case GoGoError.NoSuchEntityError:
		return http.StatusNotFound
	case GoGoError.FailedCommitingRequest:
		return http.StatusNotFound
	case GoGoError.EntityAlreadyExists:
		return http.StatusBadRequest
	case GoGoError.TechnicalFailrue:
		return http.StatusNotFound
	case GoGoError.InvalidInput:
		return http.StatusBadRequest
	}
	return http.StatusNotFound
}