package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"

	"github.com/stretchr/testify/assert"
)

/* =========================
   SETUP HELPERS
========================= */

func setupPICHandler() (*PICHandler, repository.PICRepositoryInterface, *services.SSEService) {
	mockRepo := testhelpers.NewMockPICRepository()
	sseService := services.NewSSEService()
	picService := services.NewPICService(mockRepo)
	handler := NewPICHandler(picService, sseService)
	return handler, mockRepo, sseService
}

/* =========================
   TEST GET ALL
========================= */

func TestPICHandler_GetAll_Success(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	// Create test data
	mockRepo.Create(dto.CreatePICRequest{
		Nama:    stringPtr("PIC 1"),
		Telepon: stringPtr("081234567890"),
	}, "id-1")
	mockRepo.Create(dto.CreatePICRequest{
		Nama:    stringPtr("PIC 2"),
		Telepon: stringPtr("081234567891"),
	}, "id-2")

	req := httptest.NewRequest(http.MethodGet, "/api/pic", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.PICResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestPICHandler_GetAll_EmptyResult(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/pic", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.PICResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 0)
}

func TestPICHandler_GetAll_VerifyContentType(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/pic", nil)
	w := httptest.NewRecorder()

	handler.handleGetAll(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

/* =========================
   TEST GET BY ID
========================= */

func TestPICHandler_GetByID_Success(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Test PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	req := httptest.NewRequest(http.MethodGet, "/api/pic/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "test-id")

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.PICResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "test-id", response.ID)
	assert.Equal(t, "Test PIC", response.Nama)
}

func TestPICHandler_GetByID_NotFound(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/pic/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.handleGetByID(w, req, "nonexistent")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Data tidak ditemukan")
}

func TestPICHandler_GetByID_VariousIDFormats(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	testIDs := []string{
		"uuid-1234-5678-90ab-cdef",
		"simple-id",
		"id-with-numbers-123",
		"CamelCaseID",
	}

	for _, id := range testIDs {
		t.Run(id, func(t *testing.T) {
			mockRepo.Create(dto.CreatePICRequest{
				Nama:         stringPtr("Test PIC"),
				Telepon:      stringPtr("081234567890"),
				IDPerusahaan: stringPtr("perusahaan-1"),
			}, id)

			req := httptest.NewRequest(http.MethodGet, "/api/pic/"+id, nil)
			w := httptest.NewRecorder()

			handler.handleGetByID(w, req, id)

			assert.Equal(t, http.StatusOK, w.Code)

			var response dto.PICResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, id, response.ID)
		})
	}
}

/* =========================
   TEST CREATE
========================= */

func TestPICHandler_Create_Success_MinimalData(t *testing.T) {
	handler, _, _ := setupPICHandler()

	reqBody := dto.CreatePICRequest{
		Nama:         stringPtr("New PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.handleCreate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.PICResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "New PIC", response.Nama)
}

func TestPICHandler_Create_Success_CompleteData(t *testing.T) {
	handler, _, _ := setupPICHandler()

	reqBody := dto.CreatePICRequest{
		Nama:         stringPtr("Complete PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-123"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.handleCreate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.PICResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Complete PIC", response.Nama)
}

func TestPICHandler_Create_Success_WithoutUserContext(t *testing.T) {
	handler, _, _ := setupPICHandler()

	reqBody := dto.CreatePICRequest{
		Nama:         stringPtr("New PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.handleCreate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPICHandler_Create_InvalidBody(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request body")
}

func TestPICHandler_Create_EmptyBody(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	// Should fail validation in service layer
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPICHandler_Create_MissingContentType(t *testing.T) {
	handler, _, _ := setupPICHandler()

	reqBody := dto.CreatePICRequest{
		Nama:         stringPtr("New PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
	// No Content-Type header
	w := httptest.NewRecorder()

	handler.handleCreate(w, req)

	// Should still work as Go's json decoder is lenient
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPICHandler_Create_WithIDInURL_ShouldFail(t *testing.T) {
	handler, _, _ := setupPICHandler()

	reqBody := dto.CreatePICRequest{
		Nama:         stringPtr("New PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/pic/some-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID tidak diperlukan")
}

func TestPICHandler_Create_MissingRequiredFields(t *testing.T) {
	handler, _, _ := setupPICHandler()

	testCases := []struct {
		name    string
		reqBody dto.CreatePICRequest
	}{
		{
			name: "Missing Nama",
			reqBody: dto.CreatePICRequest{
				Telepon:      stringPtr("081234567890"),
				IDPerusahaan: stringPtr("perusahaan-1"),
			},
		},
		{
			name: "Missing IDPerusahaan",
			reqBody: dto.CreatePICRequest{
				Nama:    stringPtr("Test PIC"),
				Telepon: stringPtr("081234567890"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.reqBody)

			req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.handleCreate(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

/* =========================
   TEST UPDATE
========================= */

func TestPICHandler_Update_Success_PartialUpdate(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Old Name"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	updateReq := dto.UpdatePICRequest{
		Nama: stringPtr("New Name"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/pic/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.PICResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", response.Nama)
}

func TestPICHandler_Update_Success_CompleteUpdate(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Old Name"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	updateReq := dto.UpdatePICRequest{
		Nama:         stringPtr("Updated Name"),
		Telepon:      stringPtr("089876543210"),
		IDPerusahaan: stringPtr("perusahaan-2"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/pic/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.PICResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", response.Nama)
	assert.Equal(t, "089876543210", response.Telepon)
}

func TestPICHandler_Update_NotFound(t *testing.T) {
	handler, _, _ := setupPICHandler()

	updateReq := dto.UpdatePICRequest{
		Nama: stringPtr("New Name"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/pic/nonexistent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "nonexistent")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPICHandler_Update_InvalidBody(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/pic/test-id", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request body")
}

func TestPICHandler_Update_EmptyBody(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Test"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	req := httptest.NewRequest(http.MethodPut, "/api/pic/test-id", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	// Empty update should still return OK
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPICHandler_Update_WithoutID_ShouldFail(t *testing.T) {
	handler, _, _ := setupPICHandler()

	updateReq := dto.UpdatePICRequest{
		Nama: stringPtr("New Name"),
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/pic", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID wajib")
}

/* =========================
   TEST DELETE
========================= */

func TestPICHandler_Delete_Success(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Test PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/pic/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Delete success", response["message"])
}

func TestPICHandler_Delete_NotFound(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/pic/nonexistent", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "nonexistent")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPICHandler_Delete_WithoutUserContext(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Test PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/pic/test-id", nil)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPICHandler_Delete_WithoutID_ShouldFail(t *testing.T) {
	handler, _, _ := setupPICHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/pic", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID wajib")
}

func TestPICHandler_Delete_VerifyDeleted(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Test PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	req := httptest.NewRequest(http.MethodDelete, "/api/pic/test-id", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleDelete(w, req, "test-id")

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify it's deleted
	req2 := httptest.NewRequest(http.MethodGet, "/api/pic/test-id", nil)
	w2 := httptest.NewRecorder()
	handler.handleGetByID(w2, req2, "test-id")

	assert.Equal(t, http.StatusNotFound, w2.Code)
}

/* =========================
   TEST METHOD NOT ALLOWED
========================= */

func TestPICHandler_MethodNotAllowed(t *testing.T) {
	handler, _, _ := setupPICHandler()

	methods := []string{http.MethodPatch, http.MethodOptions, http.MethodHead}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/pic", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}

/* =========================
   TEST SERVE HTTP ROUTING
========================= */

func TestPICHandler_ServeHTTP_Routes(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	testCases := []struct {
		name           string
		method         string
		path           string
		setupData      func()
		expectedStatus int
	}{
		{
			name:           "GET all - empty",
			method:         http.MethodGet,
			path:           "/api/pic",
			setupData:      func() {},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET all - with data",
			method: http.MethodGet,
			path:   "/api/pic",
			setupData: func() {
				mockRepo.Create(dto.CreatePICRequest{
					Nama:         stringPtr("PIC 1"),
					Telepon:      stringPtr("081234567890"),
					IDPerusahaan: stringPtr("perusahaan-1"),
				}, "id-1")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET by ID - found",
			method: http.MethodGet,
			path:   "/api/pic/test-id",
			setupData: func() {
				mockRepo.Create(dto.CreatePICRequest{
					Nama:         stringPtr("Test PIC"),
					Telepon:      stringPtr("081234567890"),
					IDPerusahaan: stringPtr("perusahaan-1"),
				}, "test-id")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET by ID - not found",
			method:         http.MethodGet,
			path:           "/api/pic/nonexistent",
			setupData:      func() {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "POST with ID - should fail",
			method:         http.MethodPost,
			path:           "/api/pic/test-id",
			setupData:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "PUT without ID - should fail",
			method:         http.MethodPut,
			path:           "/api/pic",
			setupData:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "DELETE without ID - should fail",
			method:         http.MethodDelete,
			path:           "/api/pic",
			setupData:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupData()

			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

/* =========================
   TEST EDGE CASES
========================= */

func TestPICHandler_Create_SpecialCharactersInName(t *testing.T) {
	handler, _, _ := setupPICHandler()

	specialNames := []string{
		"Name with spaces",
		"Name-with-dashes",
		"Name_with_underscores",
		"Name.with.dots",
		"Name (with parentheses)",
		"Ñame wïth åccents",
	}

	for _, name := range specialNames {
		t.Run(name, func(t *testing.T) {
			reqBody := dto.CreatePICRequest{
				Nama:         stringPtr(name),
				Telepon:      stringPtr("081234567890"),
				IDPerusahaan: stringPtr("perusahaan-1"),
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.handleCreate(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
		})
	}
}

func TestPICHandler_Update_ClearOptionalFields(t *testing.T) {
	handler, mockRepo, _ := setupPICHandler()

	mockRepo.Create(dto.CreatePICRequest{
		Nama:         stringPtr("Test PIC"),
		Telepon:      stringPtr("081234567890"),
		IDPerusahaan: stringPtr("perusahaan-1"),
	}, "test-id")

	// Update to clear optional fields
	emptyString := ""
	updateReq := dto.UpdatePICRequest{
		Telepon: &emptyString,
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/pic/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.handleUpdate(w, req, "test-id")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPICHandler_ConcurrentRequests(t *testing.T) {
	handler, _, _ := setupPICHandler()

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(index int) {
			reqBody := dto.CreatePICRequest{
				Nama:         stringPtr("PIC Concurrent"),
				Telepon:      stringPtr("081234567890"),
				IDPerusahaan: stringPtr("perusahaan-1"),
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/api/pic", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.handleCreate(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

/* =========================
   TEST RESPONSE HEADERS
========================= */

func TestPICHandler_VerifyJSONHeaders(t *testing.T) {
	handler, _, _ := setupPICHandler()

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/pic"},
		{http.MethodGet, "/api/pic/nonexistent"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			req := httptest.NewRequest(ep.method, ep.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}
