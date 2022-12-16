package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/Nelwhix/mech-attendance/Controllers"
	"github.com/go-faker/faker/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

type Record struct {
	gorm.Model
	Course string
	Name string `valid:"required"`
	Matric string `valid:"numeric,required"`
}

func main() {
	Enverr := godotenv.Load()
	if Enverr != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to MECH25 Attendance Web App by Isioma Nelson")
	})

	// Record Controller
	router.HandleFunc("/attendance/{course}", Controllers.RenderAttendanceForm).Methods("GET")
	router.HandleFunc("/attendance/{course}", Controllers.SubmitAttendance).Methods("POST")
	router.HandleFunc("/records/{course}", Controllers.GetRecords).Methods("GET")
	router.HandleFunc("/records/delete/{record}", Controllers.DeleteRecord).Methods("POST")

	// User Controller
	router.HandleFunc("/admin", Controllers.RenderLogin).Methods("GET")
	router.HandleFunc("/dashboard", Controllers.RenderDashboard).Methods("GET")
	router.HandleFunc("/admin", Controllers.Login).Methods("POST")
	
	// Course Controller
	router.HandleFunc("/courses/add", Controllers.RenderCourseForm).Methods("GET")
	router.HandleFunc("/courses/add", Controllers.AddCourse).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/resources", http.FileServer(http.Dir("resources/"))))

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
				Name: faker.Name(),
				Course: courses[n],
				Matric: "180404018",
			})
		}
	}

	log.Printf("Server starting on port %v\n", CONN_PORT)
	err := http.ListenAndServe(CONN_HOST + ":" + CONN_PORT, router)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}