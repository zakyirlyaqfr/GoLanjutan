package service

import (
	"errors"
	"time"

	"golanjutan/app/model"
	"golanjutan/app/repository"
	"golanjutan/utils"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: repo,
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

func (s *AuthService) Register(username, password, role string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = s.UserRepo.Create(username, hashed, role)
	return err
}
