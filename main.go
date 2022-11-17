package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

const (
	MEG315 = "Applied Thermodynamics"
	MEG313 = "Fluid Dynamics"
	GEG311 = "Calculus of Several Variable"
	MEG314 = "Numerical Methods"
	CEG311 = "Civil Engineering Technology"
)

type Attendance struct {
	Course string `valid:"alpha,required"`
	Name string `valid:"alpha,required"`
	Matric string `valid:"alpha,required"`
}

func renderAttendanceForm(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("./templates/index.html")
	parsedTemplate.Execute(w, nil)
}

func submitAttendance(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	record := new(Attendance)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(record, r.PostForm)

	if decodeErr != nil {
		log.Printf("error mapping form data to struct: ", decodeErr)
	}	

	valid, validationErrorMessage := validateInput(w, r, record)
}

func validateInput(w http.ResponseWriter, r *http.Request, record *Attendance) (bool, string) {
	valid, validationError := govalidator.ValidateStruct(Attendance)

	if !valid {
		courseError := govalidator.ErrorByField(validationError, "Course")
		nameError := govalidator.ErrorByField(validationError, "Name")
		matricError := govalidator.ErrorByField(validationError, "Matric")
	}
	return valid, "Validation Error"
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", renderAttendanceForm).Methods("GET")
	router.HandleFunc("/", submitAttendance).Methods("POST")
	log.Printf("Server starting on port %v\n", CONN_PORT)
	err := http.ListenAndServe(CONN_HOST + ":" + CONN_PORT, router)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}