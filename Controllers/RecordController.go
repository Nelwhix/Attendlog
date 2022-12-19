package Controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Records struct {
	Records []Record
}

func RenderAttendanceForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	data := map[string]string{
		"Course": vars["course"],
	}

	parsedTemplate, _ := template.ParseFiles("views/index.html")
	err := parsedTemplate.Execute(w, data)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	record := new(Record)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(record, r.PostForm)

	if decodeErr != nil {
		log.Printf("error mapping form data to struct: %v", decodeErr)
	}	

	valid, validationErrorMessage := validateInput(w, r, record)

	if !valid {
		fmt.Fprint(w, validationErrorMessage)
		return
	}

	dataURI := record.Signature
	encodedImg := strings.Split(dataURI, ",")[1]
	decodedImg, _ := base64.StdEncoding.DecodeString(encodedImg)

	_, err := os.Stat("./resources/uploads/" + record.Matric + ".png")

	if !errors.Is(err, os.ErrNotExist) {
		os.Remove("./resources/uploads/" + record.Matric + ".png")
	}

	out, err := os.Create("./resources/uploads/" + record.Matric + ".png")

	if err != nil {
		log.Printf("error creating a file for writing %v", err)
		return
	}
	defer out.Close()
	os.WriteFile(out.Name(), decodedImg, 0644)

	record.Signature = filepath.Base(out.Name())
	

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var result Record

	beforeTime := time.Now().Add(-time.Hour * 6)
	db.Where("Course == ? AND Name == ? AND Matric == ? AND created_at BETWEEN ? AND ?",
	 record.Course, record.Name, record.Matric, beforeTime, time.Now()).First(&result)

	if result.Name != "" {
		fmt.Fprint(w, "Duplicate entries are not allowed")
		return
	 }

	db.AutoMigrate(&Record{})
	db.Create(&record)
	fmt.Fprintf(w, "Record submitted successfully");
}

func validateInput(w http.ResponseWriter, r *http.Request, record *Record) (bool, string) {
	valid, validationError := govalidator.ValidateStruct(record)

	if len(record.Matric) != 9 {
		valid = false
	}

	if !valid {
		nameError := govalidator.ErrorByField(validationError, "Name")
		matricError := govalidator.ErrorByField(validationError, "Matric")

		if nameError != "" {
			return valid, "Please fill in a valid name"
		}

		if matricError != "" {
			return valid, "Please fill in a valid matric number"
		}

		return valid, "Your Matric number must be 9 digits long"
	}


	return valid, "Validation Error"
}

func GetRecords(w http.ResponseWriter, r *http.Request) {
	isActive := hasActiveSession(r)

	if !isActive {
		fmt.Fprintln(w, "You are not authorized to view this page")
		return
	}

	vars := mux.Vars(r)

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var records []Record
	beforeTime := time.Now().Add(-time.Hour * 6)
	db.Where("Course == ? AND created_at BETWEEN ? AND ?", vars["course"], beforeTime, time.Now()).Find(&records)

	data := Records{
		Records: records,
	}

	parsedTemplates, parseErr := template.ParseFiles("views/records.html")

	if parseErr != nil {
		log.Printf("Error parsing html: %v", parseErr)
	}

	err = parsedTemplates.Execute(w, data)

	if err != nil {
		log.Printf("Error occured while executing the template or writing its output : %v", err)
		return
	}
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var record Record
	db.First(&record, vars["record"])
	course := record.Course

	db.Delete(&Record{}, vars["record"])
	
	http.Redirect(w, r, "/records/" + course, 200)
}

// func ExportRecords(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	var result []Record

// 	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	db.Where("Course == ?", vars["course"]).Find(&result)

// 	f, err := os.Create(vars["course"] + "-records.csv")
// 	defer f.Close()

// 	if err != nil {
// 		log.Fatalln("failed to open file", err)
// 	}


// 	// csvWriter := csv.NewWriter(f)
// 	// err = csvWriter.WriteAll(result)
// }