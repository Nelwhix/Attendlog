package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJwt(userId string) (string, error) {
	tokenBytes, err := os.ReadFile("./cert/id_rsa")
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"iss":     os.Getenv("APP_NAME"),
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"user_id": userId,
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("could not parse key: %w", err)
	}

	tokenInstance := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := tokenInstance.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateJwt(cToken string) (string, error) {
	tokenBytes, err := os.ReadFile("./cert/id_rsa.pub")
	if err != nil {
		return "", fmt.Errorf("error opening public key file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(tokenBytes)
	if err != nil {
		return "", err
	}

	tok, err := jwt.Parse(cToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("could not validate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return "", fmt.Errorf("validate: invalid token")
	}

	issuer, ok := claims["iss"].(string)
	if !ok {
		return "", fmt.Errorf("validate: invalid token")
	}

	if issuer != os.Getenv("APP_NAME") {
		return "", fmt.Errorf("validate: invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("validate: invalid token")
	}

	if time.Now().Unix() > int64(exp) {
		return "", fmt.Errorf("expired token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("validate: invalid token")
	}
	return userID, nil
}
