package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	DB, err = gorm.Open(sqlite.Open("gyo-en.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB.AutoMigrate(&User{}, &MonitoredURL{}, &CheckResult{})
}