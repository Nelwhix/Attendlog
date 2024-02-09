package models

import (
	"time"
)

type Record struct {
	ID        string `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"column:email_address;unique"`
	Signature string
	LinkID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
