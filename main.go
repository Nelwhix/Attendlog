package main

import (
	"os"
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

func init() {
	file, err := os.OpenFile("applogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := log.New(file, "App Log", log.LstdFlags)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Controllers.RenderAttendanceForm).Methods("GET")
	router.HandleFunc("/admin", Controllers.RenderLogin).Methods("GET")
	router.HandleFunc("/admin", Controllers.Login).Methods("POST")
	router.HandleFunc("/", Controllers.SubmitAttendance).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/resources", http.FileServer(http.Dir("resources/"))))

	handlers.CompressHandler(router)
	logger.Printf("Server starting on port %v\n", CONN_PORT)
	err := http.ListenAndServe(CONN_HOST + ":" + CONN_PORT, router)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}