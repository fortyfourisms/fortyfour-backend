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
	repo     repository.SEEditRequestRepositoryInterface
	seRepo   repository.SERepositoryInterface
	seSvc    SEService
	userRepo repository.UserRepositoryInterface
	sseSvc   *SSEService
}

func NewSEEditRequestService(
	repo repository.SEEditRequestRepositoryInterface,
	seRepo repository.SERepositoryInterface,
	seSvc SEService,
	userRepo repository.UserRepositoryInterface,
	sseSvc *SSEService,
) SEEditRequestService {
	return &seEditRequestService{
		repo:     repo,
		seRepo:   seRepo,
		seSvc:    seSvc,
		userRepo: userRepo,
		sseSvc:   sseSvc,
	}
}

// CreateRequest — user submit request edit SE
func (s *seEditRequestService) CreateRequest(userID, idSE string, req dto.CreateSEEditRequestDTO) (*dto.SEEditRequestResponse, error) {
	// Pastikan SE ada
	se, err := s.seRepo.GetByID(idSE)
	if err != nil {
		return nil, errors.New("SE tidak ditemukan")
	}

	catatanUser := req.Catatan
	if catatanUser == nil {
		catatanUser = req.CatatanUser
	}

	// Serialize data perubahan ke JSON
	dataJSON, err := json.Marshal(req.DataPerubahan)
	if err != nil {
		return nil, errors.New("data perubahan tidak valid")
	}

	now := time.Now()
	editReq := &models.SEEditRequest{
		ID:            uuid.NewString(),
		IDSE:          idSE,
		IDUser:        userID,
		Status:        models.SEEditRequestPending,
		CatatanUser:   catatanUser,
		DataPerubahan: string(dataJSON),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(editReq); err != nil {
		return nil, err
	}

	s.notifyAdminsOnCreate(userID, se, editReq)

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

func (s *seEditRequestService) notifyAdminsOnCreate(requesterUserID string, se *dto.SEResponse, editReq *models.SEEditRequest) {
	if s.userRepo == nil || s.sseSvc == nil {
		return
	}

	admins, err := s.userRepo.FindAllAdmins()
	if err != nil || len(admins) == 0 {
		return
	}

	requesterName := requesterUserID
	if requester, err := s.userRepo.FindByID(requesterUserID); err == nil && requester != nil {
		switch {
		case requester.DisplayName != nil && *requester.DisplayName != "":
			requesterName = *requester.DisplayName
		case requester.Username != "":
			requesterName = requester.Username
		}
	}

	seName := editReq.IDSE
	if se != nil && se.NamaSE != "" {
		seName = se.NamaSE
	}

	message := fmt.Sprintf("Request edit data SE %s diajukan oleh %s", seName, requesterName)
	if editReq.CatatanUser != nil && *editReq.CatatanUser != "" {
		message += ". Catatan: " + *editReq.CatatanUser
	}

	adminIDs := make([]string, 0, len(admins))
	for _, admin := range admins {
		if admin.ID != "" {
			adminIDs = append(adminIDs, admin.ID)
		}
	}

	eventData := map[string]interface{}{
		"id":             editReq.ID,
		"id_se":          editReq.IDSE,
		"id_user":        editReq.IDUser,
		"status":         editReq.Status,
		"catatan_user":   editReq.CatatanUser,
		"data_perubahan": mapSEEditRequestToResponse(editReq).DataPerubahan,
		"nama_se":        seName,
		"requester_name": requesterName,
	}

	s.sseSvc.NotifyCreateToUsers("se_request", eventData, adminIDs, models.NotifSEEditRequested, message)
}
