package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/* =========================
   MOCK SERVICE
========================= */

type mockSektorService struct {
	mock.Mock
}

func (m *mockSektorService) GetAll() ([]dto.SektorResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SektorResponse), args.Error(1)
}

func (m *mockSektorService) GetByID(id string) (*dto.SektorResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SektorResponse), args.Error(1)
}

/* =========================
   SETUP HELPERS
========================= */

func setupSektorHandler() (*handlers.SektorHandler, *mockSektorService) {
	mockSvc := new(mockSektorService)
	handler := handlers.NewSektorHandler(mockSvc)
	return handler, mockSvc
}

/* =========================
   TEST GET ALL
========================= */

func TestSektorHandler_GetAll_Success(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	expectedData := []dto.SektorResponse{
		{
			ID:         "1",
			NamaSektor: "Sektor A",
		},
		{
			ID:         "2",
			NamaSektor: "Sektor B",
		},
	}

	mockSvc.On("GetAll").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []dto.SektorResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Sektor A", response[0].NamaSektor)
	assert.Equal(t, "Sektor B", response[1].NamaSektor)

	mockSvc.AssertExpectations(t)
}

func TestSektorHandler_GetAll_EmptyResult(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	mockSvc.On("GetAll").Return([]dto.SektorResponse{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []dto.SektorResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 0)

	mockSvc.AssertExpectations(t)
}

func TestSektorHandler_GetAll_ServiceError(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	mockSvc.On("GetAll").Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "database error")

	mockSvc.AssertExpectations(t)
}

func TestSektorHandler_GetAll_VerifyContentType(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	mockSvc.On("GetAll").Return([]dto.SektorResponse{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	mockSvc.AssertExpectations(t)
}

func TestSektorHandler_GetAll_MultipleSektors(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	// Create 10 sektors
	sektors := make([]dto.SektorResponse, 10)
	for i := 0; i < 10; i++ {
		sektors[i] = dto.SektorResponse{
			ID:         string(rune('A' + i)),
			NamaSektor: "Sektor " + string(rune('A'+i)),
		}
	}

	mockSvc.On("GetAll").Return(sektors, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []dto.SektorResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 10)

	mockSvc.AssertExpectations(t)
}

/* =========================
   TEST GET BY ID
========================= */

func TestSektorHandler_GetByID_Success(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	expectedData := &dto.SektorResponse{
		ID:         "test-id-123",
		NamaSektor: "Sektor Test",
	}

	mockSvc.On("GetByID", "test-id-123").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor/test-id-123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response dto.SektorResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "test-id-123", response.ID)
	assert.Equal(t, "Sektor Test", response.NamaSektor)

	mockSvc.AssertExpectations(t)
}

func TestSektorHandler_GetByID_NotFound(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	mockSvc.On("GetByID", "nonexistent-id").Return(nil, errors.New("sektor not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/sektor/nonexistent-id", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Data tidak ditemukan")

	mockSvc.AssertExpectations(t)
}

func TestSektorHandler_GetByID_WithSpecialCharacters(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	specialIDs := []string{
		"id-with-dashes",
		"id_with_underscores",
		"123-456-789",
		"uuid-1234-5678-90ab-cdef",
	}

	for _, id := range specialIDs {
		t.Run(id, func(t *testing.T) {
			expectedData := &dto.SektorResponse{
				ID:         id,
				NamaSektor: "Test Sektor",
			}

			mockSvc.On("GetByID", id).Return(expectedData, nil).Once()

			req := httptest.NewRequest(http.MethodGet, "/api/sektor/"+id, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestSektorHandler_GetByID_VerifyContentType(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	expectedData := &dto.SektorResponse{
		ID:         "1",
		NamaSektor: "Sektor A",
	}

	mockSvc.On("GetByID", "1").Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/sektor/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	mockSvc.AssertExpectations(t)
}

/* =========================
   TEST METHOD NOT ALLOWED
========================= */

func TestSektorHandler_MethodNotAllowed(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/sektor", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

			mockSvc.AssertNotCalled(t, "GetAll")
			mockSvc.AssertNotCalled(t, "GetByID")
		})
	}
}

/* =========================
   TEST ROUTING
========================= */

func TestSektorHandler_Routing(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		path           string
		setupMock      func(*mockSektorService)
		expectedStatus int
	}{
		{
			name:   "GET all",
			method: http.MethodGet,
			path:   "/api/sektor",
			setupMock: func(ms *mockSektorService) {
				ms.On("GetAll").Return([]dto.SektorResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET by ID - found",
			method: http.MethodGet,
			path:   "/api/sektor/123",
			setupMock: func(ms *mockSektorService) {
				ms.On("GetByID", "123").Return(&dto.SektorResponse{ID: "123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET by ID - not found",
			method: http.MethodGet,
			path:   "/api/sektor/999",
			setupMock: func(ms *mockSektorService) {
				ms.On("GetByID", "999").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "POST - not allowed",
			method:         http.MethodPost,
			path:           "/api/sektor",
			setupMock:      func(ms *mockSektorService) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "PUT - not allowed",
			method:         http.MethodPut,
			path:           "/api/sektor/123",
			setupMock:      func(ms *mockSektorService) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "DELETE - not allowed",
			method:         http.MethodDelete,
			path:           "/api/sektor/123",
			setupMock:      func(ms *mockSektorService) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, mockSvc := setupSektorHandler()
			tc.setupMock(mockSvc)

			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

/* =========================
   TEST EDGE CASES
========================= */

func TestSektorHandler_GetByID_EmptyID(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	mockSvc.On("GetAll").Return([]dto.SektorResponse{}, nil)

	// Empty ID should route to GetAll
	req := httptest.NewRequest(http.MethodGet, "/api/sektor/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should either return 404 or redirect to GetAll
	// Adjust based on actual implementation
	assert.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusNotFound)

	mockSvc.AssertExpectations(t)
}

func TestSektorHandler_ConcurrentRequests(t *testing.T) {
	handler, mockSvc := setupSektorHandler()

	mockSvc.On("GetAll").Return([]dto.SektorResponse{
		{ID: "1", NamaSektor: "Sektor A"},
	}, nil)

	done := make(chan bool)

	// Simulate 10 concurrent requests
	for i := 0; i < 10; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/api/sektor", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	mockSvc.AssertExpectations(t)
}
