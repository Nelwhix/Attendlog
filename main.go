package main

import (
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/Nelwhix/mech-attendance/Controllers"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)


func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Controllers.Create).Methods("GET")
	router.HandleFunc("/admin", Controllers.Admin_Create).Methods("GET")
	router.HandleFunc("/", Controllers.Store).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/resources", http.FileServer(http.Dir("resources/"))))

	handlers.CompressHandler(router)
	log.Printf("Server starting on port %v\n", CONN_PORT)
	err := http.ListenAndServe(CONN_HOST + ":" + CONN_PORT, router)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}