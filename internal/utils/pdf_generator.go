package utils

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"

	"github.com/go-pdf/fpdf"
)

const (
	pageMargin   = 15.0
	lineHeight   = 7.0
	colLabelW    = 70.0
	colValueW    = 110.0
	headerHeight = 10.0
)

// toSafe converts a UTF-8 string to a safe Latin-1 representation for fpdf.
// fpdf built-in fonts (Arial, Helvetica, etc.) only support Latin-1 / ISO-8859-1.
// Characters outside that range — like the en dash (–) stored as UTF-8 — are
// silently corrupted into â€" and similar garbage. This function replaces the
// most common offenders with plain ASCII equivalents before they reach fpdf.
func toSafe(s string) string {
	r := strings.NewReplacer(
		// Dashes
		"\u2013", "-", // en dash  –
		"\u2014", "--", // em dash  —
		"\u2012", "-", // figure dash
		"\u2015", "--", // horizontal bar

		// Quotes
		"\u2018", "'", // left single quotation mark  '
		"\u2019", "'", // right single quotation mark  '
		"\u201A", ",", // single low-9 quotation mark  ‚
		"\u201C", "\"", // left double quotation mark  "
		"\u201D", "\"", // right double quotation mark  "
		"\u201E", "\"", // double low-9 quotation mark  „

		// Spaces & separators
		"\u00A0", " ", // non-breaking space
		"\u2009", " ", // thin space
		"\u200B", "", // zero-width space

		// Symbols
		"\u2026", "...", // ellipsis  …
		"\u2022", "*", // bullet  •
		"\u00B7", "*", // middle dot  ·
		"\u00A9", "(c)", // copyright  ©
		"\u00AE", "(R)", // registered  ®
		"\u2122", "(TM)", // trademark  ™
		"\u00B0", " deg", // degree  °
		"\u00D7", "x", // multiplication sign  ×
		"\u00F7", "/", // division sign  ÷

		// Arrows
		"\u2192", "->", // →
		"\u2190", "<-", // ←
		"\u2194", "<->", // ↔
	)
	return r.Replace(s)
}

// ════════════════════════════════════════════════════════════════════════════
// SE
// ════════════════════════════════════════════════════════════════════════════

