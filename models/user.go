package models

import (
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
