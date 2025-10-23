package model


import "time"


type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"user"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	AlumniID  *int       `json:"alumni_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

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
	 Role     string `json:"role"`
}