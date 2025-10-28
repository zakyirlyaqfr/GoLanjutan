package service

import (
	"context"
	"errors"
	"time"

	"golanjutan/app/model"
	"golanjutan/app/repository"
	"golanjutan/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	// "go.mongodb.org/mongo-driver/bson/primitive" // Tidak perlu
)

// Helper context
func getAuthCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

type AuthService struct {
	UserRepo   repository.IUserRepository
	AlumniRepo repository.IAlumniRepository
}

func NewAuthService(userRepo repository.IUserRepository, alumniRepo repository.IAlumniRepository) *AuthService {
	return &AuthService{
		UserRepo:   userRepo,
		AlumniRepo: alumniRepo,
	}
}

// =============================
// HANDLER (dipanggil dari route)
// =============================
// ... HandleLogin dan HandleRegister tidak berubah ...
func (s *AuthService) HandleLogin(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	ctx, cancel := getAuthCtx()
	defer cancel()
	res, err := s.Login(ctx, req.Username, req.Password)
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

	ctx, cancel := getAuthCtx()
	defer cancel()
	user, err := s.Register(ctx, req)
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

func (s *AuthService) Login(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	user, err := s.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if !utils.CheckPassword(user.Password, password) {
		return nil, errors.New("username atau password salah")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID, // DIUBAH: Kirim int64 langsung (bukan .Hex())
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	if user.AlumniID != nil {
		claims["alumni_id"] = *user.AlumniID // DIUBAH: Kirim int64 langsung (bukan .Hex())
	}

	token, err := utils.GenerateJWTWithClaims(claims)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
		Role:  user.Role,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 1. Buat Model Alumni
	alumniReq := model.CreateAlumniRequest{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  &req.NoTelepon,
		Alamat:     &req.Alamat,
	}
	alumni := model.Alumni{
		// ID akan di-set oleh repository (auto-increment)
		NIM:        alumniReq.NIM,
		Nama:       alumniReq.Nama,
		Jurusan:    alumniReq.Jurusan,
		Angkatan:   alumniReq.Angkatan,
		TahunLulus: alumniReq.TahunLulus,
		Email:      alumniReq.Email,
		NoTelepon:  alumniReq.NoTelepon,
		Alamat:     alumniReq.Alamat,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	createdAlumni, err := s.AlumniRepo.Create(ctx, alumni)
	if err != nil {
		return nil, err
	}

	// 2. Buat Model User
	user := model.User{
		// ID akan di-set oleh repository (auto-increment)
		Username:  req.Username,
		Password:  string(hashedPassword),
		Role:      req.Role,
		AlumniID:  &createdAlumni.ID, // Hubungkan ID Alumni (ini sudah benar, int64)
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.UserRepo.Create(ctx, user)
	if err != nil {
		// Rollback?
		return nil, err
	}

	return createdUser, nil
}