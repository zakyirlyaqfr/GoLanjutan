package service

import (
	"errors"
	"golanjutan/app/model"
	"golanjutan/app/repository"

	"github.com/gofiber/fiber/v2"
)

type TrashService struct {
	AlumniRepo    *repository.AlumniRepository
	PekerjaanRepo *repository.PekerjaanRepository
}

func NewTrashService(alumniRepo *repository.AlumniRepository, pekerjaanRepo *repository.PekerjaanRepository) *TrashService {
	return &TrashService{
		AlumniRepo:    alumniRepo,
		PekerjaanRepo: pekerjaanRepo,
	}
}

// ==================== HANDLER (from routes) ====================

func (s *TrashService) HandleGetTrash(c *fiber.Ctx) error {
	user, err := s.getUserFromContext(c) // Menggunakan helper
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := s.GetTrash(user.Role, user.AlumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "trash data retrieved successfully",
		"data":    data,
	})
}

// ==================== SERVICE LOGIC ====================

func (s *TrashService) GetTrash(role string, alumniID *int) (map[string]interface{}, error) {
	if role == "admin" {
		alumni, err := s.AlumniRepo.GetTrashed()
		if err != nil {
			return nil, err
		}
		pekerjaan, err := s.PekerjaanRepo.GetTrashed()
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"alumni":    alumni,
			"pekerjaan": pekerjaan,
		}, nil
	}

	// role user
	if alumniID == nil {
		return map[string]interface{}{
			// "alumni":    []model.Alumni{}, // User tidak bisa melihat trash alumni
			"pekerjaan": []model.PekerjaanAlumni{},
		}, nil
	}

	pekerjaan, err := s.PekerjaanRepo.GetTrashedByAlumniID(*alumniID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pekerjaan": pekerjaan,
	}, nil
}

// ==================== HELPER FUNCTION (copy from pekerjaan_service) ====================
func (s *TrashService) getUserFromContext(c *fiber.Ctx) (*model.User, error) {
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