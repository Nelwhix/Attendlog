package models

import (
	"github.com/Nelwhix/Attendlog/database"
	"time"
)

type User struct {
	ID               string `gorm:"primaryKey"`
	FirstName        string
	LastName         string
	UserName         string `gorm:"unique"`
	Email            string `gorm:"column:email_address;unique"`
	Password         string
	SecurityQuestion string
	Answer           string
	Links            []Link
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func GetUserById(userID string) (User, error) {
	var cUser User
	db, err := database.New()
	if err != nil {
		return User{}, err
	}

	db.First(&cUser, "id = ?", userID)
	if cUser.Email == "" {
		return User{}, err
	}

	return cUser, nil
}
