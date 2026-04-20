package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
)

type SEEditRequestService interface {
	CreateRequest(userID, idSE string, req dto.CreateSEEditRequestDTO) (*dto.SEEditRequestResponse, error)
	GetPending() ([]dto.SEEditRequestResponse, error)
	GetByUser(userID string) ([]dto.SEEditRequestResponse, error)
	Review(id string, req dto.ReviewSEEditRequestDTO) (*dto.SEEditRequestResponse, error)
}

type seEditRequestService struct {
	repo    repository.SEEditRequestRepositoryInterface
	seRepo  repository.SERepositoryInterface
	seSvc   SEService
}

func NewSEEditRequestService(
	repo repository.SEEditRequestRepositoryInterface,
	seRepo repository.SERepositoryInterface,
	seSvc SEService,
) SEEditRequestService {
	return &seEditRequestService{
		repo:   repo,
		seRepo: seRepo,
		seSvc:  seSvc,
	}
}

// CreateRequest — user submit request edit SE
func (s *seEditRequestService) CreateRequest(userID, idSE string, req dto.CreateSEEditRequestDTO) (*dto.SEEditRequestResponse, error) {
	// Pastikan SE ada
	se, err := s.seRepo.GetByID(idSE)
	if err != nil {
		return nil, errors.New("SE tidak ditemukan")
	}
	_ = se

	// Serialize data perubahan ke JSON
	dataJSON, err := json.Marshal(req.DataPerubahan)
	if err != nil {
		return nil, errors.New("data perubahan tidak valid")
	}

	editReq := &models.SEEditRequest{
		ID:            uuid.NewString(),
		IDSE:          idSE,
		IDUser:        userID,
		Status:        models.SEEditRequestPending,
		CatatanUser:   req.Catatan,
		DataPerubahan: string(dataJSON),
	}

	if err := s.repo.Create(editReq); err != nil {
		return nil, err
	}

	return mapSEEditRequestToResponse(editReq), nil
}

// GetPending — admin list semua request pending
func (s *seEditRequestService) GetPending() ([]dto.SEEditRequestResponse, error) {
	requests, err := s.repo.FindAllPending()
	if err != nil {
		return nil, err
	}

	result := make([]dto.SEEditRequestResponse, 0, len(requests))
	for _, r := range requests {
		r := r
		result = append(result, *mapSEEditRequestToResponse(&r))
	}
	return result, nil
}

// GetByUser — user list semua request miliknya
func (s *seEditRequestService) GetByUser(userID string) ([]dto.SEEditRequestResponse, error) {
	requests, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.SEEditRequestResponse, 0, len(requests))
	for _, r := range requests {
		r := r
		result = append(result, *mapSEEditRequestToResponse(&r))
	}
	return result, nil
}

// Review — admin approve/reject, jika approve otomatis update SE
func (s *seEditRequestService) Review(id string, req dto.ReviewSEEditRequestDTO) (*dto.SEEditRequestResponse, error) {
	editReq, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("request tidak ditemukan: %v", err)
	}

	if editReq.Status != models.SEEditRequestPending {
		return nil, errors.New("request sudah diproses sebelumnya")
	}

	status := models.SEEditRequestStatus(strings.ToLower(req.Status))
	if status != models.SEEditRequestApproved && status != models.SEEditRequestRejected {
		return nil, errors.New("status harus 'approved' atau 'rejected'")
	}

	// Jika approve, terapkan perubahan ke SE
	if status == models.SEEditRequestApproved {
		var updateReq dto.UpdateSERequest
		if err := json.Unmarshal([]byte(editReq.DataPerubahan), &updateReq); err != nil {
			return nil, errors.New("gagal membaca data perubahan")
		}

		if _, err := s.seSvc.Update(editReq.IDSE, updateReq); err != nil {
			return nil, errors.New("gagal menerapkan perubahan: " + err.Error())
		}
	}

	if err := s.repo.UpdateStatus(id, status, req.Catatan); err != nil {
		return nil, err
	}

	editReq.Status = status
	editReq.Catatan = req.Catatan

	return mapSEEditRequestToResponse(editReq), nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func mapSEEditRequestToResponse(r *models.SEEditRequest) *dto.SEEditRequestResponse {
	resp := &dto.SEEditRequestResponse{
		ID:          r.ID,
		IDSE:        r.IDSE,
		IDUser:      r.IDUser,
		Status:      r.Status,
		CatatanUser: r.CatatanUser,
		Catatan:     r.Catatan,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
	}

	// Parse data perubahan JSON ke DTO
	var updateReq dto.UpdateSERequest
	if err := json.Unmarshal([]byte(r.DataPerubahan), &updateReq); err == nil {
		resp.DataPerubahan = &updateReq
	}

	return resp
}
