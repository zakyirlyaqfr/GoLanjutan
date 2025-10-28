package model

import (
	"time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	ID         int64      `bson:"_id" json:"id"` // Diubah dari ObjectID ke int64 dan omitempty dihapus
	NIM        string     `bson:"nim" json:"nim"`
	Nama       string     `bson:"nama" json:"nama"`
	Jurusan    string     `bson:"jurusan" json:"jurusan"`
	Angkatan   int        `bson:"angkatan" json:"angkatan"`
	TahunLulus int        `bson:"tahun_lulus" json:"tahun_lulus"`
	Email      string     `bson:"email" json:"email"`
	NoTelepon  *string    `bson:"no_telepon,omitempty" json:"no_telepon,omitempty"`
	Alamat     *string    `bson:"alamat,omitempty" json:"alamat,omitempty"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// CreateAlumniRequest dan UpdateAlumniRequest tidak perlu diubah
// karena mereka merepresentasikan JSON request body, bukan skema DB.
type CreateAlumniRequest struct {
	NIM        string  `json:"nim"`
	Nama       string  `json:"nama"`
	Jurusan    string  `json:"jurusan"`
	Angkatan   int     `json:"angkatan"`
	TahunLulus int     `json:"tahun_lulus"`
	Email      string  `json:"email"`
	NoTelepon  *string `json:"no_telepon,omitempty"`
	Alamat     *string `json:"alamat,omitempty"`
}

type UpdateAlumniRequest struct {
	Nama       string  `json:"nama"`
	Jurusan    string  `json:"jurusan"`
	Angkatan   int     `json:"angkatan"`
	TahunLulus int     `json:"tahun_lulus"`
	Email      string  `json:"email"`
	NoTelepon  *string `json:"no_telepon,omitempty"`
	Alamat     *string `json:"alamat,omitempty"`
}

type TrashAlumni struct {
	ID         int64      `bson:"_id" json:"id"` // Diubah
	NIM        string     `bson:"nim" json:"nim"`
	Nama       string     `bson:"nama" json:"nama"`
	Jurusan    string     `bson:"jurusan" json:"jurusan"`
	Angkatan   int        `bson:"angkatan" json:"angkatan"`
	TahunLulus int        `bson:"tahun_lulus" json:"tahun_lulus"`
	Email      string     `bson:"email" json:"email"`
	NoTelepon  *string    `bson:"no_telepon,omitempty" json:"no_telepon,omitempty"`
	Alamat     *string    `bson:"alamat,omitempty" json:"alamat,omitempty"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}