package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type IkasService struct {
	repo *repository.IkasRepository
}

func NewIkasService(repo *repository.IkasRepository) *IkasService {
	return &IkasService{repo: repo}
}

// Update method Create di IkasService
func (s *IkasService) Create(req dto.CreateIkasRequest, id string) error {
	var idIdentifikasi, idProteksi, idDeteksi, idGulih string
	var nilaiIden, nilaiProt, nilaiDet, nilaiGul float64
	var err error

	// Create Identifikasi jika ada data
	if req.Identifikasi != nil {
		idIdentifikasi = uuid.New().String() // Generate UUID di sini
		nilaiIden, err = s.repo.CreateIdentifikasi(idIdentifikasi, req.Identifikasi)
		if err != nil {
			rollbar.Error(err)
			return err
		}
	} else if req.IDIdentifikasi != "" {
		// Backward compatibility: gunakan ID manual
		idIdentifikasi = req.IDIdentifikasi
	}

	// Create Proteksi jika ada data
	if req.Proteksi != nil {
		idProteksi = uuid.New().String() // Generate UUID di sini
		nilaiProt, err = s.repo.CreateProteksi(idProteksi, req.Proteksi)
		if err != nil {
			rollbar.Error(err)
			return err
		}
	} else if req.IDProteksi != "" {
		idProteksi = req.IDProteksi
	}

	// Create Deteksi jika ada data
	if req.Deteksi != nil {
		idDeteksi = uuid.New().String() // Generate UUID di sini
		nilaiDet, err = s.repo.CreateDeteksi(idDeteksi, req.Deteksi)
		if err != nil {
			rollbar.Error(err)
			return err
		}
	} else if req.IDDeteksi != "" {
		idDeteksi = req.IDDeteksi
	}

	// Create Gulih jika ada data
	if req.Gulih != nil {
		idGulih = uuid.New().String() // Generate UUID di sini
		nilaiGul, err = s.repo.CreateGulih(idGulih, req.Gulih)
		if err != nil {
			rollbar.Error(err)
			return err
		}
	} else if req.IDGulih != "" {
		idGulih = req.IDGulih
	}

	// Hitung nilai kematangan (rata-rata dari 4 nilai)
	nilaiKematangan := (nilaiIden + nilaiProt + nilaiDet + nilaiGul) / 4.0

	// Create IKAS record
	return s.repo.Create(req, id, nilaiKematangan,
		idIdentifikasi, idProteksi, idDeteksi, idGulih)
}

func (s *IkasService) GetAll() ([]dto.IkasResponse, error) {
	return s.repo.GetAll()
}

func (s *IkasService) GetByID(id string) (*dto.IkasResponse, error) {
	return s.repo.GetByID(id)
}

func (s *IkasService) Update(id string, req dto.UpdateIkasRequest) (*dto.IkasResponse, error) {
	// Ambil data existing untuk mendapatkan ID nested tables
	existing, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	var needRecalculateKematangan bool
	var nilaiIden, nilaiProt, nilaiDet, nilaiGul float64

	// Update Identifikasi jika ada data
	if req.Identifikasi != nil {
		if existing.Identifikasi == nil {
			return nil, errors.New("identifikasi record not found for this ikas")
		}
		nilai, err := s.repo.UpdateIdentifikasi(existing.Identifikasi.ID, req.Identifikasi)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		nilaiIden = nilai
		needRecalculateKematangan = true
	} else if existing.Identifikasi != nil {
		nilaiIden = existing.Identifikasi.NilaiIdentifikasi
	}

	// Update Proteksi jika ada data
	if req.Proteksi != nil {
		if existing.Proteksi == nil {
			return nil, errors.New("proteksi record not found for this ikas")
		}
		nilai, err := s.repo.UpdateProteksi(existing.Proteksi.ID, req.Proteksi)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		nilaiProt = nilai
		needRecalculateKematangan = true
	} else if existing.Proteksi != nil {
		nilaiProt = existing.Proteksi.NilaiProteksi
	}

	// Update Deteksi jika ada data
	if req.Deteksi != nil {
		if existing.Deteksi == nil {
			return nil, errors.New("deteksi record not found for this ikas")
		}
		nilai, err := s.repo.UpdateDeteksi(existing.Deteksi.ID, req.Deteksi)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		nilaiDet = nilai
		needRecalculateKematangan = true
	} else if existing.Deteksi != nil {
		nilaiDet = existing.Deteksi.NilaiDeteksi
	}

	// Update Gulih jika ada data
	if req.Gulih != nil {
		if existing.Gulih == nil {
			return nil, errors.New("gulih record not found for this ikas")
		}
		nilai, err := s.repo.UpdateGulih(existing.Gulih.ID, req.Gulih)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		nilaiGul = nilai
		needRecalculateKematangan = true
	} else if existing.Gulih != nil {
		nilaiGul = existing.Gulih.NilaiGulih
	}

	// Recalculate nilai_kematangan jika ada perubahan pada nested data
	if needRecalculateKematangan {
		nilaiKematangan := (nilaiIden + nilaiProt + nilaiDet + nilaiGul) / 4.0
		req.NilaiKematangan = &nilaiKematangan
	}

	// Update data IKAS utama
	if err := s.repo.Update(id, req); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data terbaru dengan JOIN
	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *IkasService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *IkasService) ImportFromExcel(fileData []byte) (*dto.IkasResponse, error) {
	// Parse Excel - semua data sudah diambil dari Excel
	excelData, err := s.repo.ParseExcelForImport(fileData)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Generate ID baru
	newID := uuid.New().String()

	// Create menggunakan service Create yang sudah ada
	if err := s.Create(*excelData, newID); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data yang baru dibuat
	resp, err := s.GetByID(newID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}
