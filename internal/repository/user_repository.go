package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	query := `UPDATE users SET email = $1, password_hash = $2, role = $3 WHERE id = $4`
	_, err := r.db.Exec(query, user.Email, user.PasswordHash, user.Role, user.ID)
	return err
}

