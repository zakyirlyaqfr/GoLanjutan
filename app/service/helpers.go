package service

import (
	"errors"
	"golanjutan/app/model"

	"github.com/gofiber/fiber/v2"
)

// ==================== HELPER FUNCTION ====================

// Fungsi ini mengambil data user dari context (biasanya diset lewat middleware JWT)
// Didefinisikan SATU KALI di sini untuk digunakan oleh semua service.
func getUserFromContext(c *fiber.Ctx) (*model.User, error) {
	userData := c.Locals("user")
	if userData == nil {
		return nil, errors.New("user tidak ditemukan di context")
	}

	user, ok := userData.(*model.User)
	if !ok {
		return nil, errors.New("format user context tidak valid")
	}

	return user, nil
}