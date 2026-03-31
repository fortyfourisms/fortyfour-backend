package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ════════════════════════════════════════════════════════════════════════════
// MOCK SE EXPORT SERVICE
// ════════════════════════════════════════════════════════════════════════════

type mockSEExportService struct {
	mock.Mock
}

func (m *mockSEExportService) ExportAllPDF() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockSEExportService) ExportByPerusahaanPDF(idPerusahaan string) ([]byte, error) {
	args := m.Called(idPerusahaan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockSEExportService) ExportByIDPDF(id string) (*dto.SEResponse, []byte, error) {
	args := m.Called(id)
	var se *dto.SEResponse
	if args.Get(0) != nil {
		se = args.Get(0).(*dto.SEResponse)
	}
	var pdf []byte
	if args.Get(1) != nil {
		pdf = args.Get(1).([]byte)
	}
	return se, pdf, args.Error(2)
}

// ════════════════════════════════════════════════════════════════════════════
// MOCK CSIRT EXPORT SERVICE
// ════════════════════════════════════════════════════════════════════════════

type mockCsirtExportService struct {
	mock.Mock
}

func (m *mockCsirtExportService) ExportAllPDF() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCsirtExportService) ExportByPerusahaanPDF(idPerusahaan string) ([]byte, error) {
	args := m.Called(idPerusahaan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCsirtExportService) ExportByIDPDF(id string) (*dto.CsirtResponse, []byte, error) {
	args := m.Called(id)
	var csirt *dto.CsirtResponse
	if args.Get(0) != nil {
		csirt = args.Get(0).(*dto.CsirtResponse)
	}
	var pdf []byte
	if args.Get(1) != nil {
		pdf = args.Get(1).([]byte)
	}
	return csirt, pdf, args.Error(2)
}

// ════════════════════════════════════════════════════════════════════════════
// CONTEXT HELPERS
// ════════════════════════════════════════════════════════════════════════════

func withExportAdminCtx(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.RoleKey, "admin")
	return r.WithContext(ctx)
}

func withExportUserCtx(r *http.Request, idPerusahaan string) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.RoleKey, "user")
	ctx = context.WithValue(ctx, middleware.IDPerusahaanKey, idPerusahaan)
	return r.WithContext(ctx)
}

func withExportUserCtxNoPerusahaan(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.RoleKey, "user")
	return r.WithContext(ctx)
}

var fakePDF = []byte("%PDF-1.4 fake pdf content")

// ════════════════════════════════════════════════════════════════════════════
// SE EXPORT HANDLER — ExportAll
// ════════════════════════════════════════════════════════════════════════════

func TestSEExportHandler_ExportAll_Admin_Success(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportAllPDF").Return(fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "laporan-se.pdf")
	mockSvc.AssertExpectations(t)
}

func TestSEExportHandler_ExportAll_Admin_FilterPerusahaan(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-abc").Return(fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf?id_perusahaan=p-abc", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	mockSvc.AssertExpectations(t)
}

func TestSEExportHandler_ExportAll_Admin_ServiceError(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportAllPDF").Return(nil, assert.AnError)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSEExportHandler_ExportAll_User_Success(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-user").Return(fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	mockSvc.AssertExpectations(t)
}

func TestSEExportHandler_ExportAll_User_QueryParamIgnored(t *testing.T) {
	// query param id_perusahaan harus diabaikan untuk user — selalu pakai dari JWT
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-user").Return(fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf?id_perusahaan=p-lain", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// pastikan ExportByPerusahaanPDF dipanggil dengan p-user, bukan p-lain
	mockSvc.AssertExpectations(t)
}

