package service

import (
	"errors"
	"strings"
	"golanjutan/app/model"
	"golanjutan/app/repository"
)

type AlumniService struct {
	Repo *repository.AlumniRepository
}

func NewAlumniService(repo *repository.AlumniRepository) *AlumniService {
	return &AlumniService{Repo: repo}
}

func (s *AlumniService) GetAll() ([]model.Alumni, error) {
	return s.Repo.GetAll()
}

func (s *AlumniService) GetByID(id int) (*model.Alumni, error) {
	return s.Repo.GetByID(id)
}

func (s *AlumniService) Create(req model.CreateAlumniRequest) (int, error) {
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return 0, errors.New("nim, nama, jurusan, dan email harus diisi")
	}
	return s.Repo.Create(req)
}

func (s *AlumniService) Update(id int, req model.UpdateAlumniRequest) error {
	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return errors.New("nama, jurusan, dan email harus diisi")
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
	offset := (page - 1) * limit

	// whitelist kolom sort
	allowedSort := map[string]bool{
		"id": true, "nim": true, "nama": true, "jurusan": true,
		"angkatan": true, "tahun_lulus": true, "created_at": true,
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if strings.ToUpper(sortOrder) != "ASC" {
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

	return model.AlumniResponse{
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

// ✅ hanya superadmin (id=1) yang boleh soft delete alumni
func (s *AlumniService) SoftDeleteAlumni(user *model.User, alumniID int) error {
	if user.ID != 1 {
		return errors.New("hanya superadmin yang bisa menghapus alumni")
	}
	return s.Repo.SoftDelete(alumniID)
}

// ✅ hanya superadmin (id=1) yang boleh restore alumni
func (s *AlumniService) RestoreAlumni(user *model.User, alumniID int) error {
	if user.ID != 1 {
		return errors.New("hanya superadmin yang bisa restore alumni")
	}
	return s.Repo.Restore(alumniID)
}
