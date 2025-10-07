package service

import (
	"errors"
	"time"

	"golanjutan/app/model"
	"golanjutan/app/repository"
	"golanjutan/utils"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo   repository.UserRepository
	AlumniRepo repository.AlumniRepository
}

func NewAuthService(userRepo *repository.UserRepository, alumniRepo *repository.AlumniRepository) *AuthService {
	return &AuthService{
		UserRepo:   *userRepo,
		AlumniRepo: *alumniRepo,
	}
}

func (s *AuthService) Login(username, password string) (*model.LoginResponse, error) {
	user, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if !utils.CheckPassword(user.Password, password) {
		return nil, errors.New("username atau password salah")
	}

	// Buat claims sesuai middleware (pakai user_id dan role)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	// Generate JWT token
	token, err := utils.GenerateJWTWithClaims(claims)
	if err != nil {
		return nil, err
	}

	// Balikkan token + user info
	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) Register(username, password, role, nim, nama, jurusan string, angkatan int, tahunLulus int, email, noTelepon, alamat string) (*model.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Step 1: Insert User
	user, err := s.UserRepo.Create(username, string(hashedPassword), role)
	if err != nil {
		return nil, err
	}

	// Step 2: Insert Alumni
	alumniReq := model.CreateAlumniRequest{
		NIM:        nim,
		Nama:       nama,
		Jurusan:    jurusan,
		Angkatan:   angkatan,
		TahunLulus: tahunLulus,
		Email:      email,
		NoTelepon:  &noTelepon, // pakai pointer
		Alamat:     &alamat,    // pakai pointer
	}

	alumniID, err := s.AlumniRepo.Create(alumniReq)
	if err != nil {
		return nil, err
	}

	// Step 3: Update User → set alumni_id
	user.AlumniID = &alumniID
	if err := s.UserRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
