package Controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Course string
	Name string `valid:"required"`
	Matric string `valid:"numeric,required"`
}

type Records struct {
	Records []Record
}

func RenderAttendanceForm(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("views/index.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	record := new(Record)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(record, r.PostForm)

	if decodeErr != nil {
		log.Printf("error mapping form data to struct: %v", decodeErr)
	}	

	valid, validationErrorMessage := validateInput(w, r, record)

	if !valid {
		fmt.Fprint(w, validationErrorMessage)
		return
	}
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Record{})
	db.Create(&record)
	fmt.Fprintf(w, "Record submitted successfully");
}

func validateInput(w http.ResponseWriter, r *http.Request, record *Record) (bool, string) {
	valid, validationError := govalidator.ValidateStruct(record)

	if !valid {
		nameError := govalidator.ErrorByField(validationError, "Name")
		matricError := govalidator.ErrorByField(validationError, "Matric")

		if nameError != "" {
			return valid, "Please fill in a valid name"
		}

		if matricError != "" {
			return valid, "Please fill in a valid matric number"
		}
	}


	return valid, "Validation Error"
}

func GetRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var records []Record
	db.Where("Course <> ?", vars["course"]).Find(&records)

	data := Records{
		Records: records,
	}

	parsedTemplates, parseErr := template.ParseFiles("views/records.html")

	if parseErr != nil {
		log.Printf("Error parsing html: %v", parseErr)
	}

	err = parsedTemplates.Execute(w, data)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}