package Controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

const (
	UserName = "Nelwhix"
	Password = "admin"
)

type Courses struct {
	Courses []Course
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
	}
	encoded, err := cookieHandler.Encode("session", value)

	if err == nil {
		cookie := &http.Cookie{
			Name: "session",
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(response, cookie)
	}
}

func hasActiveSession(request *http.Request) bool {
	cookie, err := request.Cookie("session")

	if err == nil {
		cookieValue := make(map[string]string)

		err = cookieHandler.Decode("session", cookie.Value, &cookieValue)

		var userName string
		if err == nil {
			userName = cookieValue["username"]
		}
		
		if userName != "" {
			return true
		}
	}
	return false
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name: "session",
		Value: "",
		Path: "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
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
	isActive := hasActiveSession(r)

	if !isActive {
		fmt.Fprintln(w, "You are not authorized to view this page")
		return
	}

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var courses []Course
	db.Find(&courses)

	data := Courses{
		Courses: courses,
	}

	parsedTemplate, parseErr := template.ParseFiles("views/dashboard.html")
	if parseErr != nil {
		log.Printf("Error parsing html: %v", parseErr)
	}
	err = parsedTemplate.Execute(w, data)

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
		setSession(user.Username, w)
		target = "/dashboard"
	} else {
		fmt.Fprintf(w, "Bad credentials")
		return
	}

	http.Redirect(w, r, target, http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/admin", http.StatusFound)
}