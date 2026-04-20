package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/middleware"
	"ikas/internal/repository"
	"ikas/internal/services"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockIkasProducer implements services.IkasProducerInterface
type mockIkasProducer struct {
	mock.Mock
}

func (m *mockIkasProducer) PublishIkasCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishIkasUpdated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishIkasDeleted(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishIkasAuditLog(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishIkasImported(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishJawabanIdentifikasiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishJawabanProteksiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishJawabanDeteksiCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockIkasProducer) PublishJawabanGulihCreated(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Check interfaces compliance
var _ repository.IkasRepositoryInterface = (*mockIkasRepository)(nil)
var _ services.IkasProducerInterface = (*mockIkasProducer)(nil)


func setupIkasHandler(repo repository.IkasRepositoryInterface, producer services.IkasProducerInterface) *IkasHandler {
	service := services.NewIkasService(
		repo,
		nil, // identifikasiRepo
		nil, // proteksiRepo
		nil, // deteksiRepo
		nil, // gulihRepo
		nil, // jawabanIdentifikasiRepo
		nil, // jawabanProteksiRepo
		nil, // jawabanDeteksiRepo
		nil, // jawabanGulihRepo
		producer,
	)
	return NewIkasHandler(service)
}

func TestIkasHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	expectedData := []dto.IkasResponse{{ID: "123", Responden: "Test"}}
	repo.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ikas", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Berhasil mengambil data", response["message"])
}

func TestIkasHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	repo.On("GetAll").Return([]dto.IkasResponse{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ikas", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestIkasHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	expectedData := &dto.IkasResponse{ID: "123", Responden: "Test"}
	repo.On("GetByID", "123").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ikas/123", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIkasHandler_ServeHTTP_GetByID_NotFound(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	repo.On("GetByID", "123").Return((*dto.IkasResponse)(nil), errors.New("data tidak ditemukan"))

	req := httptest.NewRequest(http.MethodGet, "/api/maturity/ikas/123", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIkasHandler_ServeHTTP_Create_Success(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	createReq := dto.CreateIkasRequest{IDPerusahaan: "1", Responden: "User", Tanggal: "2026-01-01"}
	repo.On("CheckExistsByPerusahaanIDAndYear", "1", 2026).Return(false, nil)

	producer.On("PublishIkasCreated", mock.Anything, mock.MatchedBy(func(e dto_event.IkasCreatedEvent) bool {
		return e.IDPerusahaan == "1" && e.UserID == "user-123"
	})).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.MatchedBy(func(e dto_event.IkasAuditLogEvent) bool {
		return e.Action == "CREATE_IKAS"
	})).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ikas", bytes.NewBuffer(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(context.WithValue(ctx, middleware.UserIDKey, "user-123"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIkasHandler_ServeHTTP_Create_InvalidJSON(t *testing.T) {
	handler := setupIkasHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ikas", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Create_ExistsOrError(t *testing.T) {
	repo := new(mockIkasRepository)
	handler := setupIkasHandler(repo, nil)

	createReq := dto.CreateIkasRequest{IDPerusahaan: "1", Tanggal: "2026-01-01"}
	repo.On("CheckExistsByPerusahaanIDAndYear", "1", 2026).Return(true, nil) // Conflict

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ikas", bytes.NewBuffer(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Update_Success(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	perusahaanID := "2"
	updateReq := dto.UpdateIkasRequest{IDPerusahaan: &perusahaanID}

	current := &dto.IkasResponse{
		ID:         "123",
		Tanggal:    "2026-01-01",
		Perusahaan: &dto.PerusahaanInIkas{ID: "1"},
	}
	repo.On("GetByID", "123").Return(current, nil)
	repo.On("GetLatestByPerusahaan", "1").Return((*dto.IkasResponse)(nil), nil)

	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasUpdated", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ikas/123", bytes.NewBuffer(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(context.WithValue(ctx, middleware.UserIDKey, "user-123"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIkasHandler_ServeHTTP_Update_InvalidJSON(t *testing.T) {
	handler := setupIkasHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ikas/123", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Update_NotFound(t *testing.T) {
	repo := new(mockIkasRepository)
	handler := setupIkasHandler(repo, nil)

	repo.On("GetByID", "123").Return((*dto.IkasResponse)(nil), errors.New("sql: no rows in result set"))

	body, _ := json.Marshal(dto.UpdateIkasRequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ikas/123", bytes.NewBuffer(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIkasHandler_ServeHTTP_Update_SystemError(t *testing.T) {
	repo := new(mockIkasRepository)
	handler := setupIkasHandler(repo, nil)

	repo.On("GetByID", "123").Return((*dto.IkasResponse)(nil), errors.New("db connection failed"))

	body, _ := json.Marshal(dto.UpdateIkasRequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/maturity/ikas/123", bytes.NewBuffer(body))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Delete_Success(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	repo.On("GetByID", "123").Return(&dto.IkasResponse{ID: "123"}, nil)

	producer.On("PublishIkasDeleted", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ikas/123", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(context.WithValue(ctx, middleware.UserIDKey, "user-123"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIkasHandler_ServeHTTP_Delete_NotFound(t *testing.T) {
	repo := new(mockIkasRepository)
	handler := setupIkasHandler(repo, nil)

	repo.On("GetByID", "123").Return((*dto.IkasResponse)(nil), errors.New("sql: no rows in result set"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ikas/123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIkasHandler_ServeHTTP_Delete_SystemError(t *testing.T) {
	repo := new(mockIkasRepository)
	handler := setupIkasHandler(repo, nil)

	repo.On("GetByID", "123").Return((*dto.IkasResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ikas/123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Import_InvalidForm(t *testing.T) {
	handler := setupIkasHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ikas/import", nil)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=invalid")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func createMultipartRequest(t *testing.T, field string, filename string, content []byte) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	if filename != "" {
		part, err := writer.CreateFormFile(field, filename)
		assert.NoError(t, err)
		part.Write(content)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/maturity/ikas/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestIkasHandler_ServeHTTP_Import_NoFile(t *testing.T) {
	handler := setupIkasHandler(nil, nil)

	req := createMultipartRequest(t, "wrong_field", "test.xlsx", []byte("dummy"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Import_WrongExtension(t *testing.T) {
	handler := setupIkasHandler(nil, nil)

	req := createMultipartRequest(t, "file", "test.txt", []byte("dummy"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Import_ParseError(t *testing.T) {
	repo := new(mockIkasRepository)
	handler := setupIkasHandler(repo, nil)

	repo.On("ParseExcelForImport", mock.Anything).Return((*dto.ParsedExcelData)(nil), errors.New("parse error"))

	req := createMultipartRequest(t, "file", "test.xlsx", []byte("dummy data"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_Import_Success(t *testing.T) {
	repo := new(mockIkasRepository)
	producer := new(mockIkasProducer)
	handler := setupIkasHandler(repo, producer)

	importData := &dto.ParsedExcelData{
		IkasRequest: dto.CreateIkasRequest{
			IDPerusahaan: "1",
		},
		JawabanIdentifikasi: []dto.ExcelSubdomainAnswer{{PertanyaanID: 1, Jawaban: 1.0}},
	}

	repo.On("ParseExcelForImport", mock.Anything).Return(importData, nil)
	repo.On("CheckExistsByPerusahaanIDAndYear", "1", 2026).Return(false, nil) // Create IKAS

	producer.On("PublishIkasCreated", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishIkasAuditLog", mock.Anything, mock.Anything).Return(nil)
	producer.On("PublishJawabanIdentifikasiCreated", mock.Anything, mock.Anything).Return(nil)

	req := createMultipartRequest(t, "file", "test.xlsx", []byte("dummy data"))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIkasHandler_ServeHTTP_Import_SystemError(t *testing.T) {
	repo := new(mockIkasRepository)
	handler := setupIkasHandler(repo, nil)

	importData := &dto.ParsedExcelData{
		IkasRequest: dto.CreateIkasRequest{
			IDPerusahaan: "1",
		},
		JawabanIdentifikasi: []dto.ExcelSubdomainAnswer{{PertanyaanID: 1, Jawaban: 1.0}},
	}

	repo.On("ParseExcelForImport", mock.Anything).Return(importData, nil)
	repo.On("CheckExistsByPerusahaanIDAndYear", "1", 2026).Return(false, errors.New("db error"))

	req := createMultipartRequest(t, "file", "test.xlsx", []byte("dummy data"))
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIkasHandler_ServeHTTP_MethodValidation(t *testing.T) {
	handler := setupIkasHandler(nil, nil)

	t.Run("POST with ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/maturity/ikas/1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/maturity/ikas", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE without ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/maturity/ikas", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/maturity/ikas", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
