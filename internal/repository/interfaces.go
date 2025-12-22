package repository

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

// UserRepositoryInterface defines methods for user data access
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	EmailExists(email string, excludeID *string) (bool, error)
	UsernameExists(username string, excludeID *string) (bool, error)
}

// PostRepositoryInterface defines methods for post data access
type PostRepositoryInterface interface {
	Create(post *models.Post) error
	FindAll() ([]*models.Post, error)
	FindByID(id int) (*models.Post, error)
	FindByAuthorID(authorID string) ([]*models.Post, error)
	Update(post *models.Post) error
	Delete(id int) error
}

// JabatanRepositoryInterface defines methods for jabatan data access
type JabatanRepositoryInterface interface {
	Create(req dto.CreateJabatanRequest, id string) error
	GetAll() ([]dto.JabatanResponse, error)
	GetByID(id string) (*dto.JabatanResponse, error)
	Update(id string, jabatan dto.JabatanResponse) error
	Delete(id string) error
}
