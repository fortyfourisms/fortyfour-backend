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

// mockPertanyaanDeteksiRepository implements repository.PertanyaanDeteksiRepositoryInterface
type mockPertanyaanDeteksiRepository struct {
	mock.Mock
}

func (m *mockPertanyaanDeteksiRepository) Create(req dto.CreatePertanyaanDeteksiRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPertanyaanDeteksiRepository) GetAll() ([]dto.PertanyaanDeteksiResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.PertanyaanDeteksiResponse), args.Error(1)
}

func (m *mockPertanyaanDeteksiRepository) GetByID(id int) (*dto.PertanyaanDeteksiResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PertanyaanDeteksiResponse), args.Error(1)
}

func (m *mockPertanyaanDeteksiRepository) GetTotalCount() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *mockPertanyaanDeteksiRepository) Update(id int, req dto.UpdatePertanyaanDeteksiRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockPertanyaanDeteksiRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockPertanyaanDeteksiRepository) CheckSubKategoriExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockPertanyaanDeteksiRepository) CheckRuangLingkupExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

// mockPertanyaanDeteksiProducer implements services.PertanyaanDeteksiProducerInterface
type mockPertanyaanDeteksiProducer struct {
	mock.Mock
}

func (m *mockPertanyaanDeteksiProducer) PublishPertanyaanDeteksiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanDeteksiProducer) PublishPertanyaanDeteksiUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanDeteksiProducer) PublishPertanyaanDeteksiDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.PertanyaanDeteksiRepositoryInterface = (*mockPertanyaanDeteksiRepository)(nil)
var _ services.PertanyaanDeteksiProducerInterface = (*mockPertanyaanDeteksiProducer)(nil)

func setupPertanyaanDeteksiHandler(repo *mockPertanyaanDeteksiRepository, producer *mockPertanyaanDeteksiProducer) *PertanyaanDeteksiHandler {
	service := services.NewPertanyaanDeteksiService(repo, producer)
	return NewPertanyaanDeteksiHandler(service)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	producer := new(mockPertanyaanDeteksiProducer)
	handler := setupPertanyaanDeteksiHandler(repo, producer)

	expectedData := []dto.PertanyaanDeteksiResponse{{ID: 1}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-deteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	producer := new(mockPertanyaanDeteksiProducer)
	handler := setupPertanyaanDeteksiHandler(repo, producer)

	repo.On("GetAll").Return([]dto.PertanyaanDeteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-deteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	producer := new(mockPertanyaanDeteksiProducer)
	handler := setupPertanyaanDeteksiHandler(repo, producer)

	expectedData := &dto.PertanyaanDeteksiResponse{ID: 1}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanDeteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanDeteksiResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	producer := new(mockPertanyaanDeteksiProducer)
	handler := setupPertanyaanDeteksiHandler(repo, producer)

	createReq := dto.CreatePertanyaanDeteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanDeteksi: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)
	repo.On("Create", mock.Anything).Return(int64(1), nil)
	
	producer.On("PublishPertanyaanDeteksiCreated", mock.Anything, mock.MatchedBy(func(e dto_event.PertanyaanDeteksiCreatedEvent) bool {
		return e.Request.PertanyaanDeteksi == "Pertanyaan Test"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-deteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanDeteksiHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-deteksi", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Create_SubKategoriNotFound(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	createReq := dto.CreatePertanyaanDeteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanDeteksi: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-deteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	createReq := dto.CreatePertanyaanDeteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanDeteksi: ""}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)
	repo.On("Create", mock.Anything).Return(int64(0), errors.New("pertanyaan_deteksi tidak boleh kosong"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-deteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Create_SystemError(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	createReq := dto.CreatePertanyaanDeteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanDeteksi: "Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, errors.New("db error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-deteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	producer := new(mockPertanyaanDeteksiProducer)
	handler := setupPertanyaanDeteksiHandler(repo, producer)

	updateReq := dto.UpdatePertanyaanDeteksiRequest{PertanyaanDeteksi: deteksiStrPtr("Pertanyaan Ubah")}
	
	repo.On("GetByID", 1).Return(&dto.PertanyaanDeteksiResponse{ID: 1}, nil)
	repo.On("Update", 1, mock.Anything).Return(nil)
	
	producer.On("PublishPertanyaanDeteksiUpdated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-deteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanDeteksiHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-deteksi/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanDeteksiResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdatePertanyaanDeteksiRequest{PertanyaanDeteksi: deteksiStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-deteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.PertanyaanDeteksiResponse{ID: 1}, nil)
	repo.On("Update", 1, mock.Anything).Return(errors.New("pertanyaan_deteksi tidak boleh kosong"))
	
	body, _ := json.Marshal(dto.UpdatePertanyaanDeteksiRequest{PertanyaanDeteksi: deteksiStrPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-deteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Update_SystemError(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanDeteksiResponse)(nil), errors.New("db error"))

	body, _ := json.Marshal(dto.UpdatePertanyaanDeteksiRequest{PertanyaanDeteksi: deteksiStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-deteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	producer := new(mockPertanyaanDeteksiProducer)
	handler := setupPertanyaanDeteksiHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.PertanyaanDeteksiResponse{ID: 1}, nil)
	repo.On("Delete", 1).Return(nil)
	producer.On("PublishPertanyaanDeteksiDeleted", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanDeteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_Delete_SystemError(t *testing.T) {
	repo := new(mockPertanyaanDeteksiRepository)
	handler := setupPertanyaanDeteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanDeteksiResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-deteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanDeteksiHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupPertanyaanDeteksiHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-deteksi/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-deteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-deteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/pertanyaan-deteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func deteksiStrPtr(s string) *string {
	return &s
}
