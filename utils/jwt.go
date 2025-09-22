package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret string

// InitJWT untuk inisialisasi secret dari config
func InitJWT(secret string) {
	jwtSecret = secret
}

// GenerateJWT untuk membuat JWT Token
func GenerateJWT(username string, role string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// VerifyJWT untuk memverifikasi JWT Token
func VerifyJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
}
