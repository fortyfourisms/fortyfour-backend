package utils

import (
	"errors"
	"strings"
	"survey/internal/dto"
)

// CREATE VALIDATION
func ValidateCreateResponden(req dto.CreateRespondenRequest) error {

	// normalize input
	req.UserID = strings.TrimSpace(req.UserID)
	req.NoTelepon = strings.TrimSpace(req.NoTelepon)
	req.Sektor = strings.TrimSpace(req.Sektor)
	req.SektorLainnya = strings.TrimSpace(req.SektorLainnya)
	req.SertifikatTraining = strings.TrimSpace(req.SertifikatTraining)

	// USER ID
	if req.UserID == "" {
		return errors.New("user_id wajib diisi")
	}
	if len(req.UserID) < 3 {
		return errors.New("user_id minimal 3 karakter")
	}

	// NO TELEPON
	if req.NoTelepon == "" {
		return errors.New("no_telepon wajib diisi")
	}
	if !isPhone(req.NoTelepon) {
		return errors.New("format no_telepon tidak valid")
	}
	if len(req.NoTelepon) < 8 {
		return errors.New("no_telepon terlalu pendek")
	}

	// SEKTOR
	if req.Sektor == "" {
		return errors.New("sektor wajib diisi")
	}
	if !isValidSektor(req.Sektor) {
		return errors.New("sektor tidak valid")
	}

	// SEKTOR LAINNYA
	if req.Sektor == "Lainnya" {
		if req.SektorLainnya == "" {
			return errors.New("sektor_lainnya wajib diisi jika sektor = Lainnya")
		}
	} else {
		req.SektorLainnya = ""
	}

	// SERTIFIKAT
	if req.SertifikatTraining == "" {
		return errors.New("sertifikat_training wajib diisi")
	}

	return nil
}

// UPDATE VALIDATION
func ValidateUpdateResponden(req dto.UpdateRespondenRequest) error {
	return ValidateCreateResponden(dto.CreateRespondenRequest{
		UserID:             req.UserID,
		NoTelepon:          req.NoTelepon,
		Sektor:             req.Sektor,
		SektorLainnya:      req.SektorLainnya,
		SertifikatTraining: req.SertifikatTraining,
	})
}

// HELPER FUNCTIONS

// Validasi nomor telepon (hanya angka + optional '+' di depan)
func isPhone(phone string) bool {
	if phone == "" {
		return false
	}

	for i, r := range phone {
		if i == 0 && r == '+' {
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// Validasi sektor (pakai mapping cepat)
func isValidSektor(sektor string) bool {
	validSektor := map[string]bool{
		"Keamanan Siber": true,
		"Jasa Konsultan dan Sertifikasi Keamanan Informasi":        true,
		"Industri Pulp dan Kertas":                                 true,
		"Jasa Konstruksi":                                          true,
		"Jasa Sertifikasi, Inspeksi, Pengujian, dan Survey":        true,
		"Jasa Industri":                                            true,
		"Komponen Kendaraan Listrik":                               true,
		"Alat Kesehatan":                                           true,
		"Peralatan Listrik":                                        true,
		"Industri Kecil dan Menengah Furnitur, dan Bahan Bangunan": true,
		"Logam":                               true,
		"Permesinan dan Alat Mesin Pertanian": true,
		"Maritim, Alat Transportasi, dan Alat Pertahanan":      true,
		"Elektronika dan Telematika":                           true,
		"Hasil Hutan dan Perkebunan":                           true,
		"Makanan, Hasil Laut, dan Perikanan":                   true,
		"Minuman, Tembakau, dan Bahan Penyegar":                true,
		"Kemurgi, Oleokimia, dan Pakan":                        true,
		"Kimia hulu":                                           true,
		"Kawasan Industri":                                     true,
		"Semen, Keramik, dan Pengolahan Bahan Galian Nonlogam": true,
		"Tekstil, Kulit, dan Alas Kaki":                        true,
		"Industri Aneka":                                       true,
		"Industri Bahan dan Produk Farmasi":                    true,
		"Kimia Hilir":                                          true,
		"Lainnya":                                              true,
	}

	return validSektor[sektor]
}
