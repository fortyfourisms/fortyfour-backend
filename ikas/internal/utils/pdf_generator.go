package utils

import (
	"bytes"
	"fmt"
	"ikas/internal/dto"
	"strings"
	"time"

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
func toSafe(s string) string {
	r := strings.NewReplacer(
		"\u2013", "-", "\u2014", "--", "\u2012", "-", "\u2015", "--",
		"\u2018", "'", "\u2019", "'", "\u201A", ",", "\u201C", "\"", "\u201D", "\"", "\u201E", "\"",
		"\u00A0", " ", "\u2009", " ", "\u200B", "", "\u2026", "...", "\u2022", "*", "\u00B7", "*",
		"\u00A9", "(c)", "\u00AE", "(R)", "\u2122", "(TM)", "\u00B0", " deg", "\u00D7", "x", "\u00F7", "/",
		"\u2192", "->", "\u2190", "<-", "\u2194", "<->",
	)
	return r.Replace(s)
}

// GenerateIkasPDF generates a professional report for IKAS assessment.
func GenerateIkasPDF(data *dto.IkasResponse) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.SetAutoPageBreak(true, pageMargin)

	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Laporan Indeks Kematangan Keamanan Siber (IKAS)", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(0, 8, toSafe(data.Perusahaan.NamaPerusahaan), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6,
		fmt.Sprintf("Digenerate pada: %s", time.Now().Format("02 January 2006, 15:04:05 WIB")),
		"", 1, "C", false, 0, "")
	pdf.Ln(6)

	// Summary Section
	printSectionTitle(pdf, "Ringkasan Hasil Kematangan")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 245, 255)

	pdf.CellFormat(colLabelW, 12, "  Skor Total", "1", 0, "L", true, 0, "")
	pdf.SetTextColor(41, 98, 255)
	pdf.CellFormat(colValueW, 12, fmt.Sprintf("  %.2f", data.NilaiKematangan), "1", 1, "L", true, 0, "")
	pdf.SetTextColor(0, 0, 0)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(colLabelW, 10, "  Kategori", "1", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(colValueW, 10, "  "+toSafe(data.KategoriKematanganKeamananSiber), "1", 1, "L", false, 0, "")

	pdf.Ln(6)

	// Assessment Info
	printSectionTitle(pdf, "Informasi Asesmen")
	infoRows := [][]string{
		{"Responden", data.Responden},
		{"Jabatan", data.Jabatan},
		{"Telepon", data.Telepon},
		{"Tanggal Data", data.Tanggal},
		{"Target Skor", fmt.Sprintf("%.2f", data.TargetNilai)},
		{"Status Validasi", map[bool]string{true: "Terverifikasi", false: "Belum Diverifikasi"}[data.IsValidated]},
	}
	printRows(pdf, infoRows)

	pdf.Ln(8)

	// Domain Breakdown
	printSectionTitle(pdf, "Rincian Skor Per Domain")

	// Table Header for domains
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(230, 235, 245)
	pdf.CellFormat(80, 8, "  Domain", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "  Skor", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 8, "  Level Kematangan", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)

	domains := []struct {
		Nama  string
		Skor  float64
		Level string
	}{
		{"Identifikasi (Identify)", 0, "-"},
		{"Proteksi (Protect)", 0, "-"},
		{"Deteksi (Detect)", 0, "-"},
		{"Penanggulangan & Pemulihan (Respond & Recover)", 0, "-"},
	}

	if data.Identifikasi != nil {
		domains[0].Skor = data.Identifikasi.NilaiIdentifikasi
		domains[0].Level = data.Identifikasi.KategoriTingkatKematanganDomain
	}
	if data.Proteksi != nil {
		domains[1].Skor = data.Proteksi.NilaiProteksi
		domains[1].Level = data.Proteksi.KategoriTingkatKematanganDomain
	}
	if data.Deteksi != nil {
		domains[2].Skor = data.Deteksi.NilaiDeteksi
		domains[2].Level = data.Deteksi.KategoriTingkatKematanganDomain
	}
	if data.Gulih != nil {
		domains[3].Skor = data.Gulih.NilaiGulih
		domains[3].Level = data.Gulih.KategoriTingkatKematanganDomain
	}

	for _, d := range domains {
		pdf.CellFormat(80, 8, "  "+toSafe(d.Nama), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("  %.2f", d.Skor), "1", 0, "C", false, 0, "")
		pdf.CellFormat(70, 8, "  "+toSafe(d.Level), "1", 1, "L", false, 0, "")
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 100, 100)
	pdf.MultiCell(0, 4, toSafe("Catatan: Laporan ini dihasilkan secara otomatis oleh sistem. Nilai kematangan dihitung berdasarkan metodologi standar Indeks Kematangan Keamanan Siber (IKAS)."), "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF IKAS: %w", err)
	}
	return buf.Bytes(), nil
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
