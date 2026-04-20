package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/middleware"
	"ikas/internal/repository"
	"ikas/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper functions for jawaban proteksi tests
func jpStrPtr(s string) *string       { return &s }
func jpFloat64Ptr(f float64) *float64 { return &f }

// ─── Mock Producer ───────────────────────────────────────────────────────────

type mockJawabanProteksiProducer struct {
	mock.Mock
}

func (m *mockJawabanProteksiProducer) PublishJawabanProteksiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanProteksiProducer) PublishJawabanProteksiUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanProteksiProducer) PublishJawabanProteksiDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanProteksiProducer) PublishIkasAuditLog(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ─── Mock Repository ─────────────────────────────────────────────────────────

type mockJawabanProteksiRepository struct {
	mock.Mock
}

func (m *mockJawabanProteksiRepository) Create(req dto.CreateJawabanProteksiRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockJawabanProteksiRepository) GetAll() ([]dto.JawabanProteksiResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.JawabanProteksiResponse), args.Error(1)
}
func (m *mockJawabanProteksiRepository) GetByID(id int) (*dto.JawabanProteksiResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JawabanProteksiResponse), args.Error(1)
}
func (m *mockJawabanProteksiRepository) GetByIkasID(perusahaanID string) ([]dto.JawabanProteksiResponse, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).([]dto.JawabanProteksiResponse), args.Error(1)
}
func (m *mockJawabanProteksiRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanProteksiResponse, error) {
	args := m.Called(pertanyaanID)
	return args.Get(0).([]dto.JawabanProteksiResponse), args.Error(1)
}
func (m *mockJawabanProteksiRepository) Update(id int, req dto.UpdateJawabanProteksiRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}
func (m *mockJawabanProteksiRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockJawabanProteksiRepository) CheckPertanyaanExists(pertanyaanID int) (bool, error) {
	args := m.Called(pertanyaanID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanProteksiRepository) CheckIkasExists(perusahaanID string) (bool, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanProteksiRepository) CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error) {
	args := m.Called(perusahaanID, pertanyaanID, excludeID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanProteksiRepository) RecalculateProteksi(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}
func (m *mockJawabanProteksiRepository) UpsertToBuffer(req dto.CreateJawabanProteksiRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *mockJawabanProteksiRepository) GetBufferCount(perusahaanID string) (int, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).(int), args.Error(1)
}
func (m *mockJawabanProteksiRepository) FlushBuffer(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}
func (m *mockJawabanProteksiRepository) GetByPerusahaanID(perusahaanID string) ([]dto.JawabanProteksiResponse, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).([]dto.JawabanProteksiResponse), args.Error(1)
}

func (m *mockJawabanProteksiRepository) CloneByIkasID(oldIkasID string, newIkasID string) error {
	return nil
}

// Interface compliance
var _ repository.JawabanProteksiRepositoryInterface = (*mockJawabanProteksiRepository)(nil)
var _ services.JawabanProteksiProducerInterface = (*mockJawabanProteksiProducer)(nil)

func setupJawabanProteksiHandler(repo *mockJawabanProteksiRepository, ikasRepo *mockIkasRepository, producer *mockJawabanProteksiProducer) *JawabanProteksiHandler {
	if ikasRepo != nil {
		ikasRepo.On("IsLocked", mock.Anything).Return(false, nil).Maybe()
		ikasRepo.On("GetByID", mock.Anything).Return(&dto.IkasResponse{
			ID: "uuid1",
			Perusahaan: &dto.PerusahaanInIkas{ID: "1"},
		}, nil).Maybe()
		ikasRepo.On("CheckOwnership", mock.Anything, mock.Anything).Return(true, nil).Maybe()
	}
	service := services.NewJawabanProteksiService(repo, ikasRepo, producer)
	return NewJawabanProteksiHandler(service)
}

// ─── GET ALL ─────────────────────────────────────────────────────────────────

