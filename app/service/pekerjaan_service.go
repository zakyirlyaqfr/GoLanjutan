package service

import (
	"errors"
	"fmt"
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ==================== STRUCT ====================

type PekerjaanService struct {
	Repo *repository.PekerjaanRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{Repo: repo}
}

// ==================== HANDLER ====================
// (Handler sudah ada, tidak perlu diubah)
func (s *PekerjaanService) HandleGetAll(c *fiber.Ctx) error {
	data, err := s.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleGetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	data, err := s.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleGetByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid alumni_id"})
	}
	data, err := s.GetByAlumniID(alumniID)
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

	res, err := s.GetAllWithFilter(page, limit, sortBy, sortOrder, search)
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

	id, err := s.Create(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	data, _ := s.GetByID(id)
	return c.Status(201).JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleUpdate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := s.Update(id, req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	data, _ := s.GetByID(id)
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) HandleSoftDelete(c *fiber.Ctx) error {
	user, _ := getUserFromContext(c)
	id, _ := strconv.Atoi(c.Params("id"))

	if err := s.SoftDeletePekerjaan(user, id); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil di-soft delete"})
}

func (s *PekerjaanService) HandleHardDelete(c *fiber.Ctx) error {
	user, _ := getUserFromContext(c)
	id, _ := strconv.Atoi(c.Params("id"))

	if err := s.HardDeletePekerjaan(user, id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus permanen"})
}

func (s *PekerjaanService) HandleRestore(c *fiber.Ctx) error {
	user, _ := getUserFromContext(c)
	id, _ := strconv.Atoi(c.Params("id"))

	if err := s.RestorePekerjaan(user, id); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil di-restore"})
}


// ==================== SERVICE LOGIC ====================

func (s *PekerjaanService) GetAll() ([]model.PekerjaanAlumni, error) {
	return s.Repo.GetAll()
}

func (s *PekerjaanService) GetByID(id int) (*model.PekerjaanAlumni, error) {
	return s.Repo.GetByID(id)
}

func (s *PekerjaanService) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	return s.Repo.GetByAlumniID(alumniID)
}

// Diubah: Logika bisnis dipindahkan ke sini dari repo
func (s *PekerjaanService) Create(req model.CreatePekerjaanRequest) (int, error) {
	// 1. Validasi Input
	if req.AlumniID == 0 ||
		req.NamaPerusahaan == "" ||
		req.PosisiJabatan == "" ||
		req.BidangIndustri == "" ||
		req.LokasiKerja == "" ||
		req.TanggalMulaiKerja == "" {
		return 0, errors.New("field required tidak lengkap")
	}

	// 2. Logika Bisnis (Default & Parsing)
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return 0, errors.New("tanggal_mulai_kerja harus dalam format YYYY-MM-DD")
	}

	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return 0, errors.New("tanggal_selesai_kerja harus dalam format YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}
	
	status := "aktif"
	if req.StatusPekerjaan != nil && *req.StatusPekerjaan != "" {
		status = *req.StatusPekerjaan
	}

	// 3. Buat domain model
	pekerjaan := model.PekerjaanAlumni{
		AlumniID:          req.AlumniID,
		NamaPerusahaan:    req.NamaPerusahaan,
		PosisiJabatan:     req.PosisiJabatan,
		BidangIndustri:    req.BidangIndustri,
		LokasiKerja:       req.LokasiKerja,
		GajiRange:         req.GajiRange,
		TanggalMulaiKerja: tanggalMulai,
		TanggalSelesaiKerja: tanggalSelesai,
		StatusPekerjaan:   status,
		DeskripsiPekerjaan: req.DeskripsiPekerjaan,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// 4. Kirim ke Repo
	return s.Repo.Create(pekerjaan)
}

// Diubah: Logika bisnis dipindahkan ke sini dari repo
func (s *PekerjaanService) Update(id int, req model.UpdatePekerjaanRequest) error {
	// 1. Validasi Input
	if req.NamaPerusahaan == "" ||
		req.PosisiJabatan == "" ||
		req.BidangIndustri == "" ||
		req.LokasiKerja == "" ||
		req.TanggalMulaiKerja == "" {
		return errors.New("field required tidak lengkap")
	}

	// Cek dulu apakah data ada
	p, err := s.Repo.GetByID(id)
	if err != nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	// 2. Logika Bisnis (Default & Parsing)
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return errors.New("tanggal_mulai_kerja harus dalam format YYYY-MM-DD")
	}

	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return errors.New("tanggal_selesai_kerja harus dalam format YYYY-MM-DD")
		}
		tanggalSelesai = &t
	}
	
	status := "aktif"
	if req.StatusPekerjaan != nil && *req.StatusPekerjaan != "" {
		status = *req.StatusPekerjaan
	}

	// 3. Update domain model
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

	// 4. Kirim ke Repo
	return s.Repo.Update(id, *p)
}


