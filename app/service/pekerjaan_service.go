package service

import (
	"context"
	"errors"
	"fmt"
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"strconv" // DITAMBAHKAN
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	// "go.mongodb.org/mongo-driver/bson/primitive" // DIHAPUS
)

// ==================== STRUCT ====================

type PekerjaanService struct {
	Repo repository.IPekerjaanRepository
}

func NewPekerjaanService(repo repository.IPekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{Repo: repo}
}

// Helper context
func getPekerjaanCtx() (context.Context, context.CancelFunc) {
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

// ==================== HANDLER ====================
// ... Semua handler (HandleGetAll, HandleGetByID, dst.) sudah benar ...
// ... Tidak ada perubahan di sisi handler ...
func (s *PekerjaanService) HandleGetAll(c *fiber.Ctx) error {
	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	data, err := s.GetAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	data, err := s.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleGetByAlumniID(c *fiber.Ctx) error {
	alumniID := c.Params("alumni_id")
	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	data, err := s.GetByAlumniID(ctx, alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleGetAllWithFilter(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	sortBy := c.Query("sortBy", "created_at")
	sortOrder := c.Query("sortOrder", "DESC")
	search := c.Query("search", "")

	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	res, err := s.GetAllWithFilter(ctx, page, limit, sortBy, sortOrder, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    res.Data,
		"meta":    res.Meta,
	})
}

func (s *PekerjaanService) HandleCreate(c *fiber.Ctx) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	data, err := s.Create(ctx, req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	data, err := s.Update(ctx, id, req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleSoftDelete(c *fiber.Ctx) error {
	user, _ := getUserFromContext(c)
	id := c.Params("id")

	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	if err := s.SoftDeletePekerjaan(ctx, user, id); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil di-soft delete"})
}

func (s *PekerjaanService) HandleHardDelete(c *fiber.Ctx) error {
	user, _ := getUserFromContext(c)
	id := c.Params("id")

	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	if err := s.HardDeletePekerjaan(ctx, user, id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus permanen"})
}

func (s *PekerjaanService) HandleRestore(c *fiber.Ctx) error {
	user, _ := getUserFromContext(c)
	id := c.Params("id")

	ctx, cancel := getPekerjaanCtx()
	defer cancel()
	if err := s.RestorePekerjaan(ctx, user, id); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil di-restore"})
}
// ==================== SERVICE LOGIC ====================
// ... GetAll, GetByID, GetByAlumniID tidak berubah ...
func (s *PekerjaanService) GetAll(ctx context.Context) ([]model.PekerjaanAlumni, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PekerjaanService) GetByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *PekerjaanService) GetByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error) {
	return s.Repo.GetByAlumniID(ctx, alumniID)
}


// Diubah: Logika bisnis dipindahkan ke sini dari repo
func (s *PekerjaanService) Create(ctx context.Context, req model.CreatePekerjaanRequest) (*model.PekerjaanAlumni, error) {
	// 1. Validasi Input
	if req.AlumniID == "" ||
		req.NamaPerusahaan == "" ||
		req.PosisiJabatan == "" ||
		req.BidangIndustri == "" ||
		req.LokasiKerja == "" ||
		req.TanggalMulaiKerja == "" {
		return nil, errors.New("field required tidak lengkap")
	}

	// DIUBAH: Konversi AlumniID string ke int64
	alumniIntID, err := strconv.ParseInt(req.AlumniID, 10, 64)
	if err != nil {
		return nil, errors.New("alumni_id tidak valid (bukan int64)")
	}

	// 2. Logika Bisnis (Default & Parsing)
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return nil, errors.New("tanggal_mulai_kerja harus dalam format YYYY-MM-DD")
	}

	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return nil, errors.New("tanggal_selesai_kerja harus dalam format YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}

	status := "aktif"
	if req.StatusPekerjaan != nil && *req.StatusPekerjaan != "" {
		status = *req.StatusPekerjaan
	}

	// 3. Buat domain model
	pekerjaan := model.PekerjaanAlumni{
		// ID akan di-set oleh repository (auto-increment)
		AlumniID:            alumniIntID, // DIUBAH
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tanggalMulai,
		TanggalSelesaiKerja: tanggalSelesai,
		StatusPekerjaan:     status,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// 4. Kirim ke Repo
	return s.Repo.Create(ctx, pekerjaan)
}

// ... Update, GetAllWithFilter, SoftDelete, HardDelete, Restore sudah benar ...
// ... Perbandingan `p.AlumniID != *user.AlumniID` sudah benar untuk int64 ...
func (s *PekerjaanService) Update(ctx context.Context, id string, req model.UpdatePekerjaanRequest) (*model.PekerjaanAlumni, error) {
	if req.NamaPerusahaan == "" ||
		req.PosisiJabatan == "" ||
		req.BidangIndustri == "" ||
		req.LokasiKerja == "" ||
		req.TanggalMulaiKerja == "" {
		return nil, errors.New("field required tidak lengkap")
	}

	p, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("pekerjaan tidak ditemukan")
	}

	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return nil, errors.New("tanggal_mulai_kerja harus dalam format YYYY-MM-DD")
	}

	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return nil, errors.New("tanggal_selesai_kerja harus dalam format YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}

	status := "aktif"
	if req.StatusPekerjaan != nil && *req.StatusPekerjaan != "" {
		status = *req.StatusPekerjaan
	}

	p.NamaPerusahaan = req.NamaPerusahaan
	p.PosisiJabatan = req.PosisiJabatan
	p.BidangIndustri = req.BidangIndustri
	p.LokasiKerja = req.LokasiKerja
	p.GajiRange = req.GajiRange
	p.TanggalMulaiKerja = tanggalMulai
	p.TanggalSelesaiKerja = tanggalSelesai
	p.StatusPekerjaan = status
	p.DeskripsiPekerjaan = req.DeskripsiPekerjaan
	p.UpdatedAt = time.Now()

	err = s.Repo.Update(ctx, id, *p)
	return p, err
}

func (s *PekerjaanService) GetAllWithFilter(ctx context.Context, page, limit int, sortBy, sortOrder, search string) (model.PekerjaanResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	allowedSort := map[string]bool{
		"_id":             true,
		"id":              true,
		"alumni_id":       true,
		"nama_perusahaan": true,
		"posisi_jabatan":  true,
		"bidang_industri": true,
		"lokasi_kerja":    true,
		"created_at":      true,
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
		return model.PekerjaanResponse{}, err
	}

	total, err := s.Repo.Count(ctx, search)
	if err != nil {
		return model.PekerjaanResponse{}, err
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return model.PekerjaanResponse{
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

// ==================== DELETE & RESTORE LOGIC ====================

func (s *PekerjaanService) SoftDeletePekerjaan(ctx context.Context, user *model.User, pekerjaanID string) error {
	role := strings.ToLower(user.Role)

	if role == "admin" {
		return s.Repo.SoftDelete(ctx, pekerjaanID)
	}

	if role == "user" {
		if user.AlumniID == nil {
			return errors.New("akun belum terhubung dengan data alumni")
		}

		p, err := s.Repo.GetByID(ctx, pekerjaanID)
		if err != nil {
			return errors.New("pekerjaan tidak ditemukan")
		}
		// Perbandingan int64 (Sudah Benar)
		if p.AlumniID != *user.AlumniID {
			return errors.New("tidak bisa hapus pekerjaan orang lain")
		}
		return s.Repo.SoftDelete(ctx, pekerjaanID)
	}

	return errors.New("akses ditolak")
}

func (s *PekerjaanService) HardDeletePekerjaan(ctx context.Context, user *model.User, pekerjaanID string) error {
	role := strings.ToLower(user.Role)

	p, err := s.Repo.GetByIDIncludeDeleted(ctx, pekerjaanID)
	if err != nil {
		return fmt.Errorf("pekerjaan tidak ditemukan")
	}

	if p.DeletedAt == nil {
		return fmt.Errorf("harus soft delete dulu sebelum hard delete")
	}

	if role == "admin" {
		return s.Repo.HardDelete(ctx, pekerjaanID)
	}

	if role == "user" {
		if user.AlumniID == nil {
			return errors.New("akun belum terhubung dengan data alumni")
		}
		// Perbandingan int64 (Sudah Benar)
		if p.AlumniID != *user.AlumniID {
			return fmt.Errorf("tidak bisa hapus pekerjaan orang lain")
		}
		return s.Repo.HardDelete(ctx, pekerjaanID)
	}

	return fmt.Errorf("role tidak valid")
}

func (s *PekerjaanService) RestorePekerjaan(ctx context.Context, user *model.User, pekerjaanID string) error {
	role := strings.ToLower(user.Role)

	p, err := s.Repo.GetByIDIncludeDeleted(ctx, pekerjaanID)
	if err != nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	if role == "admin" {
		return s.Repo.Restore(ctx, pekerjaanID)
	}

	if role == "user" {
		if user.AlumniID == nil {
			return errors.New("akun belum terhubung dengan data alumni")
		}
		// Perbandingan int64 (Sudah Benar)
		if p.AlumniID != *user.AlumniID {
			return errors.New("tidak bisa restore pekerjaan orang lain")
		}
		return s.Repo.Restore(ctx, pekerjaanID)
	}

	return errors.New("akses ditolak")
}