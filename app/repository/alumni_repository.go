package repository

import (
	"database/sql"
	"fmt"
	"golanjutan/app/model"
	"time"
)

type AlumniRepository struct {
	DB *sql.DB
}

func NewAlumniRepository(db *sql.DB) *AlumniRepository {
	return &AlumniRepository{DB: db}
}

// ... (GetAll, GetByID, GetTrashed - tidak berubah) ...

func (r *AlumniRepository) GetAll() ([]model.Alumni, error) {
	rows, err := r.DB.Query(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		var noTel sql.NullString
		var alamat sql.NullString
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &noTel, &alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		if noTel.Valid {
			a.NoTelepon = &noTel.String
		}
		if alamat.Valid {
			a.Alamat = &alamat.String
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AlumniRepository) GetByID(id int) (*model.Alumni, error) {
	row := r.DB.QueryRow(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE id = $1 AND deleted_at IS NULL
    `, id)

	var a model.Alumni
	var noTel sql.NullString
	var alamat sql.NullString
	if err := row.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &noTel, &alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return nil, err
	}
	if noTel.Valid {
		a.NoTelepon = &noTel.String
	}
	if alamat.Valid {
		a.Alamat = &alamat.String
	}
	return &a, nil
}

// Ditambahkan: GetByID termasuk yang sudah di soft-delete (untuk pengecekan di service)
func (r *AlumniRepository) GetByIDIncludeDeleted(id int) (*model.Alumni, error) {
	row := r.DB.QueryRow(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at, deleted_at
        FROM alumni
        WHERE id = $1
    `, id)

	var a model.Alumni
	var noTel, alamat sql.NullString
	var deletedAt sql.NullTime
	if err := row.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &noTel, &alamat, &a.CreatedAt, &a.UpdatedAt, &deletedAt); err != nil {
		return nil, err
	}
	if noTel.Valid {
		a.NoTelepon = &noTel.String
	}
	if alamat.Valid {
		a.Alamat = &alamat.String
	}
	if deletedAt.Valid {
		a.DeletedAt = &deletedAt.Time
	}
	return &a, nil
}

func (r *AlumniRepository) GetTrashed() ([]model.Alumni, error) {
	rows, err := r.DB.Query(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at, deleted_at
        FROM alumni
        WHERE deleted_at IS NOT NULL
        ORDER BY deleted_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		var noTel, alamat sql.NullString
		var deletedAt sql.NullTime
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
			&a.Email, &noTel, &alamat, &a.CreatedAt, &a.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		if noTel.Valid {
			a.NoTelepon = &noTel.String
		}
		if alamat.Valid {
			a.Alamat = &alamat.String
		}
		if deletedAt.Valid {
			a.DeletedAt = &deletedAt.Time
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AlumniRepository) Create(req model.CreateAlumniRequest) (int, error) {
	var id int
	err := r.DB.QueryRow(`
        INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `, req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), time.Now()).Scan(&id)
	return id, err
}

func (r *AlumniRepository) Update(id int, req model.UpdateAlumniRequest) error {
	_, err := r.DB.Exec(`
        UPDATE alumni SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, email=$5, no_telepon=$6, alamat=$7, updated_at=$8
        WHERE id=$9
    `, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), id)
	return err
}

func (r *AlumniRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM alumni WHERE id = $1`, id)
	return err
}

// ... (GetAllWithFilter, Count - tidak berubah) ...
func (r *AlumniRepository) GetAllWithFilter(limit, offset int, sortBy, sortOrder, search string) ([]model.Alumni, error) {
	query := `
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE (nama ILIKE $1 OR jurusan ILIKE $1 OR nim ILIKE $1) AND deleted_at IS NULL
    `

	// Hati-hati dengan SQL Injection, tapi ini mengikuti kode asli Anda.
	// Sebaiknya, kolom sortBy dan sortOrder divalidasi di service.
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", sortBy, sortOrder)

	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		var noTel, alamat sql.NullString
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
			&a.Email, &noTel, &alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		if noTel.Valid {
			a.NoTelepon = &noTel.String
		}
		if alamat.Valid {
			a.Alamat = &alamat.String
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AlumniRepository) Count(search string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
        SELECT COUNT(*) 
        FROM alumni
        WHERE (nama ILIKE $1 OR jurusan ILIKE $1 OR nim ILIKE $1) AND deleted_at IS NULL
    `, "%"+search+"%").Scan(&total)
	return total, err
}


func (r *AlumniRepository) SoftDelete(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// soft delete alumni
	_, err = tx.Exec(`UPDATE alumni SET deleted_at = NOW() WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// soft delete pekerjaan terkait
	_, err = tx.Exec(`UPDATE pekerjaan_alumni SET deleted_at = NOW() WHERE alumni_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// HardDelete menghapus permanen alumni dan semua pekerjaan terkait
// Logika pengecekan 'deleted_at' dipindah ke service
func (r *AlumniRepository) HardDelete(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// hapus semua pekerjaan milik alumni
	_, err = tx.Exec(`DELETE FROM pekerjaan_alumni WHERE alumni_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// hapus permanen alumni
	_, err = tx.Exec(`DELETE FROM alumni WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit()
}

func (r *AlumniRepository) Restore(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// restore alumni
	_, err = tx.Exec(`UPDATE alumni SET deleted_at = NULL WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// restore pekerjaan terkait
	_, err = tx.Exec(`UPDATE pekerjaan_alumni SET deleted_at = NULL WHERE alumni_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *AlumniRepository) GetAllActive() ([]model.Alumni, error) {
    // Implementasi ini identik dengan GetAll(). Anda bisa hapus salah satu jika sama.
	return r.GetAll()
}