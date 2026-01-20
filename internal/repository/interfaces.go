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
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	UpdateWithPhoto(user *models.User) error
	UpdatePassword(id, hashedPassword string) error
	GetPasswordByID(id string) (string, error)
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

// PerusahaanRepositoryInterface
type PerusahaanRepositoryInterface interface {
	Create(req dto.CreatePerusahaanRequest, id string) error
	GetByID(id string) (*dto.PerusahaanResponse, error)
	GetAll() ([]dto.PerusahaanResponse, error)
	Update(id string, perusahaan dto.PerusahaanResponse) error
	Delete(id string) error
}

// PICPerusahaanRepositoryInterface
type PICRepositoryInterface interface {
	Create(req dto.CreatePICRequest, id string) error
	GetByID(id string) (*dto.PICResponse, error)
	GetAll() ([]dto.PICResponse, error)
	Update(id string, req dto.UpdatePICRequest) error
	Delete(id string) error
}

// IdentifikasiRepositoryInterface
type IdentifikasiRepositoryInterface interface {
	Create(req dto.CreateIdentifikasiRequest, id string) error
	GetAll() ([]models.Identifikasi, error)
	GetByID(id string) (*models.Identifikasi, error)
	Update(id string, identifikasi models.Identifikasi) error
	Delete(id string) error
}

// ProteksiRepositoryInterface
type ProteksiRepositoryInterface interface {
	Create(req dto.CreateProteksiRequest, id string) error
	GetAll() ([]models.Proteksi, error)
	GetByID(id string) (*models.Proteksi, error)
	Update(id string, proteksi models.Proteksi) error
	Delete(id string) error
}

// DeteksiRepositoryInterface
type DeteksiRepositoryInterface interface {
	Create(req dto.CreateDeteksiRequest, id string) error
	GetAll() ([]models.Deteksi, error)
	GetByID(id string) (*models.Deteksi, error)
	Update(id string, deteksi models.Deteksi) error
	Delete(id string) error
}

// GulihRepositoryInterface
type GulihRepositoryInterface interface {
	Create(req dto.CreateGulihRequest, id string) error
	GetAll() ([]models.Gulih, error)
	GetByID(id string) (*models.Gulih, error)
	Update(id string, gulih models.Gulih) error
	Delete(id string) error
}

// SektorRepositoryInterface
type SektorRepositoryInterface interface {
	GetAll() ([]dto.SektorResponse, error)
	GetByID(id string) (*dto.SektorResponse, error)
}

// SubSektorRepositoryInterface
type SubSektorRepositoryInterface interface {
	GetAll() ([]dto.SubSektorResponse, error)
	GetByID(id string) (*dto.SubSektorResponse, error)
	GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error)
}
// SERepositoryInterface
type SERepositoryInterface interface {
	Create(req dto.CreateSERequest, id string, totalBobot int, kategori string) error
	GetAll() ([]dto.SEResponse, error)
	GetByID(id string) (*dto.SEResponse, error)
	Update(id string, req dto.UpdateSERequest, totalBobot int, kategori string) error
	Delete(id string) error
}