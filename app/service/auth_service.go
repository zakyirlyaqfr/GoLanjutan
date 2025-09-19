package service


import (
"errors"
"time"


"golanjutan/app/model"
"golanjutan/app/repository"
"golanjutan/config"


"golang.org/x/crypto/bcrypt"
"github.com/golang-jwt/jwt/v5"
)


type AuthService struct {
UserRepo *repository.UserRepository
}


func NewAuthService(ur *repository.UserRepository) *AuthService {
return &AuthService{UserRepo: ur}
}


func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
user, err := s.UserRepo.GetByUsername(req.Username)
if err != nil {
return nil, errors.New("username atau password salah")
}
// password stored may be pgcrypto crypt result or bcrypt; try bcrypt first
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
// try pgcrypto crypt check is not possible here; assume bcrypt only for Go-based seeder
return nil, errors.New("username atau password salah")
}


// create token
secret := config.AppEnv.JWTSecret
claims := jwt.MapClaims{
"user_id": user.ID,
"username": user.Username,
"role": user.Role,
"exp": time.Now().Add(24 * time.Hour).Unix(),
}
t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signed, err := t.SignedString([]byte(secret))
if err != nil {
return nil, err
}
user.Password = ""
return &model.LoginResponse{Token: signed, User: user}, nil
}


// helper to create user with bcrypt hashed password (for seeding via Go)
func (s *AuthService) CreateUser(username, password, role string) (int, error) {
h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
return 0, err
}
return s.UserRepo.Create(username, string(h), role)
}