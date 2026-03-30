package repository

import (
	"ikas/internal/dto"
	"ikas/internal/models"
)

// IkasRepositoryInterface
type IkasRepositoryInterface interface {
	//CREATE IKAS
	Create(
		req dto.CreateIkasRequest,
		id string,
		nilaiKematangan float64,
	) error

	//READ
	GetAll() ([]dto.IkasResponse, error)
	GetByID(id string) (*dto.IkasResponse, error)

	//UPDATE IKAS
	Update(id string, req dto.UpdateIkasRequest) error

	//DELETE
	Delete(id string) error

	//IMPORT
	ParseExcelForImport(fileData []byte) (*dto.ParsedExcelData, error)

	//HELPER
	FindPerusahaanByName(namaPerusahaan string) (string, error)
	CheckExistsByPerusahaanID(idPerusahaan string) (bool, error)
	GetIDByPerusahaanID(idPerusahaan string) (string, error)
}

// RuangLingkupRepositoryInterface
type RuangLingkupRepositoryInterface interface {
	// CREATE
	Create(req dto.CreateRuangLingkupRequest) (int64, error)

	// READ
	GetAll() ([]dto.RuangLingkupResponse, error)
	GetByID(id int) (*dto.RuangLingkupResponse, error)

	// UPDATE
	Update(id int, req dto.UpdateRuangLingkupRequest) error

	// DELETE
	Delete(id int) error

	// HELPER
	CheckDuplicateName(nama string, excludeID int) (bool, error)
}

// DomainRepositoryInterface
type DomainRepositoryInterface interface {
	Create(req dto.CreateDomainRequest) (int64, error)
	GetAll() ([]dto.DomainResponse, error)
	GetByID(id int) (*dto.DomainResponse, error)
	Update(id int, req dto.UpdateDomainRequest) error
	Delete(id int) error
	CheckDuplicateName(nama string, excludeID int) (bool, error)
}

// KategoriRepositoryInterface
type KategoriRepositoryInterface interface {
	Create(req dto.CreateKategoriRequest) (int64, error)
	GetAll() ([]dto.KategoriResponse, error)
	GetByID(id int) (*dto.KategoriResponse, error)
	Update(id int, req dto.UpdateKategoriRequest) error
	Delete(id int) error
	CheckDuplicateName(domainID int, namaKategori string, excludeID int) (bool, error)
	CheckDomainExists(domainID int) (bool, error)
}

// SubKategoriRepositoryInterface
type SubKategoriRepositoryInterface interface {
	Create(req dto.CreateSubKategoriRequest) (int64, error)
	GetAll() ([]dto.SubKategoriResponse, error)
	GetByID(id int) (*dto.SubKategoriResponse, error)
	Update(id int, req dto.UpdateSubKategoriRequest) error
	Delete(id int) error
	CheckDuplicateName(kategoriID int, namaSubKategori string, excludeID int) (bool, error)
	CheckKategoriExists(kategoriID int) (bool, error)
}

// IdentifikasiRepositoryInterface
type IdentifikasiRepositoryInterface interface {
	GetAll() ([]models.Identifikasi, error)
	GetByID(id string) (*models.Identifikasi, error)
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
