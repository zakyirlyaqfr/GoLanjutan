package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golanjutan/config"
)

// Protected middleware
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
		}

		tokenStr := parts[1]
		t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(config.AppEnv.JWTSecret), nil
		})
		if err != nil || !t.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
		}

		// set locals
		if uid, ok := claims["user_id"]; ok {
			switch v := uid.(type) {
			case float64:
				c.Locals("user_id", int(v))
			case int:
				c.Locals("user_id", v)
			case string:
				if i, err := strconv.Atoi(v); err == nil {
					c.Locals("user_id", i)
				}
			}
		}

		if role, ok := claims["role"].(string); ok {
			c.Locals("role", strings.ToLower(role))
		}

		if alumniID, ok := claims["alumni_id"]; ok {
			switch v := alumniID.(type) {
			case float64:
				c.Locals("alumni_id", int(v))
			case int:
				c.Locals("alumni_id", v)
			}
		}

		return c.Next()
	}
}

// RequireRole middleware
func RequireRole(r string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		}
		if roleStr, ok := role.(string); ok {
			if roleStr != strings.ToLower(r) {
				return fiber.NewError(fiber.StatusForbidden, "forbidden: insufficient role")
			}
		}
		return c.Next()
	}
}
