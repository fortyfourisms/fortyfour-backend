package repository

import "ikas/internal/dto"

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

// RuangLingkupRepositoryInterface
type RuangLingkupRepositoryInterface interface {
	// CREATE
	Create(req dto.CreateRuangLingkupRequest, id string) error

	// READ
	GetAll() ([]dto.RuangLingkupResponse, error)
	GetByID(id string) (*dto.RuangLingkupResponse, error)

	// UPDATE
	Update(id string, req dto.UpdateRuangLingkupRequest) error

	// DELETE
	Delete(id string) error

	// HELPER
	CheckDuplicateName(nama string, excludeID string) (bool, error)
}

// DomainRepositoryInterface
type DomainRepositoryInterface interface {
	Create(req dto.CreateDomainRequest, id string) error
	GetAll() ([]dto.DomainResponse, error)
	GetByID(id string) (*dto.DomainResponse, error)
	Update(id string, req dto.UpdateDomainRequest) error
	Delete(id string) error
	CheckDuplicateName(nama string, excludeID string) (bool, error)
}
