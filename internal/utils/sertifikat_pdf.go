package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"fortyfour-backend/internal/models"

	"github.com/go-pdf/fpdf"
)

// GenerateSertifikatPDF membuat sertifikat PDF standar dan menyimpannya ke disk.
// Mengembalikan path file PDF yang berhasil dibuat.
func GenerateSertifikatPDF(sertifikat *models.Sertifikat) (string, error) {
	pdf := fpdf.New("L", "mm", "A4", "") // Landscape
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pageW, pageH := pdf.GetPageSize()
	marginX := 20.0

	// ── Border dekoratif ──────────────────────────────────────────────────
	pdf.SetDrawColor(0, 102, 204)
	pdf.SetLineWidth(2)
	pdf.Rect(10, 10, pageW-20, pageH-20, "D")

	pdf.SetDrawColor(0, 102, 204)
	pdf.SetLineWidth(0.5)
	pdf.Rect(15, 15, pageW-30, pageH-30, "D")

	// ── Header ────────────────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(100, 100, 100)
	pdf.SetXY(marginX, 25)
	pdf.CellFormat(pageW-2*marginX, 10, "LEARNING MANAGEMENT SYSTEM", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "B", 32)
	pdf.SetTextColor(0, 51, 102)
	pdf.SetXY(marginX, 42)
	pdf.CellFormat(pageW-2*marginX, 15, "SERTIFIKAT", "", 1, "C", false, 0, "")

	// ── Garis dekoratif ───────────────────────────────────────────────────
	lineY := 60.0
	pdf.SetDrawColor(0, 102, 204)
	pdf.SetLineWidth(1)
	pdf.Line(pageW/2-50, lineY, pageW/2+50, lineY)

	// ── Body ──────────────────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(marginX, 68)
	pdf.CellFormat(pageW-2*marginX, 8, "Diberikan kepada:", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(0, 51, 102)
	pdf.SetXY(marginX, 80)
	pdf.CellFormat(pageW-2*marginX, 14, toSafe(sertifikat.NamaPeserta), "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(marginX, 100)
	pdf.CellFormat(pageW-2*marginX, 8, "Atas keberhasilan menyelesaikan kelas:", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(marginX, 112)
	pdf.CellFormat(pageW-2*marginX, 10, toSafe(sertifikat.NamaKelas), "", 1, "C", false, 0, "")

	// ── Tanggal & Nomor ───────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.SetXY(marginX, 135)
	tanggal := sertifikat.TanggalTerbit.Format("02 January 2006")
	pdf.CellFormat(pageW-2*marginX, 7, fmt.Sprintf("Tanggal Terbit: %s", tanggal), "", 1, "C", false, 0, "")

	pdf.SetXY(marginX, 143)
	pdf.CellFormat(pageW-2*marginX, 7, fmt.Sprintf("No. %s", sertifikat.NomorSertifikat), "", 1, "C", false, 0, "")

	// ── Garis bawah dekoratif ─────────────────────────────────────────────
	lineY2 := 158.0
	pdf.SetDrawColor(0, 102, 204)
	pdf.SetLineWidth(1)
	pdf.Line(pageW/2-50, lineY2, pageW/2+50, lineY2)

	// ── Footer ────────────────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(marginX, pageH-25)
	pdf.CellFormat(pageW-2*marginX, 5, "Sertifikat ini diterbitkan secara otomatis oleh sistem LMS dan sah tanpa tanda tangan.", "", 1, "C", false, 0, "")

	// ── Simpan ke file ────────────────────────────────────────────────────
	dir := "uploads/sertifikat"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.pdf", sertifikat.ID)
	fullPath := filepath.Join(dir, filename)

	if err := pdf.OutputFileAndClose(fullPath); err != nil {
		return "", err
	}

	return fullPath, nil
}
