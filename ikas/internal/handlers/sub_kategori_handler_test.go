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

// mockSubKategoriRepository implements repository.SubKategoriRepositoryInterface
type mockSubKategoriRepository struct {
	mock.Mock
}

func (m *mockSubKategoriRepository) Create(req dto.CreateSubKategoriRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockSubKategoriRepository) GetAll() ([]dto.SubKategoriResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.SubKategoriResponse), args.Error(1)
}

func (m *mockSubKategoriRepository) GetByID(id int) (*dto.SubKategoriResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SubKategoriResponse), args.Error(1)
}

func (m *mockSubKategoriRepository) Update(id int, req dto.UpdateSubKategoriRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockSubKategoriRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockSubKategoriRepository) CheckDuplicateName(kategoriID int, namaSubKategori string, excludeID int) (bool, error) {
	args := m.Called(kategoriID, namaSubKategori, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockSubKategoriRepository) CheckKategoriExists(kategoriID int) (bool, error) {
	args := m.Called(kategoriID)
	return args.Get(0).(bool), args.Error(1)
}

// mockSubKategoriProducer implements services.SubKategoriProducerInterface
type mockSubKategoriProducer struct {
	mock.Mock
}

func (m *mockSubKategoriProducer) PublishSubKategoriCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockSubKategoriProducer) PublishSubKategoriUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockSubKategoriProducer) PublishSubKategoriDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.SubKategoriRepositoryInterface = (*mockSubKategoriRepository)(nil)
var _ services.SubKategoriProducerInterface = (*mockSubKategoriProducer)(nil)

func setupSubKategoriHandler(repo *mockSubKategoriRepository, producer *mockSubKategoriProducer) *SubKategoriHandler {
	service := services.NewSubKategoriService(repo, producer)
	return NewSubKategoriHandler(service)
}

func TestSubKategoriHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	expectedData := []dto.SubKategoriResponse{{ID: 1, NamaSubKategori: "Sub Kategori Test"}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/sub-kategori", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Berhasil mengambil data", response["message"])
}

func TestSubKategoriHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	repo.On("GetAll").Return([]dto.SubKategoriResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/sub-kategori", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	expectedData := &dto.SubKategoriResponse{ID: 1, NamaSubKategori: "Sub Kategori Test"}
	repo.On("GetByID", 1).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/sub-kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.SubKategoriResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/sub-kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return((*dto.SubKategoriResponse)(nil), errors.New("system error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/sub-kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	createReq := dto.CreateSubKategoriRequest{KategoriID: 1, NamaSubKategori: "Sub Kategori Finance"}
	repo.On("CheckKategoriExists", 1).Return(true, nil)
	repo.On("CheckDuplicateName", 1, "Sub Kategori Finance", 0).Return(false, nil)
	producer.On("PublishSubKategoriCreated", mock.Anything, mock.MatchedBy(func(e dto_event.SubKategoriCreatedEvent) bool {
		return e.Request.NamaSubKategori == "Sub Kategori Finance"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/sub-kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupSubKategoriHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/sub-kategori", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Create_ValidationError(t *testing.T) {
	handler := setupSubKategoriHandler(nil, nil)

	createReq := dto.CreateSubKategoriRequest{KategoriID: 1, NamaSubKategori: ""} // Invalid (empty)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/sub-kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Create_KategoriNotFound(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	handler := setupSubKategoriHandler(repo, nil)

	createReq := dto.CreateSubKategoriRequest{KategoriID: 99, NamaSubKategori: "Sub Kategori Finance"}
	repo.On("CheckKategoriExists", 99).Return(false, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/sub-kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Create_Duplicate(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	handler := setupSubKategoriHandler(repo, nil)

	createReq := dto.CreateSubKategoriRequest{KategoriID: 1, NamaSubKategori: "Sub Kategori Kembar"}
	repo.On("CheckKategoriExists", 1).Return(true, nil)
	repo.On("CheckDuplicateName", 1, "Sub Kategori Kembar", 0).Return(true, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/sub-kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Create_PublishError(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	createReq := dto.CreateSubKategoriRequest{KategoriID: 1, NamaSubKategori: "Sub Kategori Test"}
	repo.On("CheckKategoriExists", 1).Return(true, nil)
	repo.On("CheckDuplicateName", 1, "Sub Kategori Test", 0).Return(false, nil)
	producer.On("PublishSubKategoriCreated", mock.Anything, mock.Anything).Return(errors.New("publish error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/sub-kategori", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	updateReq := dto.UpdateSubKategoriRequest{NamaSubKategori: subKatStrPtr("Sub Kategori Finance")}
	repo.On("GetByID", 1).Return(&dto.SubKategoriResponse{ID: 1, KategoriID: 1}, nil)
	repo.On("CheckDuplicateName", 1, "Sub Kategori Finance", 1).Return(false, nil)
	producer.On("PublishSubKategoriUpdated", mock.Anything, mock.MatchedBy(func(e dto_event.SubKategoriUpdatedEvent) bool {
		return e.ID == 1 && *e.Request.NamaSubKategori == "Sub Kategori Finance"
	})).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	handler := setupSubKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.SubKategoriResponse)(nil), errors.New("data tidak ditemukan"))

	body, _ := json.Marshal(dto.UpdateSubKategoriRequest{NamaSubKategori: subKatStrPtr("x")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupSubKategoriHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori/1", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Update_ValidationError(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	handler := setupSubKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.SubKategoriResponse{ID: 1, KategoriID: 1}, nil)

	// Trigger validation error by passing empty string
	body, _ := json.Marshal(dto.UpdateSubKategoriRequest{NamaSubKategori: subKatStrPtr("")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Update_KategoriNotFound(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	handler := setupSubKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return(&dto.SubKategoriResponse{ID: 1, KategoriID: 1}, nil)
	repo.On("CheckKategoriExists", 99).Return(false, nil)

	kategoriID := 99
	body, _ := json.Marshal(dto.UpdateSubKategoriRequest{KategoriID: &kategoriID, NamaSubKategori: subKatStrPtr("Sub Kategori Finance")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Update_Conflict(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.SubKategoriResponse{ID: 1, KategoriID: 1}, nil)
	repo.On("CheckDuplicateName", 1, "Sub Kategori Kembar", 1).Return(true, nil)

	body, _ := json.Marshal(dto.UpdateSubKategoriRequest{NamaSubKategori: subKatStrPtr("Sub Kategori Kembar")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Update_PublishError(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.SubKategoriResponse{ID: 1, KategoriID: 1}, nil)
	repo.On("CheckDuplicateName", 1, "Sub Kategori Finance", 1).Return(false, nil)
	producer.On("PublishSubKategoriUpdated", mock.Anything, mock.Anything).Return(errors.New("db error"))

	body, _ := json.Marshal(dto.UpdateSubKategoriRequest{NamaSubKategori: subKatStrPtr("Sub Kategori Finance")})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	producer := new(mockSubKategoriProducer)
	handler := setupSubKategoriHandler(repo, producer)

	repo.On("GetByID", 1).Return(&dto.SubKategoriResponse{ID: 1}, nil)
	producer.On("PublishSubKategoriDeleted", mock.Anything, mock.MatchedBy(func(e dto_event.SubKategoriDeletedEvent) bool {
		return e.ID == 1
	})).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/sub-kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	handler := setupSubKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.SubKategoriResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/sub-kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_Delete_Error(t *testing.T) {
	repo := new(mockSubKategoriRepository)
	handler := setupSubKategoriHandler(repo, nil)

	repo.On("GetByID", 1).Return((*dto.SubKategoriResponse)(nil), errors.New("internal error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/sub-kategori/1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSubKategoriHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupSubKategoriHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/sub-kategori/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/sub-kategori", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/sub-kategori", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/sub-kategori", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func subKatStrPtr(s string) *string {
	return &s
}
