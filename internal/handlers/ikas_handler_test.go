package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupIkasHandler() (*IkasHandler, repository.IkasRepositoryInterface, *services.SSEService) {
	mockRepo := testhelpers.NewMockIkasRepository()
	sseSvc := services.NewSSEService()
	ikasSvc := services.NewIkasService(mockRepo)
	handler := NewIkasHandler(ikasSvc, sseSvc)
	return handler, mockRepo, sseSvc
}

func TestIkasHandler_handleGetAll(t *testing.T) {
	h, _, _ := setupIkasHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ikas", nil)
	w := httptest.NewRecorder()
	h.handleGetAll(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestIkasHandler_handleGetByID(t *testing.T) {
	h, repo, _ := setupIkasHandler()
	_ = repo.Create(dto.CreateIkasRequest{
		IDPerusahaan:   "p-1",
		IDIdentifikasi: "i-1",
		IDProteksi:     "pro-1",
		IDDeteksi:      "det-1",
		IDGulih:        "g-1",
		Tanggal:        "2025-06-01",
		Responden:      "Budi",
		Telepon:        "081234",
		Jabatan:        "Manager",
		TargetNilai:    5.0,
	}, "test-id", 0, "i-1", "pro-1", "det-1", "g-1")

	req := httptest.NewRequest(http.MethodGet, "/api/ikas/test-id", nil)
	w := httptest.NewRecorder()
	h.handleGetByID(w, req, "test-id")
	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestIkasHandler_handleGetByID_NotFound(t *testing.T) {
	h, _, _ := setupIkasHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ikas/ghost", nil)
	w := httptest.NewRecorder()
	h.handleGetByID(w, req, "ghost")
	if w.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", w.Code)
	}
}

func TestIkasHandler_handleCreate(t *testing.T) {
	h, _, _ := setupIkasHandler()
	body, _ := json.Marshal(dto.CreateIkasRequest{
		IDPerusahaan:   "p-baru",
		IDIdentifikasi: "i-baru",
		IDProteksi:     "pro-baru",
		IDDeteksi:      "det-baru",
		IDGulih:        "g-baru",
		Tanggal:        "2025-06-02",
		Responden:      "Ani",
		Telepon:        "082222",
		Jabatan:        "Direktur",
		TargetNilai:    5.0,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/ikas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.handleCreate(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("want 201, got %d", w.Code)
	}
}

func TestIkasHandler_handleCreate_InvalidBody(t *testing.T) {
	h, _, _ := setupIkasHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ikas", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.handleCreate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestIkasHandler_handleUpdate(t *testing.T) {
	h, repo, _ := setupIkasHandler()
	_ = repo.Create(dto.CreateIkasRequest{
		IDPerusahaan:   "p-lama",
		IDIdentifikasi: "i-lama",
		IDProteksi:     "pro-lama",
		IDDeteksi:      "det-lama",
		IDGulih:        "g-lama",
		Tanggal:        "2025-06-01",
		Responden:      "Rudi",
		Telepon:        "08333",
		Jabatan:        "Staff",
		TargetNilai:    4.5,
	}, "upd-id", 0, "i-1", "pro-1", "det-1", "g-1")

	update := dto.UpdateIkasRequest{
		Responden: strPtr("Rudy"),
		TargetNilai: float64Ptr(4.9),
	}
	b, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPut, "/api/ikas/upd-id", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.handleUpdate(w, req, "upd-id")
	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestIkasHandler_handleDelete(t *testing.T) {
	h, repo, _ := setupIkasHandler()
	_ = repo.Create(dto.CreateIkasRequest{
		IDPerusahaan:   "p-hps",
		IDIdentifikasi: "i-hps",
		IDProteksi:     "pro-hps",
		IDDeteksi:      "det-hps",
		IDGulih:        "g-hps",
		Tanggal:        "2025-06-01",
		Responden:      "Dede",
		Telepon:        "08444",
		Jabatan:        "OB",
		TargetNilai:    4.0,
	}, "del-id", 0, "i-1", "pro-1", "det-1", "g-1")

	req := httptest.NewRequest(http.MethodDelete, "/api/ikas/del-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.handleDelete(w, req, "del-id")
	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestIkasHandler_ServeHTTP(t *testing.T) {
	h, _, _ := setupIkasHandler()
	cases := []struct {
		name string
		m    string
		path string
		code int
	}{
		{"GET all", http.MethodGet, "/api/ikas", http.StatusOK},
		{"GET by ID", http.MethodGet, "/api/ikas/xx", http.StatusNotFound},
		{"POST dg ID", http.MethodPost, "/api/ikas/123", http.StatusBadRequest},
		{"PUT tanpa ID", http.MethodPut, "/api/ikas", http.StatusBadRequest},
		{"DELETE tanpa ID", http.MethodDelete, "/api/ikas", http.StatusBadRequest},
		{"PATCH not allowed", http.MethodPatch, "/api/ikas", http.StatusMethodNotAllowed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.m, tc.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			if w.Code != tc.code {
				t.Errorf("%s: want %d, got %d", tc.name, tc.code, w.Code)
			}
		})
	}
}

/* helper */
func strPtr(s string) *string    { return &s }
func float64Ptr(f float64) *float64 { return &f }