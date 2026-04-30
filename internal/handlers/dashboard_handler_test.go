package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDashboardHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	h := NewDashboardHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/dashboard/sektor", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestDashboardHandler_ServeHTTP_NotFound(t *testing.T) {
	h := NewDashboardHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/unknown", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDashboardHandler_InvalidKategoriSE(t *testing.T) {
	h := NewDashboardHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/kse?kategori_se=Invalid", nil)
	w := httptest.NewRecorder()
	h.SummarySE(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["error"], "kategori_se tidak valid")
}

func TestPtrStr_EmptyString_ReturnsNil(t *testing.T) {
	result := ptrStr("")
	assert.Nil(t, result)
}

func TestPtrStr_NonEmptyString_ReturnsPointer(t *testing.T) {
	result := ptrStr("hello")
	assert.NotNil(t, result)
	assert.Equal(t, "hello", *result)
}

func TestNewDashboardHandler_ReturnsNonNil(t *testing.T) {
	h := NewDashboardHandler(nil)
	assert.NotNil(t, h)
}

