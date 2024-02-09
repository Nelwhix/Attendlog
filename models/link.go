package models

import "time"

type Link struct {
	ID        string `gorm:"primaryKey"`
	Latitude  float64
	Longitude float64
	UserID    string
	Records   []Record
	CreatedAt time.Time
	UpdatedAt time.Time
}
