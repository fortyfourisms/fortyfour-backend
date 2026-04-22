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

// mockDomainRepository implements repository.DomainRepositoryInterface
type mockDomainRepository struct {
	mock.Mock
}

func (m *mockDomainRepository) Create(req dto.CreateDomainRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockDomainRepository) GetAll() ([]dto.DomainResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.DomainResponse), args.Error(1)
}

func (m *mockDomainRepository) GetByID(id int) (*dto.DomainResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DomainResponse), args.Error(1)
}

func (m *mockDomainRepository) Update(id int, req dto.UpdateDomainRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockDomainRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockDomainRepository) CheckDuplicateName(nama string, excludeID int) (bool, error) {
	args := m.Called(nama, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

// mockDomainProducer implements services.DomainProducerInterface
type mockDomainProducer struct {
	mock.Mock
}

func (m *mockDomainProducer) PublishDomainCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockDomainProducer) PublishDomainUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockDomainProducer) PublishDomainDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

var _ repository.DomainRepositoryInterface = (*mockDomainRepository)(nil)
var _ services.DomainProducerInterface = (*mockDomainProducer)(nil)

func setupDomainHandler(repo *mockDomainRepository, producer *mockDomainProducer) *DomainHandler {
	service := services.NewDomainService(repo, producer)
	return NewDomainHandler(service)
}

func TestDomainHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	expectedData := []dto.DomainResponse{{ID: 1, NamaDomain: "Test Domain"}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/domain", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Berhasil mengambil data", response["message"])
}

func TestDomainHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	repo.On("GetAll").Return([]dto.DomainResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/domain", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDomainHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	expectedData := &dto.DomainResponse{ID: 1, NamaDomain: "Domain Finance"}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/domain/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDomainHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.DomainResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/domain/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDomainHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.DomainResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/domain/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDomainHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	createReq := dto.CreateDomainRequest{NamaDomain: "Domain Finance"}
	repo.On("CheckDuplicateName", "Domain Finance", 0).Return(false, nil)
	repo.On("Create", createReq).Return(int64(1), nil)
	repo.On("GetByID", 1).Return(&dto.DomainResponse{ID: 1, NamaDomain: "Domain Finance"}, nil)
	producer.On("PublishDomainCreated", mock.Anything, mock.MatchedBy(func(e dto_event.DomainCreatedEvent) bool {
		return e.Request.NamaDomain == "Domain Finance"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/domain", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDomainHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupDomainHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/domain", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	handler := setupDomainHandler(nil, nil)

	createReq := dto.CreateDomainRequest{NamaDomain: ""} // Invalid (empty)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/domain", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_ServeHTTP_Create_Duplicate(t *testing.T) {
	repo := new(mockDomainRepository)
	handler := setupDomainHandler(repo, nil)

	createReq := dto.CreateDomainRequest{NamaDomain: "Duplicate"}
	repo.On("CheckDuplicateName", "Duplicate", 0).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/domain", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDomainHandler_ServeHTTP_Create_PublishError(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	createReq := dto.CreateDomainRequest{NamaDomain: "Domain"}
	repo.On("CheckDuplicateName", "Domain", 0).Return(false, nil)
	repo.On("Create", createReq).Return(int64(1), nil)
	repo.On("GetByID", 1).Return(&dto.DomainResponse{ID: 1, NamaDomain: "Domain"}, nil)
	producer.On("PublishDomainCreated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/domain", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDomainHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	updateReq := dto.UpdateDomainRequest{NamaDomain: strPtr("Domain Finance")}
	repo.On("GetByID", 1).Return(&dto.DomainResponse{ID: 1}, nil)
	repo.On("CheckDuplicateName", "Domain Finance", 1).Return(false, nil)
	repo.On("Update", 1, updateReq).Return(nil)
	producer.On("PublishDomainUpdated", mock.Anything, mock.MatchedBy(func(e dto_event.DomainUpdatedEvent) bool {
		return e.ID == 1 && *e.Request.NamaDomain == "Domain Finance"
	})).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/domain/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDomainHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockDomainRepository)
	handler := setupDomainHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.DomainResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdateDomainRequest{NamaDomain: strPtr("x")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/domain/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDomainHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupDomainHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/domain/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_ServeHTTP_Update_Duplicate(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	updateReq := dto.UpdateDomainRequest{NamaDomain: strPtr("Domain Finance")}
	repo.On("GetByID", 1).Return(&dto.DomainResponse{ID: 1}, nil)
	repo.On("CheckDuplicateName", "Domain Finance", 1).Return(false, nil)
	repo.On("Update", 1, updateReq).Return(nil)
	producer.On("PublishDomainUpdated", mock.Anything, mock.Anything).Return(errors.New("db error"))

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/domain/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDomainHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockDomainRepository)
	handler := setupDomainHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.DomainResponse{ID: 1}, nil)

	// Trigger validation error by passing empty string (which is restricted in service)
	body, _ := json.Marshal(dto.UpdateDomainRequest{NamaDomain: strPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/domain/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_ServeHTTP_Update_Conflict(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.DomainResponse{ID: 1}, nil)
	repo.On("CheckDuplicateName", "Duplicate", 1).Return(true, nil)

	body, _ := json.Marshal(dto.UpdateDomainRequest{NamaDomain: strPtr("Duplicate")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/domain/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDomainHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockDomainRepository)
	producer := new(mockDomainProducer)
	handler := setupDomainHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.DomainResponse{ID: 1}, nil)
	repo.On("Delete", 1).Return(nil)
	producer.On("PublishDomainDeleted", mock.Anything, mock.MatchedBy(func(e dto_event.DomainDeletedEvent) bool {
		return e.ID == 1
	})).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/domain/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDomainHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockDomainRepository)
	handler := setupDomainHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.DomainResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/domain/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDomainHandler_ServeHTTP_Delete_Error(t *testing.T) {
	repo := new(mockDomainRepository)
	handler := setupDomainHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.DomainResponse)(nil), errors.New("internal error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/domain/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDomainHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupDomainHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/domain/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/domain", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/domain", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/domain", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func strPtr(s string) *string {
	return &s
}
