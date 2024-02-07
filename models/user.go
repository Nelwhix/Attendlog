package models

import (
	"time"
)

type User struct {
	ID               string `gorm:"primaryKey"`
	FirstName        string `gorm:"column:first_name"`
	LastName         string `gorm:"column:last_name"`
	UserName         string `gorm:"column:user_name;unique_index"`
	Email            string `gorm:"column:email_address;unique_index"`
	Password         string `gorm:"column:password"`
	SecurityQuestion string `gorm:"column:security_question"`
	Answer           string `gorm:"column:answer"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
