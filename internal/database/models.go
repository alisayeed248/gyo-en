package database

import (
	"time"
)

// users are people who sign up
type User struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// monitoredUrl is a URl that a user wants to monitor
type MonitoredURL struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:not null"`
	URL       string    `json:"url" gorm:"not null"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CheckResult struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	URL          string        `json:"url" gorm:"not null"`
	IsUp         bool          `json:"is_up"`
	ResponseTime time.Duration `json:"response_time"`
	StatusCode   int           `json:"status_code"`
	ErrorMessage string        `json:"error_message"`
	CheckedAt    time.Time     `json:"checked_at"`
}
