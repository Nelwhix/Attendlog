package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Nelwhix/Attendlog/controllers"
	"github.com/Nelwhix/Attendlog/models"
	"github.com/Nelwhix/Attendlog/services"
	"github.com/gorilla/csrf"
	"gorm.io/driver/sqlite"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const (
	ConnHost = "localhost"
	ConnPort = "8080"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cCookie, err := r.Cookie("accessToken")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userID, err := services.ValidateJwt(cCookie.Value)
		if err != nil {
			log.Printf("error validating jwt: %v", err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		cUser, err := models.GetUserById(userID)
		if err != nil {
			log.Printf("error retrieving user: %v", err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "currentUser", cUser)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func main() {
	Enverr := godotenv.Load()
	if Enverr != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Attendlog %v by Isioma Nelson", os.Getenv("APP_VERSION"))
	})
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/public"))))
	router.Handle("/favicon.ico", http.FileServer(http.Dir("./resources/public")))
	router.HandleFunc("/auth/login", controllers.RenderLogin).Methods("GET")
	router.HandleFunc("/auth/login", controllers.Login).Methods("POST")
	router.HandleFunc("/auth/signup", controllers.RenderSignUp).Methods("GET")
	router.HandleFunc("/auth/signup", controllers.SignUp).Methods("POST")
	router.HandleFunc("/generate-qrcode", controllers.GenerateQrCode).Methods("GET")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)
	protected.HandleFunc("/dashboard", controllers.RenderDashboard).Methods("GET")
	protected.HandleFunc("/attendance", controllers.CreateNewLink).Methods("POST")
	protected.HandleFunc("/attendance-success", controllers.RenderAttendanceSuccess).Methods("GET")

	compressed := handlers.CompressHandler(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, compressed)
	csrfMiddleware := csrf.Protect([]byte(os.Getenv("APP_KEY")))

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("./storage/app-%v.db", os.Getenv("APP_ENV"))), &gorm.Config{})
	if err != nil {
		log.Fatal("error opening db: ", err.Error())
		return
	}

	err = db.AutoMigrate(&models.User{}, &models.Record{}, &models.Link{})
	if err != nil {
		log.Fatal("error migrating models : ", err.Error())
		return
	}

	log.Printf("Server starting on port %v\n", ConnPort)
	srv := &http.Server{
		Handler:      csrfMiddleware(loggedRouter),
		Addr:         fmt.Sprintf("%s:%s", ConnHost, ConnPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
