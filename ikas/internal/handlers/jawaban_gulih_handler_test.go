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

// Helper functions for jawaban gulih tests
func jgStrPtr(s string) *string       { return &s }
func jgFloat64Ptr(f float64) *float64 { return &f }

// ─── Mock Producer ───────────────────────────────────────────────────────────

type mockJawabanGulihProducer struct {
	mock.Mock
}

func (m *mockJawabanGulihProducer) PublishJawabanGulihCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanGulihProducer) PublishJawabanGulihUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanGulihProducer) PublishJawabanGulihDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *mockJawabanGulihProducer) PublishIkasAuditLog(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ─── Mock Repository ─────────────────────────────────────────────────────────

type mockJawabanGulihRepository struct {
	mock.Mock
}

func (m *mockJawabanGulihRepository) Create(req dto.CreateJawabanGulihRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockJawabanGulihRepository) GetAll() ([]dto.JawabanGulihResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.JawabanGulihResponse), args.Error(1)
}
func (m *mockJawabanGulihRepository) GetByID(id int) (*dto.JawabanGulihResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JawabanGulihResponse), args.Error(1)
}
func (m *mockJawabanGulihRepository) GetByPerusahaan(perusahaanID string) ([]dto.JawabanGulihResponse, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).([]dto.JawabanGulihResponse), args.Error(1)
}
func (m *mockJawabanGulihRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanGulihResponse, error) {
	args := m.Called(pertanyaanID)
	return args.Get(0).([]dto.JawabanGulihResponse), args.Error(1)
}
func (m *mockJawabanGulihRepository) Update(id int, req dto.UpdateJawabanGulihRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}
func (m *mockJawabanGulihRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockJawabanGulihRepository) CheckPertanyaanExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanGulihRepository) CheckPerusahaanExists(id string) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanGulihRepository) CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error) {
	args := m.Called(perusahaanID, pertanyaanID, excludeID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanGulihRepository) RecalculateGulih(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}
func (m *mockJawabanGulihRepository) UpsertToBuffer(req dto.CreateJawabanGulihRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *mockJawabanGulihRepository) GetBufferCount(perusahaanID string) (int, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).(int), args.Error(1)
}
func (m *mockJawabanGulihRepository) FlushBuffer(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}

// Interface compliance
var _ repository.JawabanGulihRepositoryInterface = (*mockJawabanGulihRepository)(nil)
var _ services.JawabanGulihProducerInterface = (*mockJawabanGulihProducer)(nil)

func setupJawabanGulihHandler(repo *mockJawabanGulihRepository, ikasRepo *mockIkasRepository, producer *mockJawabanGulihProducer) *JawabanGulihHandler {
	service := services.NewJawabanGulihService(repo, ikasRepo, producer)
	return NewJawabanGulihHandler(service)
}

// ─── GET ALL (no query params) ───────────────────────────────────────────────

