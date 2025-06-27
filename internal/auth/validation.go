package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/alisayeed248/gyo-en/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(username, password string) (*database.User, error) {
	// make a user Struct, check the db to see if we have a user that is this parameter
	var user database.User
	result := database.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	// compare password with stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// password doesn't match
		return nil, nil
	}

	return &user, nil
}

type Claims struct {
	UserID uint `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("some-secret-key")

func GenerateJWT(userID uint, username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}