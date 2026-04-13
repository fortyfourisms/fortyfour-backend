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

// mockIdentifikasiRepository implements repository.IdentifikasiRepositoryInterface for testing purposes.
type mockIdentifikasiRepository struct {
	GetAllFn      func() ([]models.Identifikasi, error)
	GetByIDFn     func(id string) (*models.Identifikasi, error)
	GetByIkasIDFn func(ikasID string) ([]models.Identifikasi, error)
}

func (m *mockIdentifikasiRepository) GetAll() ([]models.Identifikasi, error) {
	return m.GetAllFn()
}

func (m *mockIdentifikasiRepository) GetByID(id string) (*models.Identifikasi, error) {
	return m.GetByIDFn(id)
}

func (m *mockIdentifikasiRepository) GetByIkasID(ikasID string) ([]models.Identifikasi, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

var _ repository.IdentifikasiRepositoryInterface = (*mockIdentifikasiRepository)(nil)

func setupIdentifikasiHandler(repo repository.IdentifikasiRepositoryInterface, ikasRepo repository.IkasRepositoryInterface) *IdentifikasiHandler {
	service := services.NewIdentifikasiService(repo, ikasRepo)
	return NewIdentifikasiHandler(service)
}

func TestIdentifikasiHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetAllFn: func() ([]models.Identifikasi, error) {
			return []models.Identifikasi{
				{ID: "1"},
				{ID: "2"},
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupIdentifikasiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/identifikasi", nil)
	// Inject admin role so GetAll works without ikas_id
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

func TestIdentifikasiHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetAllFn: func() ([]models.Identifikasi, error) {
			return nil, errors.New("database error")
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupIdentifikasiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/identifikasi", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestIdentifikasiHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetByIDFn: func(id string) (*models.Identifikasi, error) {
			return &models.Identifikasi{
				ID: "uuid-test",
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupIdentifikasiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/identifikasi/uuid-test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Berhasil mengambil data", response["message"])
	assert.NotNil(t, response["data"])
}

func TestIdentifikasiHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := &mockIdentifikasiRepository{
		GetByIDFn: func(id string) (*models.Identifikasi, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupIdentifikasiHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/identifikasi/invalid-id", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIdentifikasiHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	repo := &mockIdentifikasiRepository{}
	ikasRepo := new(mockIkasRepository)
	handler := setupIdentifikasiHandler(repo, ikasRepo)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/identifikasi", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}
