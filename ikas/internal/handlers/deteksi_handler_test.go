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

// mockDeteksiRepository implements repository.DeteksiRepositoryInterface for testing purposes.
type mockDeteksiRepository struct {
	GetAllFn            func() ([]models.Deteksi, error)
	GetByIDFn           func(id string) (*models.Deteksi, error)
	GetByIkasIDFn       func(ikasID string) ([]models.Deteksi, error)
	GetByPerusahaanIDFn func(perusahaanID string) ([]models.Deteksi, error)
}

func (m *mockDeteksiRepository) GetAll() ([]models.Deteksi, error) {
	return m.GetAllFn()
}

func (m *mockDeteksiRepository) GetByID(id string) (*models.Deteksi, error) {
	return m.GetByIDFn(id)
}

func (m *mockDeteksiRepository) GetByIkasID(ikasID string) ([]models.Deteksi, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

func (m *mockDeteksiRepository) GetByPerusahaanID(perusahaanID string) ([]models.Deteksi, error) {
	if m.GetByPerusahaanIDFn != nil {
		return m.GetByPerusahaanIDFn(perusahaanID)
	}
	return nil, nil
}

var _ repository.DeteksiRepositoryInterface = (*mockDeteksiRepository)(nil)

func setupDeteksiHandler(repo repository.DeteksiRepositoryInterface, ikasRepo repository.IkasRepositoryInterface) *DeteksiHandler {
	service := services.NewDeteksiService(repo, ikasRepo)
	return NewDeteksiHandler(service)
}

func TestDeteksiHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetAllFn: func() ([]models.Deteksi, error) {
			return []models.Deteksi{
				{ID: "1", NilaiDeteksi: 80.5},
				{ID: "2", NilaiDeteksi: 90.0},
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupDeteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/deteksi", nil)
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
	assert.Equal(t, float64(2), response["total"])
	assert.NotNil(t, response["data"])
}

func TestDeteksiHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetAllFn: func() ([]models.Deteksi, error) {
			return nil, errors.New("database error")
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupDeteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/deteksi", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeteksiHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return &models.Deteksi{
				ID:           "uuid-test",
				NilaiDeteksi: 85.0,
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupDeteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/deteksi/uuid-test", nil)
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

func TestDeteksiHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := &mockDeteksiRepository{
		GetByIDFn: func(id string) (*models.Deteksi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupDeteksiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/deteksi/invalid-id", nil)
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeteksiHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	repo := &mockDeteksiRepository{}
	ikasRepo := &mockIkasRepository{}
	handler := setupDeteksiHandler(repo, ikasRepo)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/deteksi", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}
