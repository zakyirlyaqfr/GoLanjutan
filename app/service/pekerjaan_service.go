package service

import (
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"errors"
	"time"
	"strings"
)

type PekerjaanService struct {
	Repo *repository.PekerjaanRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository) *PekerjaanService {
    return &PekerjaanService{Repo: repo}
}

func (s *PekerjaanService) GetAll() ([]model.PekerjaanAlumni, error) {
	return s.Repo.GetAll()
}

func (s *PekerjaanService) GetByID(id int) (*model.PekerjaanAlumni, error) {
	return s.Repo.GetByID(id)
}

func (s *PekerjaanService) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	return s.Repo.GetByAlumniID(alumniID)
}

func (s *PekerjaanService) Create(req model.CreatePekerjaanRequest) (int, error) {
	// minimal validation: tanggal valid
	if req.AlumniID == 0 || req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" || req.TanggalMulaiKerja == "" {
		return 0, errors.New("field required tidak lengkap")
	}
	// validate date format
	if _, err := time.Parse("2006-01-02", req.TanggalMulaiKerja); err != nil {
		return 0, errors.New("tanggal_mulai_kerja harus dalam format YYYY-MM-DD")
	}
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		if _, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja); err != nil {
			return 0, errors.New("tanggal_selesai_kerja harus dalam format YYYY-MM-DD")
		}
	}
	return s.Repo.Create(req)
}

func (s *PekerjaanService) Update(id int, req model.UpdatePekerjaanRequest) error {
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" || req.TanggalMulaiKerja == "" {
		return errors.New("field required tidak lengkap")
	}
	if _, err := time.Parse("2006-01-02", req.TanggalMulaiKerja); err != nil {
		return errors.New("tanggal_mulai_kerja harus dalam format YYYY-MM-DD")
	}
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		if _, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja); err != nil {
			return errors.New("tanggal_selesai_kerja harus dalam format YYYY-MM-DD")
		}
	}
	return s.Repo.Update(id, req)
}

func (s *PekerjaanService) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *PekerjaanService) GetAllWithFilter(page, limit int, sortBy, sortOrder, search string) (model.PekerjaanResponse, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// whitelist kolom sort
	allowedSort := map[string]bool{
		"id": true, "alumni_id": true, "nama_perusahaan": true,
		"posisi_jabatan": true, "bidang_industri": true,
		"lokasi_kerja": true, "created_at": true,
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if strings.ToUpper(sortOrder) != "ASC" {
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

	return model.PekerjaanResponse{
		Data: data,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  sortOrder,
			Search: search,
		},
	}, nil
}