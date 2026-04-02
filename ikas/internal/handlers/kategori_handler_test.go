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

// mockKategoriRepository implements repository.KategoriRepositoryInterface
type mockKategoriRepository struct {
	mock.Mock
}

func (m *mockKategoriRepository) Create(req dto.CreateKategoriRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockKategoriRepository) GetAll() ([]dto.KategoriResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.KategoriResponse), args.Error(1)
}

func (m *mockKategoriRepository) GetByID(id int) (*dto.KategoriResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.KategoriResponse), args.Error(1)
}

func (m *mockKategoriRepository) Update(id int, req dto.UpdateKategoriRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockKategoriRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockKategoriRepository) CheckDuplicateName(domainID int, namaKategori string, excludeID int) (bool, error) {
	args := m.Called(domainID, namaKategori, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockKategoriRepository) CheckDomainExists(domainID int) (bool, error) {
	args := m.Called(domainID)
	return args.Get(0).(bool), args.Error(1)
}

// mockKategoriProducer implements services.KategoriProducerInterface
type mockKategoriProducer struct {
	mock.Mock
}

func (m *mockKategoriProducer) PublishKategoriCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockKategoriProducer) PublishKategoriUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockKategoriProducer) PublishKategoriDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.KategoriRepositoryInterface = (*mockKategoriRepository)(nil)
var _ services.KategoriProducerInterface = (*mockKategoriProducer)(nil)

func setupKategoriHandler(repo *mockKategoriRepository, producer *mockKategoriProducer) *KategoriHandler {
	service := services.NewKategoriService(repo, producer)
	return NewKategoriHandler(service)
}

func TestKategoriHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	expectedData := []dto.KategoriResponse{{ID: 1, NamaKategori: "Kategori Test"}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/kategori", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Berhasil mengambil data", response["message"])
}

func TestKategoriHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	repo.On("GetAll").Return([]dto.KategoriResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/kategori", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestKategoriHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	expectedData := &dto.KategoriResponse{ID: 1, NamaKategori: "Kategori Test"}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestKategoriHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.KategoriResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKategoriHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.KategoriResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestKategoriHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	createReq := dto.CreateKategoriRequest{DomainID: 1, NamaKategori: "Kategori Finance"}
	repo.On("CheckDomainExists", 1).Return(true, nil)
	repo.On("CheckDuplicateName", 1, "Kategori Finance", 0).Return(false, nil)
	producer.On("PublishKategoriCreated", mock.Anything, mock.MatchedBy(func(e dto_event.KategoriCreatedEvent) bool {
		return e.Request.NamaKategori == "Kategori Finance"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestKategoriHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupKategoriHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/kategori", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestKategoriHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	handler := setupKategoriHandler(nil, nil)

	createReq := dto.CreateKategoriRequest{DomainID: 1, NamaKategori: ""} // Invalid (empty)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestKategoriHandler_ServeHTTP_Create_DomainNotFound(t *testing.T) {
	repo := new(mockKategoriRepository)
	handler := setupKategoriHandler(repo, nil)

	createReq := dto.CreateKategoriRequest{DomainID: 99, NamaKategori: "Kategori Finance"}
	repo.On("CheckDomainExists", 99).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKategoriHandler_ServeHTTP_Create_Duplicate(t *testing.T) {
	repo := new(mockKategoriRepository)
	handler := setupKategoriHandler(repo, nil)

	createReq := dto.CreateKategoriRequest{DomainID: 1, NamaKategori: "Kategori Kembar"}
	repo.On("CheckDomainExists", 1).Return(true, nil)
	repo.On("CheckDuplicateName", 1, "Kategori Kembar", 0).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestKategoriHandler_ServeHTTP_Create_PublishError(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	createReq := dto.CreateKategoriRequest{DomainID: 1, NamaKategori: "Kategori Test"}
	repo.On("CheckDomainExists", 1).Return(true, nil)
	repo.On("CheckDuplicateName", 1, "Kategori Test", 0).Return(false, nil)
	producer.On("PublishKategoriCreated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestKategoriHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	updateReq := dto.UpdateKategoriRequest{NamaKategori: katStrPtr("Kategori Finance")}
	repo.On("GetByID", 1).Return(&dto.KategoriResponse{ID: 1, DomainID: 1}, nil)
	repo.On("CheckDuplicateName", 1, "Kategori Finance", 1).Return(false, nil)
	producer.On("PublishKategoriUpdated", mock.Anything, mock.MatchedBy(func(e dto_event.KategoriUpdatedEvent) bool {
		return e.ID == 1 && *e.Request.NamaKategori == "Kategori Finance"
	})).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestKategoriHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockKategoriRepository)
	handler := setupKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.KategoriResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdateKategoriRequest{NamaKategori: katStrPtr("x")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKategoriHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupKategoriHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestKategoriHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockKategoriRepository)
	handler := setupKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.KategoriResponse{ID: 1, DomainID: 1}, nil)

	// Trigger validation error by passing empty string
	body, _ := json.Marshal(dto.UpdateKategoriRequest{NamaKategori: katStrPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestKategoriHandler_ServeHTTP_Update_DomainNotFound(t *testing.T) {
	repo := new(mockKategoriRepository)
	handler := setupKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.KategoriResponse{ID: 1, DomainID: 1}, nil)
	repo.On("CheckDomainExists", 99).Return(false, nil)

	domainID := 99
	body, _ := json.Marshal(dto.UpdateKategoriRequest{DomainID: &domainID, NamaKategori: katStrPtr("Kategori Finance")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKategoriHandler_ServeHTTP_Update_Conflict(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.KategoriResponse{ID: 1, DomainID: 1}, nil)
	repo.On("CheckDuplicateName", 1, "Kategori Kembar", 1).Return(true, nil)

	body, _ := json.Marshal(dto.UpdateKategoriRequest{NamaKategori: katStrPtr("Kategori Kembar")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestKategoriHandler_ServeHTTP_Update_PublishError(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.KategoriResponse{ID: 1, DomainID: 1}, nil)
	repo.On("CheckDuplicateName", 1, "Kategori Finance", 1).Return(false, nil)
	producer.On("PublishKategoriUpdated", mock.Anything, mock.Anything).Return(errors.New("db error"))

	body, _ := json.Marshal(dto.UpdateKategoriRequest{NamaKategori: katStrPtr("Kategori Finance")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestKategoriHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockKategoriRepository)
	producer := new(mockKategoriProducer)
	handler := setupKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.KategoriResponse{ID: 1}, nil)
	producer.On("PublishKategoriDeleted", mock.Anything, mock.MatchedBy(func(e dto_event.KategoriDeletedEvent) bool {
		return e.ID == 1
	})).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestKategoriHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockKategoriRepository)
	handler := setupKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.KategoriResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestKategoriHandler_ServeHTTP_Delete_Error(t *testing.T) {
	repo := new(mockKategoriRepository)
	handler := setupKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.KategoriResponse)(nil), errors.New("internal error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestKategoriHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupKategoriHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/kategori/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/kategori", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/kategori", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/kategori", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func katStrPtr(s string) *string {
	return &s
}
