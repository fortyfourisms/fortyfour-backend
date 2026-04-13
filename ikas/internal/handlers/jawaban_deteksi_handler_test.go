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

// Helper functions for jawaban deteksi tests
func jdStrPtr(s string) *string       { return &s }
func jdFloat64Ptr(f float64) *float64 { return &f }

// ─── Mock Producer ───────────────────────────────────────────────────────────

type mockJawabanDeteksiProducer struct {
	mock.Mock
}

func (m *mockJawabanDeteksiProducer) PublishJawabanDeteksiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanDeteksiProducer) PublishJawabanDeteksiUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanDeteksiProducer) PublishJawabanDeteksiDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanDeteksiProducer) PublishIkasAuditLog(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ─── Mock Repository ─────────────────────────────────────────────────────────

type mockJawabanDeteksiRepository struct {
	mock.Mock
}

func (m *mockJawabanDeteksiRepository) Create(req dto.CreateJawabanDeteksiRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) GetAll() ([]dto.JawabanDeteksiResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.JawabanDeteksiResponse), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) GetByID(id int) (*dto.JawabanDeteksiResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JawabanDeteksiResponse), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) GetByIkasID(perusahaanID string) ([]dto.JawabanDeteksiResponse, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).([]dto.JawabanDeteksiResponse), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanDeteksiResponse, error) {
	args := m.Called(pertanyaanID)
	return args.Get(0).([]dto.JawabanDeteksiResponse), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) Update(id int, req dto.UpdateJawabanDeteksiRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}
func (m *mockJawabanDeteksiRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockJawabanDeteksiRepository) CheckPertanyaanExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) CheckIkasExists(id string) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error) {
	args := m.Called(perusahaanID, pertanyaanID, excludeID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) RecalculateDeteksi(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}
func (m *mockJawabanDeteksiRepository) UpsertToBuffer(req dto.CreateJawabanDeteksiRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *mockJawabanDeteksiRepository) GetBufferCount(perusahaanID string) (int, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).(int), args.Error(1)
}
func (m *mockJawabanDeteksiRepository) FlushBuffer(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}

// Interface compliance
var _ repository.JawabanDeteksiRepositoryInterface = (*mockJawabanDeteksiRepository)(nil)
var _ services.JawabanDeteksiProducerInterface = (*mockJawabanDeteksiProducer)(nil)

func setupJawabanDeteksiHandler(repo *mockJawabanDeteksiRepository, ikasRepo *mockIkasRepository, producer *mockJawabanDeteksiProducer) *JawabanDeteksiHandler {
	service := services.NewJawabanDeteksiService(repo, ikasRepo, producer)
	return NewJawabanDeteksiHandler(service)
}

// ─── GET ALL (no query params) ───────────────────────────────────────────────