func TestJawabanProteksiHandler_GetAll_Success(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanProteksiResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanProteksiHandler_GetAll_Error(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanProteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET BY PERUSAHAAN ───────────────────────────────────────────────────────

func TestJawabanProteksiHandler_GetByIkasID_Success(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByIkasID", "550e8400-e29b-41d4-a716-446655440000").Return([]dto.JawabanProteksiResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi?ikas_id=550e8400-e29b-41d4-a716-446655440000", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanProteksiHandler_GetByIkasID_InvalidUUID(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi?ikas_id=invalid-uuid", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_GetByIkasID_Error(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByIkasID", "550e8400-e29b-41d4-a716-446655440000").Return([]dto.JawabanProteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi?ikas_id=550e8400-e29b-41d4-a716-446655440000", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET BY PERTANYAAN ───────────────────────────────────────────────────────

func TestJawabanProteksiHandler_GetByPertanyaan_Success(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByPertanyaan", 1).Return([]dto.JawabanProteksiResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi?pertanyaan_proteksi_id=1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanProteksiHandler_GetByPertanyaan_InvalidFormat(t *testing.T) {
	handler := setupJawabanProteksiHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi?pertanyaan_proteksi_id=abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_GetByPertanyaan_Error(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByPertanyaan", 1).Return([]dto.JawabanProteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi?pertanyaan_proteksi_id=1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanProteksiHandler_GetByPertanyaan_InvalidID_Zero(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	// pertanyaan_proteksi_id=0 → service returns "pertanyaan_proteksi_id tidak valid" → handler 500 (default)
	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi?pertanyaan_proteksi_id=0", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET BY ID ───────────────────────────────────────────────────────────────

func TestJawabanProteksiHandler_GetByID_Success(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return(&dto.JawabanProteksiResponse{ID: 1, IkasID: "ikas1"}, nil)
	ikasRepo.On("GetByID", "ikas1").Return(&dto.IkasResponse{ID: "ikas1", Perusahaan: &dto.PerusahaanInIkas{ID: "p1"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi/1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanProteksiHandler_GetByID_NotFound(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanProteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi/1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanProteksiHandler_GetByID_Error(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanProteksiResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-proteksi/1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── CREATE ──────────────────────────────────────────────────────────────────

func TestJawabanProteksiHandler_Create_Success(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanProteksiProducer)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, producer)

	createReq := dto.CreateJawabanProteksiRequest{
		PertanyaanProteksiID: 1,
		IkasID:               "550e8400-e29b-41d4-a716-446655440000",
		JawabanProteksi:      jpFloat64Ptr(3.0),
	}

	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(false, nil)
	producer.On("PublishJawabanProteksiCreated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestJawabanProteksiHandler_Create_InvalidJSON(t *testing.T) {
	handler := setupJawabanProteksiHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_Create_ValidationError(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	// jawaban_proteksi is nil → triggers "jawaban_proteksi tidak boleh kosong" → 400
	createReq := dto.CreateJawabanProteksiRequest{
		PertanyaanProteksiID: 1,
		IkasID:               "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_Create_PertanyaanNotFound(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanProteksiRequest{
		PertanyaanProteksiID: 1,
		IkasID:               "550e8400-e29b-41d4-a716-446655440000",
		JawabanProteksi:      jpFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanProteksiHandler_Create_PerusahaanNotFound(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanProteksiRequest{
		PertanyaanProteksiID: 1,
		IkasID:               "550e8400-e29b-41d4-a716-446655440000",
		JawabanProteksi:      jpFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanProteksiHandler_Create_Duplicate(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanProteksiRequest{
		PertanyaanProteksiID: 1,
		IkasID:               "550e8400-e29b-41d4-a716-446655440000",
		JawabanProteksi:      jpFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestJawabanProteksiHandler_Create_ServerError(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanProteksiProducer)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, producer)

	createReq := dto.CreateJawabanProteksiRequest{
		PertanyaanProteksiID: 1,
		IkasID:               "550e8400-e29b-41d4-a716-446655440000",
		JawabanProteksi:      jpFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(false, nil)
	producer.On("PublishJawabanProteksiCreated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanProteksiHandler_Create_InvalidPertanyaanType(t *testing.T) {
	handler := setupJawabanProteksiHandler(nil, nil, nil)
	// pertanyaan_proteksi_id as string triggers JSON decode error containing "pertanyaan_proteksi_id"
	body := []byte(`{"pertanyaan_proteksi_id":"not-int","perusahaan_id":"uuid","jawaban_proteksi":3}`)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_Create_InvalidJawabanType(t *testing.T) {
	handler := setupJawabanProteksiHandler(nil, nil, nil)
	// jawaban_proteksi as string triggers JSON decode error containing "jawaban_proteksi"
	body := []byte(`{"pertanyaan_proteksi_id":1,"perusahaan_id":"uuid","jawaban_proteksi":"not-number"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── UPDATE ──────────────────────────────────────────────────────────────────

func TestJawabanProteksiHandler_Update_Success(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanProteksiProducer)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, producer)

	updateReq := dto.UpdateJawabanProteksiRequest{
		JawabanProteksi: jpFloat64Ptr(4.0),
	}

	existing := &dto.JawabanProteksiResponse{ID: 1, IkasID: "uuid1", JawabanProteksi: jpFloat64Ptr(3.0)}
	repo.On("GetByID", 1).Return(existing, nil)
	ikasRepo.On("GetByID", "uuid1").Return(&dto.IkasResponse{ID: "uuid1", Perusahaan: &dto.PerusahaanInIkas{ID: "ikas1"}}, nil)
	producer.On("PublishJawabanProteksiUpdated", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-proteksi/1", bytes.NewReader(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(context.WithValue(ctx, middleware.UserIDKey, "user1"))

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanProteksiHandler_Update_NotFound(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	handler := setupJawabanProteksiHandler(repo, nil, nil)

	updateReq := dto.UpdateJawabanProteksiRequest{
		JawabanProteksi: jpFloat64Ptr(4.0),
	}
	repo.On("GetByID", 1).Return((*dto.JawabanProteksiResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-proteksi/1", bytes.NewReader(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanProteksiHandler_Update_InvalidJSON(t *testing.T) {
	handler := setupJawabanProteksiHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-proteksi/1", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_Update_ValidationError(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanProteksiProducer)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, producer)

	existing := &dto.JawabanProteksiResponse{ID: 1, IkasID: "uuid1"}
	repo.On("GetByID", 1).Return(existing, nil)

	// Validasi only without evidence
	updateReq := dto.UpdateJawabanProteksiRequest{
		Validasi: jpStrPtr("yes"),
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-proteksi/1", bytes.NewReader(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_Update_CannotUnmarshal(t *testing.T) {
	handler := setupJawabanProteksiHandler(nil, nil, nil)
	// jawaban_proteksi as string triggers "cannot unmarshal" error
	body := []byte(`{"jawaban_proteksi":"not-a-number"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-proteksi/1", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanProteksiHandler_Update_ServerError(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanProteksiProducer)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, producer)

	existing := &dto.JawabanProteksiResponse{ID: 1, IkasID: "uuid1", JawabanProteksi: jpFloat64Ptr(3.0)}
	repo.On("GetByID", 1).Return(existing, nil)
	ikasRepo.On("GetByID", "uuid1").Return(&dto.IkasResponse{ID: "uuid1", Perusahaan: &dto.PerusahaanInIkas{ID: "ikas1"}}, nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishJawabanProteksiUpdated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	updateReq := dto.UpdateJawabanProteksiRequest{
		JawabanProteksi: jpFloat64Ptr(4.0),
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-proteksi/1", bytes.NewReader(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── DELETE ──────────────────────────────────────────────────────────────────

func TestJawabanProteksiHandler_Delete_Success(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanProteksiProducer)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, producer)

	repo.On("GetByID", 1).Return(&dto.JawabanProteksiResponse{ID: 1, IkasID: "uuid1"}, nil)
	ikasRepo.On("GetByID", "uuid1").Return(&dto.IkasResponse{ID: "uuid1", Perusahaan: &dto.PerusahaanInIkas{ID: "ikas1"}}, nil)
	producer.On("PublishJawabanProteksiDeleted", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-proteksi/1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(context.WithValue(ctx, middleware.UserIDKey, "user1"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanProteksiHandler_Delete_NotFound(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	handler := setupJawabanProteksiHandler(repo, nil, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanProteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-proteksi/1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanProteksiHandler_Delete_ServerError(t *testing.T) {
	repo := new(mockJawabanProteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanProteksiProducer)
	handler := setupJawabanProteksiHandler(repo, ikasRepo, producer)

	repo.On("GetByID", 1).Return(&dto.JawabanProteksiResponse{ID: 1, IkasID: "uuid1"}, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishJawabanProteksiDeleted", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-proteksi/1", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── METHOD VALIDATION ──────────────────────────────────────────────────────

func TestJawabanProteksiHandler_MethodValidation(t *testing.T) {
	handler := setupJawabanProteksiHandler(nil, nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-proteksi/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-proteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-proteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/jawaban-proteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
