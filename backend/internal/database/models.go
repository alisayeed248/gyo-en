package database

import (
	"time"
)

// users are people who sign up
type User struct {
	Id	uint 	`json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Email string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// monitoredUrl is a URl that a user wants to monitor
type MonitoredUrl struct {
	ID uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"url" gorm:not null"`
	URL string `json:"url" gorm:"not null"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

