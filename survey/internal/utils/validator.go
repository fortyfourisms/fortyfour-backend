package utils

import (
	"errors"
	"regexp"
	"strings"
	"survey/internal/dto"
)

// CREATE VALIDATION (UPSERT USER)
func ValidateCreateResponden(req dto.CreateRespondenRequest) error {

	// NORMALIZATION
	req.IdPerusahaan = strings.TrimSpace(req.IdPerusahaan)
	req.NamaLengkap = strings.TrimSpace(req.NamaLengkap)
	req.Jabatan = strings.TrimSpace(req.Jabatan)
	req.Email = strings.TrimSpace(req.Email)
	req.NoTelepon = strings.TrimSpace(req.NoTelepon)

	// SAFE POINTER HANDLING
	var sertifikat string
	if req.SertifikatTraining != nil {
		sertifikat = strings.TrimSpace(*req.SertifikatTraining)
	}

	// VALIDATION

	if req.IdPerusahaan == "" {
		return errors.New("id_perusahaan wajib diisi")
	}

	if req.NamaLengkap == "" {
		return errors.New("nama_lengkap wajib diisi")
	}
	if len(req.NamaLengkap) < 3 {
		return errors.New("nama_lengkap minimal 3 karakter")
	}

	if req.Jabatan == "" {
		return errors.New("jabatan wajib diisi")
	}

	if req.Email == "" {
		return errors.New("email wajib diisi")
	}
	if !isValidEmail(req.Email) {
		return errors.New("format email tidak valid")
	}

	if req.NoTelepon == "" {
		return errors.New("no_telepon wajib diisi")
	}
	if len(req.NoTelepon) < 8 {
		return errors.New("no_telepon terlalu pendek")
	}
	if !isPhone(req.NoTelepon) {
		return errors.New("format no_telepon tidak valid")
	}

	_ = sertifikat

	return nil
}

// PHONE VALIDATION
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

// EMAIL VALIDATION
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
