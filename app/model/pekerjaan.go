package model

import (
	"time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanAlumni struct {
	ID                  int64      `bson:"_id" json:"id"`             // Diubah
	AlumniID            int64      `bson:"alumni_id" json:"alumni_id"` // Diubah
	NamaPerusahaan      string     `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string     `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string     `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string     `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           *string    `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	TanggalMulaiKerja   time.Time  `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string     `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string    `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
	CreatedAt           time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type CreatePekerjaanRequest struct {
	AlumniID            string   `json:"alumni_id"` // Tetap string, akan di-parse di service/handler
	NamaPerusahaan      string   `json:"nama_perusahaan"`
	PosisiJabatan       string   `json:"posisi_jabatan"`
	BidangIndustri      string   `json:"bidang_industri"`
	LokasiKerja         string   `json:"lokasi_kerja"`
	GajiRange           *string  `json:"gaji_range,omitempty"`
	TanggalMulaiKerja   string   `json:"tanggal_mulai_kerja"` 
	TanggalSelesaiKerja *string  `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     *string  `json:"status_pekerjaan,omitempty"`
	DeskripsiPekerjaan  *string  `json:"deskripsi_pekerjaan,omitempty"`
}

type UpdatePekerjaanRequest struct {
	NamaPerusahaan    string   `json:"nama_perusahaan"`
	PosisiJabatan     string   `json:"posisi_jabatan"`
	BidangIndustri    string   `json:"bidang_industri"`
	LokasiKerja       string   `json:"lokasi_kerja"`
	GajiRange         *string  `json:"gaji_range,omitempty"`
	TanggalMulaiKerja string   `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan   *string  `json:"status_pekerjaan,omitempty"`
	DeskripsiPekerjaan *string `json:"deskripsi_pekerjaan,omitempty"`
}

type TrashPekerjaan struct {
	ID                  int64      `bson:"_id" json:"id"`             // Diubah
	AlumniID            int64      `bson:"alumni_id" json:"alumni_id"` // Diubah
	NamaPerusahaan      string     `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string     `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string     `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string     `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           *string    `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	TanggalMulaiKerja   time.Time  `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string     `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string    `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
	DeletedAt           *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}