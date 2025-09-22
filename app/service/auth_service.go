package service

import (
	"errors"
	"time"

	"golanjutan/app/repository"
	"golanjutan/utils"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
    return &AuthService{
        UserRepo: repo,
    }
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	if !utils.CheckPassword(user.Password, password) {
		return "", errors.New("username atau password salah")
	}

	token, err := utils.GenerateJWT(user.Username, user.Role, time.Hour*24)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(username, password, role string) error {
    hashed, err := utils.HashPassword(password)
    if err != nil {
        return err
    }
	// Sekarang asumsi UserRepo.Create() return-nya hanya error
	_, err = s.UserRepo.Create(username, hashed, role)
	return err
}

