package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"ikas/internal/middleware"
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
	GetAllFn      func() ([]models.Proteksi, error)
	GetByIDFn     func(id string) (*models.Proteksi, error)
	GetByIkasIDFn        func(ikasID string) ([]models.Proteksi, error)
	GetByPerusahaanIDFn  func(perusahaanID string) ([]models.Proteksi, error)
}

func (m *mockProteksiRepository) GetAll() ([]models.Proteksi, error) {
	return m.GetAllFn()
}

func (m *mockProteksiRepository) GetByID(id string) (*models.Proteksi, error) {
	return m.GetByIDFn(id)
}

func (m *mockProteksiRepository) GetByIkasID(ikasID string) ([]models.Proteksi, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

func (m *mockProteksiRepository) GetByPerusahaanID(perusahaanID string) ([]models.Proteksi, error) {
	if m.GetByPerusahaanIDFn != nil {
		return m.GetByPerusahaanIDFn(perusahaanID)
	}
	return nil, nil
}

var _ repository.ProteksiRepositoryInterface = (*mockProteksiRepository)(nil)

func setupProteksiHandler(repo repository.ProteksiRepositoryInterface, ikasRepo repository.IkasRepositoryInterface) *ProteksiHandler {
	service := services.NewProteksiService(repo, ikasRepo)
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
	ikasRepo := new(mockIkasRepository)
	handler := setupProteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

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
	ikasRepo := new(mockIkasRepository)
	handler := setupProteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

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
	ikasRepo := new(mockIkasRepository)
	handler := setupProteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi/uuid-test", nil)
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
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
	ikasRepo := new(mockIkasRepository)
	handler := setupProteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/proteksi/invalid-id", nil)
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProteksiHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	repo := &mockProteksiRepository{}
	ikasRepo := new(mockIkasRepository)
	handler := setupProteksiHandler(repo, ikasRepo)

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