// GenerateSEPDF generates a PDF for a slice of SEResponse.
// If namaFilter is non-empty it is shown as a subtitle filter.
func GenerateSEPDF(data []dto.SEResponse, namaFilter string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.SetAutoPageBreak(true, pageMargin)

	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Laporan Sistem Elektronik (SE)", "", 1, "C", false, 0, "")

	if namaFilter != "" {
		pdf.SetFont("Arial", "I", 11)
		pdf.CellFormat(0, 7, toSafe("Perusahaan: "+namaFilter), "", 1, "C", false, 0, "")
	}

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6,
		fmt.Sprintf("Digenerate pada: %s", time.Now().Format("02 January 2006, 15:04:05 WIB")),
		"", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, fmt.Sprintf("Total data: %d", len(data)), "", 1, "C", false, 0, "")

	pdf.Ln(4)

	for i, se := range data {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetFillColor(41, 98, 255)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(0, headerHeight,
			toSafe(fmt.Sprintf("  SE #%d - %s", i+1, se.NamaSE)),
			"", 1, "L", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(1)

		printSEBlock(pdf, se)
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// GenerateSEByIDPDF generates a PDF for a single SEResponse.
func GenerateSEByIDPDF(se *dto.SEResponse) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.SetAutoPageBreak(true, pageMargin)

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Laporan Sistem Elektronik (SE)", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(0, 8, toSafe(se.NamaSE), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6,
		fmt.Sprintf("Digenerate pada: %s", time.Now().Format("02 January 2006, 15:04:05 WIB")),
		"", 1, "C", false, 0, "")
	pdf.Ln(4)

	printSEBlock(pdf, *se)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// ════════════════════════════════════════════════════════════════════════════
// CSIRT
// ════════════════════════════════════════════════════════════════════════════

// GenerateCsirtPDF generates a PDF for a slice of CsirtResponse.
// If namaFilter is non-empty it is shown as a subtitle.
func GenerateCsirtPDF(data []dto.CsirtResponse, namaFilter string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.SetAutoPageBreak(true, pageMargin)

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Laporan CSIRT", "", 1, "C", false, 0, "")

	if namaFilter != "" {
		pdf.SetFont("Arial", "I", 11)
		pdf.CellFormat(0, 7, toSafe("Perusahaan: "+namaFilter), "", 1, "C", false, 0, "")
	}

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6,
		fmt.Sprintf("Digenerate pada: %s", time.Now().Format("02 January 2006, 15:04:05 WIB")),
		"", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, fmt.Sprintf("Total data: %d", len(data)), "", 1, "C", false, 0, "")

	pdf.Ln(4)

	for i, csirt := range data {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetFillColor(41, 98, 255)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(0, headerHeight,
			toSafe(fmt.Sprintf("  CSIRT #%d - %s", i+1, csirt.NamaCsirt)),
			"", 1, "L", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(1)

		printCsirtBlock(pdf, csirt)
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// GenerateCsirtByIDPDF generates a PDF for a single CsirtResponse.
func GenerateCsirtByIDPDF(csirt *dto.CsirtResponse) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.SetAutoPageBreak(true, pageMargin)

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Laporan CSIRT", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(0, 8, toSafe(csirt.NamaCsirt), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6,
		fmt.Sprintf("Digenerate pada: %s", time.Now().Format("02 January 2006, 15:04:05 WIB")),
		"", 1, "C", false, 0, "")
	pdf.Ln(4)

	printCsirtBlock(pdf, *csirt)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// ── Internal helpers ──────────────────────────────────────────────────────────

func printSEBlock(pdf *fpdf.Fpdf, se dto.SEResponse) {
	printSectionTitle(pdf, "Informasi SE")

	namaPerusahaan := se.IDPerusahaan
	if se.Perusahaan != nil {
		namaPerusahaan = se.Perusahaan.NamaPerusahaan
	}
	// Gunakan " - " (ASCII) bukan " - " (UTF-8 en dash) agar tidak corrupt
	namaSubSektor := se.IDSubSektor
	if se.SubSektor != nil {
		namaSubSektor = fmt.Sprintf("%s - %s", se.SubSektor.NamaSektor, se.SubSektor.NamaSubSektor)
	}
	namaCsirt := se.IDCsirt
	if se.Csirt != nil {
		namaCsirt = se.Csirt.NamaCsirt
	}

	rows := [][]string{
		{"Perusahaan", namaPerusahaan},
		{"Sub Sektor", namaSubSektor},
		{"CSIRT", namaCsirt},
		{"Nama SE", se.NamaSE},
		{"IP SE", se.IpSE},
		{"AS Number SE", se.AsNumberSE},
		{"Pengelola SE", se.PengelolaSE},
		{"Fitur SE", se.FiturSE},
		{"Dibuat", se.CreatedAt},
		{"Diupdate", se.UpdatedAt},
	}
	printRows(pdf, rows)

	pdf.Ln(3)

	printSectionTitle(pdf, "Karakteristik Instansi")
	karakteristik := [][]string{
		{"Nilai Investasi", se.NilaiInvestasi},
		{"Anggaran Operasional", se.AnggaranOperasional},
		{"Kepatuhan Peraturan", se.KepatuhanPeraturan},
		{"Teknik Kriptografi", se.TeknikKriptografi},
		{"Jumlah Pengguna", se.JumlahPengguna},
		{"Data Pribadi", se.DataPribadi},
		{"Klasifikasi Data", se.KlasifikasiData},
		{"Kekritisan Proses", se.KekritisanProses},
		{"Dampak Kegagalan", se.DampakKegagalan},
		{"Potensi Kerugian & Dampak Negatif", se.PotensiKerugiandanDampakNegatif},
	}
	printRows(pdf, karakteristik)

	pdf.Ln(3)

	printSectionTitle(pdf, "Hasil Kategorisasi")
	hasil := [][]string{
		{"Total Bobot", fmt.Sprintf("%d", se.TotalBobot)},
		{"Kategori SE", se.KategoriSE},
	}
	printRows(pdf, hasil)
}

func printSectionTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(220, 230, 255)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(colLabelW+colValueW, lineHeight, "  "+title, "1", 1, "L", true, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func printRows(pdf *fpdf.Fpdf, rows [][]string) {
	for i, row := range rows {
		if i%2 == 0 {
			pdf.SetFillColor(245, 247, 255)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(colLabelW, lineHeight, toSafe("  "+row[0]), "1", 0, "L", true, 0, "")

		pdf.SetFont("Arial", "", 9)
		pdf.CellFormat(colValueW, lineHeight, toSafe("  "+row[1]), "1", 1, "L", true, 0, "")
	}
}

func statusAda(s string) string {
	if s != "" {
		return "Ada"
	}
	return "-"
}

func printCsirtBlock(pdf *fpdf.Fpdf, csirt dto.CsirtResponse) {
	printSectionTitle(pdf, "Informasi CSIRT")

	namaPerusahaan := csirt.Perusahaan.NamaPerusahaan
	if namaPerusahaan == "" {
		namaPerusahaan = csirt.Perusahaan.ID
	}

	email := "-"
	if csirt.EmailCsirt != nil && *csirt.EmailCsirt != "" {
		email = *csirt.EmailCsirt
	}

	telepon := "-"
	if csirt.TeleponCsirt != nil && *csirt.TeleponCsirt != "" {
		telepon = *csirt.TeleponCsirt
	}

	tglReg := "-"
	if csirt.TanggalRegistrasi != "" {
		tglReg = csirt.TanggalRegistrasi
	}

	tglKadaluarsa := "-"
	if csirt.TanggalKadaluarsa != "" {
		tglKadaluarsa = csirt.TanggalKadaluarsa
	}

	rows := [][]string{
		{"Perusahaan", namaPerusahaan},
		{"Nama CSIRT", csirt.NamaCsirt},
		{"Website CSIRT", csirt.WebCsirt},
		{"Email CSIRT", email},
		{"Telepon CSIRT", telepon},
		{"Photo CSIRT", statusAda(csirt.PhotoCsirt)},
		{"File RFC2350", statusAda(csirt.FileRFC2350)},
		{"File Public Key PGP", statusAda(csirt.FilePublicKeyPGP)},
		{"File STR (Sertifikat)", statusAda(csirt.FileStr)},
	}
	printRows(pdf, rows)

	pdf.Ln(3)

	printSectionTitle(pdf, "Tanggal Registrasi")
	tanggalRows := [][]string{
		{"Tanggal Registrasi", tglReg},
		{"Tanggal Kadaluarsa", tglKadaluarsa},
	}
	printRows(pdf, tanggalRows)
}
