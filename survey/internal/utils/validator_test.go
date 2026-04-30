package utils

import (
	"testing"

	"survey/internal/dto"
)

func TestValidateCreateResponden_Success(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateCreateResponden_IdPerusahaanEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "id_perusahaan wajib diisi" {
		t.Fatalf("expected id_perusahaan error, got %v", err)
	}
}

func TestValidateCreateResponden_NamaLengkapEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "nama_lengkap wajib diisi" {
		t.Fatalf("expected nama_lengkap error, got %v", err)
	}
}

func TestValidateCreateResponden_NamaLengkapTooShort(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "ab",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "nama_lengkap minimal 3 karakter" {
		t.Fatalf("expected min length error, got %v", err)
	}
}

func TestValidateCreateResponden_JabatanEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "",
		Email:              "email@mail.com",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "jabatan wajib diisi" {
		t.Fatalf("expected jabatan error, got %v", err)
	}
}

func TestValidateCreateResponden_EmailEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "email wajib diisi" {
		t.Fatalf("expected email error, got %v", err)
	}
}

func TestValidateCreateResponden_EmailInvalid(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "invalid-email",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "format email tidak valid" {
		t.Fatalf("expected invalid email error, got %v", err)
	}
}

func TestValidateCreateResponden_PhoneEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "no_telepon wajib diisi" {
		t.Fatalf("expected phone required error, got %v", err)
	}
}

func TestValidateCreateResponden_PhoneInvalidFormat(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "abc12345",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "format no_telepon tidak valid" {
		t.Fatalf("expected invalid phone format, got %v", err)
	}
}

func TestValidateCreateResponden_PhoneTooShort(t *testing.T) {
	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "123",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "no_telepon terlalu pendek" {
		t.Fatalf("expected phone too short, got %v", err)
	}
}

func TestValidateUpdateResponden_Success(t *testing.T) {
	req := dto.UpdateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "+62812345678",
		SertifikatTraining: "yes",
	}

	err := ValidateUpdateResponden(req)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestIsPhone(t *testing.T) {
	if !isPhone("+628123") {
		t.Error("expected valid phone")
	}

	if isPhone("abc123") {
		t.Error("expected invalid phone")
	}
}
