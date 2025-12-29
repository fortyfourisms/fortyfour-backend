package services

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
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
	// Update data
	if err := s.repo.Update(id, req); err != nil {
		return nil, err
	}

	// Ambil data terbaru dengan JOIN
	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *IkasService) Delete(id string) error {
	return s.repo.Delete(id)
}
