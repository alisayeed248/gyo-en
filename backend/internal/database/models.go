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

