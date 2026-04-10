package models

import (
	"strings"
	"time"
)

// STRExpiryWarnDays adalah batas peringatan sebelum STR kadaluarsa (180 hari).
// Ubah nilai ini jika ingin mengubah berapa hari sebelum kadaluarsa notifikasi dikirim.
const STRExpiryWarnDays = 180

type Csirt struct {
	ID                     string  `json:"id"`
	IdPerusahaan           string  `json:"id_perusahaan"`
	NamaCsirt              string  `json:"nama_csirt"`
	WebCsirt               string  `json:"web_csirt"`
	TeleponCsirt           *string `json:"telepon_csirt"`
	PhotoCsirt             *string `json:"photo_csirt"`
	FileRFC2350            *string `json:"file_rfc2350"`
	FilePublicKeyPGP       *string `json:"file_public_key_pgp"`
	FileStr                *string `json:"file_str"`
	TanggalRegistrasi      *string `json:"tanggal_registrasi"`
	TanggalKadaluarsa      *string `json:"tanggal_kadaluarsa"`
	TanggalRegistrasiUlang *string `json:"tanggal_registrasi_ulang"`
}

// parseDate memparsing string tanggal dari database ke time.Time.
// Mendukung berbagai format yang mungkin dikembalikan oleh MySQL driver
// dengan parseTime=true maupun tanpa parseTime.
func parseDate(s *string) (time.Time, bool) {
	if s == nil || *s == "" {
		return time.Time{}, false
	}

	val := strings.TrimSpace(*s)
	if val == "" {
		return time.Time{}, false
	}

	// Coba beberapa format yang mungkin dari database
	// Urutan: dari yang paling spesifik ke yang paling umum
	formats := []string{
		"2006-01-02",                              // YYYY-MM-DD (standar DATE MySQL tanpa parseTime)
		"2006-01-02 15:04:05 -0700 MST",           // time.Time.String() — output parseTime=true dengan timezone
		"2006-01-02 15:04:05 -0700 -07",           // variasi time.Time.String() tanpa named timezone
		time.RFC3339,                              // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05Z",                   // ISO 8601 dengan Z
		"2006-01-02T15:04:05-07:00",              // ISO 8601 dengan offset
		"2006-01-02 15:04:05",                    // datetime tanpa timezone
		"2006-01-02 15:04:05.999999999 -0700 MST", // time.Time.String() dengan nanosecond
	}

	for _, format := range formats {
		t, err := time.Parse(format, val)
		if err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// ──────────────────────────────────────────────────────────
//  Tanggal Kadaluarsa helpers
// ──────────────────────────────────────────────────────────

// IsSTRExpired mengecek apakah tanggal_kadaluarsa sudah lewat
func (c *Csirt) IsSTRExpired() bool {
	t, ok := parseDate(c.TanggalKadaluarsa)
	if !ok {
		return false
	}
	return time.Now().After(t)
}

// IsSTRExpiringSoon mengecek apakah tanggal_kadaluarsa akan jatuh tempo
// dalam STRExpiryWarnDays hari ke depan (belum lewat, tapi sudah dekat)
func (c *Csirt) IsSTRExpiringSoon() bool {
	t, ok := parseDate(c.TanggalKadaluarsa)
	if !ok {
		return false
	}
	now := time.Now()
	warnDate := t.AddDate(0, 0, -STRExpiryWarnDays)
	return now.After(warnDate) && now.Before(t)
}

// DaysUntilSTRExpiry mengembalikan sisa hari sebelum STR kadaluarsa.
// Mengembalikan 0 jika sudah expired atau tanggal tidak valid.
func (c *Csirt) DaysUntilSTRExpiry() int {
	t, ok := parseDate(c.TanggalKadaluarsa)
	if !ok {
		return 0
	}
	remaining := time.Until(t)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Hours() / 24)
}

// ──────────────────────────────────────────────────────────
//  Tanggal Registrasi Ulang helpers
// ──────────────────────────────────────────────────────────

// IsRegistrasiUlangPassed mengecek apakah tanggal_registrasi_ulang sudah lewat
func (c *Csirt) IsRegistrasiUlangPassed() bool {
	t, ok := parseDate(c.TanggalRegistrasiUlang)
	if !ok {
		return false
	}
	return time.Now().After(t)
}

// IsRegistrasiUlangSoon mengecek apakah tanggal_registrasi_ulang akan jatuh tempo
// dalam STRExpiryWarnDays hari ke depan (belum lewat, tapi sudah dekat)
func (c *Csirt) IsRegistrasiUlangSoon() bool {
	t, ok := parseDate(c.TanggalRegistrasiUlang)
	if !ok {
		return false
	}
	now := time.Now()
	warnDate := t.AddDate(0, 0, -STRExpiryWarnDays)
	return now.After(warnDate) && now.Before(t)
}

// DaysUntilRegistrasiUlang mengembalikan sisa hari sebelum registrasi ulang.
// Mengembalikan 0 jika sudah lewat atau tanggal tidak valid.
func (c *Csirt) DaysUntilRegistrasiUlang() int {
	t, ok := parseDate(c.TanggalRegistrasiUlang)
	if !ok {
		return 0
	}
	remaining := time.Until(t)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Hours() / 24)
}