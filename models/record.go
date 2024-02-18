package models

import (
	"time"
)

type Record struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `gorm:"column:email_address" json:"email_address"`
	Signature string    `json:"signature"`
	LinkID    string    `json:"link_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
