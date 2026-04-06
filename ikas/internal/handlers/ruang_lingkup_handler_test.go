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

// mockRuangLingkupRepository implements repository.RuangLingkupRepositoryInterface
type mockRuangLingkupRepository struct {
	mock.Mock
}

func (m *mockRuangLingkupRepository) Create(req dto.CreateRuangLingkupRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRuangLingkupRepository) GetAll() ([]dto.RuangLingkupResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.RuangLingkupResponse), args.Error(1)
}

func (m *mockRuangLingkupRepository) GetByID(id int) (*dto.RuangLingkupResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RuangLingkupResponse), args.Error(1)
}

func (m *mockRuangLingkupRepository) Update(id int, req dto.UpdateRuangLingkupRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockRuangLingkupRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockRuangLingkupRepository) CheckDuplicateName(nama string, excludeID int) (bool, error) {
	args := m.Called(nama, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

// mockRuangLingkupProducer implements services.RuangLingkupProducerInterface
type mockRuangLingkupProducer struct {
	mock.Mock
}

func (m *mockRuangLingkupProducer) PublishRuangLingkupCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockRuangLingkupProducer) PublishRuangLingkupUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockRuangLingkupProducer) PublishRuangLingkupDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.RuangLingkupRepositoryInterface = (*mockRuangLingkupRepository)(nil)
var _ services.RuangLingkupProducerInterface = (*mockRuangLingkupProducer)(nil)

func setupRuangLingkupHandler(repo *mockRuangLingkupRepository, producer *mockRuangLingkupProducer) *RuangLingkupHandler {
	service := services.NewRuangLingkupService(repo, producer)
	return NewRuangLingkupHandler(service)
}

func TestRuangLingkupHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	expectedData := []dto.RuangLingkupResponse{{ID: 1, NamaRuangLingkup: "Lingkup Test"}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ruang-lingkup", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Berhasil mengambil data", response["message"])
}

func TestRuangLingkupHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	repo.On("GetAll").Return([]dto.RuangLingkupResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ruang-lingkup", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	expectedData := &dto.RuangLingkupResponse{ID: 1, NamaRuangLingkup: "Lingkup Test"}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ruang-lingkup/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.RuangLingkupResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ruang-lingkup/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.RuangLingkupResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ruang-lingkup/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	createReq := dto.CreateRuangLingkupRequest{NamaRuangLingkup: "Lingkup Finance"}
	repo.On("CheckDuplicateName", "Lingkup Finance", 0).Return(false, nil)
	producer.On("PublishRuangLingkupCreated", mock.Anything, mock.MatchedBy(func(e dto_event.RuangLingkupCreatedEvent) bool {
		return e.Request.NamaRuangLingkup == "Lingkup Finance"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ruang-lingkup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupRuangLingkupHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ruang-lingkup", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	handler := setupRuangLingkupHandler(nil, nil)

	createReq := dto.CreateRuangLingkupRequest{NamaRuangLingkup: ""} // Invalid (empty)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ruang-lingkup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Create_Duplicate(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	handler := setupRuangLingkupHandler(repo, nil)

	createReq := dto.CreateRuangLingkupRequest{NamaRuangLingkup: "Lingkup Kembar"}
	repo.On("CheckDuplicateName", "Lingkup Kembar", 0).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ruang-lingkup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Create_PublishError(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	createReq := dto.CreateRuangLingkupRequest{NamaRuangLingkup: "Lingkup Test"}
	repo.On("CheckDuplicateName", "Lingkup Test", 0).Return(false, nil)
	producer.On("PublishRuangLingkupCreated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ruang-lingkup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	updateReq := dto.UpdateRuangLingkupRequest{NamaRuangLingkup: rlStrPtr("Lingkup Finance")}
	repo.On("GetByID", 1).Return(&dto.RuangLingkupResponse{ID: 1}, nil)
	repo.On("CheckDuplicateName", "Lingkup Finance", 1).Return(false, nil)
	producer.On("PublishRuangLingkupUpdated", mock.Anything, mock.MatchedBy(func(e dto_event.RuangLingkupUpdatedEvent) bool {
		return e.ID == 1 && *e.Request.NamaRuangLingkup == "Lingkup Finance"
	})).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ruang-lingkup/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	handler := setupRuangLingkupHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.RuangLingkupResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdateRuangLingkupRequest{NamaRuangLingkup: rlStrPtr("x")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ruang-lingkup/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupRuangLingkupHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ruang-lingkup/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Update_Duplicate(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.RuangLingkupResponse{ID: 1}, nil)
	repo.On("CheckDuplicateName", "Lingkup Finance", 1).Return(false, nil)
	producer.On("PublishRuangLingkupUpdated", mock.Anything, mock.Anything).Return(errors.New("db error"))

	body, _ := json.Marshal(dto.UpdateRuangLingkupRequest{NamaRuangLingkup: rlStrPtr("Lingkup Finance")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ruang-lingkup/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	handler := setupRuangLingkupHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.RuangLingkupResponse{ID: 1}, nil)

	// Trigger validation error by passing empty string
	body, _ := json.Marshal(dto.UpdateRuangLingkupRequest{NamaRuangLingkup: rlStrPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ruang-lingkup/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Update_Conflict(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.RuangLingkupResponse{ID: 1}, nil)
	repo.On("CheckDuplicateName", "Lingkup Kembar", 1).Return(true, nil)

	body, _ := json.Marshal(dto.UpdateRuangLingkupRequest{NamaRuangLingkup: rlStrPtr("Lingkup Kembar")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ruang-lingkup/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	producer := new(mockRuangLingkupProducer)
	handler := setupRuangLingkupHandler(repo, producer)

	// In the service:
	// _, err := s.repo.GetByID(id)
	// return s.producer.PublishRuangLingkupDeleted(...)
	repo.On("GetByID", 1).Return(&dto.RuangLingkupResponse{ID: 1}, nil)
	producer.On("PublishRuangLingkupDeleted", mock.Anything, mock.MatchedBy(func(e dto_event.RuangLingkupDeletedEvent) bool {
		return e.ID == 1
	})).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ruang-lingkup/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	handler := setupRuangLingkupHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.RuangLingkupResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ruang-lingkup/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_Delete_Error(t *testing.T) {
	repo := new(mockRuangLingkupRepository)
	handler := setupRuangLingkupHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.RuangLingkupResponse)(nil), errors.New("internal error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ruang-lingkup/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRuangLingkupHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupRuangLingkupHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/ruang-lingkup/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/ruang-lingkup", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ruang-lingkup", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/ruang-lingkup", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func rlStrPtr(s string) *string {
	return &s
}
