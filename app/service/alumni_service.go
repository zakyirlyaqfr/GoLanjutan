package service

import (
	"context"
	"errors"
	"fmt"
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ==================== STRUCT ====================

type AlumniService struct {
	Repo repository.IAlumniRepository
}

func NewAlumniService(repo repository.IAlumniRepository) *AlumniService {
	return &AlumniService{Repo: repo}
}

// Helper context
func getCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// BARU: Menambahkan helper yang hilang
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

// ==================== HANDLER (from routes) ====================
// ... Semua handler (HandleGetAll, HandleGetByID, HandleCreate, HandleUpdate) sudah benar ...
// ... Saya sertakan HandleSoftDelete, HandleHardDelete, HandleRestore untuk menunjukkan perubahannya ...

func (s *AlumniService) HandleGetAllWithFilter(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	sortBy := c.Query("sortBy", "created_at")
	sortOrder := c.Query("sortOrder", "DESC")
	search := c.Query("search", "")

	ctx, cancel := getCtx()
	defer cancel()

	res, err := s.GetAllWithFilter(ctx, page, limit, sortBy, sortOrder, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": res.Data, "meta": res.Meta})
}

func (s *AlumniService) HandleGetAll(c *fiber.Ctx) error {
	ctx, cancel := getCtx()
	defer cancel()
	res, err := s.GetAll(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (s *AlumniService) HandleGetByID(c *fiber.Ctx) error {
	id := c.Params("id") 

	ctx, cancel := getCtx()
	defer cancel()
	res, err := s.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (s *AlumniService) HandleCreate(c *fiber.Ctx) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
	}

	ctx, cancel := getCtx()
	defer cancel()
	newAlumni, err := s.Create(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": newAlumni})
}

func (s *AlumniService) HandleUpdate(c *fiber.Ctx) error {
	id := c.Params("id") 
	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
	}

	ctx, cancel := getCtx()
	defer cancel()
	updated, err := s.Update(ctx, id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": updated})
}


func (s *AlumniService) HandleSoftDelete(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	id := c.Params("id")

	ctx, cancel := getCtx()
	defer cancel()
	if err := s.SoftDeleteAlumni(ctx, user, id); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Alumni berhasil di-soft delete"})
}

func (s *AlumniService) HandleHardDelete(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	id := c.Params("id")

	ctx, cancel := getCtx()
	defer cancel()
	if err := s.HardDeleteAlumni(ctx, user, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"message": "alumni dan semua pekerjaan terkait berhasil dihapus permanen",
	})
}

func (s *AlumniService) HandleRestore(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	id := c.Params("id")

	ctx, cancel := getCtx()
	defer cancel()
	if err := s.RestoreAlumni(ctx, user, id); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Alumni berhasil di-restore"})
}

// ==================== SERVICE LOGIC ====================
// ... GetAll, GetByID, Create, Update, GetAllWithFilter sudah benar ...
// ... Saya sertakan SoftDelete, HardDelete, Restore untuk menunjukkan perubahannya ...

func (s *AlumniService) GetAll(ctx context.Context) ([]model.Alumni, error) {
	return s.Repo.GetAll(ctx)
}

func (s *AlumniService) GetByID(ctx context.Context, id string) (*model.Alumni, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *AlumniService) Create(ctx context.Context, req model.CreateAlumniRequest) (*model.Alumni, error) {
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return nil, errors.New("nim, nama, jurusan, dan email harus diisi")
	}

	alumni := model.Alumni{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.Repo.Create(ctx, alumni)
}

func (s *AlumniService) Update(ctx context.Context, id string, req model.UpdateAlumniRequest) (*model.Alumni, error) {
	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return nil, errors.New("nama, jurusan, dan email harus diisi")
	}

	alumni, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("alumni tidak ditemukan")
	}

	alumni.Nama = req.Nama
	alumni.Jurusan = req.Jurusan
	alumni.Angkatan = req.Angkatan
	alumni.TahunLulus = req.TahunLulus
	alumni.Email = req.Email
	alumni.NoTelepon = req.NoTelepon
	alumni.Alamat = req.Alamat
	alumni.UpdatedAt = time.Now()

	err = s.Repo.Update(ctx, id, *alumni)
	return alumni, err
}

func (s *AlumniService) GetAllWithFilter(ctx context.Context, page, limit int, sortBy, sortOrder, search string) (model.AlumniResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	allowedSort := map[string]bool{
		"_id":         true,
		"id":          true, 
		"nim":         true,
		"nama":        true,
		"jurusan":     true,
		"angkatan":    true,
		"tahun_lulus": true,
		"created_at":  true,
	}
	if sortBy == "id" {
		sortBy = "_id"
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if strings.ToUpper(sortOrder) != "ASC" && strings.ToUpper(sortOrder) != "DESC" {
		sortOrder = "DESC"
	}

	data, err := s.Repo.GetAllWithFilter(ctx, limit, offset, sortBy, sortOrder, search)
	if err != nil {
		return model.AlumniResponse{}, err
	}

	total, err := s.Repo.Count(ctx, search)
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


// Logika otorisasi (hanya superadmin)
func (s *AlumniService) SoftDeleteAlumni(ctx context.Context, user *model.User, alumniID string) error {
	// DIUBAH: Cek superadmin menggunakan ID int64. Ganti '1' dengan ID superadmin Anda.
	if user.ID != 1 {
		return errors.New("hanya superadmin yang bisa menghapus alumni")
	}
	return s.Repo.SoftDelete(ctx, alumniID)
}

// Logika otorisasi (hanya admin) + Logika bisnis (cek soft delete)
func (s *AlumniService) HardDeleteAlumni(ctx context.Context, user *model.User, alumniID string) error {
	if strings.ToLower(user.Role) != "admin" {
		return errors.New("hanya admin yang bisa hard delete alumni")
	}

	alumni, err := s.Repo.GetByIDIncludeDeleted(ctx, alumniID)
	if err != nil {
		return errors.New("alumni tidak ditemukan")
	}

	if alumni.DeletedAt == nil {
		return fmt.Errorf("alumni belum dihapus (soft delete dulu sebelum hard delete)")
	}

	return s.Repo.HardDelete(ctx, alumniID)
}

// Logika otorisasi (hanya superadmin)
func (s *AlumniService) RestoreAlumni(ctx context.Context, user *model.User, alumniID string) error {
	// DIUBAH: Cek superadmin menggunakan ID int64. Ganti '1' dengan ID superadmin Anda.
	if user.ID != 1 {
		return errors.New("hanya superadmin yang bisa restore alumni")
	}
	return s.Repo.Restore(ctx, alumniID)
}