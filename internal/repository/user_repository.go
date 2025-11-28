package repository

import (
	"database/sql"
	"errors"
	"fortyfour-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, password, email) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, email, created_at, updated_at 
	          FROM users WHERE username = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
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

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `SELECT id, username, password, email, created_at, updated_at 
	          FROM users WHERE id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
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
	query := `UPDATE users SET username = ?, email = ? WHERE id = ?`

	_, err := r.db.Exec(query, user.Username, user.Email, user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
