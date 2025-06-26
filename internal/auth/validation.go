package auth

import (
	"github.com/alisayeed248/gyo-en/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(username, password string) (bool, error) {
	// make a user Struct, check the db to see if we have a user that is this parameter
	var user database.User
	result := database.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return false, result.Error
	}

	// compare password with stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// password doesn't match
		return false, nil
	}

	return true, nil
}