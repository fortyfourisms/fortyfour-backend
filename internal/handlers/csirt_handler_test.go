package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

//
// MOCK SERVICE
//

type mockCsirtService struct {
	GetAllFn            func() ([]dto.CsirtResponse, error)
	GetByIDFn           func(string) (*dto.CsirtResponse, error)
	GetByPerusahaanFn   func(string) ([]dto.CsirtResponse, error)
	CreateFn            func(dto.CreateCsirtRequest) (*models.Csirt, error)
	UpdateFn            func(string, dto.UpdateCsirtRequest) (*models.Csirt, error)
	DeleteFn            func(string) error
}

func (m *mockCsirtService) GetAll() ([]dto.CsirtResponse, error) {
	return m.GetAllFn()
}
func (m *mockCsirtService) GetByID(id string) (*dto.CsirtResponse, error) {
	return m.GetByIDFn(id)
}
func (m *mockCsirtService) GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error) {
	return m.GetByPerusahaanFn(idPerusahaan)
}
func (m *mockCsirtService) Create(req dto.CreateCsirtRequest) (*models.Csirt, error) {
	return m.CreateFn(req)
}
func (m *mockCsirtService) Update(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
	return m.UpdateFn(id, req)
}
func (m *mockCsirtService) Delete(id string) error {
	return m.DeleteFn(id)
}

// withCsirtAdminContext set role admin di context
func withCsirtAdminContext(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "admin")
	return req.WithContext(ctx)
}

//
// HELPERS
//

func createMultipartRequest(
	t *testing.T,
	method, url string,
	fields map[string]string,
	files map[string]string,
) (*http.Request, func()) {

	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	for field, filename := range files {
		part, err := writer.CreateFormFile(field, filename)
		if err != nil {
			t.Fatalf("failed create form file: %v", err)
		}
		_, _ = io.WriteString(part, "dummy content")
	}

	writer.Close()

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	cleanup := func() {
		_ = os.RemoveAll("uploads")
	}

	return req, cleanup
}

//
// TESTS — GET ALL
//

func TestCsirtHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetAllFn: func() ([]dto.CsirtResponse, error) {
			return []dto.CsirtResponse{
				{ID: "1", NamaCsirt: "CSIRT A"},
			}, nil
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt", nil)
	req = withCsirtAdminContext(req)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []dto.CsirtResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "CSIRT A", response[0].NamaCsirt)
}

func TestCsirtHandler_GetAll_ServiceError(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetAllFn: func() ([]dto.CsirtResponse, error) {
			return nil, errors.New("db error")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt", nil)
	req = withCsirtAdminContext(req)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

//
// TESTS — GET BY ID
//

func TestCsirtHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID:        id,
				NamaCsirt: "CSIRT A",
			}, nil
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response dto.CsirtResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID)
	assert.Equal(t, "CSIRT A", response.NamaCsirt)
}

func TestCsirtHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt/123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

//
// TESTS — CREATE
//

func TestCsirtHandler_Create_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		CreateFn: func(req dto.CreateCsirtRequest) (*models.Csirt, error) {
			assert.Equal(t, "CSIRT A", req.NamaCsirt)
			return &models.Csirt{
				ID:        "new-id",
				NamaCsirt: req.NamaCsirt,
			}, nil
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(
		t,
		http.MethodPost,
		"/api/csirt",
		map[string]string{
			"nama_csirt":    "CSIRT A",
			"id_perusahaan": "1",
		},
		map[string]string{
			"photo_csirt": "image.jpg",
		},
	)
	t.Cleanup(cleanup)
	req = withCsirtAdminContext(req)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Csirt
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "new-id", response.ID)
	assert.Equal(t, "CSIRT A", response.NamaCsirt)
}

