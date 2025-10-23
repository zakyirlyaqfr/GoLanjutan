package service

import (
	"errors"
	"fmt"
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ==================== STRUCT ====================

type AlumniService struct {
	Repo *repository.AlumniRepository
}

func NewAlumniService(repo *repository.AlumniRepository) *AlumniService {
	return &AlumniService{Repo: repo}
}

// ==================== HANDLER (from routes) ====================

func (s *AlumniService) HandleGetAllWithFilter(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	sortBy := c.Query("sortBy", "created_at")
	sortOrder := c.Query("sortOrder", "DESC")
	search := c.Query("search", "")

	res, err := s.GetAllWithFilter(page, limit, sortBy, sortOrder, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": res.Data, "meta": res.Meta})
}

func (s *AlumniService) HandleGetAll(c *fiber.Ctx) error {
	res, err := s.GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (s *AlumniService) HandleGetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
	}
	res, err := s.GetByID(id)
	if err != nil {
		// Asumsi error dari GetByID adalah not found
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "alumni not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (s *AlumniService) HandleCreate(c *fiber.Ctx) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
	}
	id, err := s.Create(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	newAlumni, _ := s.GetByID(id) // Ambil data baru untuk ditampilkan
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": newAlumni})
}

func (s *AlumniService) HandleUpdate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
	}
	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
	}
	if err := s.Update(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	updated, _ := s.GetByID(id) // Ambil data baru untuk ditampilkan
	return c.JSON(fiber.Map{"success": true, "data": updated})
}

func (s *AlumniService) HandleSoftDelete(c *fiber.Ctx) error {
	user, err := getUserFromContext(c) // Menggunakan helper
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	id, _ := strconv.Atoi(c.Params("id"))

	if err := s.SoftDeleteAlumni(user, id); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Alumni berhasil di-soft delete"})
}

func (s *AlumniService) HandleHardDelete(c *fiber.Ctx) error {
	user, err := getUserFromContext(c) // Menggunakan helper
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id tidak valid"})
	}

	if err := s.HardDeleteAlumni(user, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "alumni dan semua pekerjaan terkait berhasil dihapus permanen",
	})
}

func (s *AlumniService) HandleRestore(c *fiber.Ctx) error {
	user, err := getUserFromContext(c) // Menggunakan helper
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	id, _ := strconv.Atoi(c.Params("id"))

	if err := s.RestoreAlumni(user, id); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Alumni berhasil di-restore"})
}

// ==================== SERVICE LOGIC ====================

func (s *AlumniService) GetAll() ([]model.Alumni, error) {
	return s.Repo.GetAll()
}

func (s *AlumniService) GetByID(id int) (*model.Alumni, error) {
	return s.Repo.GetByID(id)
}

func (s *AlumniService) Create(req model.CreateAlumniRequest) (int, error) {
	// Validasi bisnis
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return 0, errors.New("nim, nama, jurusan, dan email harus diisi")
	}
	return s.Repo.Create(req)
}

func (s *AlumniService) Update(id int, req model.UpdateAlumniRequest) error {
	// Validasi bisnis
	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return errors.New("nama, jurusan, dan email harus diisi")
	}
	
	// Cek dulu apakah alumni ada
	_, err := s.Repo.GetByID(id)
	if err != nil {
		return errors.New("alumni tidak ditemukan")
	}
	
	return s.Repo.Update(id, req)
}

func (s *AlumniService) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *AlumniService) GetAllWithFilter(page, limit int, sortBy, sortOrder, search string) (model.AlumniResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Whitelist kolom sort
	allowedSort := map[string]bool{
		"id":          true,
		"nim":         true,
		"nama":        true,
		"jurusan":     true,
		"angkatan":    true,
		"tahun_lulus": true,
		"created_at":  true,
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if strings.ToUpper(sortOrder) != "ASC" && strings.ToUpper(sortOrder) != "DESC" {
		sortOrder = "DESC"
	}

	data, err := s.Repo.GetAllWithFilter(limit, offset, sortBy, sortOrder, search)
	if err != nil {
		return model.AlumniResponse{}, err
	}

	total, err := s.Repo.Count(search)
	if err != nil {
		return model.AlumniResponse{}, err
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	
	return model.AlumniResponse{
		Data: data,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  totalPages,
			SortBy: sortBy,
			Order:  sortOrder,
			Search: search,
		},
	}, nil
}

// ✅ Logika otorisasi (hanya superadmin)
func (s *AlumniService) SoftDeleteAlumni(user *model.User, alumniID int) error {
	if user.ID != 1 {
		return errors.New("hanya superadmin yang bisa menghapus alumni")
	}
	return s.Repo.SoftDelete(alumniID)
}

// ✅ Logika otorisasi (hanya admin) + Logika bisnis (cek soft delete)
func (s *AlumniService) HardDeleteAlumni(user *model.User, alumniID int) error {
	if strings.ToLower(user.Role) != "admin" {
		return errors.New("hanya admin yang bisa hard delete alumni")
	}

	// Logika bisnis dipindahkan dari repo ke service
	alumni, err := s.Repo.GetByIDIncludeDeleted(alumniID)
	if err != nil {
		return errors.New("alumni tidak ditemukan")
	}

	if alumni.DeletedAt == nil {
		return fmt.Errorf("alumni belum dihapus (soft delete dulu sebelum hard delete)")
	}

	return s.Repo.HardDelete(alumniID)
}

// ✅ Logika otorisasi (hanya superadmin)
func (s *AlumniService) RestoreAlumni(user *model.User, alumniID int) error {
	if user.ID != 1 {
		return errors.New("hanya superadmin yang bisa restore alumni")
	}
	return s.Repo.Restore(alumniID)
}

// // ==================== HELPER FUNCTION (copy from pekerjaan_service) ====================
// func getUserFromContext(c *fiber.Ctx) (*model.User, error) {
// 	userData := c.Locals("user")
// 	if userData == nil {
// 		return nil, errors.New("user tidak ditemukan di context")
// 	}

// 	user, ok := userData.(*model.User)
// 	if !ok {
// 		return nil, errors.New("format user context tidak valid")
// 	}

// 	return user, nil
// }