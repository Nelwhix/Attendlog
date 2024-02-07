package main

import (
	"errors"
	"fmt"
	"github.com/Nelwhix/Attendlog/controllers"
	"github.com/gorilla/csrf"
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

type Record struct {
	gorm.Model
	Course string
	Name   string `valid:"required"`
	Matric string `valid:"numeric,required"`
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("accessToken")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	Enverr := godotenv.Load()
	if Enverr != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/public"))))
	router.HandleFunc("/auth/login", controllers.RenderLogin).Methods("GET")
	router.HandleFunc("/auth/signup", controllers.RenderSignUp).Methods("GET")
	router.HandleFunc("/auth/signup", controllers.SignUp).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)
	protected.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Attendlog API %v by Isioma Nelson", os.Getenv("APP_VERSION"))
	})

	//router.HandleFunc("/attendance/{course}", controllers.RenderAttendanceForm).Methods("GET")
	//router.HandleFunc("/attendance/{course}", controllers.SubmitAttendance).Methods("POST")
	//router.HandleFunc("/records/{course}", controllers.GetRecords).Methods("GET")
	//router.HandleFunc("/records/delete/{record}", controllers.DeleteRecord).Methods("POST")
	//router.HandleFunc("/records/export/{course}", controllers.ExportRecords).Methods("GET")
	//
	//router.HandleFunc("/admin", controllers.Login).Methods("POST")
	//
	//// Course Controller
	//router.HandleFunc("/courses/add", controllers.RenderCourseForm).Methods("GET")
	//router.HandleFunc("/courses/add", controllers.AddCourse).Methods("POST")
	//router.PathPrefix("/").Handler(http.StripPrefix("/resources", http.FileServer(http.Dir("resources/"))))

	compressed := handlers.CompressHandler(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, compressed)
	csrfMiddleware := csrf.Protect([]byte(os.Getenv("APP_KEY")))

	log.Printf("Server starting on port %v\n", ConnPort)
	srv := &http.Server{
		Handler:      csrfMiddleware(loggedRouter),
		Addr:         fmt.Sprintf("%s:%s", ConnHost, ConnPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
