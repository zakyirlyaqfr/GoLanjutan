package service

import (
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"errors"
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
	// simple validation
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
