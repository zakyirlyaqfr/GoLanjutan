package repository

import (
	"database/sql"
	"fmt"
	"golanjutan/app/model"
)

type PekerjaanRepository struct {
	DB *sql.DB
}

func NewPekerjaanRepository(db *sql.DB) *PekerjaanRepository {
	return &PekerjaanRepository{DB: db}
}

// ... (GetAll, GetByID, GetByAlumniID - tidak berubah) ...
func (r *PekerjaanRepository) GetAll() ([]model.PekerjaanAlumni, error) {
	rows, err := r.DB.Query(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
        tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
        FROM pekerjaan_alumni
        WHERE deleted_at IS NULL
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
        WHERE id = $1 AND deleted_at IS NULL
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
        WHERE alumni_id = $1 AND deleted_at IS NULL
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


// Diubah: Menerima model.PekerjaanAlumni yang sudah divalidasi dari service
func (r *PekerjaanRepository) Create(p model.PekerjaanAlumni) (int, error) {
	var id int
	err := r.DB.QueryRow(`
        INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
         tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
        RETURNING id
    `, p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange,
		p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, p.CreatedAt, p.UpdatedAt).Scan(&id)
	return id, err
}

// Diubah: Menerima model.PekerjaanAlumni yang sudah divalidasi dari service
func (r *PekerjaanRepository) Update(id int, p model.PekerjaanAlumni) error {
	_, err := r.DB.Exec(`
        UPDATE pekerjaan_alumni
        SET nama_perusahaan=$1, posisi_jabatan=$2, bidang_industri=$3, lokasi_kerja=$4, gaji_range=$5, 
        tanggal_mulai_kerja=$6, tanggal_selesai_kerja=$7, status_pekerjaan=$8, deskripsi_pekerjaan=$9, updated_at=$10
        WHERE id=$11
    `, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange,
		p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, p.UpdatedAt, id)
	return err
}

func (r *PekerjaanRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id = $1`, id)
	return err
}

// ... (GetAllWithFilter, GetTrashed, GetTrashedByAlumniID, Count - tidak berubah) ...
func (r *PekerjaanRepository) GetAllWithFilter(limit, offset int, sortBy, sortOrder, search string) ([]model.PekerjaanAlumni, error) {
	query := `
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
        tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
        FROM pekerjaan_alumni
        WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1) AND deleted_at IS NULL
    `
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", sortBy, sortOrder)

	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var gaji, deskripsi sql.NullString
		var tanggalSelesai sql.NullTime
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &gaji, &p.TanggalMulaiKerja, &tanggalSelesai,
			&p.StatusPekerjaan, &deskripsi, &p.CreatedAt, &p.UpdatedAt); err != nil {
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

func (r *PekerjaanRepository) GetTrashed() ([]model.PekerjaanAlumni, error) {
	rows, err := r.DB.Query(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
        tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
        created_at, updated_at, deleted_at
        FROM pekerjaan_alumni
        WHERE deleted_at IS NOT NULL
        ORDER BY deleted_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var gaji, deskripsi sql.NullString
		var tanggalSelesai, deletedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &gaji, &p.TanggalMulaiKerja, &tanggalSelesai,
			&p.StatusPekerjaan, &deskripsi, &p.CreatedAt, &p.UpdatedAt, &deletedAt); err != nil {
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
		if deletedAt.Valid {
			p.DeletedAt = &deletedAt.Time
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) GetTrashedByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	rows, err := r.DB.Query(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
        tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
        created_at, updated_at, deleted_at
        FROM pekerjaan_alumni
        WHERE deleted_at IS NOT NULL AND alumni_id = $1
        ORDER BY deleted_at DESC
    `, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var gaji, deskripsi sql.NullString
		var tanggalSelesai, deletedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &gaji, &p.TanggalMulaiKerja, &tanggalSelesai,
			&p.StatusPekerjaan, &deskripsi, &p.CreatedAt, &p.UpdatedAt, &deletedAt); err != nil {
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
		if deletedAt.Valid {
			p.DeletedAt = &deletedAt.Time
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Count(search string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
        SELECT COUNT(*) 
        FROM pekerjaan_alumni
        WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1) AND deleted_at IS NULL
    `, "%"+search+"%").Scan(&total)
	return total, err
}


func (r *PekerjaanRepository) SoftDelete(id int) error {
	_, err := r.DB.Exec(`UPDATE pekerjaan_alumni SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// HardDelete menghapus permanen satu pekerjaan
// Logika pengecekan 'deleted_at' dipindah ke service
func (r *PekerjaanRepository) HardDelete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id = $1`, id)
	return err
}

func (r *PekerjaanRepository) Restore(id int) error {
	_, err := r.DB.Exec(`UPDATE pekerjaan_alumni SET deleted_at = NULL WHERE id = $1`, id)
	return err
}

func (r *PekerjaanRepository) GetByIDRow(id int) *sql.Row {
	return r.DB.QueryRow(`SELECT id, alumni_id FROM pekerjaan_alumni WHERE id = $1`, id)
}

func (r *PekerjaanRepository) GetActive() (*sql.Rows, error) {
	return r.DB.Query(`SELECT * FROM pekerjaan_alumni WHERE deleted_at IS NULL`)
}

// GetByIDIncludeDeleted sudah ada dan benar
func (r *PekerjaanRepository) GetByIDIncludeDeleted(id int) (*model.PekerjaanAlumni, error) {
	row := r.DB.QueryRow(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
        tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at
        FROM pekerjaan_alumni
        WHERE id = $1
    `, id)

	var p model.PekerjaanAlumni
	var gaji sql.NullString
	var tanggalSelesai sql.NullTime
	var deskripsi sql.NullString
	var deletedAt sql.NullTime

	if err := row.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &gaji,
		&p.TanggalMulaiKerja, &tanggalSelesai, &p.StatusPekerjaan, &deskripsi, &p.CreatedAt, &p.UpdatedAt, &deletedAt); err != nil {
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
	if deletedAt.Valid {
		p.DeletedAt = &deletedAt.Time
	}

	return &p, nil
}