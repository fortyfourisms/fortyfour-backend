package repository

import "ikas/internal/dto"

// IkasRepositoryInterface
type IkasRepositoryInterface interface {
	//CREATE DOMAIN
	CreateIdentifikasi(data *dto.CreateIdentifikasiData) (int64, float64, error)
	CreateProteksi(data *dto.CreateProteksiData) (int64, float64, error)
	CreateDeteksi(data *dto.CreateDeteksiData) (int64, float64, error)
	CreateGulih(data *dto.CreateGulihData) (int64, float64, error)

	//CREATE IKAS
	Create(
		req dto.CreateIkasRequest,
		id string,
		nilaiKematangan float64,
		idIden, idProt, idDet, idGul int,
	) error

	//READ
	GetAll() ([]dto.IkasResponse, error)
	GetByID(id string) (*dto.IkasResponse, error)

	//UPDATE IKAS
	Update(id string, req dto.UpdateIkasRequest) error

	//UPDATE DOMAIN
	UpdateIdentifikasi(id int, data *dto.UpdateIdentifikasiData) (float64, error)
	UpdateProteksi(id int, data *dto.UpdateProteksiData) (float64, error)
	UpdateDeteksi(id int, data *dto.UpdateDeteksiData) (float64, error)
	UpdateGulih(id int, data *dto.UpdateGulihData) (float64, error)

	//DELETE
	Delete(id string) error

	//IMPORT
	ParseExcelForImport(fileData []byte) (*dto.CreateIkasRequest, error)

	//HELPER
	FindPerusahaanByName(namaPerusahaan string) (string, error)
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
