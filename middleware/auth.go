package middleware


import (
"fmt"
"strings"
"strconv"


"github.com/gofiber/fiber/v2"
"github.com/golang-jwt/jwt/v5"
"golanjutan/config"
)


// Protected middleware: cek token, set locals user_id dan role
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
secret := config.AppEnv.JWTSecret
t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
return nil, fmt.Errorf("unexpected signing method")
}
return []byte(secret), nil
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
// jwt lib may have float64 numbers
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
c.Locals("role", role)
}
return c.Next()
}
}


// RequireRole middleware generator
func RequireRole(r string) fiber.Handler {
return func(c *fiber.Ctx) error {
role := c.Locals("role")
if role == nil {
return fiber.NewError(fiber.StatusForbidden, "forbidden")
}
if roleStr, ok := role.(string); ok {
if roleStr != r {
return fiber.NewError(fiber.StatusForbidden, "forbidden: insufficient role")
}
}
return c.Next()
}
}