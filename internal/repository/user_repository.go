package repository

import (
	"database/sql"
	"errors"
	"fortyfour-backend/internal/models"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String() // <- generate UUID
	}

	query := `INSERT INTO users (id, username, password, email, id_jabatan) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Password, user.Email, user.IDJabatan)
	if err != nil {
		return err
	}

	return nil
}
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, email, id_jabatan, created_at, updated_at 
	          FROM users WHERE username = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.IDJabatan,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	query := `SELECT id, username, password, email, id_jabatan, created_at, updated_at 
	          FROM users WHERE id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.IDJabatan,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET username = ?, email = ?, id_jabatan = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Username, user.Email, user.IDJabatan, user.ID)
	return err
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
