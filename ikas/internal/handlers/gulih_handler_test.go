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

// mockGulihRepository implements repository.GulihRepositoryInterface for testing purposes.
type mockGulihRepository struct {
	GetAllFn      func() ([]models.Gulih, error)
	GetByIDFn     func(id string) (*models.Gulih, error)
	GetByIkasIDFn func(ikasID string) ([]models.Gulih, error)
}

func (m *mockGulihRepository) GetAll() ([]models.Gulih, error) {
	return m.GetAllFn()
}

func (m *mockGulihRepository) GetByID(id string) (*models.Gulih, error) {
	return m.GetByIDFn(id)
}

func (m *mockGulihRepository) GetByIkasID(ikasID string) ([]models.Gulih, error) {
	if m.GetByIkasIDFn != nil {
		return m.GetByIkasIDFn(ikasID)
	}
	return nil, nil
}

var _ repository.GulihRepositoryInterface = (*mockGulihRepository)(nil)

func setupGulihHandler(repo repository.GulihRepositoryInterface, ikasRepo repository.IkasRepositoryInterface) *GulihHandler {
	service := services.NewGulihService(repo, ikasRepo)
	return NewGulihHandler(service)
}

func TestGulihHandler_ServeHTTP_GetAll_Success(t *testing.T) {
	repo := &mockGulihRepository{
		GetAllFn: func() ([]models.Gulih, error) {
			return []models.Gulih{
				{ID: "1", NilaiGulih: 3.5},
				{ID: "2", NilaiGulih: 4.0},
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupGulihHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/gulih", nil)
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

func TestGulihHandler_ServeHTTP_GetAll_Error(t *testing.T) {
	repo := &mockGulihRepository{
		GetAllFn: func() ([]models.Gulih, error) {
			return nil, errors.New("database error")
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupGulihHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/gulih", nil)
	// Inject admin role
	ctx := context.WithValue(req.Context(), middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGulihHandler_ServeHTTP_GetByID_Success(t *testing.T) {
	repo := &mockGulihRepository{
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return &models.Gulih{
				ID: "uuid-test",
			}, nil
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupGulihHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/gulih/uuid-test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Berhasil mengambil data", response["message"])
	assert.NotNil(t, response["data"])
}

func TestGulihHandler_ServeHTTP_GetByID_Error(t *testing.T) {
	repo := &mockGulihRepository{
		GetByIDFn: func(id string) (*models.Gulih, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}
	ikasRepo := new(mockIkasRepository)
	handler := setupGulihHandler(repo, ikasRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/gulih/invalid-id", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGulihHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	repo := &mockGulihRepository{}
	ikasRepo := new(mockIkasRepository)
	handler := setupGulihHandler(repo, ikasRepo)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/gulih", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}
