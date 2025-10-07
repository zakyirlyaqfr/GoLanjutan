package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret string

func InitJWT(secret string) {
	jwtSecret = secret
}

// GenerateJWTWithClaims membuat JWT langsung dari MapClaims
func GenerateJWTWithClaims(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func GenerateJWT(userID int, alumniID *int, username, role string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"username":  username,
		"role":      role,
		"exp":       time.Now().Add(expiry).Unix(),
	}
	if alumniID != nil {
		claims["alumni_id"] = *alumniID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func VerifyJWT(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, nil, fmt.Errorf("invalid token claims")
	}
	return t, claims, nil
}