func TestSEExportHandler_ExportAll_User_NoPerusahaan_Forbidden(t *testing.T) {
	h := NewSEExportHandler(new(mockSEExportService))
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf", nil)
	req = withExportUserCtxNoPerusahaan(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSEExportHandler_ExportAll_MethodNotAllowed(t *testing.T) {
	h := NewSEExportHandler(new(mockSEExportService))
	req := httptest.NewRequest(http.MethodPost, "/api/se/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// SE EXPORT HANDLER — ExportByID
// ════════════════════════════════════════════════════════════════════════════

func TestSEExportHandler_ExportByID_Admin_Success(t *testing.T) {
	se := &dto.SEResponse{ID: "se-1", IDPerusahaan: "p-abc"}
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByIDPDF", "se-1").Return(se, fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/se-1/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "se-1")
	mockSvc.AssertExpectations(t)
}

func TestSEExportHandler_ExportByID_User_OwnData_Success(t *testing.T) {
	se := &dto.SEResponse{ID: "se-1", IDPerusahaan: "p-user"}
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByIDPDF", "se-1").Return(se, fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/se-1/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSEExportHandler_ExportByID_User_OtherPerusahaan_Forbidden(t *testing.T) {
	se := &dto.SEResponse{ID: "se-1", IDPerusahaan: "p-lain"}
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByIDPDF", "se-1").Return(se, fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/se-1/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSEExportHandler_ExportByID_NotFound(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByIDPDF", "not-exist").Return(nil, nil, assert.AnError)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/not-exist/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	// assert.AnError message != "data tidak ditemukan" → 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSEExportHandler_ExportByID_DataTidakDitemukan(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByIDPDF", "not-exist").Return(
		(*dto.SEResponse)(nil),
		([]byte)(nil),
		errDataTidakDitemukan,
	)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/not-exist/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// CSIRT EXPORT HANDLER — ExportAll
// ════════════════════════════════════════════════════════════════════════════

func TestCsirtExportHandler_ExportAll_Admin_Success(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportAllPDF").Return(fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "laporan-csirt.pdf")
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportHandler_ExportAll_Admin_FilterPerusahaan(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-abc").Return(fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf?id_perusahaan=p-abc", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportHandler_ExportAll_Admin_ServiceError(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportAllPDF").Return(nil, assert.AnError)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCsirtExportHandler_ExportAll_User_Success(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-user").Return(fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportHandler_ExportAll_User_QueryParamIgnored(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-user").Return(fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf?id_perusahaan=p-lain", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportHandler_ExportAll_User_NoPerusahaan_Forbidden(t *testing.T) {
	h := NewCsirtExportHandler(new(mockCsirtExportService))
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf", nil)
	req = withExportUserCtxNoPerusahaan(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCsirtExportHandler_ExportAll_MethodNotAllowed(t *testing.T) {
	h := NewCsirtExportHandler(new(mockCsirtExportService))
	req := httptest.NewRequest(http.MethodPost, "/api/csirt/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// CSIRT EXPORT HANDLER — ExportByID
// ════════════════════════════════════════════════════════════════════════════

func TestCsirtExportHandler_ExportByID_Admin_Success(t *testing.T) {
	csirt := &dto.CsirtResponse{
		ID:         "csirt-1",
		Perusahaan: dto.PerusahaanResponse{ID: "p-abc"},
	}
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByIDPDF", "csirt-1").Return(csirt, fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/csirt-1/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "csirt-1")
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportHandler_ExportByID_User_OwnData_Success(t *testing.T) {
	csirt := &dto.CsirtResponse{
		ID:         "csirt-1",
		Perusahaan: dto.PerusahaanResponse{ID: "p-user"},
	}
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByIDPDF", "csirt-1").Return(csirt, fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/csirt-1/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCsirtExportHandler_ExportByID_User_OtherPerusahaan_Forbidden(t *testing.T) {
	csirt := &dto.CsirtResponse{
		ID:         "csirt-1",
		Perusahaan: dto.PerusahaanResponse{ID: "p-lain"},
	}
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByIDPDF", "csirt-1").Return(csirt, fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/csirt-1/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCsirtExportHandler_ExportByID_DataTidakDitemukan(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByIDPDF", "not-exist").Return(
		(*dto.CsirtResponse)(nil),
		([]byte)(nil),
		errDataTidakDitemukan,
	)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/not-exist/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCsirtExportHandler_ExportByID_ServiceError(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByIDPDF", "csirt-1").Return(nil, nil, assert.AnError)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/csirt-1/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// SE EXPORT HANDLER — cabang ServeHTTP yang belum dicakup
// ════════════════════════════════════════════════════════════════════════════
func TestSEExportHandler_ExportByID_MethodNotAllowed(t *testing.T) {
	h := NewSEExportHandler(new(mockSEExportService))

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/se/se-1/export-pdf", nil)
			req = withExportAdminCtx(req)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}

// Path tidak dikenali → 404 "Route tidak ditemukan"
func TestSEExportHandler_ServeHTTP_RouteNotFound(t *testing.T) {
	h := NewSEExportHandler(new(mockSEExportService))

	req := httptest.NewRequest(http.MethodGet, "/api/se/unknown-path", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Route tidak ditemukan")
}

// ID kosong pada path ByID (/api/se//export-pdf) → 400 "ID tidak valid"
func TestSEExportHandler_ExportByID_EmptyID_Returns400(t *testing.T) {
	h := NewSEExportHandler(new(mockSEExportService))

	// Simulasikan path yang menghasilkan id="" setelah trim prefix "/api/se/" dan suffix "/export-pdf"
	req := httptest.NewRequest(http.MethodGet, "/api/se//export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID tidak valid")
}

// User - ExportAll service error
func TestSEExportHandler_ExportAll_User_ServiceError(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-user").Return(nil, assert.AnError)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// servePDF — verifikasi Content-Length dan Content-Disposition header dengan nama file benar
func TestSEExportHandler_ExportByID_ResponseHeaders(t *testing.T) {
	se := &dto.SEResponse{ID: "se-abc", IDPerusahaan: "p-1"}
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByIDPDF", "se-abc").Return(se, fakePDF, nil)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/se-abc/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "laporan-se-se-abc.pdf")
	assert.Equal(t, len(fakePDF), w.Body.Len())
}

// Admin - ExportAll filter perusahaan service error
func TestSEExportHandler_ExportAll_Admin_FilterPerusahaan_ServiceError(t *testing.T) {
	mockSvc := new(mockSEExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-fail").Return(nil, assert.AnError)

	h := NewSEExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/se/export-pdf?id_perusahaan=p-fail", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ════════════════════════════════════════════════════════════════════════════
// CSIRT EXPORT HANDLER — cabang ServeHTTP yang belum dicakup
// ════════════════════════════════════════════════════════════════════════════
func TestCsirtExportHandler_ExportByID_MethodNotAllowed(t *testing.T) {
	h := NewCsirtExportHandler(new(mockCsirtExportService))

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/csirt/csirt-1/export-pdf", nil)
			req = withExportAdminCtx(req)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}

// Path tidak dikenali → 404 "Route tidak ditemukan"
func TestCsirtExportHandler_ServeHTTP_RouteNotFound(t *testing.T) {
	h := NewCsirtExportHandler(new(mockCsirtExportService))

	req := httptest.NewRequest(http.MethodGet, "/api/csirt/unknown-path", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Route tidak ditemukan")
}

// ID kosong pada path ByID → 400 "ID tidak valid"
func TestCsirtExportHandler_ExportByID_EmptyID_Returns400(t *testing.T) {
	h := NewCsirtExportHandler(new(mockCsirtExportService))

	req := httptest.NewRequest(http.MethodGet, "/api/csirt//export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID tidak valid")
}

// User - ExportAll service error
func TestCsirtExportHandler_ExportAll_User_ServiceError(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-user").Return(nil, assert.AnError)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf", nil)
	req = withExportUserCtx(req, "p-user")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// servePDF — verifikasi Content-Length dan Content-Disposition dengan nama file benar
func TestCsirtExportHandler_ExportByID_ResponseHeaders(t *testing.T) {
	csirt := &dto.CsirtResponse{
		ID:         "csirt-xyz",
		Perusahaan: dto.PerusahaanResponse{ID: "p-1"},
	}
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByIDPDF", "csirt-xyz").Return(csirt, fakePDF, nil)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/csirt-xyz/export-pdf", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "laporan-csirt-csirt-xyz.pdf")
	assert.Equal(t, len(fakePDF), w.Body.Len())
}

// Admin - ExportAll filter perusahaan service error
func TestCsirtExportHandler_ExportAll_Admin_FilterPerusahaan_ServiceError(t *testing.T) {
	mockSvc := new(mockCsirtExportService)
	mockSvc.On("ExportByPerusahaanPDF", "p-fail").Return(nil, assert.AnError)

	h := NewCsirtExportHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/csirt/export-pdf?id_perusahaan=p-fail", nil)
	req = withExportAdminCtx(req)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ════════════════════════════════════════════════════════════════════════════
// HELPERS untuk error "data tidak ditemukan"
// ════════════════════════════════════════════════════════════════════════════

type customErr struct{ msg string }

func (e *customErr) Error() string { return e.msg }

var errDataTidakDitemukan = &customErr{"data tidak ditemukan"}