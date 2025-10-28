package service

import (
	"context"
	"errors"
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"strconv" // DITAMBAHKAN
	"time"

	"github.com/gofiber/fiber/v2"
	// "go.mongodb.org/mongo-driver/bson/primitive" // DIHAPUS
)

// Helper context
func getTrashCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

type TrashService struct {
	AlumniRepo    repository.IAlumniRepository
	PekerjaanRepo repository.IPekerjaanRepository
}

func NewTrashService(alumniRepo repository.IAlumniRepository, pekerjaanRepo repository.IPekerjaanRepository) *TrashService {
	return &TrashService{
		AlumniRepo:    alumniRepo,
		PekerjaanRepo: pekerjaanRepo,
	}
}

// ==================== HANDLER (from routes) ====================

func (s *TrashService) HandleGetTrash(c *fiber.Ctx) error {
	user, err := s.getUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := getTrashCtx()
	defer cancel()
	// user.AlumniID sekarang *int64
	data, err := s.GetTrash(ctx, user.Role, user.AlumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "trash data retrieved successfully",
		"data":    data,
	})
}

// ==================== SERVICE LOGIC ====================

// DIUBAH: Tanda tangan fungsi menerima *int64
func (s *TrashService) GetTrash(ctx context.Context, role string, alumniID *int64) (map[string]interface{}, error) {
	if role == "admin" {
		alumni, err := s.AlumniRepo.GetTrashed(ctx)
		if err != nil {
			return nil, err
		}
		pekerjaan, err := s.PekerjaanRepo.GetTrashed(ctx)
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
			"pekerjaan": []model.PekerjaanAlumni{},
		}, nil
	}

	// DIUBAH: Konversi *int64 ke string untuk dikirim ke repo
	alumniIDString := strconv.FormatInt(*alumniID, 10)
	pekerjaan, err := s.PekerjaanRepo.GetTrashedByAlumniID(ctx, alumniIDString)
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