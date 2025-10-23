package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golanjutan/app/model"
	"golanjutan/config"
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
		var userID int
		switch v := claims["user_id"].(type) {
		case float64:
			userID = int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				userID = i
			}
		}

		role, _ := claims["role"].(string)
		role = strings.ToLower(role)

		var alumniID *int
		if aID, ok := claims["alumni_id"]; ok {
			switch v := aID.(type) {
			case float64:
				val := int(v)
				alumniID = &val
			case string:
				if i, err := strconv.Atoi(v); err == nil {
					alumniID = &i
				}
			}
		}

		user := &model.User{
			ID:       userID,
			Role:     role,
			AlumniID: alumniID,
		}

		// ✅ Simpan user struct dan field individual di context
		c.Locals("user", user)
		c.Locals("user_id", userID)
		c.Locals("role", role)
		c.Locals("alumni_id", alumniID)

		return c.Next()
	}
}

// RequireRole memastikan user memiliki role tertentu
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
