package repository

import (
	"golanjutan/app/model"
	"database/sql"
	"time"
)

type PekerjaanRepository struct {
	DB *sql.DB
}

func NewPekerjaanRepository(db *sql.DB) *PekerjaanRepository {
    return &PekerjaanRepository{DB: db}
}

func (r *PekerjaanRepository) GetAll() ([]model.PekerjaanAlumni, error) {
	rows, err := r.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var gaji sql.NullString
		var tanggalSelesai sql.NullTime
		var deskripsi sql.NullString
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, 
			&p.LokasiKerja, &gaji, &p.TanggalMulaiKerja, &tanggalSelesai, &p.StatusPekerjaan, &deskripsi, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if gaji.Valid {
			p.GajiRange = &gaji.String
		}
		if tanggalSelesai.Valid {
			p.TanggalSelesaiKerja = &tanggalSelesai.Time
		}
		if deskripsi.Valid {
			p.DeskripsiPekerjaan = &deskripsi.String
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) GetByID(id int) (*model.PekerjaanAlumni, error) {
	row := r.DB.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE id = $1
	`, id)

	var p model.PekerjaanAlumni
	var gaji sql.NullString
	var tanggalSelesai sql.NullTime
	var deskripsi sql.NullString
	if err := row.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &gaji,
		 &p.TanggalMulaiKerja, &tanggalSelesai, &p.StatusPekerjaan, &deskripsi, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	if gaji.Valid {
		p.GajiRange = &gaji.String
	}
	if tanggalSelesai.Valid {
		p.TanggalSelesaiKerja = &tanggalSelesai.Time
	}
	if deskripsi.Valid {
		p.DeskripsiPekerjaan = &deskripsi.String
	}
	return &p, nil
}

func (r *PekerjaanRepository) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	rows, err := r.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE alumni_id = $1
		ORDER BY created_at DESC
	`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var gaji sql.NullString
		var tanggalSelesai sql.NullTime
		var deskripsi sql.NullString
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &gaji,
			 &p.TanggalMulaiKerja, &tanggalSelesai, &p.StatusPekerjaan, &deskripsi, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if gaji.Valid {
			p.GajiRange = &gaji.String
		}
		if tanggalSelesai.Valid {
			p.TanggalSelesaiKerja = &tanggalSelesai.Time
		}
		if deskripsi.Valid {
			p.DeskripsiPekerjaan = &deskripsi.String
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Create(req model.CreatePekerjaanRequest) (int, error) {
	var id int
	var tanggalSelesai interface{} = nil
	if req.TanggalSelesaiKerja != nil {
		tanggalSelesai = *req.TanggalSelesaiKerja
	}
	status := "aktif"
	if req.StatusPekerjaan != nil && *req.StatusPekerjaan != "" {
		status = *req.StatusPekerjaan
	}
	err := r.DB.QueryRow(`
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		 tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id
	`, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja, req.GajiRange, 
	req.TanggalMulaiKerja, tanggalSelesai, status, req.DeskripsiPekerjaan, time.Now(), time.Now()).Scan(&id)
	return id, err
}

func (r *PekerjaanRepository) Update(id int, req model.UpdatePekerjaanRequest) error {
	var tanggalSelesai interface{} = nil
	if req.TanggalSelesaiKerja != nil {
		tanggalSelesai = *req.TanggalSelesaiKerja
	}
	status := "aktif"
	if req.StatusPekerjaan != nil && *req.StatusPekerjaan != "" {
		status = *req.StatusPekerjaan
	}
	_, err := r.DB.Exec(`
		UPDATE pekerjaan_alumni
		SET nama_perusahaan=$1, posisi_jabatan=$2, bidang_industri=$3, lokasi_kerja=$4, gaji_range=$5, 
		tanggal_mulai_kerja=$6, tanggal_selesai_kerja=$7, status_pekerjaan=$8, deskripsi_pekerjaan=$9, updated_at=$10
		WHERE id=$11
	`, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja, req.GajiRange, 
	req.TanggalMulaiKerja, tanggalSelesai, status, req.DeskripsiPekerjaan, time.Now(), id)
	return err
}

func (r *PekerjaanRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id = $1`, id)
	return err
}