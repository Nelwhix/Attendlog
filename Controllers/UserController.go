package Controllers

import (
	"fmt"
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
	parsedTemplate, _ := template.ParseFiles("../views/login.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func RenderDashboard(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("../views/dashboard.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func readLoginForm(r *http.Request) *User {
	r.ParseForm()
	user := new(User)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(user, r.PostForm)

	if decodeErr != nil {
		log.Printf("error mapping parsed form data to struct : %v", decodeErr)
	}

	return user
}

func validateUser(w http.ResponseWriter, r *http.Request, user *User) (bool, string) {
	valid, validationErr := govalidator.ValidateStruct(user)

	if !valid {
		userNameErr := govalidator.ErrorByField(validationErr, "Username")
		passwordErr := govalidator.ErrorByField(validationErr, "Password")

		if userNameErr != "" {
			log.Printf("Username validation error : %v", userNameErr)
			return valid, "Please fill in a valid Username"
		}

		if passwordErr != "" {
			log.Printf("Password validation error : %v", passwordErr)
			return valid, "Please fill in a valid password"
		}
	}

	return valid, "Validation error"
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := readLoginForm(r)
	valid, validationErr := validateUser(w, r, user)

	if !valid {
		fmt.Fprint(w, validationErr)
		return
	}

	if (user.Username == UserName && user.Password == Password) {
		
	} else {
		fmt.Fprintf(w, "Bad credentials")
	}
}