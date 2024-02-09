package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

var Db *gorm.DB

func New() (*gorm.DB, error) {
	if Db != nil {
		return Db, nil
	}

	Db, err := gorm.Open(sqlite.Open(fmt.Sprintf("./storage/app-%v.db", os.Getenv("APP_ENV"))), &gorm.Config{})
	if err != nil {
		log.Fatal("error opening db: ", err.Error())
		return nil, err
	}

	return Db, nil
}
