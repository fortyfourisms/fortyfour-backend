package services

import (
	"context"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/rabbitmq"
	"ikas/internal/repository"
	"time"

	"github.com/google/uuid"
)

type IkasService struct {
	repo     repository.IkasRepositoryInterface
	producer *rabbitmq.Producer
}

func NewIkasService(repo repository.IkasRepositoryInterface, producer *rabbitmq.Producer) *IkasService {
	return &IkasService{
		repo:     repo,
		producer: producer,
	}
}

func (s *IkasService) Create(req dto.CreateIkasRequest, id string) error {
	var idIdentifikasi, idProteksi, idDeteksi, idGulih string
	var nilaiIden, nilaiProt, nilaiDet, nilaiGul float64
	var err error

	// Create Identifikasi jika ada data
	if req.Identifikasi != nil {
		idIdentifikasi = uuid.New().String()
		nilaiIden, err = s.repo.CreateIdentifikasi(idIdentifikasi, req.Identifikasi)
		if err != nil {
			return err
		}
	} else if req.IDIdentifikasi != "" {
		idIdentifikasi = req.IDIdentifikasi
	}

	// Create Proteksi jika ada data
	if req.Proteksi != nil {
		idProteksi = uuid.New().String()
		nilaiProt, err = s.repo.CreateProteksi(idProteksi, req.Proteksi)
		if err != nil {
			return err
		}
	} else if req.IDProteksi != "" {
		idProteksi = req.IDProteksi
	}

	// Create Deteksi jika ada data
	if req.Deteksi != nil {
		idDeteksi = uuid.New().String()
		nilaiDet, err = s.repo.CreateDeteksi(idDeteksi, req.Deteksi)
		if err != nil {
			return err
		}
	} else if req.IDDeteksi != "" {
		idDeteksi = req.IDDeteksi
	}

	// Create Gulih jika ada data
	if req.Gulih != nil {
		idGulih = uuid.New().String()
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
	err = s.repo.Create(req, id, nilaiKematangan,
		idIdentifikasi, idProteksi, idDeteksi, idGulih)
	if err != nil {
		return err
	}

	go s.publishIkasCreatedEvent(id, req, nilaiKematangan)

	return nil
}

func (s *IkasService) publishIkasCreatedEvent(ikasID string, req dto.CreateIkasRequest, nilaiKematangan float64) {
	if s.producer == nil {
		return // Skip jika producer tidak ada
	}

	event := dto_event.IkasCreatedEvent{
		IkasID:          ikasID,
		IDPerusahaan:    req.IDPerusahaan,
		Tanggal:         req.Tanggal,
		Responden:       req.Responden,
		NilaiKematangan: nilaiKematangan,
		TargetNilai:     req.TargetNilai,
		CreatedAt:       time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasCreated(ctx, event); err != nil {
		// Log error tapi jangan fail create operation
		// logger.Error(err, "operation failed")
	}
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
		return nil, err
	}

	// Simpan old nilai untuk event
	oldNilaiKematangan := existing.NilaiKematangan

	var needRecalculateKematangan bool
	var nilaiIden, nilaiProt, nilaiDet, nilaiGul float64

	// Update Identifikasi jika ada data
	if req.Identifikasi != nil {
		if existing.Identifikasi == nil {
			return nil, errors.New("identifikasi record not found for this ikas")
		}
		nilai, err := s.repo.UpdateIdentifikasi(existing.Identifikasi.ID, req.Identifikasi)
		if err != nil {
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
			return nil, err
		}
		nilaiGul = nilai
		needRecalculateKematangan = true
	} else if existing.Gulih != nil {
		nilaiGul = existing.Gulih.NilaiGulih
	}

	// Recalculate nilai_kematangan jika ada perubahan pada nested data
	var newNilaiKematangan float64
	if needRecalculateKematangan {
		newNilaiKematangan = (nilaiIden + nilaiProt + nilaiDet + nilaiGul) / 4.0
		req.NilaiKematangan = &newNilaiKematangan
	}

	// Update data IKAS utama
	if err := s.repo.Update(id, req); err != nil {
		return nil, err
	}

	// Ambil data terbaru dengan JOIN
	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	go s.publishIkasUpdatedEvent(id, oldNilaiKematangan, newNilaiKematangan, req)

	return updated, nil
}

func (s *IkasService) publishIkasUpdatedEvent(ikasID string, oldNilai, newNilai float64, req dto.UpdateIkasRequest) {
	if s.producer == nil {
		return
	}

	// Collect updated fields
	updatedFields := []string{}
	if req.IDPerusahaan != nil {
		updatedFields = append(updatedFields, "id_perusahaan")
	}
	if req.Tanggal != nil {
		updatedFields = append(updatedFields, "tanggal")
	}
	if req.Responden != nil {
		updatedFields = append(updatedFields, "responden")
	}
	if req.NilaiKematangan != nil {
		updatedFields = append(updatedFields, "nilai_kematangan")
	}
	if req.Identifikasi != nil {
		updatedFields = append(updatedFields, "identifikasi")
	}
	if req.Proteksi != nil {
		updatedFields = append(updatedFields, "proteksi")
	}
	if req.Deteksi != nil {
		updatedFields = append(updatedFields, "deteksi")
	}
	if req.Gulih != nil {
		updatedFields = append(updatedFields, "gulih")
	}

	event := dto_event.IkasUpdatedEvent{
		IkasID:             ikasID,
		OldNilaiKematangan: oldNilai,
		NewNilaiKematangan: newNilai,
		UpdatedFields:      updatedFields,
		UpdatedAt:          time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasUpdated(ctx, event); err != nil {
		// Log error tapi jangan fail update operation
		// logger.Error(err, "operation failed")
	}
}

func (s *IkasService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	go s.publishIkasDeletedEvent(id)

	return nil
}

func (s *IkasService) publishIkasDeletedEvent(ikasID string) {
	if s.producer == nil {
		return
	}

	event := dto_event.IkasDeletedEvent{
		IkasID:    ikasID,
		DeletedAt: time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasDeleted(ctx, event); err != nil {
		// Log error tapi jangan fail delete operation
		// logger.Error(err, "operation failed")
	}
}

func (s *IkasService) ImportFromExcel(fileData []byte) (*dto.IkasResponse, error) {
	// Parse Excel - semua data sudah diambil dari Excel
	excelData, err := s.repo.ParseExcelForImport(fileData)
	if err != nil {
		return nil, err
	}

	// Generate ID baru
	newID := uuid.New().String()

	// Create menggunakan service Create yang sudah ada
	if err := s.Create(*excelData, newID); err != nil {
		return nil, err
	}

	// Ambil data yang baru dibuat
	resp, err := s.GetByID(newID)
	if err != nil {
		return nil, err
	}

	go s.publishIkasImportedEvent(newID, excelData, resp.NilaiKematangan)

	return resp, nil
}

func (s *IkasService) publishIkasImportedEvent(ikasID string, req *dto.CreateIkasRequest, nilaiKematangan float64) {
	if s.producer == nil {
		return
	}

	event := dto_event.IkasImportedEvent{
		IkasID:          ikasID,
		IDPerusahaan:    req.IDPerusahaan,
		NilaiKematangan: nilaiKematangan,
		ImportedAt:      time.Now(),
	}

	ctx := context.Background()
	if err := s.producer.PublishIkasImported(ctx, event); err != nil {
		// Log error tapi jangan fail import operation
		// logger.Error(err, "operation failed")
	}
}
