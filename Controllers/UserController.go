package Controllers

import (
	"gorm.io/gorm"
	"net/http"
	"http/template"
)

type User struct {
	gorm.Model
	Name string
	email string
	password string
}

func Admin_Register(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("views/register.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func Admin_Login(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("views/login.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}