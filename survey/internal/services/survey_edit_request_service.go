package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
	"survey/internal/repository"

	"github.com/google/uuid"
)

type EditRequestService interface {
	Create(userID string, req dto.CreateEditRequestDTO) (*dto.EditRequestResponse, error)
	GetPending() ([]dto.EditRequestResponse, error)
	GetByUser(userID string) ([]dto.EditRequestResponse, error)
	Review(id string, req dto.ReviewEditRequestDTO) (*dto.EditRequestResponse, error)
}

type editRequestService struct {
	repo       repository.EditRequestRepositoryInterface
	risikoRepo repository.RisikoRepositoryInterface
}

func NewEditRequestService(
	repo repository.EditRequestRepositoryInterface,
	risikoRepo repository.RisikoRepositoryInterface,
) EditRequestService {
	return &editRequestService{
		repo:       repo,
		risikoRepo: risikoRepo,
	}
}

// CREATE
func (s *editRequestService) Create(userID string, req dto.CreateEditRequestDTO) (*dto.EditRequestResponse, error) {

	// VALIDASI INPUT
	if req.RespondenID <= 0 {
		return nil, errors.New("responden_id wajib diisi")
	}

	if req.RisikoID <= 0 {
		return nil, errors.New("risiko_id wajib diisi")
	}

	// VALIDASI RISIKO ADA
	if _, err := s.risikoRepo.GetByID(req.RisikoID); err != nil {
		return nil, errors.New("risiko tidak ditemukan")
	}

	// CEK DUPLICATE REQUEST (pending)
	existing, err := s.repo.FindPendingByRespondenRisiko(req.RespondenID, req.RisikoID)
	if err != nil {
		return nil, err
	}

	if len(existing) > 0 {
		return nil, errors.New("masih ada request pending untuk risiko ini")
	}

	// SERIALIZE DATA
	dataJSON, err := json.Marshal(req.DataPerubahan)
	if err != nil {
		return nil, errors.New("data_perubahan tidak valid")
	}

	now := time.Now()

	model := &models.EditRequest{
		ID:            uuid.NewString(),
		RespondenID:   req.RespondenID,
		RisikoID:      req.RisikoID,
		UserID:        userID,
		Status:        models.EditRequestPending,
		CatatanUser:   req.CatatanUser,
		DataPerubahan: string(dataJSON),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	return mapToResponse(model), nil
}

// GET PENDING (ADMIN)
func (s *editRequestService) GetPending() ([]dto.EditRequestResponse, error) {

	data, err := s.repo.FindAllPending()
	if err != nil {
		return nil, err
	}

	result := make([]dto.EditRequestResponse, 0, len(data))
	for i := range data {
		result = append(result, *mapToResponse(&data[i]))
	}

	return result, nil
}

// GET BY USER
func (s *editRequestService) GetByUser(userID string) ([]dto.EditRequestResponse, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user_id tidak valid")
	}

	data, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.EditRequestResponse, 0, len(data))
	for i := range data {
		result = append(result, *mapToResponse(&data[i]))
	}

	return result, nil
}

// REVIEW (ADMIN)
func (s *editRequestService) Review(id string, req dto.ReviewEditRequestDTO) (*dto.EditRequestResponse, error) {

	// AMBIL DATA
	editReq, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("request tidak ditemukan")
	}

	// VALIDASI STATUS
	if editReq.Status != models.EditRequestPending {
		return nil, errors.New("request sudah diproses sebelumnya")
	}

	status := models.EditRequestStatus(strings.ToLower(req.Status))

	if status != models.EditRequestApproved && status != models.EditRequestRejected {
		return nil, errors.New("status harus approved atau rejected")
	}

	// APPLY UPDATE JIKA APPROVED
	if status == models.EditRequestApproved {

		var updateMap map[string]interface{}

		if err := json.Unmarshal([]byte(editReq.DataPerubahan), &updateMap); err != nil {
			return nil, errors.New("gagal membaca data_perubahan")
		}

		// UPDATE RISIKO (dynamic)
		if err := s.risikoRepo.UpdatePartial(editReq.RisikoID, updateMap); err != nil {
			return nil, fmt.Errorf("gagal update risiko: %v", err)
		}
	}

	// UPDATE STATUS REQUEST
	if err := s.repo.UpdateStatus(id, status, req.Catatan); err != nil {
		return nil, err
	}

	editReq.Status = status
	editReq.Catatan = req.Catatan
	editReq.UpdatedAt = time.Now()

	return mapToResponse(editReq), nil
}

// MAPPING
func mapToResponse(m *models.EditRequest) *dto.EditRequestResponse {

	resp := &dto.EditRequestResponse{
		ID:          m.ID,
		RespondenID: m.RespondenID,
		RisikoID:    m.RisikoID,
		UserID:      m.UserID,
		Status:      m.Status,
		CatatanUser: m.CatatanUser,
		Catatan:     m.Catatan,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}

	// PARSE JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(m.DataPerubahan), &data); err == nil {
		resp.DataPerubahan = data
	}

	return resp
}
