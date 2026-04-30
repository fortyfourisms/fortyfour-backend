package utils

import (
	"errors"
	"regexp"
	"strings"
	"survey/internal/dto"
)

// CREATE VALIDATION
func ValidateCreateResponden(req dto.CreateRespondenRequest) error {

	// normalize
	req.IdPerusahaan = strings.TrimSpace(req.IdPerusahaan)
	req.NamaLengkap = strings.TrimSpace(req.NamaLengkap)
	req.Jabatan = strings.TrimSpace(req.Jabatan)
	req.Email = strings.TrimSpace(req.Email)
	req.NoTelepon = strings.TrimSpace(req.NoTelepon)
	req.SertifikatTraining = strings.TrimSpace(req.SertifikatTraining)

	// ID PERUSAHAAN
	if req.IdPerusahaan == "" {
		return errors.New("id_perusahaan wajib diisi")
	}

	// NAMA LENGKAP
	if req.NamaLengkap == "" {
		return errors.New("nama_lengkap wajib diisi")
	}
	if len(req.NamaLengkap) < 3 {
		return errors.New("nama_lengkap minimal 3 karakter")
	}

	// JABATAN
	if req.Jabatan == "" {
		return errors.New("jabatan wajib diisi")
	}

	// EMAIL
	if req.Email == "" {
		return errors.New("email wajib diisi")
	}
	if !isValidEmail(req.Email) {
		return errors.New("format email tidak valid")
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

	return nil
}

// UPDATE VALIDATION
func ValidateUpdateResponden(req dto.UpdateRespondenRequest) error {
	return ValidateCreateResponden(dto.CreateRespondenRequest{
		IdPerusahaan:       req.IdPerusahaan,
		NamaLengkap:        req.NamaLengkap,
		Jabatan:            req.Jabatan,
		Email:              req.Email,
		NoTelepon:          req.NoTelepon,
		SertifikatTraining: req.SertifikatTraining,
	})
}

// Validasi nomor telepon (angka + optional '+')
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

// Validasi email sederhana
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
