package utils

import (
	"bytes"
	"testing"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================================================
// HELPERS: builder data dummy
// ================================================================

func makeSEResponse(namaSE string) dto.SEResponse {
	perusahaan := &dto.PerusahaanMiniResponse{
		ID:             "perusahaan-1",
		NamaPerusahaan: "PT Contoh Teknologi",
	}
	subSektor := &dto.SubSektorMiniResponse{
		ID:            "sub-1",
		NamaSubSektor: "Teknologi Informasi",
		NamaSektor:    "Informatika",
	}
	csirt := &dto.CsirtMiniResponse{
		ID:        "csirt-1",
		NamaCsirt: "CSIRT Contoh",
	}
	return dto.SEResponse{
		ID:                              "se-1",
		IDPerusahaan:                    "perusahaan-1",
		IDSubSektor:                     "sub-1",
		IDCsirt:                         "csirt-1",
		NamaSE:                          namaSE,
		IpSE:                            "192.168.1.1",
		AsNumberSE:                      "AS12345",
		PengelolaSE:                     "Tim IT",
		FiturSE:                         "Autentikasi, Enkripsi",
		NilaiInvestasi:                  "A",
		AnggaranOperasional:             "B",
		KepatuhanPeraturan:              "C",
		TeknikKriptografi:               "A",
		JumlahPengguna:                  "B",
		DataPribadi:                     "C",
		KlasifikasiData:                 "A",
		KekritisanProses:                "B",
		DampakKegagalan:                 "C",
		PotensiKerugiandanDampakNegatif: "A",
		TotalBobot:                      25,
		KategoriSE:                      "Tinggi",
		CreatedAt:                       "2024-01-01",
		UpdatedAt:                       "2024-06-01",
		Perusahaan:                      perusahaan,
		SubSektor:                       subSektor,
		Csirt:                           csirt,
	}
}

func makeCsirtResponse(namaCsirt string) dto.CsirtResponse {
	telepon := "+62812345678"
	return dto.CsirtResponse{
		ID:                     "csirt-1",
		NamaCsirt:              namaCsirt,
		WebCsirt:               "https://csirt.contoh.id",
		TeleponCsirt:           &telepon,
		PhotoCsirt:             "photo.jpg",
		FileRFC2350:            "rfc2350.pdf",
		FilePublicKeyPGP:       "key.pgp",
		FileStr:                "str.pdf",
		TanggalRegistrasi:      "2024-01-01",
		TanggalKadaluarsa:      "2025-01-01",
		TanggalRegistrasiUlang: "2025-01-15",
		Perusahaan: dto.PerusahaanResponse{
			ID:             "perusahaan-1",
			NamaPerusahaan: "PT Contoh Teknologi",
		},
	}
}

// isPDFBytes memeriksa bahwa bytes yang diberikan merupakan PDF valid
// dengan mengecek magic bytes "%PDF" di awal file.
func isPDFBytes(b []byte) bool {
	return len(b) > 4 && bytes.HasPrefix(b, []byte("%PDF"))
}

// ================================================================
// TEST: toSafe — pure function, tidak butuh fpdf
// ================================================================

func TestToSafe_DashesConversion(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"\u2013", "-"},   // en dash → -
		{"\u2014", "--"},  // em dash → --
		{"\u2012", "-"},   // figure dash → -
		{"\u2015", "--"},  // horizontal bar → --
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, toSafe(c.input), "input: %q", c.input)
	}
}

func TestToSafe_QuotesConversion(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"\u2018", "'"},  // left single quote
		{"\u2019", "'"},  // right single quote
		{"\u201C", "\""}, // left double quote
		{"\u201D", "\""}, // right double quote
		{"\u201A", ","},  // single low-9
		{"\u201E", "\""}, // double low-9
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, toSafe(c.input), "input: %q", c.input)
	}
}

func TestToSafe_SpacesConversion(t *testing.T) {
	assert.Equal(t, " ", toSafe("\u00A0"))  // non-breaking space
	assert.Equal(t, " ", toSafe("\u2009"))  // thin space
	assert.Equal(t, "", toSafe("\u200B"))   // zero-width space dihapus
}

func TestToSafe_SymbolsConversion(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"\u2026", "..."},   // ellipsis
		{"\u2022", "*"},     // bullet
		{"\u00B7", "*"},     // middle dot
		{"\u00A9", "(c)"},   // copyright
		{"\u00AE", "(R)"},   // registered
		{"\u2122", "(TM)"},  // trademark
		{"\u00B0", " deg"},  // degree
		{"\u00D7", "x"},     // multiplication
		{"\u00F7", "/"},     // division
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, toSafe(c.input), "input: %q", c.input)
	}
}

func TestToSafe_ArrowsConversion(t *testing.T) {
	assert.Equal(t, "->", toSafe("\u2192"))   // →
	assert.Equal(t, "<-", toSafe("\u2190"))   // ←
	assert.Equal(t, "<->", toSafe("\u2194"))  // ↔
}