func TestJawabanGulihHandler_GetAll_Success(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanGulihResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanGulihHandler_GetAll_Error(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanGulihResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET ALL filtered by perusahaan_id ───────────────────────────────────────

func TestJawabanGulihHandler_GetByPerusahaan_Success(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetByPerusahaan", "550e8400-e29b-41d4-a716-446655440000").Return([]dto.JawabanGulihResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih?perusahaan_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanGulihHandler_GetByPerusahaan_InvalidUUID(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih?perusahaan_id=invalid-uuid", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanGulihHandler_GetByPerusahaan_Error(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetByPerusahaan", "550e8400-e29b-41d4-a716-446655440000").Return([]dto.JawabanGulihResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih?perusahaan_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET ALL filtered by pertanyaan_gulih_id ──────────────────────────────────

func TestJawabanGulihHandler_GetByPertanyaan_Success(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetByPertanyaan", 1).Return([]dto.JawabanGulihResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih?pertanyaan_gulih_id=1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanGulihHandler_GetByPertanyaan_Error(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetByPertanyaan", 1).Return([]dto.JawabanGulihResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih?pertanyaan_gulih_id=1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET BY ID ───────────────────────────────────────────────────────────────

func TestJawabanGulihHandler_GetByID_Success(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return(&dto.JawabanGulihResponse{ID: 1}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanGulihHandler_GetByID_NotFound(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanGulihResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanGulihHandler_GetByID_Error(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanGulihResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanGulihHandler_GetByID_InvalidID(t *testing.T) {
	handler := setupJawabanGulihHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-gulih/abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── CREATE ──────────────────────────────────────────────────────────────────

func TestJawabanGulihHandler_Create_Success(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanGulihProducer)
	handler := setupJawabanGulihHandler(repo, ikasRepo, producer)

	createReq := dto.CreateJawabanGulihRequest{
		PertanyaanGulihID: 1,
		PerusahaanID:      "550e8400-e29b-41d4-a716-446655440000",
		JawabanGulih:      jgFloat64Ptr(3.0),
	}

	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckPerusahaanExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(false, nil)
	producer.On("PublishJawabanGulihCreated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestJawabanGulihHandler_Create_InvalidJSON(t *testing.T) {
	handler := setupJawabanGulihHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanGulihHandler_Create_ValidationError(t *testing.T) {
	handler := setupJawabanGulihHandler(nil, nil, nil)
	createReq := dto.CreateJawabanGulihRequest{}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanGulihHandler_Create_PertanyaanNotFound(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanGulihRequest{
		PertanyaanGulihID: 1,
		PerusahaanID:      "550e8400-e29b-41d4-a716-446655440000",
		JawabanGulih:      jgFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanGulihHandler_Create_PerusahaanNotFound(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanGulihRequest{
		PertanyaanGulihID: 1,
		PerusahaanID:      "550e8400-e29b-41d4-a716-446655440000",
		JawabanGulih:      jgFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckPerusahaanExists", "550e8400-e29b-41d4-a716-446655440000").Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanGulihHandler_Create_Duplicate(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanGulihHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanGulihRequest{
		PertanyaanGulihID: 1,
		PerusahaanID:      "550e8400-e29b-41d4-a716-446655440000",
		JawabanGulih:      jgFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckPerusahaanExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestJawabanGulihHandler_Create_ServerError(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanGulihProducer)
	handler := setupJawabanGulihHandler(repo, ikasRepo, producer)

	createReq := dto.CreateJawabanGulihRequest{
		PertanyaanGulihID: 1,
		PerusahaanID:      "550e8400-e29b-41d4-a716-446655440000",
		JawabanGulih:      jgFloat64Ptr(3.0),
	}
	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckPerusahaanExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(false, nil)
	producer.On("PublishJawabanGulihCreated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── UPDATE ──────────────────────────────────────────────────────────────────

func TestJawabanGulihHandler_Update_Success(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanGulihProducer)
	handler := setupJawabanGulihHandler(repo, ikasRepo, producer)

	updateReq := dto.UpdateJawabanGulihRequest{
		JawabanGulih: jgFloat64Ptr(4.0),
	}

	existing := &dto.JawabanGulihResponse{ID: 1, PerusahaanID: "uuid1", JawabanGulih: jgFloat64Ptr(3.0)}
	repo.On("GetByID", 1).Return(existing, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishJawabanGulihUpdated", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-gulih/1", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user1"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanGulihHandler_Update_NotFound(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	handler := setupJawabanGulihHandler(repo, nil, nil)

	updateReq := dto.UpdateJawabanGulihRequest{
		JawabanGulih: jgFloat64Ptr(4.0),
	}
	repo.On("GetByID", 1).Return((*dto.JawabanGulihResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-gulih/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanGulihHandler_Update_InvalidJSON(t *testing.T) {
	handler := setupJawabanGulihHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-gulih/1", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanGulihHandler_Update_InvalidID(t *testing.T) {
	handler := setupJawabanGulihHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-gulih/abc", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanGulihHandler_Update_ValidationError(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanGulihProducer)
	handler := setupJawabanGulihHandler(repo, ikasRepo, producer)

	existing := &dto.JawabanGulihResponse{ID: 1, PerusahaanID: "uuid1"}
	repo.On("GetByID", 1).Return(existing, nil)

	// Validasi only without evidence
	updateReq := dto.UpdateJawabanGulihRequest{
		Validasi: jgStrPtr("yes"),
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-gulih/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanGulihHandler_Update_ServerError(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanGulihProducer)
	handler := setupJawabanGulihHandler(repo, ikasRepo, producer)

	existing := &dto.JawabanGulihResponse{ID: 1, PerusahaanID: "uuid1", JawabanGulih: jgFloat64Ptr(3.0)}
	repo.On("GetByID", 1).Return(existing, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishJawabanGulihUpdated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	updateReq := dto.UpdateJawabanGulihRequest{
		JawabanGulih: jgFloat64Ptr(4.0),
	}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-gulih/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── DELETE ──────────────────────────────────────────────────────────────────

func TestJawabanGulihHandler_Delete_Success(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanGulihProducer)
	handler := setupJawabanGulihHandler(repo, ikasRepo, producer)

	repo.On("GetByID", 1).Return(&dto.JawabanGulihResponse{ID: 1, PerusahaanID: "uuid1"}, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishJawabanGulihDeleted", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-gulih/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user1"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanGulihHandler_Delete_NotFound(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	handler := setupJawabanGulihHandler(repo, nil, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanGulihResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanGulihHandler_Delete_InvalidID(t *testing.T) {
	handler := setupJawabanGulihHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-gulih/abc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanGulihHandler_Delete_ServerError(t *testing.T) {
	repo := new(mockJawabanGulihRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanGulihProducer)
	handler := setupJawabanGulihHandler(repo, ikasRepo, producer)

	repo.On("GetByID", 1).Return(&dto.JawabanGulihResponse{ID: 1, PerusahaanID: "uuid1"}, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishJawabanGulihDeleted", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── ROUTING / ENDPOINT NOT FOUND ───────────────────────────────────────────

func TestJawabanGulihHandler_EndpointNotFound(t *testing.T) {
	handler := setupJawabanGulihHandler(nil, nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-gulih/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-gulih", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-gulih", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/jawaban-gulih", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