func TestJawabanDeteksiHandler_GetAll_Success(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanDeteksiResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanDeteksiHandler_GetAll_Error(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanDeteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET ALL filtered by perusahaan_id ───────────────────────────────────────

func TestJawabanDeteksiHandler_GetByIkasID_Success(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByIkasID", "550e8400-e29b-41d4-a716-446655440000").Return([]dto.JawabanDeteksiResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi?perusahaan_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanDeteksiHandler_GetByIkasID_InvalidUUID(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	// The service validates UUIDs; invalid UUID triggers error
	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi?perusahaan_id=invalid-uuid", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanDeteksiHandler_GetByIkasID_Error(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByIkasID", "550e8400-e29b-41d4-a716-446655440000").Return([]dto.JawabanDeteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi?perusahaan_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET ALL filtered by pertanyaan_deteksi_id ───────────────────────────────

func TestJawabanDeteksiHandler_GetByPertanyaan_Success(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByPertanyaan", 1).Return([]dto.JawabanDeteksiResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi?pertanyaan_deteksi_id=1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanDeteksiHandler_GetByPertanyaan_Error(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByPertanyaan", 1).Return([]dto.JawabanDeteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi?pertanyaan_deteksi_id=1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET BY ID ───────────────────────────────────────────────────────────────

func TestJawabanDeteksiHandler_GetByID_Success(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return(&dto.JawabanDeteksiResponse{ID: 1}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanDeteksiHandler_GetByID_NotFound(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanDeteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanDeteksiHandler_GetByID_Error(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanDeteksiResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanDeteksiHandler_GetByID_InvalidID(t *testing.T) {
	handler := setupJawabanDeteksiHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-deteksi/abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── CREATE ──────────────────────────────────────────────────────────────────

func TestJawabanDeteksiHandler_Create_Success(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanDeteksiProducer)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, producer)

	createReq := dto.CreateJawabanDeteksiRequest{
		PertanyaanDeteksiID: 1,
		IkasID:        "550e8400-e29b-41d4-a716-446655440000",
		JawabanDeteksi:      jdFloat64Ptr(3.0),
	}

	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(false, nil)
	producer.On("PublishJawabanDeteksiCreated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestJawabanDeteksiHandler_Create_InvalidJSON(t *testing.T) {
	handler := setupJawabanDeteksiHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanDeteksiHandler_Create_ValidationError(t *testing.T) {
	handler := setupJawabanDeteksiHandler(nil, nil, nil)
	createReq := dto.CreateJawabanDeteksiRequest{}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanDeteksiHandler_Create_PertanyaanNotFound(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanDeteksiRequest{
		PertanyaanDeteksiID: 1,
		IkasID:        "550e8400-e29b-41d4-a716-446655440000",
		JawabanDeteksi:      jdFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanDeteksiHandler_Create_PerusahaanNotFound(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanDeteksiRequest{
		PertanyaanDeteksiID: 1,
		IkasID:        "550e8400-e29b-41d4-a716-446655440000",
		JawabanDeteksi:      jdFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanDeteksiHandler_Create_Duplicate(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanDeteksiRequest{
		PertanyaanDeteksiID: 1,
		IkasID:        "550e8400-e29b-41d4-a716-446655440000",
		JawabanDeteksi:      jdFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestJawabanDeteksiHandler_Create_ServerError(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanDeteksiProducer)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, producer)

	createReq := dto.CreateJawabanDeteksiRequest{
		PertanyaanDeteksiID: 1,
		IkasID:        "550e8400-e29b-41d4-a716-446655440000",
		JawabanDeteksi:      jdFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(false, nil)
	producer.On("PublishJawabanDeteksiCreated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── UPDATE ──────────────────────────────────────────────────────────────────

func TestJawabanDeteksiHandler_Update_Success(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanDeteksiProducer)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, producer)

	updateReq := dto.UpdateJawabanDeteksiRequest{
		JawabanDeteksi: jdFloat64Ptr(4.0),
	}

	existing := &dto.JawabanDeteksiResponse{ID: 1, IkasID: "uuid1", JawabanDeteksi: jdFloat64Ptr(3.0)}
	repo.On("GetByID", 1).Return(existing, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishJawabanDeteksiUpdated", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-deteksi/1", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user1"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanDeteksiHandler_Update_NotFound(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	handler := setupJawabanDeteksiHandler(repo, nil, nil)

	updateReq := dto.UpdateJawabanDeteksiRequest{
		JawabanDeteksi: jdFloat64Ptr(4.0),
	}
	repo.On("GetByID", 1).Return((*dto.JawabanDeteksiResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-deteksi/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanDeteksiHandler_Update_InvalidJSON(t *testing.T) {
	handler := setupJawabanDeteksiHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-deteksi/1", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanDeteksiHandler_Update_InvalidID(t *testing.T) {
	handler := setupJawabanDeteksiHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-deteksi/abc", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanDeteksiHandler_Update_ValidationError(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanDeteksiProducer)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, producer)

	existing := &dto.JawabanDeteksiResponse{ID: 1, IkasID: "uuid1"}
	repo.On("GetByID", 1).Return(existing, nil)

	// Validasi only without evidence
	updateReq := dto.UpdateJawabanDeteksiRequest{
		Validasi: jdStrPtr("yes"),
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-deteksi/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanDeteksiHandler_Update_ServerError(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanDeteksiProducer)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, producer)

	existing := &dto.JawabanDeteksiResponse{ID: 1, IkasID: "uuid1", JawabanDeteksi: jdFloat64Ptr(3.0)}
	repo.On("GetByID", 1).Return(existing, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishJawabanDeteksiUpdated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	updateReq := dto.UpdateJawabanDeteksiRequest{
		JawabanDeteksi: jdFloat64Ptr(4.0),
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-deteksi/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── DELETE ──────────────────────────────────────────────────────────────────

func TestJawabanDeteksiHandler_Delete_Success(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanDeteksiProducer)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, producer)

	repo.On("GetByID", 1).Return(&dto.JawabanDeteksiResponse{ID: 1, IkasID: "uuid1"}, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishJawabanDeteksiDeleted", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-deteksi/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user1"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanDeteksiHandler_Delete_NotFound(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	handler := setupJawabanDeteksiHandler(repo, nil, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanDeteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanDeteksiHandler_Delete_InvalidID(t *testing.T) {
	handler := setupJawabanDeteksiHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-deteksi/abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanDeteksiHandler_Delete_ServerError(t *testing.T) {
	repo := new(mockJawabanDeteksiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanDeteksiProducer)
	handler := setupJawabanDeteksiHandler(repo, ikasRepo, producer)

	repo.On("GetByID", 1).Return(&dto.JawabanDeteksiResponse{ID: 1, IkasID: "uuid1"}, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishJawabanDeteksiDeleted", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── ROUTING / ENDPOINT NOT FOUND ───────────────────────────────────────────

func TestJawabanDeteksiHandler_EndpointNotFound(t *testing.T) {
	handler := setupJawabanDeteksiHandler(nil, nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-deteksi/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-deteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-deteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/jawaban-deteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
