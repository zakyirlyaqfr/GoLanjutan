package model

import (
	"time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        int64      `bson:"_id" json:"id"`                                // Diubah
	Username  string     `bson:"username" json:"user"`
	Password  string     `bson:"password" json:"-"`
	Role      string     `bson:"role" json:"role"`
	AlumniID  *int64     `bson:"alumni_id,omitempty" json:"alumni_id,omitempty"` // Diubah
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// Request struct (LoginRequest, RegisterRequest, LoginResponse) tidak perlu tag bson
// ... (sisa file user.go tidak berubah) ...

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	NIM        string `json:"nim"`
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	Email      string `json:"email"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
	Role  string `json:"role"`
}