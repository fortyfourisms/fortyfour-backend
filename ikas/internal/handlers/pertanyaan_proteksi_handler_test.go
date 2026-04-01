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

// mockPertanyaanProteksiRepository implements repository.PertanyaanProteksiRepositoryInterface
type mockPertanyaanProteksiRepository struct {
	mock.Mock
}

func (m *mockPertanyaanProteksiRepository) Create(req dto.CreatePertanyaanProteksiRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPertanyaanProteksiRepository) GetAll() ([]dto.PertanyaanProteksiResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.PertanyaanProteksiResponse), args.Error(1)
}

func (m *mockPertanyaanProteksiRepository) GetByID(id int) (*dto.PertanyaanProteksiResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PertanyaanProteksiResponse), args.Error(1)
}

func (m *mockPertanyaanProteksiRepository) GetTotalCount() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *mockPertanyaanProteksiRepository) Update(id int, req dto.UpdatePertanyaanProteksiRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockPertanyaanProteksiRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockPertanyaanProteksiRepository) CheckSubKategoriExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockPertanyaanProteksiRepository) CheckRuangLingkupExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

// mockPertanyaanProteksiProducer implements services.PertanyaanProteksiProducerInterface
type mockPertanyaanProteksiProducer struct {
	mock.Mock
}

func (m *mockPertanyaanProteksiProducer) PublishPertanyaanProteksiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanProteksiProducer) PublishPertanyaanProteksiUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanProteksiProducer) PublishPertanyaanProteksiDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.PertanyaanProteksiRepositoryInterface = (*mockPertanyaanProteksiRepository)(nil)
var _ services.PertanyaanProteksiProducerInterface = (*mockPertanyaanProteksiProducer)(nil)

func setupPertanyaanProteksiHandler(repo *mockPertanyaanProteksiRepository, producer *mockPertanyaanProteksiProducer) *PertanyaanProteksiHandler {
	service := services.NewPertanyaanProteksiService(repo, producer)
	return NewPertanyaanProteksiHandler(service)
}

func TestPertanyaanProteksiHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	producer := new(mockPertanyaanProteksiProducer)
	handler := setupPertanyaanProteksiHandler(repo, producer)

	expectedData := []dto.PertanyaanProteksiResponse{{ID: 1}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-proteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	producer := new(mockPertanyaanProteksiProducer)
	handler := setupPertanyaanProteksiHandler(repo, producer)

	repo.On("GetAll").Return([]dto.PertanyaanProteksiResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-proteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	producer := new(mockPertanyaanProteksiProducer)
	handler := setupPertanyaanProteksiHandler(repo, producer)

	expectedData := &dto.PertanyaanProteksiResponse{ID: 1}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-proteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanProteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-proteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanProteksiResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-proteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	producer := new(mockPertanyaanProteksiProducer)
	handler := setupPertanyaanProteksiHandler(repo, producer)

	createReq := dto.CreatePertanyaanProteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanProteksi: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)
	repo.On("Create", mock.Anything).Return(int64(1), nil)
	
	producer.On("PublishPertanyaanProteksiCreated", mock.Anything, mock.MatchedBy(func(e dto_event.PertanyaanProteksiCreatedEvent) bool {
		return e.Request.PertanyaanProteksi == "Pertanyaan Test"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-proteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanProteksiHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-proteksi", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Create_SubKategoriNotFound(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	createReq := dto.CreatePertanyaanProteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanProteksi: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-proteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	createReq := dto.CreatePertanyaanProteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanProteksi: ""}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)
	repo.On("Create", mock.Anything).Return(int64(0), errors.New("pertanyaan_proteksi tidak boleh kosong"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-proteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Create_SystemError(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	createReq := dto.CreatePertanyaanProteksiRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanProteksi: "Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, errors.New("db error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-proteksi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	producer := new(mockPertanyaanProteksiProducer)
	handler := setupPertanyaanProteksiHandler(repo, producer)

	updateReq := dto.UpdatePertanyaanProteksiRequest{PertanyaanProteksi: proteksiStrPtr("Pertanyaan Ubah")}
	
	repo.On("GetByID", 1).Return(&dto.PertanyaanProteksiResponse{ID: 1}, nil)
	repo.On("Update", 1, mock.Anything).Return(nil)
	
	producer.On("PublishPertanyaanProteksiUpdated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-proteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	t.Log("RESPONSE BODY:", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanProteksiHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-proteksi/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanProteksiResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdatePertanyaanProteksiRequest{PertanyaanProteksi: proteksiStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-proteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.PertanyaanProteksiResponse{ID: 1}, nil)
	repo.On("Update", 1, mock.Anything).Return(errors.New("pertanyaan_proteksi tidak boleh kosong"))
	
	body, _ := json.Marshal(dto.UpdatePertanyaanProteksiRequest{PertanyaanProteksi: proteksiStrPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-proteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Update_SystemError(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanProteksiResponse)(nil), errors.New("db error"))

	body, _ := json.Marshal(dto.UpdatePertanyaanProteksiRequest{PertanyaanProteksi: proteksiStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-proteksi/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	producer := new(mockPertanyaanProteksiProducer)
	handler := setupPertanyaanProteksiHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.PertanyaanProteksiResponse{ID: 1}, nil)
	repo.On("Delete", 1).Return(nil)
	producer.On("PublishPertanyaanProteksiDeleted", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-proteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanProteksiResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-proteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_Delete_SystemError(t *testing.T) {
	repo := new(mockPertanyaanProteksiRepository)
	handler := setupPertanyaanProteksiHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanProteksiResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-proteksi/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanProteksiHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupPertanyaanProteksiHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-proteksi/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-proteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-proteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/pertanyaan-proteksi", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func proteksiStrPtr(s string) *string {
	return &s
}
