package Controllers

import (
	"net/http"
	"html/template"
	"log"

	"gorm.io/gorm"
	"github.com/gorilla/schema"
)

type Subject struct {
	gorm.Model
	Semester string `valid:"required"`
	Name string	`valid:"required"`
	Code string	`valid:"required"`
}

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
	subject := new(Subject)

}