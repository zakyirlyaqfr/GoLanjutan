package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	// "go.mongodb.org/mongo-driver/bson/primitive" // DIHAPUS
)

var jwtSecret string

// InitJWT menginisialisasi secret key
func InitJWT(secret string) {
	jwtSecret = secret
}

// GenerateJWTWithClaims membuat JWT langsung dari MapClaims
// (Tidak berubah, karena auth_service sudah mengirim int64)
func GenerateJWTWithClaims(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GenerateJWT membuat token JWT dengan payload lengkap
// (Fungsi ini tidak lagi dipakai oleh auth_service, tapi tetap saya perbaiki)
func GenerateJWT(userID int64, alumniID *int64, username, role string, expiry time.Duration) (string, error) { // DIUBAH
	claims := jwt.MapClaims{
		"user_id":  userID, // DIUBAH (bukan .Hex())
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(expiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	if alumniID != nil {
		claims["alumni_id"] = *alumniID // DIUBAH (bukan .Hex())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// VerifyJWT memverifikasi token JWT dan mengembalikan claims
// (Tidak berubah)
func VerifyJWT(tokenStr string) (jwt.MapClaims, error) {
	// ... (isi fungsi sama) ...
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}