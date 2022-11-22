package Controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
)

const (
	UserName = "Nelwhix"
	Password = "admin"
)

type User struct {
	Username string `valid:"alpha,required"`
	Password string `valid:"alpha,required"`
}

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("views/login.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	// map form to struct
	r.ParseForm()
	user := new(User)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(user, r.PostForm)

	// server-side validation
	valid, validationErr := govalidator.ValidateStruct(user)
	
	if !valid {
		userNameErr := govalidator.ErrorByField(validationErr, "Username")
		passwordErr := govalidator.ErrorByField(validationErr, "Password")

		if userNameErr != "" {
			log.Printf("Username validation error : %v", userNameErr)
		}

		if passwordErr != "" {
			log.Printf("password validation error: %v", passwordErr)
		}
	}

}