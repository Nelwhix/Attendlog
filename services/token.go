package services

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJwt(userId string) (string, error) {
	claims := jwt.MapClaims{
		"iss":     os.Getenv("APP_NAME"),
		"iat":     time.Now(),
		"exp":     time.Now().Add(24 * time.Hour),
		"user_id": userId,
	}

	tokenInstance := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	token, err := tokenInstance.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", err
	}

	return token, nil
}
