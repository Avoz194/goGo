package APIHandler

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

//	Creating an new mux server.
//	define the valid paths and methods of the server.
//	define the cors.
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
	server.HandleFunc("/api/tasks/{id}/owner", functionHandler).Methods("GET", "PUT")

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
