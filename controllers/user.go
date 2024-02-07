package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Nelwhix/Attendlog/models"
	"github.com/Nelwhix/Attendlog/requests"
	"github.com/Nelwhix/Attendlog/services"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

var decoder = schema.NewDecoder()
var validate = validator.New(validator.WithRequiredStructEnabled())

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFiles("views/login.tmpl")
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}

	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func RenderSignUp(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFiles("views/signup.tmpl")
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}

	err = parsedTemplate.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func RenderSignUpWithData(w http.ResponseWriter, r *http.Request, flashMessage map[string]string) {
	parsedTemplate, err := template.ParseFiles("views/signup.tmpl")
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}

	jsonData, _ := json.Marshal(flashMessage)
	err = parsedTemplate.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"flashMessage":   string(jsonData),
	})
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		flashMessage := map[string]string{
			"type":    "error",
			"message": "error parsing form",
		}
		RenderSignUpWithData(w, r, flashMessage)
		log.Printf("error parsing form: %v", err.Error())
		return
	}

	r.PostForm.Del("gorilla.csrf.Token")

	var signupRequest requests.SignUp
	err = decoder.Decode(&signupRequest, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		flashMessage := map[string]string{
			"type":    "error",
			"message": err.Error(),
		}
		RenderSignUpWithData(w, r, flashMessage)
		log.Printf("error decoding form: %v", err.Error())
		return
	}

	err = validate.Struct(signupRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		flashMessage := map[string]string{
			"type":    "error",
			"message": err.Error(),
		}
		RenderSignUpWithData(w, r, flashMessage)
		log.Printf("validation error: %v", err.Error())
		return
	}

	// insert user
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("./storage/app-%v.db", os.Getenv("APP_ENV"))), &gorm.Config{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
		log.Printf("error opening database: %v", err.Error())
		return
	}
	var newUser models.User
	// Migrate the schema
	db.AutoMigrate(&models.User{})

	newUser = models.User{
		ID:               ulid.Make().String(),
		FirstName:        signupRequest.FirstName,
		LastName:         signupRequest.LastName,
		UserName:         signupRequest.UserName,
		Email:            signupRequest.Email,
		Password:         services.HashPassword(signupRequest.Password),
		SecurityQuestion: signupRequest.SecurityQuestion,
		Answer:           services.HashPassword(signupRequest.Answer),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
		log.Printf("error inserting record: %v", result.Error.Error())
		return
	}

	//token, err := services.GenerateJwt(newUser.ID)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprint(w, "Internal Server Error")
	//	log.Printf("error generating token: %v", err.Error())
	//	return
	//}
	//// set it in the cookie
	//cookie := &http.Cookie{
	//	Name:     "accessToken",
	//	Value:    token,
	//	Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
	//	HttpOnly: true,
	//	Secure:   true,
	//}

	flashMessage := map[string]string{
		"type":    "success",
		"message": "account created successfully, proceed to login!",
	}
	w.WriteHeader(http.StatusCreated)
	RenderSignUpWithData(w, r, flashMessage)
}
