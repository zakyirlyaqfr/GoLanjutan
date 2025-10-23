package model

import "time"

type PekerjaanAlumni struct {
	ID                 int        `json:"id"`
	AlumniID           int        `json:"alumni_id"`
	NamaPerusahaan     string     `json:"nama_perusahaan"`
	PosisiJabatan      string     `json:"posisi_jabatan"`
	BidangIndustri     string     `json:"bidang_industri"`
	LokasiKerja        string     `json:"lokasi_kerja"`
	GajiRange          *string    `json:"gaji_range,omitempty"`
	TanggalMulaiKerja  time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan    string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan *string    `json:"deskripsi_pekerjaan,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type CreatePekerjaanRequest struct {
	AlumniID           int      `json:"alumni_id"`
	NamaPerusahaan     string   `json:"nama_perusahaan"`
	PosisiJabatan      string   `json:"posisi_jabatan"`
	BidangIndustri     string   `json:"bidang_industri"`
	LokasiKerja        string   `json:"lokasi_kerja"`
	GajiRange          *string  `json:"gaji_range,omitempty"`
	TanggalMulaiKerja  string   `json:"tanggal_mulai_kerja"` // expect "YYYY-MM-DD"
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty"` // optional "YYYY-MM-DD"
	StatusPekerjaan    *string  `json:"status_pekerjaan,omitempty"`
	DeskripsiPekerjaan *string  `json:"deskripsi_pekerjaan,omitempty"`
}

type UpdatePekerjaanRequest struct {
	NamaPerusahaan     string   `json:"nama_perusahaan"`
	PosisiJabatan      string   `json:"posisi_jabatan"`
	BidangIndustri     string   `json:"bidang_industri"`
	LokasiKerja        string   `json:"lokasi_kerja"`
	GajiRange          *string  `json:"gaji_range,omitempty"`
	TanggalMulaiKerja  string   `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan    *string  `json:"status_pekerjaan,omitempty"`
	DeskripsiPekerjaan *string  `json:"deskripsi_pekerjaan,omitempty"`
}

type TrashPekerjaan struct {
	ID                int        `json:"id"`
	AlumniID          int        `json:"alumni_id"`
	NamaPerusahaan    string     `json:"nama_perusahaan"`
	PosisiJabatan     string     `json:"posisi_jabatan"`
	BidangIndustri    string     `json:"bidang_industri"`
	LokasiKerja       string     `json:"lokasi_kerja"`
	GajiRange         *string    `json:"gaji_range,omitempty"`
	TanggalMulaiKerja time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan   string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan *string   `json:"deskripsi_pekerjaan,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at"`
}