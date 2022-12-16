package Controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func RenderCourseForm(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("views/subjects.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func AddCourse(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	course := new(Course)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(course, r.PostForm)

	if decodeErr != nil {
		log.Printf("error mapping form data to struct: %v", decodeErr)
	}

	valid, validationErrMsg := validateCourse(w, r, course)

	if !valid {
		fmt.Fprint(w, validationErrMsg)
		return
	}

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
}

func validateCourse(w http.ResponseWriter, r *http.Request, course *Course) (bool, string) {
	valid, validationErr := govalidator.ValidateStruct(course)

	if !valid {
		nameErr := govalidator.ErrorByField(validationErr, "Name")
		codeErr := govalidator.ErrorByField(validationErr, "Code")

		if nameErr != "" {
			return valid, "Please fill in a valid name"
		}

		if codeErr != "" {
			return valid, "Please fill in a valid Course Code"
		}
	}

	return valid, "Validation error"
}