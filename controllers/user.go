package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Nelwhix/Attendlog/database"
	"github.com/Nelwhix/Attendlog/models"
	"github.com/Nelwhix/Attendlog/requests"
	"github.com/Nelwhix/Attendlog/services"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/oklog/ulid/v2"
	"html/template"
	"log"
	"net/http"
	"time"
)

var decoder = schema.NewDecoder()
var validate = validator.New(validator.WithRequiredStructEnabled())

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFiles("templates/login.tmpl")
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

func RenderSignUp(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFiles("templates/signup.tmpl")
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

func RenderDashboard(w http.ResponseWriter, r *http.Request) {
	cUser := r.Context().Value("currentUser").(models.User)

	parsedTemplate, err := template.ParseFiles("templates/dashboard.tmpl")
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}

	err = parsedTemplate.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"UserName":       cUser.UserName,
	})
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func RenderDashboardWithData(w http.ResponseWriter, r *http.Request, flashMessage map[string]string) {
	cUser := r.Context().Value("currentUser").(models.User)
	jsonData, _ := json.Marshal(flashMessage)
	parsedTemplate, err := template.ParseFiles("templates/dashboard.tmpl")
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}

	err = parsedTemplate.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"UserName":       cUser.UserName,
		"flashMessage":   string(jsonData),
	})
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func RenderSignUpWithData(w http.ResponseWriter, r *http.Request, flashMessage map[string]string) {
	parsedTemplate, err := template.ParseFiles("templates/signup.tmpl")
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

func RenderLoginWithData(w http.ResponseWriter, r *http.Request, flashMessage map[string]string) {
	parsedTemplate, err := template.ParseFiles("templates/login.tmpl")
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
	db, err := database.New()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
		log.Printf("error opening database: %v", err.Error())
		return
	}
	var newUser models.User
	db.Where("email_address = ?", signupRequest.Email).First(&newUser)

	if newUser.ID != "" {
		w.WriteHeader(http.StatusBadRequest)
		flashMessage := map[string]string{
			"type":    "error",
			"message": "email is taken",
		}
		RenderSignUpWithData(w, r, flashMessage)
		return
	}
	db.Where("user_name = ?", signupRequest.UserName).First(&newUser)
	if newUser.ID != "" {
		w.WriteHeader(http.StatusBadRequest)
		flashMessage := map[string]string{
			"type":    "error",
			"message": "username is taken",
		}
		RenderSignUpWithData(w, r, flashMessage)
		return
	}

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

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Welcome to Attendlog! \n%s\n\n", mimeHeaders)))
	htmlTemplate, err := template.ParseFiles("./mail_templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
		log.Printf("error parsing email template: %v", err.Error())
		return
	}

	err = htmlTemplate.Execute(&body, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
		log.Printf("error parsing email template: %v", err.Error())
		return
	}

	err = services.SendMail(body, signupRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
		log.Printf("error sending welcome mail: %v", err.Error())
		return
	}
	flashMessage := map[string]string{
		"type":    "success",
		"message": "account created successfully, proceed to login!",
	}
	w.WriteHeader(http.StatusCreated)
	RenderSignUpWithData(w, r, flashMessage)
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		flashMessage := map[string]string{
			"type":    "error",
			"message": "error parsing form",
		}
		RenderLoginWithData(w, r, flashMessage)
		log.Printf("error parsing form: %v", err.Error())
		return
	}

	r.PostForm.Del("gorilla.csrf.Token")

	var loginRequest requests.Login
	err = decoder.Decode(&loginRequest, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		flashMessage := map[string]string{
			"type":    "error",
			"message": err.Error(),
		}
		RenderLoginWithData(w, r, flashMessage)
		log.Printf("error decoding form: %v", err.Error())
		return
	}

	err = validate.Struct(loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		flashMessage := map[string]string{
			"type":    "error",
			"message": err.Error(),
		}
		RenderLoginWithData(w, r, flashMessage)
		log.Printf("validation error: %v", err.Error())
		return
	}

	// fetch user, attempt auth
	var newUser models.User

	db, err := database.New()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error opening db: %v", err.Error())
		return
	}
	db.Where("email_address = ?", loginRequest.Email).First(&newUser)
	if newUser.ID == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)

		flashMessage := map[string]string{
			"type":    "error",
			"message": "invalid email/password",
		}
		RenderLoginWithData(w, r, flashMessage)
		return
	}

	token, err := services.GenerateJwt(newUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
		log.Printf("error generating token: %v", err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	if loginRequest.RememberMe {
		cookie.Expires = time.Now().Add(24 * time.Hour)
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
