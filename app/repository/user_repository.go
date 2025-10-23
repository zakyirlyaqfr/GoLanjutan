package repository

import (
	"database/sql"
	"golanjutan/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	row := r.DB.QueryRow(`
SELECT id, username, password, role, alumni_id, created_at, updated_at 
FROM users 
WHERE username=$1
`, username)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.AlumniID, &u.CreatedAt, &u.UpdatedAt)
	return &u, err
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	row := r.DB.QueryRow(`
SELECT id, username, password, role, alumni_id, created_at, updated_at 
FROM users 
WHERE id=$1
`, id)

	var u model.User
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.AlumniID, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(username, password, role string) (*model.User, error) {
	var id int
	err := r.DB.QueryRow(
		`INSERT INTO users (username, password, role, created_at, updated_at) 
		 VALUES ($1,$2,$3,NOW(),NOW()) 
		 RETURNING id`,
		username, password, role,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Fetch the created user (lengkap dengan alumni_id)
	var u model.User
	row := r.DB.QueryRow(`
SELECT id, username, password, role, alumni_id, created_at, updated_at 
FROM users 
WHERE id=$1
`, id)

	err = row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.AlumniID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(user *model.User) error {
	_, err := r.DB.Exec(
		`UPDATE users 
		 SET username=$1, password=$2, role=$3, alumni_id=$4, updated_at=NOW() 
		 WHERE id=$5`,
		user.Username, user.Password, user.Role, user.AlumniID, user.ID,
	)
	return err
}
