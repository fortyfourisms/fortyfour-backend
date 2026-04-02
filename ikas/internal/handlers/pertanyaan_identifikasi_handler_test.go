package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"ikas/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockPertanyaanIdentifikasiRepository implements repository.PertanyaanIdentifikasiRepositoryInterface
type mockPertanyaanIdentifikasiRepository struct {
	mock.Mock
}

func (m *mockPertanyaanIdentifikasiRepository) Create(req dto.CreatePertanyaanIdentifikasiRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPertanyaanIdentifikasiRepository) GetAll() ([]dto.PertanyaanIdentifikasiResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.PertanyaanIdentifikasiResponse), args.Error(1)
}

func (m *mockPertanyaanIdentifikasiRepository) GetByID(id int) (*dto.PertanyaanIdentifikasiResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PertanyaanIdentifikasiResponse), args.Error(1)
}

func (m *mockPertanyaanIdentifikasiRepository) GetTotalCount() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *mockPertanyaanIdentifikasiRepository) Update(id int, req dto.UpdatePertanyaanIdentifikasiRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockPertanyaanIdentifikasiRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockPertanyaanIdentifikasiRepository) CheckSubKategoriExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockPertanyaanIdentifikasiRepository) CheckRuangLingkupExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

// mockPertanyaanIdentifikasiProducer implements services.PertanyaanIdentifikasiProducerInterface
type mockPertanyaanIdentifikasiProducer struct {
	mock.Mock
}

func (m *mockPertanyaanIdentifikasiProducer) PublishPertanyaanIdentifikasiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanIdentifikasiProducer) PublishPertanyaanIdentifikasiUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanIdentifikasiProducer) PublishPertanyaanIdentifikasiDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.PertanyaanIdentifikasiRepositoryInterface = (*mockPertanyaanIdentifikasiRepository)(nil)
var _ services.PertanyaanIdentifikasiProducerInterface = (*mockPertanyaanIdentifikasiProducer)(nil)

func setupPertanyaanIdentifikasiHandler(repo *mockPertanyaanIdentifikasiRepository, producer *mockPertanyaanIdentifikasiProducer) *PertanyaanIdentifikasiHandler {
	service := services.NewPertanyaanIdentifikasiService(repo, producer)
	return NewPertanyaanIdentifikasiHandler(service)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	producer := new(mockPertanyaanIdentifikasiProducer)
	handler := setupPertanyaanIdentifikasiHandler(repo, producer)

	expectedData := []dto.PertanyaanIdentifikasiResponse{{ID: 1}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-identifikasi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	producer := new(mockPertanyaanIdentifikasiProducer)
	handler := setupPertanyaanIdentifikasiHandler(repo, producer)

	repo.On("GetAll").Return([]dto.PertanyaanIdentifikasiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-identifikasi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	producer := new(mockPertanyaanIdentifikasiProducer)
	handler := setupPertanyaanIdentifikasiHandler(repo, producer)

	expectedData := &dto.PertanyaanIdentifikasiResponse{ID: 1}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanIdentifikasiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanIdentifikasiResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	producer := new(mockPertanyaanIdentifikasiProducer)
	handler := setupPertanyaanIdentifikasiHandler(repo, producer)

	createReq := dto.CreatePertanyaanIdentifikasiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanIdentifikasi: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)

	producer.On("PublishPertanyaanIdentifikasiCreated", mock.Anything, mock.MatchedBy(func(e dto_event.PertanyaanIdentifikasiCreatedEvent) bool {
		return e.Request.PertanyaanIdentifikasi == "Pertanyaan Test"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-identifikasi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanIdentifikasiHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-identifikasi", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Create_SubKategoriNotFound(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	createReq := dto.CreatePertanyaanIdentifikasiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanIdentifikasi: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-identifikasi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	createReq := dto.CreatePertanyaanIdentifikasiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanIdentifikasi: ""}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-identifikasi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Create_SystemError(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	createReq := dto.CreatePertanyaanIdentifikasiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanIdentifikasi: "Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, errors.New("db error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-identifikasi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	producer := new(mockPertanyaanIdentifikasiProducer)
	handler := setupPertanyaanIdentifikasiHandler(repo, producer)

	updateReq := dto.UpdatePertanyaanIdentifikasiRequest{PertanyaanIdentifikasi: identStrPtr("Pertanyaan Update")}

	repo.On("GetByID", 1).Return(&dto.PertanyaanIdentifikasiResponse{ID: 1}, nil)

	producer.On("PublishPertanyaanIdentifikasiUpdated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-identifikasi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanIdentifikasiHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-identifikasi/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanIdentifikasiResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdatePertanyaanIdentifikasiRequest{PertanyaanIdentifikasi: identStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-identifikasi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.PertanyaanIdentifikasiResponse{ID: 1}, nil)

	body, _ := json.Marshal(dto.UpdatePertanyaanIdentifikasiRequest{PertanyaanIdentifikasi: identStrPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-identifikasi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Update_SystemError(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanIdentifikasiResponse)(nil), errors.New("db error"))

	body, _ := json.Marshal(dto.UpdatePertanyaanIdentifikasiRequest{PertanyaanIdentifikasi: identStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-identifikasi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	producer := new(mockPertanyaanIdentifikasiProducer)
	handler := setupPertanyaanIdentifikasiHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.PertanyaanIdentifikasiResponse{ID: 1}, nil)
	producer.On("PublishPertanyaanIdentifikasiDeleted", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanIdentifikasiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_Delete_SystemError(t *testing.T) {
	repo := new(mockPertanyaanIdentifikasiRepository)
	handler := setupPertanyaanIdentifikasiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanIdentifikasiResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-identifikasi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanIdentifikasiHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupPertanyaanIdentifikasiHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-identifikasi/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-identifikasi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-identifikasi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/pertanyaan-identifikasi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func identStrPtr(s string) *string {
	return &s
}
