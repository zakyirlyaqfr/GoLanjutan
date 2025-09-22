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
SELECT id, username, password, role, created_at, updated_at FROM users WHERE username=$1
`, username)
var u model.User
if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
return nil, err
}
return &u, nil
}


func (r *UserRepository) GetByID(id int) (*model.User, error) {
row := r.DB.QueryRow(`SELECT id, username, password, role, created_at, updated_at FROM users WHERE id=$1`, id)
var u model.User
if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
return nil, err
}
return &u, nil
}


func (r *UserRepository) Create(username, password, role string) (int, error) {
var id int
err := r.DB.QueryRow(`INSERT INTO users (username, password, role, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`, username, password, role, sql.NullTime{}, sql.NullTime{}).Scan(&id)
return id, err
}