package repository

import "fortyfour-backend/internal/models"

// UserRepositoryInterface defines methods for user data access
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

// PostRepositoryInterface defines methods for post data access
type PostRepositoryInterface interface {
	Create(post *models.Post) error
	FindAll() ([]*models.Post, error)
	FindByID(id int) (*models.Post, error)
	FindByAuthorID(authorID int) ([]*models.Post, error)
	Update(post *models.Post) error
	Delete(id int) error
}