func TestToSafe_PlainASCIIUnchanged(t *testing.T) {
	inputs := []string{
		"Hello World",
		"Nama SE: Server-01",
		"PT Contoh (c) 2024",
		"",
		"123-456",
	}
	for _, s := range inputs {
		assert.Equal(t, s, toSafe(s), "plain ASCII harus tidak berubah: %q", s)
	}
}

func TestToSafe_MixedString(t *testing.T) {
	input := "Laporan \u2013 Q1 2024\u2026"
	expected := "Laporan - Q1 2024..."
	assert.Equal(t, expected, toSafe(input))
}

func TestToSafe_MultipleCharsSameString(t *testing.T) {
	// Beberapa karakter sekaligus dalam satu string
	input := "\u2018hello\u2019 \u2013 \u2014 world\u2026"
	expected := "'hello' - -- world..."
	assert.Equal(t, expected, toSafe(input))
}

// ================================================================
// TEST: GenerateSEPDF
// ================================================================

func TestGenerateSEPDF_ReturnsPDFBytes(t *testing.T) {
	data := []dto.SEResponse{makeSEResponse("Server Utama")}
	result, err := GenerateSEPDF(data, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result), "output harus berupa PDF valid")
	assert.Greater(t, len(result), 0)
}

func TestGenerateSEPDF_WithNamaFilter(t *testing.T) {
	data := []dto.SEResponse{makeSEResponse("Server Backup")}
	result, err := GenerateSEPDF(data, "PT Maju Bersama")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateSEPDF_WithoutNamaFilter(t *testing.T) {
	data := []dto.SEResponse{makeSEResponse("Server Produksi")}
	result, err := GenerateSEPDF(data, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateSEPDF_MultipleRecords(t *testing.T) {
	data := []dto.SEResponse{
		makeSEResponse("Server A"),
		makeSEResponse("Server B"),
		makeSEResponse("Server C"),
	}
	result, err := GenerateSEPDF(data, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateSEPDF_EmptyData(t *testing.T) {
	// PDF tetap harus berhasil dibuat meski data kosong
	result, err := GenerateSEPDF([]dto.SEResponse{}, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateSEPDF_NilNestedStructs(t *testing.T) {
	// Perusahaan, SubSektor, Csirt nil → harus fallback ke ID, tidak panic
	se := makeSEResponse("Server Tanpa Relasi")
	se.Perusahaan = nil
	se.SubSektor = nil
	se.Csirt = nil

	result, err := GenerateSEPDF([]dto.SEResponse{se}, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateSEPDF_SpecialCharsInName(t *testing.T) {
	// Nama dengan karakter Unicode yang perlu di-sanitize oleh toSafe
	se := makeSEResponse("Server \u2013 Utama \u2026")
	result, err := GenerateSEPDF([]dto.SEResponse{se}, "PT Teknologi \u00AE")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

// ================================================================
// TEST: GenerateSEByIDPDF
// ================================================================

func TestGenerateSEByIDPDF_ReturnsPDFBytes(t *testing.T) {
	se := makeSEResponse("Server Detail")
	result, err := GenerateSEByIDPDF(&se)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
	assert.Greater(t, len(result), 0)
}

func TestGenerateSEByIDPDF_NilNestedStructs(t *testing.T) {
	se := makeSEResponse("Server Single")
	se.Perusahaan = nil
	se.SubSektor = nil
	se.Csirt = nil

	result, err := GenerateSEByIDPDF(&se)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateSEByIDPDF_SpecialCharsInName(t *testing.T) {
	se := makeSEResponse("Server \u2014 Detail \u2122")
	result, err := GenerateSEByIDPDF(&se)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateSEByIDPDF_EmptyOptionalFields(t *testing.T) {
	se := makeSEResponse("Server Minimal")
	se.FiturSE = ""
	se.IDCsirt = ""
	se.IDSubSektor = ""

	result, err := GenerateSEByIDPDF(&se)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

// ================================================================
// TEST: GenerateCsirtPDF
// ================================================================

func TestGenerateCsirtPDF_ReturnsPDFBytes(t *testing.T) {
	data := []dto.CsirtResponse{makeCsirtResponse("CSIRT Nasional")}
	result, err := GenerateCsirtPDF(data, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
	assert.Greater(t, len(result), 0)
}

func TestGenerateCsirtPDF_WithNamaFilter(t *testing.T) {
	data := []dto.CsirtResponse{makeCsirtResponse("CSIRT Regional")}
	result, err := GenerateCsirtPDF(data, "PT Filter Perusahaan")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtPDF_WithoutNamaFilter(t *testing.T) {
	data := []dto.CsirtResponse{makeCsirtResponse("CSIRT Lokal")}
	result, err := GenerateCsirtPDF(data, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtPDF_MultipleRecords(t *testing.T) {
	data := []dto.CsirtResponse{
		makeCsirtResponse("CSIRT A"),
		makeCsirtResponse("CSIRT B"),
		makeCsirtResponse("CSIRT C"),
	}
	result, err := GenerateCsirtPDF(data, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtPDF_EmptyData(t *testing.T) {
	result, err := GenerateCsirtPDF([]dto.CsirtResponse{}, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtPDF_NilOptionalFields(t *testing.T) {
	// TeleponCsirt nil, FileRFC, FilePGP, FileStr kosong
	csirt := makeCsirtResponse("CSIRT Minimal")
	csirt.TeleponCsirt = nil
	csirt.FileRFC2350 = ""
	csirt.FilePublicKeyPGP = ""
	csirt.FileStr = ""
	csirt.PhotoCsirt = ""
	csirt.TanggalRegistrasi = ""
	csirt.TanggalKadaluarsa = ""
	csirt.TanggalRegistrasiUlang = ""

	result, err := GenerateCsirtPDF([]dto.CsirtResponse{csirt}, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtPDF_PerusahaanNameFallbackToID(t *testing.T) {
	// NamaPerusahaan kosong → harus fallback ke ID di printCsirtBlock
	csirt := makeCsirtResponse("CSIRT Fallback")
	csirt.Perusahaan = dto.PerusahaanResponse{
		ID:             "fallback-id",
		NamaPerusahaan: "", // kosong, harus fallback ke ID
	}

	result, err := GenerateCsirtPDF([]dto.CsirtResponse{csirt}, "")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtPDF_SpecialCharsInName(t *testing.T) {
	csirt := makeCsirtResponse("CSIRT \u2013 Pusat \u2026")
	result, err := GenerateCsirtPDF([]dto.CsirtResponse{csirt}, "PT Test \u00AE")
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

// ================================================================
// TEST: GenerateCsirtByIDPDF
// ================================================================

func TestGenerateCsirtByIDPDF_ReturnsPDFBytes(t *testing.T) {
	csirt := makeCsirtResponse("CSIRT Detail")
	result, err := GenerateCsirtByIDPDF(&csirt)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
	assert.Greater(t, len(result), 0)
}

func TestGenerateCsirtByIDPDF_NilOptionalFields(t *testing.T) {
	csirt := makeCsirtResponse("CSIRT Single Minimal")
	csirt.TeleponCsirt = nil
	csirt.FileRFC2350 = ""
	csirt.FilePublicKeyPGP = ""
	csirt.FileStr = ""
	csirt.PhotoCsirt = ""

	result, err := GenerateCsirtByIDPDF(&csirt)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtByIDPDF_PerusahaanNameFallbackToID(t *testing.T) {
	csirt := makeCsirtResponse("CSIRT ID Fallback")
	csirt.Perusahaan = dto.PerusahaanResponse{
		ID:             "id-only",
		NamaPerusahaan: "",
	}

	result, err := GenerateCsirtByIDPDF(&csirt)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

func TestGenerateCsirtByIDPDF_SpecialCharsInName(t *testing.T) {
	csirt := makeCsirtResponse("CSIRT \u2014 Wilayah \u2122")
	result, err := GenerateCsirtByIDPDF(&csirt)
	require.NoError(t, err)
	assert.True(t, isPDFBytes(result))
}

// ================================================================
// TEST: konsistensi output — generate dua kali dengan data sama
//       harus menghasilkan PDF dengan ukuran yang sebanding
// ================================================================

func TestGenerateSEPDF_DeterministicSize(t *testing.T) {
	data := []dto.SEResponse{makeSEResponse("Server Konsisten")}

	result1, err1 := GenerateSEPDF(data, "")
	result2, err2 := GenerateSEPDF(data, "")

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Keduanya harus valid PDF
	assert.True(t, isPDFBytes(result1))
	assert.True(t, isPDFBytes(result2))

	// Ukuran tidak boleh berbeda jauh (toleransi 10% untuk timestamp)
	diff := len(result1) - len(result2)
	if diff < 0 {
		diff = -diff
	}
	maxAllowed := len(result1) / 10
	assert.LessOrEqual(t, diff, maxAllowed,
		"ukuran dua PDF dari data yang sama tidak boleh berbeda >10%%")
}

func TestGenerateCsirtPDF_DeterministicSize(t *testing.T) {
	data := []dto.CsirtResponse{makeCsirtResponse("CSIRT Konsisten")}

	result1, err1 := GenerateCsirtPDF(data, "")
	result2, err2 := GenerateCsirtPDF(data, "")

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.True(t, isPDFBytes(result1))
	assert.True(t, isPDFBytes(result2))

	diff := len(result1) - len(result2)
	if diff < 0 {
		diff = -diff
	}
	maxAllowed := len(result1) / 10
	assert.LessOrEqual(t, diff, maxAllowed,
		"ukuran dua PDF dari data yang sama tidak boleh berbeda >10%%")
}