package repository

import (
	"golanjutan/app/model"
	"database/sql"
	"time"
)

type AlumniRepository struct {
	DB *sql.DB
}

func NewAlumniRepository(db *sql.DB) *AlumniRepository {
	return &AlumniRepository{DB: db}
}

func (r *AlumniRepository) GetAll() ([]model.Alumni, error) {
	rows, err := r.DB.Query(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
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
		WHERE id = $1
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

// repository/alumni_repository.go
func (r *AlumniRepository) GetAllWithFilter(limit, offset int, sortBy, sortOrder, search string) ([]model.Alumni, error) {
	query := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE nama ILIKE $1 OR jurusan ILIKE $1 OR nim ILIKE $1
		ORDER BY ` + sortBy + ` ` + sortOrder + `
		LIMIT $2 OFFSET $3
	`

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
		WHERE nama ILIKE $1 OR jurusan ILIKE $1 OR nim ILIKE $1
	`, "%"+search+"%").Scan(&total)
	return total, err
}
