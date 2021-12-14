package APIHandler

import (
	"encoding/json"
	"fmt"
	"github.com/Avoz194/goGo/GoGoError"
	mod "github.com/Avoz194/goGo/Model"
	"github.com/gorilla/mux"
	"net/http"
)


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

//need to add case of not exist
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

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := mod.RemovePerson(params["id"])
	// return err in case of failure
	if err.GetError() != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(getAPIStatusForError(err))
		w.Write([]byte (err.Error()))
	}else{
		w.WriteHeader(http.StatusOK)
	}
}

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