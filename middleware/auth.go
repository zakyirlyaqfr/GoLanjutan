package middleware

import (
	"fmt"
	// "strconv" // Tidak dipakai
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golanjutan/app/model"
	"golanjutan/config"
	// "go.mongodb.org/mongo-driver/bson/primitive" // DIHAPUS
)

// Protected middleware: memvalidasi JWT & menyimpan info user ke context
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization format")
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(config.AppEnv.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
		}

		// --- Ambil data dari claims ---

		// DIUBAH: ID dari JWT adalah angka (float64 by default), bukan string hex
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid user_id in token (bukan float64)")
		}
		userID := int64(userIDFloat) // Konversi ke int64

		role, _ := claims["role"].(string)
		role = strings.ToLower(role)

		var alumniID *int64
		// DIUBAH: Cek alumni_id sebagai float64
		if aIDFloat, ok := claims["alumni_id"].(float64); ok {
			aID := int64(aIDFloat)
			alumniID = &aID
		}

		user := &model.User{
			ID:       userID, // DIUBAH
			Role:     role,
			AlumniID: alumniID, // DIUBAH
		}

		// ✅ Simpan user struct dan field individual di context
		c.Locals("user", user)
		c.Locals("user_id", userID) // Menyimpan int64
		c.Locals("role", role)
		c.Locals("alumni_id", alumniID) // Menyimpan *int64

		return c.Next()
	}
}

// RequireRole tidak perlu diubah, sudah benar
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData := c.Locals("user")
		if userData == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: missing user context")
		}

		user, ok := userData.(*model.User)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "Invalid user context type")
		}

		if !strings.EqualFold(user.Role, requiredRole) {
			return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("Forbidden: role %s required", requiredRole))
		}

		return c.Next()
	}
}