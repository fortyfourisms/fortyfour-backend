package services

import (
	"errors"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"
)

// SEExportServiceInterface defines the contract for SE export operations.
type SEExportServiceInterface interface {
	ExportAllPDF() ([]byte, error)
	ExportByPerusahaanPDF(idPerusahaan string) ([]byte, error)
	ExportByIDPDF(id string) (*dto.SEResponse, []byte, error)
}

// SEExportService handles PDF export logic for SE data.
type SEExportService struct {
	seService SEService
}

func NewSEExportService(seService SEService) *SEExportService {
	return &SEExportService{seService: seService}
}

// ExportAllPDF returns a PDF of all SE data.
// Admin only — no filter applied.
func (s *SEExportService) ExportAllPDF() ([]byte, error) {
	data, err := s.seService.GetAll()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("tidak ada data SE untuk diexport")
	}
	return utils.GenerateSEPDF(data, "")
}

// ExportByPerusahaanPDF returns a PDF of SE data filtered by perusahaan.
// Used for both: admin filtering a specific company, and regular users.
func (s *SEExportService) ExportByPerusahaanPDF(idPerusahaan string) ([]byte, error) {
	if strings.TrimSpace(idPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}

	data, err := s.seService.GetByPerusahaan(idPerusahaan)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("tidak ada data SE untuk perusahaan ini")
	}

	// Use company name from first record if available
	namaFilter := idPerusahaan
	if len(data) > 0 && data[0].Perusahaan != nil {
		namaFilter = data[0].Perusahaan.NamaPerusahaan
	}

	return utils.GenerateSEPDF(data, namaFilter)
}

// ExportByIDPDF returns a PDF for a single SE entry.
func (s *SEExportService) ExportByIDPDF(id string) (*dto.SEResponse, []byte, error) {
	if strings.TrimSpace(id) == "" {
		return nil, nil, errors.New("id wajib diisi")
	}

	se, err := s.seService.GetByID(id)
	if err != nil {
		return nil, nil, errors.New("data tidak ditemukan")
	}

	pdfBytes, err := utils.GenerateSEByIDPDF(se)
	if err != nil {
		return nil, nil, err
	}

	return se, pdfBytes, nil
}

// ════════════════════════════════════════════════════════════════════════════
// CSIRT Export
// ════════════════════════════════════════════════════════════════════════════

// CsirtExportServiceInterface defines the contract for CSIRT export operations.
type CsirtExportServiceInterface interface {
	ExportAllPDF() ([]byte, error)
	ExportByPerusahaanPDF(idPerusahaan string) ([]byte, error)
	ExportByIDPDF(id string) (*dto.CsirtResponse, []byte, error)
}

// CsirtExportService handles PDF export logic for CSIRT data.
type CsirtExportService struct {
	csirtService CsirtServiceInterface
}

func NewCsirtExportService(csirtService CsirtServiceInterface) *CsirtExportService {
	return &CsirtExportService{csirtService: csirtService}
}

// ExportAllPDF returns a PDF of all CSIRT data.
// Admin only — no filter applied.
func (s *CsirtExportService) ExportAllPDF() ([]byte, error) {
	data, err := s.csirtService.GetAll()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("tidak ada data CSIRT untuk diexport")
	}
	return utils.GenerateCsirtPDF(data, "")
}

// ExportByPerusahaanPDF returns a PDF of CSIRT data filtered by perusahaan.
func (s *CsirtExportService) ExportByPerusahaanPDF(idPerusahaan string) ([]byte, error) {
	if strings.TrimSpace(idPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}

	data, err := s.csirtService.GetByPerusahaan(idPerusahaan)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("tidak ada data CSIRT untuk perusahaan ini")
	}

	namaFilter := idPerusahaan
	if len(data) > 0 && data[0].Perusahaan.NamaPerusahaan != "" {
		namaFilter = data[0].Perusahaan.NamaPerusahaan
	}

	return utils.GenerateCsirtPDF(data, namaFilter)
}

// ExportByIDPDF returns a PDF for a single CSIRT entry.
func (s *CsirtExportService) ExportByIDPDF(id string) (*dto.CsirtResponse, []byte, error) {
	if strings.TrimSpace(id) == "" {
		return nil, nil, errors.New("id wajib diisi")
	}

	csirt, err := s.csirtService.GetByID(id)
	if err != nil {
		return nil, nil, errors.New("data tidak ditemukan")
	}

	pdfBytes, err := utils.GenerateCsirtByIDPDF(csirt)
	if err != nil {
		return nil, nil, err
	}

	return csirt, pdfBytes, nil
}
