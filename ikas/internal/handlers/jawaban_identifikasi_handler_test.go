package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"ikas/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// stringPtr for convenient testing
func jidStrPtr(s string) *string {
	return &s
}
func jidFloat64Ptr(f float64) *float64 { return &f }

type mockJawabanIdentifikasiProducer struct {
	mock.Mock
}

func (m *mockJawabanIdentifikasiProducer) PublishJawabanIdentifikasiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockJawabanIdentifikasiProducer) PublishJawabanIdentifikasiUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockJawabanIdentifikasiProducer) PublishJawabanIdentifikasiDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockJawabanIdentifikasiProducer) PublishIkasAuditLog(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type mockJawabanIdentifikasiRepository struct {
	mock.Mock
}

func (m *mockJawabanIdentifikasiRepository) Create(req dto.CreateJawabanIdentifikasiRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) GetAll() ([]dto.JawabanIdentifikasiResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.JawabanIdentifikasiResponse), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) GetByID(id int) (*dto.JawabanIdentifikasiResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JawabanIdentifikasiResponse), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) GetByIkasID(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).([]dto.JawabanIdentifikasiResponse), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanIdentifikasiResponse, error) {
	args := m.Called(pertanyaanID)
	return args.Get(0).([]dto.JawabanIdentifikasiResponse), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) Update(id int, req dto.UpdateJawabanIdentifikasiRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}
func (m *mockJawabanIdentifikasiRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockJawabanIdentifikasiRepository) CheckPertanyaanExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) CheckIkasExists(id string) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error) {
	args := m.Called(perusahaanID, pertanyaanID, excludeID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) RecalculateIdentifikasi(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}
func (m *mockJawabanIdentifikasiRepository) UpsertToBuffer(req dto.CreateJawabanIdentifikasiRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *mockJawabanIdentifikasiRepository) GetBufferCount(perusahaanID string) (int, error) {
	args := m.Called(perusahaanID)
	return args.Get(0).(int), args.Error(1)
}
func (m *mockJawabanIdentifikasiRepository) FlushBuffer(perusahaanID string) error {
	args := m.Called(perusahaanID)
	return args.Error(0)
}

// Ensure mock compatibility
var _ repository.JawabanIdentifikasiRepositoryInterface = (*mockJawabanIdentifikasiRepository)(nil)
var _ services.JawabanIdentifikasiProducerInterface = (*mockJawabanIdentifikasiProducer)(nil)

func setupJawabanIdentifikasiHandler(repo *mockJawabanIdentifikasiRepository, ikasRepo *mockIkasRepository, producer *mockJawabanIdentifikasiProducer) *JawabanIdentifikasiHandler {
	service := services.NewJawabanIdentifikasiService(repo, ikasRepo, producer)
	return NewJawabanIdentifikasiHandler(service)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanIdentifikasiResponse{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-identifikasi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, nil)

	repo.On("GetAll").Return([]dto.JawabanIdentifikasiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-identifikasi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return(&dto.JawabanIdentifikasiResponse{ID: 1}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanIdentifikasiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanIdentifikasiResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/jawaban-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanIdentifikasiProducer)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, producer)

	createReq := dto.CreateJawabanIdentifikasiRequest{
		PertanyaanIdentifikasiID: 1,
		IkasID:             "550e8400-e29b-41d4-a716-446655440000",
		JawabanIdentifikasi:      jidFloat64Ptr(3.0),
	}

	repo.On("CheckPertanyaanExists", 1).Return(true, nil)
	repo.On("CheckIkasExists", "550e8400-e29b-41d4-a716-446655440000").Return(true, nil)
	repo.On("CheckDuplicate", "550e8400-e29b-41d4-a716-446655440000", 1, 0).Return(false, nil)
	producer.On("PublishJawabanIdentifikasiCreated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-identifikasi", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupJawabanIdentifikasiHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-identifikasi", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, nil)

	createReq := dto.CreateJawabanIdentifikasiRequest{
		PertanyaanIdentifikasiID: 1,
		IkasID:             "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-identifikasi", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanIdentifikasiProducer)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, producer)

	updateReq := dto.UpdateJawabanIdentifikasiRequest{
		JawabanIdentifikasi: jidFloat64Ptr(4.0),
	}

	existing := &dto.JawabanIdentifikasiResponse{ID: 1, IkasID: "uuid1", JawabanIdentifikasi: jidFloat64Ptr(3.0)}
	repo.On("GetByID", 1).Return(existing, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishJawabanIdentifikasiUpdated", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-identifikasi/1", bytes.NewReader(body))
	req.Header.Set("X-User-ID", "user1")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	handler := setupJawabanIdentifikasiHandler(repo, nil, nil)

	updateReq := dto.UpdateJawabanIdentifikasiRequest{
		JawabanIdentifikasi: jidFloat64Ptr(4.0),
	}
	repo.On("GetByID", 1).Return((*dto.JawabanIdentifikasiResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-identifikasi/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupJawabanIdentifikasiHandler(nil, nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-identifikasi/1", bytes.NewReader([]byte("{invalid-json}")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	ikasRepo := new(mockIkasRepository)
	producer := new(mockJawabanIdentifikasiProducer)
	handler := setupJawabanIdentifikasiHandler(repo, ikasRepo, producer)

	repo.On("GetByID", 1).Return(&dto.JawabanIdentifikasiResponse{ID: 1, IkasID: "uuid1"}, nil)
	ikasRepo.On("GetIDByPerusahaanID", "uuid1").Return("ikas1", nil)
	producer.On("PublishJawabanIdentifikasiDeleted", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockJawabanIdentifikasiRepository)
	handler := setupJawabanIdentifikasiHandler(repo, nil, nil)

	repo.On("GetByID", 1).Return((*dto.JawabanIdentifikasiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJawabanIdentifikasiHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupJawabanIdentifikasiHandler(nil, nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/jawaban-identifikasi/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/jawaban-identifikasi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/jawaban-identifikasi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/jawaban-identifikasi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
