package handlers

import (
	"encoding/json"
	"errors"
	"ikas/internal/models"
	"ikas/internal/repository"
	"ikas/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockProteksiRepository implements repository.ProteksiRepositoryInterface for testing purposes.
type mockProteksiRepository struct {
	GetAllFn  func() ([]models.Proteksi, error)
	GetByIDFn func(id string) (*models.Proteksi, error)
}

func (m *mockProteksiRepository) GetAll() ([]models.Proteksi, error) {
	return m.GetAllFn()
}

func (m *mockProteksiRepository) GetByID(id string) (*models.Proteksi, error) {
	return m.GetByIDFn(id)
}

var _ repository.ProteksiRepositoryInterface = (*mockProteksiRepository)(nil)

func setupProteksiHandler(repo repository.ProteksiRepositoryInterface) *ProteksiHandler {
	service := services.NewProteksiService(repo)
	return NewProteksiHandler(service)
}

func TestProteksiHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := &mockProteksiRepository{
		GetAllFn: func() ([]models.Proteksi, error) {
			return []models.Proteksi{
				{ID: "1", NilaiProteksi: 3.5},
				{ID: "2", NilaiProteksi: 4.0},
			}, nil
		},
	}
	handler := setupProteksiHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Berhasil mengambil data", response["message"])
	assert.Equal(t, float64(2), response["total"]) // JSON unmarshals numbers to float64
	assert.NotNil(t, response["data"])
}

func TestProteksiHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := &mockProteksiRepository{
		GetAllFn: func() ([]models.Proteksi, error) {
			return nil, errors.New("database error")
		},
	}
	handler := setupProteksiHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProteksiHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := &mockProteksiRepository{
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return &models.Proteksi{
				ID: "uuid-test",
			}, nil
		},
	}
	handler := setupProteksiHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi/uuid-test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Berhasil mengambil data", response["message"])
	assert.NotNil(t, response["data"])
}

func TestProteksiHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := &mockProteksiRepository{
		GetByIDFn: func(id string) (*models.Proteksi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}
	handler := setupProteksiHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi/invalid-id", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProteksiHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	repo := &mockProteksiRepository{}
	handler := setupProteksiHandler(repo)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/proteksi", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}
