package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/services"
)

// ═══════════════════════════════════════════════════════════════════════════
// MOCK REPOSITORIES untuk AktivitasService
// ═══════════════════════════════════════════════════════════════════════════

type MockAktivitasRepo struct {
	CreateFunc          func(req dto.CreateAktivitasRequest) (int64, error)
	GetAllFunc          func() ([]dto.AktivitasResponse, error)
	GetByIDFunc         func(id int) (*dto.AktivitasResponse, error)
	GetByPerusahaanFunc func(perusahaanID string) ([]dto.AktivitasResponse, error)
	UpdateFunc          func(id int, req dto.UpdateAktivitasRequest) error
	DeleteFunc          func(id int) error
}

func (m *MockAktivitasRepo) Create(req dto.CreateAktivitasRequest) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(req)
	}
	return 1, nil
}
func (m *MockAktivitasRepo) GetAll() ([]dto.AktivitasResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}
func (m *MockAktivitasRepo) GetByID(id int) (*dto.AktivitasResponse, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}
func (m *MockAktivitasRepo) GetByPerusahaanID(perusahaanID string) ([]dto.AktivitasResponse, error) {
	if m.GetByPerusahaanFunc != nil {
		return m.GetByPerusahaanFunc(perusahaanID)
	}
	return nil, nil
}
func (m *MockAktivitasRepo) Update(id int, req dto.UpdateAktivitasRequest) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, req)
	}
	return nil
}
func (m *MockAktivitasRepo) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

type MockPerusahaanRepoForHandler struct {
	GetByIDFunc func(id string) (*dto.PerusahaanResponse, error)
}

func (m *MockPerusahaanRepoForHandler) Create(req dto.CreatePerusahaanRequest, id string) error {
	return nil
}
func (m *MockPerusahaanRepoForHandler) GetByID(id string) (*dto.PerusahaanResponse, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return &dto.PerusahaanResponse{ID: id}, nil
}
func (m *MockPerusahaanRepoForHandler) GetByNama(nama string) (*dto.PerusahaanResponse, error) {
	return nil, nil
}
func (m *MockPerusahaanRepoForHandler) GetAll() ([]dto.PerusahaanResponse, error) {
	return nil, nil
}
func (m *MockPerusahaanRepoForHandler) Update(id string, p dto.PerusahaanResponse) error {
	return nil
}
func (m *MockPerusahaanRepoForHandler) Delete(id string) error { return nil }

// helper: buat AktivitasHandler dengan mock repo
func newAktivitasHandler(aktRepo *MockAktivitasRepo, perusahaanRepo *MockPerusahaanRepoForHandler) *handlers.AktivitasHandler {
	svc := services.NewAktivitasService(aktRepo, perusahaanRepo, nil, nil)
	return handlers.NewAktivitasHandler(svc)
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GET /api/aktivitas — GetAll
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasHandler_GetAll_Success(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetAllFunc: func() ([]dto.AktivitasResponse, error) {
				return []dto.AktivitasResponse{
					{ID: 1, Judul: "Kegiatan 1"},
					{ID: 2, Judul: "Kegiatan 2"},
				}, nil
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/aktivitas", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAktivitasHandler_GetAll_Error(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetAllFunc: func() ([]dto.AktivitasResponse, error) {
				return nil, errors.New("db error")
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/aktivitas", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GET /api/aktivitas?perusahaan_id=xxx — GetByPerusahaanID
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasHandler_GetByPerusahaanID_Success(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetByPerusahaanFunc: func(pid string) ([]dto.AktivitasResponse, error) {
				return []dto.AktivitasResponse{{ID: 1, PerusahaanID: pid}}, nil
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/aktivitas?perusahaan_id=uuid-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GET /api/aktivitas/{id} — GetByID
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasHandler_GetByID_Success(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
				return &dto.AktivitasResponse{ID: id, Judul: "Test"}, nil
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/aktivitas/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAktivitasHandler_GetByID_NotFound(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
				return nil, sql.ErrNoRows
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/aktivitas/999", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: POST /api/aktivitas — Create
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasHandler_Create_Success(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{},
		&MockPerusahaanRepoForHandler{
			GetByIDFunc: func(id string) (*dto.PerusahaanResponse, error) {
				return &dto.PerusahaanResponse{ID: id}, nil
			},
		},
	)

	body := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-123",
		Judul:          "Kegiatan Baru",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas"},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/aktivitas", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d; body: %s", rec.Code, rec.Body.String())
	}
}

func TestAktivitasHandler_Create_InvalidBody(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	req := httptest.NewRequest(http.MethodPost, "/api/aktivitas", bytes.NewReader([]byte("{invalid")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAktivitasHandler_Create_ValidationError(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	// Missing required fields
	body := dto.CreateAktivitasRequest{Judul: "Test"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/aktivitas", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAktivitasHandler_Create_WithIDInPath(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	body := dto.CreateAktivitasRequest{Judul: "Test"}
	bodyBytes, _ := json.Marshal(body)

	// POST with ID in path should be rejected
	req := httptest.NewRequest(http.MethodPost, "/api/aktivitas/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for POST with ID, got %d", rec.Code)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: PUT /api/aktivitas/{id} — Update
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasHandler_Update_Success(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
				return &dto.AktivitasResponse{ID: id, PerusahaanID: "uuid-1"}, nil
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	judul := "Updated"
	body := dto.UpdateAktivitasRequest{Judul: &judul}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/aktivitas/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
}

func TestAktivitasHandler_Update_NotFound(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
				return nil, sql.ErrNoRows
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	judul := "Updated"
	body := dto.UpdateAktivitasRequest{Judul: &judul}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/aktivitas/999", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestAktivitasHandler_Update_InvalidBody(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	req := httptest.NewRequest(http.MethodPut, "/api/aktivitas/1", bytes.NewReader([]byte("{invalid")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAktivitasHandler_Update_NoID(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	req := httptest.NewRequest(http.MethodPut, "/api/aktivitas", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for PUT without ID, got %d", rec.Code)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: DELETE /api/aktivitas/{id} — Delete
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasHandler_Delete_Success(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
				return &dto.AktivitasResponse{ID: id, PerusahaanID: "uuid-1"}, nil
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/api/aktivitas/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAktivitasHandler_Delete_NotFound(t *testing.T) {
	handler := newAktivitasHandler(
		&MockAktivitasRepo{
			GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
				return nil, sql.ErrNoRows
			},
		},
		&MockPerusahaanRepoForHandler{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/api/aktivitas/999", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestAktivitasHandler_Delete_NoID(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	req := httptest.NewRequest(http.MethodDelete, "/api/aktivitas", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for DELETE without ID, got %d", rec.Code)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Method Not Allowed / Jenis Aktivitas
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasHandler_MethodNotAllowed(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	req := httptest.NewRequest(http.MethodPatch, "/api/aktivitas", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestAktivitasHandler_GetJenis(t *testing.T) {
	handler := newAktivitasHandler(&MockAktivitasRepo{}, &MockPerusahaanRepoForHandler{})

	req := httptest.NewRequest(http.MethodGet, "/api/aktivitas/jenis", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
