package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = os.Getenv("SECRET_KEY")

func GenerateJWT(adminID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin":    true,
		"admin_id": adminID,
		"exp":      time.Now().Add(24 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	log.Printf("Created token: %s", tokenString)

	return tokenString, nil
}
