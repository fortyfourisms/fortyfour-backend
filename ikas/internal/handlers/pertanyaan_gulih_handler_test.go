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

// mockPertanyaanGulihRepository implements repository.PertanyaanGulihRepositoryInterface
type mockPertanyaanGulihRepository struct {
	mock.Mock
}

func (m *mockPertanyaanGulihRepository) Create(req dto.CreatePertanyaanGulihRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPertanyaanGulihRepository) GetAll() ([]dto.PertanyaanGulihResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.PertanyaanGulihResponse), args.Error(1)
}

func (m *mockPertanyaanGulihRepository) GetByID(id int) (*dto.PertanyaanGulihResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PertanyaanGulihResponse), args.Error(1)
}

func (m *mockPertanyaanGulihRepository) GetTotalCount() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *mockPertanyaanGulihRepository) Update(id int, req dto.UpdatePertanyaanGulihRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockPertanyaanGulihRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockPertanyaanGulihRepository) CheckSubKategoriExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockPertanyaanGulihRepository) CheckRuangLingkupExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

// mockPertanyaanGulihProducer implements services.PertanyaanGulihProducerInterface
type mockPertanyaanGulihProducer struct {
	mock.Mock
}

func (m *mockPertanyaanGulihProducer) PublishPertanyaanGulihCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanGulihProducer) PublishPertanyaanGulihUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockPertanyaanGulihProducer) PublishPertanyaanGulihDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.PertanyaanGulihRepositoryInterface = (*mockPertanyaanGulihRepository)(nil)
var _ services.PertanyaanGulihProducerInterface = (*mockPertanyaanGulihProducer)(nil)

func setupPertanyaanGulihHandler(repo *mockPertanyaanGulihRepository, producer *mockPertanyaanGulihProducer) *PertanyaanGulihHandler {
	service := services.NewPertanyaanGulihService(repo, producer)
	return NewPertanyaanGulihHandler(service)
}

func TestPertanyaanGulihHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	producer := new(mockPertanyaanGulihProducer)
	handler := setupPertanyaanGulihHandler(repo, producer)

	expectedData := []dto.PertanyaanGulihResponse{{ID: 1}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-gulih", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	producer := new(mockPertanyaanGulihProducer)
	handler := setupPertanyaanGulihHandler(repo, producer)

	repo.On("GetAll").Return([]dto.PertanyaanGulihResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-gulih", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	producer := new(mockPertanyaanGulihProducer)
	handler := setupPertanyaanGulihHandler(repo, producer)

	expectedData := &dto.PertanyaanGulihResponse{ID: 1}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanGulihResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanGulihResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/pertanyaan-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	producer := new(mockPertanyaanGulihProducer)
	handler := setupPertanyaanGulihHandler(repo, producer)

	createReq := dto.CreatePertanyaanGulihRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanGulih: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)
	repo.On("Create", mock.Anything).Return(int64(1), nil)
	
	producer.On("PublishPertanyaanGulihCreated", mock.Anything, mock.MatchedBy(func(e dto_event.PertanyaanGulihCreatedEvent) bool {
		return e.Request.PertanyaanGulih == "Pertanyaan Test"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-gulih", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanGulihHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-gulih", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Create_SubKategoriNotFound(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	createReq := dto.CreatePertanyaanGulihRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanGulih: "Pertanyaan Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-gulih", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	createReq := dto.CreatePertanyaanGulihRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanGulih: ""}
	repo.On("CheckSubKategoriExists", 1).Return(true, nil)
	repo.On("CheckRuangLingkupExists", 1).Return(true, nil)
	repo.On("Create", mock.Anything).Return(int64(0), errors.New("pertanyaan_gulih tidak boleh kosong"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-gulih", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Create_SystemError(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	createReq := dto.CreatePertanyaanGulihRequest{SubKategoriID: 1, RuangLingkupID: 1, PertanyaanGulih: "Test"}
	repo.On("CheckSubKategoriExists", 1).Return(false, errors.New("db error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-gulih", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	producer := new(mockPertanyaanGulihProducer)
	handler := setupPertanyaanGulihHandler(repo, producer)

	updateReq := dto.UpdatePertanyaanGulihRequest{PertanyaanGulih: gulihStrPtr("Pertanyaan Ubah")}
	
	repo.On("GetByID", 1).Return(&dto.PertanyaanGulihResponse{ID: 1}, nil)
	repo.On("Update", 1, mock.Anything).Return(nil)
	
	producer.On("PublishPertanyaanGulihUpdated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-gulih/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupPertanyaanGulihHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-gulih/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanGulihResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdatePertanyaanGulihRequest{PertanyaanGulih: gulihStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-gulih/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.PertanyaanGulihResponse{ID: 1}, nil)
	repo.On("Update", 1, mock.Anything).Return(errors.New("pertanyaan_gulih tidak boleh kosong"))
	
	body, _ := json.Marshal(dto.UpdatePertanyaanGulihRequest{PertanyaanGulih: gulihStrPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-gulih/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Update_SystemError(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanGulihResponse)(nil), errors.New("db error"))

	body, _ := json.Marshal(dto.UpdatePertanyaanGulihRequest{PertanyaanGulih: gulihStrPtr("Test")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-gulih/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	producer := new(mockPertanyaanGulihProducer)
	handler := setupPertanyaanGulihHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.PertanyaanGulihResponse{ID: 1}, nil)
	repo.On("Delete", 1).Return(nil)
	producer.On("PublishPertanyaanGulihDeleted", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanGulihResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_Delete_SystemError(t *testing.T) {
	repo := new(mockPertanyaanGulihRepository)
	handler := setupPertanyaanGulihHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.PertanyaanGulihResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-gulih/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPertanyaanGulihHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupPertanyaanGulihHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/pertanyaan-gulih/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/pertanyaan-gulih", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/pertanyaan-gulih", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/pertanyaan-gulih", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func gulihStrPtr(s string) *string {
	return &s
}