func (s *PekerjaanService) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *PekerjaanService) GetAllWithFilter(page, limit int, sortBy, sortOrder, search string) (model.PekerjaanResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	allowedSort := map[string]bool{
		"id":              true,
		"alumni_id":       true,
		"nama_perusahaan": true,
		"posisi_jabatan":  true,
		"bidang_industri": true,
		"lokasi_kerja":    true,
		"created_at":      true,
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if strings.ToUpper(sortOrder) != "ASC" && strings.ToUpper(sortOrder) != "DESC" {
		sortOrder = "DESC"
	}

	data, err := s.Repo.GetAllWithFilter(limit, offset, sortBy, sortOrder, search)
	if err != nil {
		return model.PekerjaanResponse{}, err
	}

	total, err := s.Repo.Count(search)
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
// (Logika ini sudah benar, logika bisnis [pengecekan] sudah ada di service)

func (s *PekerjaanService) SoftDeletePekerjaan(user *model.User, pekerjaanID int) error {
	role := strings.ToLower(user.Role)

	if role == "admin" {
		return s.Repo.SoftDelete(pekerjaanID)
	}

	if role == "user" {
		if user.AlumniID == nil {
			return errors.New("akun belum terhubung dengan data alumni")
		}

		p, err := s.Repo.GetByID(pekerjaanID)
		if err != nil {
			return errors.New("pekerjaan tidak ditemukan")
		}
		if p.AlumniID != *user.AlumniID {
			return errors.New("tidak bisa hapus pekerjaan orang lain")
		}
		return s.Repo.SoftDelete(pekerjaanID)
	}

	return errors.New("akses ditolak")
}

func (s *PekerjaanService) HardDeletePekerjaan(user *model.User, pekerjaanID int) error {
	role := strings.ToLower(user.Role)

	// Logika bisnis: Pengecekan data
	p, err := s.Repo.GetByIDIncludeDeleted(pekerjaanID)
	if err != nil {
		return fmt.Errorf("pekerjaan tidak ditemukan")
	}

	// Logika bisnis: Pengecekan status
	if p.DeletedAt == nil {
		return fmt.Errorf("harus soft delete dulu sebelum hard delete")
	}

	// Logika otorisasi
	if role == "admin" {
		return s.Repo.HardDelete(pekerjaanID)
	}

	if role == "user" {
		if user.AlumniID == nil {
			return errors.New("akun belum terhubung dengan data alumni")
		}
		if p.AlumniID != *user.AlumniID {
			return fmt.Errorf("tidak bisa hapus pekerjaan orang lain")
		}
		return s.Repo.HardDelete(pekerjaanID)
	}

	return fmt.Errorf("role tidak valid")
}

func (s *PekerjaanService) RestorePekerjaan(user *model.User, pekerjaanID int) error {
	role := strings.ToLower(user.Role)

	p, err := s.Repo.GetByIDIncludeDeleted(pekerjaanID)
	if err != nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	if role == "admin" {
		return s.Repo.Restore(pekerjaanID)
	}

	if role == "user" {
		if user.AlumniID == nil {
			return errors.New("akun belum terhubung dengan data alumni")
		}
		if p.AlumniID != *user.AlumniID {
			return errors.New("tidak bisa restore pekerjaan orang lain")
		}
		return s.Repo.Restore(pekerjaanID)
	}

	return errors.New("akses ditolak")
}

// // ==================== HELPER FUNCTION ====================
// // (Helper sudah ada, tidak perlu diubah)
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