package service

import (
	"errors"
	"time"

	"golanjutan/app/model"
	"golanjutan/app/repository"
	"golanjutan/utils"

	"github.com/gofiber/fiber/v2"
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

// =============================
// HANDLER (dipanggil dari route)
// =============================
func (s *AuthService) HandleLogin(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	res, err := s.Login(req.Username, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (s *AuthService) HandleRegister(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	user, err := s.Register(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// =============================
// CORE LOGIC
// =============================

func (s *AuthService) Login(username, password string) (*model.LoginResponse, error) {
	user, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if !utils.CheckPassword(user.Password, password) {
		return nil, errors.New("username atau password salah")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	if user.AlumniID != nil {
		claims["alumni_id"] = *user.AlumniID
	}

	token, err := utils.GenerateJWTWithClaims(claims)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) Register(req model.RegisterRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepo.Create(req.Username, string(hashedPassword), req.Role)
	if err != nil {
		return nil, err
	}

	alumniReq := model.CreateAlumniRequest{
		NIM:         req.NIM,
		Nama:        req.Nama,
		Jurusan:     req.Jurusan,
		Angkatan:    req.Angkatan,
		TahunLulus:  req.TahunLulus,
		Email:       req.Email,
		NoTelepon:   &req.NoTelepon,
		Alamat:      &req.Alamat,
	}

	alumniID, err := s.AlumniRepo.Create(alumniReq)
	if err != nil {
		return nil, err
	}

	user.AlumniID = &alumniID
	if err := s.UserRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
