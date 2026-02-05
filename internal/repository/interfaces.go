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
	SetMFA(userID string, secret *string, enabled bool) error
}

type TokenRepositoryInterface interface {
	GenerateTokenPair(userID, username, role string) (*models.TokenPair, error)
	RevokeRefreshToken(refreshToken string) error
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

// IkasRepositoryInterface
type IkasRepositoryInterface interface {
	//CREATE DOMAIN
	CreateIdentifikasi(id string, data *dto.CreateIdentifikasiData) (float64, error)
	CreateProteksi(id string, data *dto.CreateProteksiData) (float64, error)
	CreateDeteksi(id string, data *dto.CreateDeteksiData) (float64, error)
	CreateGulih(id string, data *dto.CreateGulihData) (float64, error)

	//CREATE IKAS
	Create(
		req dto.CreateIkasRequest,
		id string,
		nilaiKematangan float64,
		idIden, idProt, idDet, idGul string,
	) error

	//READ
	GetAll() ([]dto.IkasResponse, error)
	GetByID(id string) (*dto.IkasResponse, error)

	//UPDATE IKAS
	Update(id string, req dto.UpdateIkasRequest) error

	//UPDATE DOMAIN
	UpdateIdentifikasi(id string, data *dto.UpdateIdentifikasiData) (float64, error)
	UpdateProteksi(id string, data *dto.UpdateProteksiData) (float64, error)
	UpdateDeteksi(id string, data *dto.UpdateDeteksiData) (float64, error)
	UpdateGulih(id string, data *dto.UpdateGulihData) (float64, error)

	//DELETE
	Delete(id string) error

	//IMPORT
	ParseExcelForImport(fileData []byte) (*dto.CreateIkasRequest, error)

	//HELPER
	FindPerusahaanByName(namaPerusahaan string) (string, error)
}

// CsirtRepositoryInterface
type CsirtRepositoryInterface interface {
	Create(req dto.CreateCsirtRequest, id string) error
	GetByID(id string) (*models.Csirt, error)
	GetAllWithPerusahaan() ([]dto.CsirtResponse, error)
	GetByIDWithPerusahaan(id string) (*dto.CsirtResponse, error)
	Update(id string, csirt models.Csirt) error
	Delete(id string) error
}

// SdmCsirtRepositoryInterface
type SdmCsirtRepositoryInterface interface {
	Create(req dto.CreateSdmCsirtRequest, id string) error
	GetAll() ([]dto.SdmCsirtResponse, error)
	GetByID(id string) (*dto.SdmCsirtResponse, error)
	Update(id string, req dto.SdmCsirtResponse) error
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