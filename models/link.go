package models

import "time"

type Link struct {
	ID                     string `gorm:"primaryKey"`
	Title                  string
	Description            string
	Latitude               float64
	Longitude              float64
	UserID                 string
	HasLocationRestriction bool
	HasSignature           bool
	Records                []Record
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
