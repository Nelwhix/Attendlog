package Controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

const (
	UserName = "Nelwhix"
	Password = "admin"
)

type User struct {
	Username string `valid:"alpha,required"`
	Password string `valid:"alpha,required"`
}

type Courses struct {
	Courses []Course
}

type Course struct {
	Name string
	Code string
}

var Store *sessions.CookieStore

func init() {
	Store = sessions.NewCookieStore([]byte("secret-key"))
}

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("views/login.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func RenderDashboard(w http.ResponseWriter, r *http.Request) {
	data := Courses{
		Courses: []Course{
			{Name: "Applied Thermodynamics", Code: "MEG315"},
			{Name: "Fluid Dynamics", Code: "MEG313"},
			{Name:"Multivariable Calculus", Code: "GEG311"},
		},
	}

	parsedTemplate, parseErr := template.ParseFiles("views/dashboard.html")
	if parseErr != nil {
		log.Printf("Error parsing html: %v", parseErr)
	}
	err := parsedTemplate.Execute(w, data)

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
	target := "/admin"
	if (user.Username == UserName && user.Password == Password) {
		session, _ := Store.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Save(r, w)
		target = "/dashboard"
	} else {
		fmt.Fprintf(w, "Bad credentials")
		return
	}

	http.Redirect(w, r, target, http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	
	http.Redirect(w, r, "/admin", http.StatusFound)
}