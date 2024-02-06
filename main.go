package main

import (
	"errors"
	"fmt"
	"github.com/Nelwhix/Attendlog/controllers"
	"github.com/gorilla/securecookie"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)

		log.Printf(
			"[%s] %s %s",
			r.Method,
			r.RequestURI,
			time.Since(startTime),
		)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashKey := []byte(os.Getenv("HASH_KEY"))
		blockKey := []byte(os.Getenv("BLOCK_KEY"))
		s := securecookie.New(hashKey, blockKey)

		cookie, err := r.Cookie("accessToken")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		value := make(map[string]string)
		if err = s.Decode("accessToken", cookie.Value, &value); err == nil {
			log.Printf("Access Token for user is: %q", value["accessToken"])
		}
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	Enverr := godotenv.Load()
	if Enverr != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/public"))))
	router.HandleFunc("/auth/login", controllers.RenderLogin).Methods("GET")
	router.HandleFunc("/auth/signup", controllers.RenderSignUp).Methods("GET")

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

	handlers.CompressHandler(router)

	// seed database if we are in dev mode
	if os.Getenv("APP_ENV") == "dev" {
		for i := 0; i < 100; i++ {
			db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

			if err != nil {
				panic("failed to connect database")
			}

			db.AutoMigrate(&Record{})

			courses := []string{"MEG315", "MEG313", "MEG314"}
			n := rand.Intn(len(courses))
			db.Create(&Record{
				Name:   faker.Name(),
				Course: courses[n],
				Matric: "180404018",
			})
		}
	}

	log.Printf("Server starting on port %v\n", ConnPort)
	srv := &http.Server{
		Handler:      router,
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
