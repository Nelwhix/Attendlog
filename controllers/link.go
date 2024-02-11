package controllers

import (
	"fmt"
	"github.com/Nelwhix/Attendlog/database"
	"github.com/Nelwhix/Attendlog/models"
	"github.com/Nelwhix/Attendlog/requests"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid/v2"
	"github.com/skip2/go-qrcode"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func CreateNewLink(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		flashMessage := map[string]string{
			"type":    "error",
			"message": "error parsing form",
		}
		RenderDashboardWithData(w, r, flashMessage)
		log.Printf("error parsing form: %v", err.Error())
		return
	}

	r.PostForm.Del("gorilla.csrf.Token")

	var createLinkRequest requests.CreateLink
	err = decoder.Decode(&createLinkRequest, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		flashMessage := map[string]string{
			"type":    "error",
			"message": err.Error(),
		}
		RenderDashboardWithData(w, r, flashMessage)
		log.Printf("error decoding form: %v", err.Error())
		return
	}

	cUser := r.Context().Value("currentUser").(models.User)

	nLink := models.Link{
		ID:                     ulid.Make().String(),
		Title:                  createLinkRequest.Title,
		UserID:                 cUser.ID,
		HasSignature:           createLinkRequest.HasSignature,
		HasLocationRestriction: createLinkRequest.HasLocationRestriction,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	if createLinkRequest.Description != "" {
		nLink.Description = createLinkRequest.Description
	}

	if createLinkRequest.HasLocationRestriction {
		nLink.HasLocationRestriction = true
		nLink.Latitude = createLinkRequest.Latitude
		nLink.Longitude = createLinkRequest.Longitude
	}

	db, err := database.New()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error opening database: %v", err.Error())
		return
	}
	result := db.Create(&nLink)
	if result.Error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error inserting record: %v", result.Error.Error())
		return
	}

	redirectUrl := fmt.Sprintf("/attendance/%v", nLink.ID)
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

func RenderAttendance(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	link := fmt.Sprintf("%v/link/%v", os.Getenv("APP_URL"), pathParams["id"])
	cUser := r.Context().Value("currentUser").(models.User)

	parsedTemplate, err := template.ParseFiles("templates/attendance.tmpl")
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}

	err = parsedTemplate.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"UserName":       cUser.UserName,
		"AttendanceLink": link,
	})
	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func GenerateQrCode(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	url := queryValues.Get("url")

	qrCode, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error generating qrcode image: %s", err.Error())
	}

	w.Header().Set("Content-Type", "image/png")

	_, err = w.Write(qrCode)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error writing qrcode image to response: %s", err.Error())
	}
}