func TestCsirtHandler_Create_InvalidMultipart(t *testing.T) {
	mockSvc := &mockCsirtService{}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/csirt", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	req = withCsirtAdminContext(req)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCsirtHandler_Create_ServiceError(t *testing.T) {
	mockSvc := &mockCsirtService{
		CreateFn: func(req dto.CreateCsirtRequest) (*models.Csirt, error) {
			return nil, errors.New("create failed")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(
		t,
		http.MethodPost,
		"/api/csirt",
		map[string]string{
			"id_perusahaan": "1",
			"nama_csirt":    "CSIRT",
		},
		nil,
	)
	t.Cleanup(cleanup)
	req = withCsirtAdminContext(req)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

//
// TESTS — UPDATE
//

func TestCsirtHandler_Update_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		UpdateFn: func(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
			assert.Equal(t, "1", id)
			assert.Equal(t, "CSIRT Updated", *req.NamaCsirt)
			return &models.Csirt{
				ID:        id,
				NamaCsirt: *req.NamaCsirt,
			}, nil
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(
		t,
		http.MethodPut,
		"/api/csirt/1",
		map[string]string{
			"nama_csirt": "CSIRT Updated",
		},
		nil,
	)
	t.Cleanup(cleanup)
	req = withCsirtAdminContext(req)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Csirt
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID)
	assert.Equal(t, "CSIRT Updated", response.NamaCsirt)
}

func TestCsirtHandler_Update_NotFound(t *testing.T) {
	mockSvc := &mockCsirtService{
		UpdateFn: func(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(
		t,
		http.MethodPut,
		"/api/csirt/999",
		map[string]string{
			"nama_csirt": "CSIRT Updated",
		},
		nil,
	)
	t.Cleanup(cleanup)
	req = withCsirtAdminContext(req)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Current implementation returns 400 for all errors in Update
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCsirtHandler_Update_Error(t *testing.T) {
	mockSvc := &mockCsirtService{
		UpdateFn: func(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
			return nil, errors.New("db error")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(
		t,
		http.MethodPut,
		"/api/csirt/1",
		map[string]string{
			"nama_csirt": "CSIRT Updated",
		},
		nil,
	)
	t.Cleanup(cleanup)
	req = withCsirtAdminContext(req)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Current implementation returns 400 for all errors in Update
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

//
// TESTS — DELETE
//

func TestCsirtHandler_Delete_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		DeleteFn: func(id string) error {
			assert.Equal(t, "1", id)
			return nil
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/csirt/1", nil)
	req = withCsirtAdminContext(req)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCsirtHandler_Delete_NotFound(t *testing.T) {
	mockSvc := &mockCsirtService{
		DeleteFn: func(id string) error {
			return errors.New("not found")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/csirt/999", nil)
	req = withCsirtAdminContext(req)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Current implementation returns 400 for all errors in Delete
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCsirtHandler_Delete_ServiceError(t *testing.T) {
	mockSvc := &mockCsirtService{
		DeleteFn: func(id string) error {
			return errors.New("delete failed")
		},
	}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/csirt/1", nil)
	req = withCsirtAdminContext(req)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

//
// TESTS — METHOD NOT ALLOWED
//

func TestCsirtHandler_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockCsirtService{}

	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPatch, "/api/csirt", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
/*
=====================================
 HELPER — USER CONTEXT (CSIRT)
=====================================
*/

func withCsirtUserContext(req *http.Request, idPerusahaan string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "user")
	ctx = context.WithValue(ctx, middleware.IDPerusahaanKey, idPerusahaan)
	return req.WithContext(ctx)
}

/*
=====================================
 TEST OWNERSHIP — GET ALL AS USER
=====================================
*/

func TestCsirtHandler_GetAll_AsUser_FilterByPerusahaan(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByPerusahaanFn: func(idPerusahaan string) ([]dto.CsirtResponse, error) {
			assert.Equal(t, "perusahaan-abc", idPerusahaan)
			return []dto.CsirtResponse{{ID: "csirt-1", NamaCsirt: "CSIRT ABC"}}, nil
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt", nil)
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.CsirtResponse
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Len(t, resp, 1)
}

func TestCsirtHandler_GetAll_AsUser_NoPerusahaan_Forbidden(t *testing.T) {
	handler := NewCsirtHandler(&mockCsirtService{})

	req := httptest.NewRequest(http.MethodGet, "/api/csirt", nil)
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCsirtHandler_GetAll_AsUser_ServiceError(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByPerusahaanFn: func(string) ([]dto.CsirtResponse, error) {
			return nil, errors.New("db error")
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt", nil)
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

/*
=====================================
 TEST OWNERSHIP — GET BY ID AS USER
=====================================
*/

func TestCsirtHandler_GetByID_AsUser_OwnData_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID: id, Perusahaan: dto.PerusahaanResponse{ID: "perusahaan-abc"},
			}, nil
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt/csirt-1", nil)
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCsirtHandler_GetByID_AsUser_OtherPerusahaan_Forbidden(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID: id, Perusahaan: dto.PerusahaanResponse{ID: "perusahaan-lain"},
			}, nil
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/csirt/csirt-1", nil)
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

/*
=====================================
 TEST OWNERSHIP — CREATE AS USER
=====================================
*/

func TestCsirtHandler_Create_AsUser_IDPerusahaanForcedFromJWT(t *testing.T) {
	var capturedReq dto.CreateCsirtRequest
	mockSvc := &mockCsirtService{
		CreateFn: func(req dto.CreateCsirtRequest) (*models.Csirt, error) {
			capturedReq = req
			return &models.Csirt{ID: "csirt-new"}, nil
		},
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{ID: id, NamaCsirt: "New"}, nil
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(t, http.MethodPost, "/api/csirt", map[string]string{
		"id_perusahaan": "perusahaan-lain", // harus di-override
		"nama_csirt":    "CSIRT Baru",
	}, nil)
	defer cleanup()
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "perusahaan-abc", capturedReq.IdPerusahaan)
}

func TestCsirtHandler_Create_AsUser_NoPerusahaan_Forbidden(t *testing.T) {
	handler := NewCsirtHandler(&mockCsirtService{})

	req, cleanup := createMultipartRequest(t, http.MethodPost, "/api/csirt", map[string]string{
		"nama_csirt": "CSIRT Baru",
	}, nil)
	defer cleanup()
	ctx := context.WithValue(req.Context(), middleware.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

/*
=====================================
 TEST OWNERSHIP — UPDATE AS USER
=====================================
*/

func TestCsirtHandler_Update_AsUser_OwnData_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID: id, Perusahaan: dto.PerusahaanResponse{ID: "perusahaan-abc"},
			}, nil
		},
		UpdateFn: func(id string, req dto.UpdateCsirtRequest) (*models.Csirt, error) {
			return &models.Csirt{ID: id}, nil
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(t, http.MethodPut, "/api/csirt/csirt-1", map[string]string{
		"nama_csirt": "Updated",
	}, nil)
	defer cleanup()
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCsirtHandler_Update_AsUser_OtherPerusahaan_Forbidden(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID: id, Perusahaan: dto.PerusahaanResponse{ID: "perusahaan-lain"},
			}, nil
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req, cleanup := createMultipartRequest(t, http.MethodPut, "/api/csirt/csirt-1", nil, nil)
	defer cleanup()
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

/*
=====================================
 TEST OWNERSHIP — DELETE AS USER
=====================================
*/

func TestCsirtHandler_Delete_AsUser_OwnData_Success(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID: id, Perusahaan: dto.PerusahaanResponse{ID: "perusahaan-abc"},
			}, nil
		},
		DeleteFn: func(id string) error { return nil },
	}
	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/csirt/csirt-1", nil)
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCsirtHandler_Delete_AsUser_OtherPerusahaan_Forbidden(t *testing.T) {
	mockSvc := &mockCsirtService{
		GetByIDFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID: id, Perusahaan: dto.PerusahaanResponse{ID: "perusahaan-lain"},
			}, nil
		},
	}
	handler := NewCsirtHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/csirt/csirt-1", nil)
	req = withCsirtUserContext(req, "perusahaan-abc")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}