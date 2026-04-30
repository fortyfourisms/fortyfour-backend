package utils

import (
	"bytes"
	"fmt"
	"ikas/internal/dto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

const (
	pageMargin = 15.0
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

// GenerateIkasPDF generates a professional report for IKAS assessment based on the BSSN template.
func GenerateIkasPDF(data *dto.IkasResponse) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.SetAutoPageBreak(true, 25)

	// Add Page
	pdf.AddPage()

	// 1. Logo BSSN at Top Center
	pwd, _ := os.Getwd()
	logoPath := filepath.Join(pwd, "internal", "assets", "bssn.png")
	pdf.ImageOptions(logoPath, (210-35)/2, 10, 35, 0, false, fpdf.ImageOptions{ReadDpi: true}, 0, "")
	pdf.Ln(32)

	// 2. Header Titles
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 6, "SURAT KETERANGAN HASIL PENILAIAN MANDIRI", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, "KEMATANGAN KEAMANAN SIBER ORGANISASI", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, "SEKTOR INDUSTRI", "", 1, "C", false, 0, "")
	pdf.Ln(6)

	// 3. Opening Text
	pdf.SetFont("Arial", "", 11)
	openingText := "Berdasarkan hasil Penilaian Mandiri Kematangan Keamanan Siber Organisasi Sektor Industri untuk Penyelenggara Sistem Elektronik (PSE) yang dilakukan oleh:"
	pdf.MultiCell(0, 5, toSafe(openingText), "", "L", false)
	pdf.Ln(5)

	// 4. PSE Details
	col1L, col1V := 35.0, 60.0
	col2L, col2V := 30.0, 55.0

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col1L, 5, "Nama PSE", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(col1V, 5, toSafe(data.Perusahaan.NamaPerusahaan), "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col2L, 5, "No. Telpon", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(col2V, 5, toSafe(data.Telepon), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col1L, 5, "Alamat", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(col1V, 5, toSafe(data.Perusahaan.Alamat), "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col2L, 5, "Email", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(col2V, 5, toSafe(data.Perusahaan.Email), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col1L, 5, "Sektor", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(col1V, 5, toSafe(data.Perusahaan.Sektor), "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col2L, 5, "Responden", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(col2V, 5, toSafe(data.Responden), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col1L, 5, "Tanggal", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	tanggalStr := data.Tanggal
	if t, err := time.Parse(time.RFC3339, tanggalStr); err == nil {
		tanggalStr = t.Format("02-01-2006")
	} else if t, err := time.Parse("2006-01-02", tanggalStr); err == nil {
		tanggalStr = t.Format("02-01-2006")
	}
	pdf.CellFormat(col1V, 5, toSafe(tanggalStr), "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(col2L, 5, "Jabatan", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(5, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(col2V, 5, toSafe(data.Jabatan), "", 1, "L", false, 0, "")

	pdf.Ln(8)

	// 5. Results Table Title
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, "HASIL PENILAIAN MANDIRI KEMATANGAN KEAMANAN SIBER MENGGUNAKAN INSTRUMEN", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "KEMATANGAN KEAMANAN SIBER (IKAS) v.1.2.1", "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// --- WATERMARK IKAS (Centered on the page) ---
	watermarkPath := filepath.Join(pwd, "internal", "assets", "ikas.png")
	wmW := 180.0
	// Centering 300mm on 210mm paper: (210 - 300) / 2 = -45
	pdf.ImageOptions(watermarkPath, -45, pdf.GetY()-15, wmW, 0, false, fpdf.ImageOptions{ReadDpi: true}, 0, "")

	// 6. Results Table
	tableW := 140.0
	startX := (210 - tableW) / 2
	pdf.SetX(startX)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(20, 7, "No.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(80, 7, "Domain", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Nilai", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	domains := []struct {
		No   string
		Name string
		Val  float64
	}{
		{"1.", "Domain Identifikasi", 0},
		{"2.", "Domain Proteksi", 0},
		{"3.", "Domain Deteksi", 0},
		{"4.", "Domain Penanggulangan dan Pemulihan", 0},
	}
	if data.Identifikasi != nil {
		domains[0].Val = data.Identifikasi.NilaiIdentifikasi
	}
	if data.Proteksi != nil {
		domains[1].Val = data.Proteksi.NilaiProteksi
	}
	if data.Deteksi != nil {
		domains[2].Val = data.Deteksi.NilaiDeteksi
	}
	if data.Gulih != nil {
		domains[3].Val = data.Gulih.NilaiGulih
	}

	for _, d := range domains {
		pdf.SetX(startX)
		pdf.CellFormat(20, 7, d.No, "1", 0, "C", false, 0, "")
		pdf.CellFormat(80, 7, " "+toSafe(d.Name), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", d.Val), "1", 1, "C", false, 0, "")
	}

	pdf.SetX(startX)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(100, 7, "Nilai Kematangan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", data.NilaiKematangan), "1", 1, "C", false, 0, "")
	pdf.SetX(startX)
	pdf.CellFormat(100, 7, "Kategori Tingkat Kematangan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 7, toSafe(data.KategoriKematanganKeamananSiber), "1", 1, "C", false, 0, "")

	// 7. Kriteria
	pdf.SetX(startX)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(140, 4, "Kriteria:", "LR", 1, "L", false, 0, "")
	pdf.SetX(startX)
	pdf.SetFont("Arial", "", 8)

	var kriteriaText string
	category := strings.ToLower(data.KategoriKematanganKeamananSiber)

	if strings.Contains(category, "level 1") {
		kriteriaText = "A. Kondisi penerapan keamanan siber dalam tahap implementasi awal\n" +
			"B. Penerapan keamanan siber belum memiliki prosedur yang terorganisir\n" +
			"C. Penerapan keamanan siber bersifat informal\n" +
			"D. Keamanan siber tidak dilakukan secara konsisten dan berkelanjutan\n" +
			"E. Dokumen manajemen risiko dan dokumen kontrol belum disusun"
	} else if strings.Contains(category, "level 2") {
		kriteriaText = "A. Kondisi penerapan keamanan siber dalam tahap implementasi yang berulang\n" +
			"B. Penerapan keamanan siber sudah memiliki prosedur yang terorganisir\n" +
			"C. Penerapan keamanan siber bersifat informal\n" +
			"D. Keamanan siber dilakukan secara berulang namun belum konsisten dan berkelanjutan\n" +
			"E. Dokumen manajemen risiko dan dokumen kontrol sudah disusun namun belum ditetapkan"
	} else if strings.Contains(category, "level 3") {
		kriteriaText = "A. Kondisi penerapan keamanan siber dalam tahap implementasi yang telah terdefinisi dengan baik\n" +
			"B. Penerapan keamanan siber terorganisir dengan jelas\n" +
			"C. Penerapan keamanan siber bersifat formal\n" +
			"D. Keamanan siber dilakukan secara berulang dan konsisten serta direviu secara berkala\n" +
			"E. Dokumen manajemen risiko dan dokumen kontrol sudah disusun dan sudah ditetapkan"
	} else if strings.Contains(category, "level 4") {
		kriteriaText = "A. Kondisi penerapan keamanan siber dalam tahap implementasi yang telah terkelola dengan baik\n" +
			"B. Penerapan keamanan siber terorganisir dengan baik namun belum dilakukan proses otomatisasi\n" +
			"C. Penerapan keamanan siber bersifat formal\n" +
			"D. Keamanan siber dilakukan secara berulang dan implementasi perbaikan dilakukan berkelanjutan\n" +
			"E. Dokumen manajemen risiko dan dokumen kontrol sudah disusun dan sudah ditetapkan"
	} else if strings.Contains(category, "level 5") {
		kriteriaText = "A. Kondisi penerapan keamanan siber telah diimplementasikan secara optimal\n" +
			"B. Penerapan keamanan siber telah terorganisir dengan baik dan telah dilakukan proses otomatisasi\n" +
			"C. Penerapan keamanan siber bersifat formal\n" +
			"D. Keamanan siber dilakukan secara berulang dan konsisten serta telah terintegrasi dan menjadi bagian budaya organisasi secara menyeluruh\n" +
			"E. Dokumen manajemen risiko dan dokumen kontrol sudah ditetapkan"
	} else {
		// Empty if unknown
		kriteriaText = ""
	}

	pdf.MultiCell(140, 4, toSafe(kriteriaText), "LRB", "L", false)

	pdf.Ln(8)

	// 8. Closing
	pdf.SetFont("Arial", "I", 9)
	keterangan := fmt.Sprintf("Keterangan:\nHasil diatas merupakan hasil penilaian mandiri IKAS %s dan masih memerlukan verifikasi dari Tim BSSN.", data.Perusahaan.NamaPerusahaan)
	pdf.MultiCell(0, 5, toSafe(keterangan), "", "L", false)
	pdf.Ln(4)

	pdf.SetFont("Arial", "", 10)
	closingText := fmt.Sprintf("Setelah mengetahui hasil penilaian mandiri IKAS sebagaimana tercantum pada tabel di atas yang dilakukan pada tanggal %s oleh %s, maka Badan Siber dan Sandi Negara dalam hal ini Direktorat Keamanan Siber dan Sandi Industri melalui Tim Asistensi IKAS %s, menerima hasil penilaian mandiri IKAS %s.", tanggalStr, data.Perusahaan.NamaPerusahaan, data.Perusahaan.NamaPerusahaan, data.Perusahaan.NamaPerusahaan)
	pdf.MultiCell(0, 5, toSafe(closingText), "", "L", false)

	pdf.Ln(10)

	// 9. Signature
	sigW := 55.0
	sigX := 210 - pageMargin - sigW

	pdf.SetFont("Arial", "B", 10)
	pdf.SetX(sigX)
	pdf.CellFormat(sigW, 5, "Ketua Tim Penilaian Kematangan", "", 1, "C", false, 0, "")
	pdf.SetX(sigX)
	pdf.CellFormat(sigW, 5, "Keamanan Siber Sektor Industri,", "", 1, "C", false, 0, "")

	// 10. Footer
	pdf.SetY(275)
	pdf.SetFont("Arial", "", 8)
	footerText := "Dokumen ini telah ditandatangani secara elektronik menggunakan sertifikat elektronik\nyang diterbitkan oleh Balai Besar Sertifikasi Elektronik (BSrE), Badan Siber dan Sandi Negara"
	pdf.MultiCell(0, 4, toSafe(footerText), "", "C", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF IKAS: %w", err)
	}
	return buf.Bytes(), nil
}
