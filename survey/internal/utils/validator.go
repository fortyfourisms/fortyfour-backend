package utils

import (
	"survey/internal/dto"
	"strings"
)

// ValidateCreateResponden validates a CreateRespondenRequest.
// Mengembalikan map field -> pesan error. Jika kosong berarti valid.
func ValidateCreateResponden(req dto.CreateRespondenRequest) map[string]string {
	errors := map[string]string{}

	// Nama lengkap
	if strings.TrimSpace(req.NamaLengkap) == "" {
		errors["nama_lengkap"] = "Nama lengkap wajib diisi"
	}

	// Jabatan
	if strings.TrimSpace(req.Jabatan) == "" {
		errors["jabatan"] = "Jabatan wajib diisi"
	}

	// Perusahaan
	if strings.TrimSpace(req.Perusahaan) == "" {
		errors["perusahaan"] = "Perusahaan wajib diisi"
	}

	// Email
	email := strings.TrimSpace(req.Email)
	switch {
	case email == "":
		errors["email"] = "Email wajib diisi"
	case !isEmail(email):
		errors["email"] = "Format email tidak valid"
	}

	// Nomor telepon
	phone := strings.TrimSpace(req.NoTelepon)
	switch {
	case phone == "":
		errors["no_telepon"] = "Nomor telepon wajib diisi"
	case !isPhone(phone):
		errors["no_telepon"] = "Nomor telepon hanya boleh berisi angka"
	}

	// Validasi sektor
	if !isValidSektor(req.Sektor) {
		errors["sektor"] = "Sektor tidak valid"
	}

	// Jika sektor lainnya
	if req.Sektor == "Lainnya" && strings.TrimSpace(req.SektorLainnya) == "" {
		errors["sektor_lainnya"] = "Keterangan sektor lainnya wajib diisi"
	}

	return errors
}

// ValidateUpdateResponden menggunakan validasi yang sama dengan create
func ValidateUpdateResponden(req dto.UpdateRespondenRequest) map[string]string {
	return ValidateCreateResponden(dto.CreateRespondenRequest{
		NamaLengkap:        req.NamaLengkap,
		Jabatan:            req.Jabatan,
		Perusahaan:         req.Perusahaan,
		Email:              req.Email,
		NoTelepon:          req.NoTelepon,
		Sektor:             req.Sektor,
		SektorLainnya:      req.SektorLainnya,
		SertifikatTraining: req.SertifikatTraining,
	})
}

// =============================
// Helper Validation Functions
// =============================

// isEmail melakukan validasi sederhana format email
func isEmail(email string) bool {
	email = strings.TrimSpace(email)

	if !strings.Contains(email, "@") {
		return false
	}

	if !strings.Contains(email, ".") {
		return false
	}

	return true
}

// isPhone memastikan nomor telepon hanya berisi angka
func isPhone(phone string) bool {
	if phone == "" {
		return false
	}

	for _, r := range phone {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// isValidSektor memastikan sektor yang dipilih valid
func isValidSektor(sektor string) bool {
	validSektor := []string{
        "Industri Makanan dan Minuman",
        "Industri Tekstil dan Pakaian",
        "Industri Kimia",
        "Industri Otomotif",
        "Industri Elektronik",
        "Industri Farmasi",
        "Industri Alat Kesehatan",
        "Jasa Konstruksi",
        "Industri Keamanan Siber",
        "Industri Pertahanan",
        "Lainnya",
	}

	for _, s := range validSektor {
		if sektor == s {
			return true
		}
	}

	return false
}