package utils

import (
	"testing"

	"survey/internal/dto"
)

func TestValidateCreateResponden_Success(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "+62812345678",
		Sektor:             "Logam",
		SektorLainnya:      "",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateCreateResponden_UserIDEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "",
		NoTelepon:          "+62812345678",
		Sektor:             "Logam",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "user_id wajib diisi" {
		t.Fatalf("expected user_id error, got %v", err)
	}
}

func TestValidateCreateResponden_UserIDTooShort(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "ab",
		NoTelepon:          "+62812345678",
		Sektor:             "Logam",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "user_id minimal 3 karakter" {
		t.Fatalf("expected min length error, got %v", err)
	}
}

func TestValidateCreateResponden_PhoneEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "",
		Sektor:             "Logam",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "no_telepon wajib diisi" {
		t.Fatalf("expected phone required error, got %v", err)
	}
}

func TestValidateCreateResponden_PhoneInvalidFormat(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "abc12345",
		Sektor:             "Logam",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "format no_telepon tidak valid" {
		t.Fatalf("expected invalid phone format, got %v", err)
	}
}

func TestValidateCreateResponden_PhoneTooShort(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "123",
		Sektor:             "Logam",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "no_telepon terlalu pendek" {
		t.Fatalf("expected phone too short, got %v", err)
	}
}

func TestValidateCreateResponden_InvalidSektor(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "+62812345678",
		Sektor:             "TidakAda",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "sektor tidak valid" {
		t.Fatalf("expected invalid sector, got %v", err)
	}
}

func TestValidateCreateResponden_LainnyaWithoutDetail(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "+62812345678",
		Sektor:             "Lainnya",
		SektorLainnya:      "",
		SertifikatTraining: "yes",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "sektor_lainnya wajib diisi jika sektor = Lainnya" {
		t.Fatalf("expected sektor_lainnya error, got %v", err)
	}
}

func TestValidateCreateResponden_SertifikatEmpty(t *testing.T) {
	req := dto.CreateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "+62812345678",
		Sektor:             "Logam",
		SertifikatTraining: "",
	}

	err := ValidateCreateResponden(req)
	if err == nil || err.Error() != "sertifikat_training wajib diisi" {
		t.Fatalf("expected sertifikat error, got %v", err)
	}
}

func TestValidateUpdateResponden_Success(t *testing.T) {
	req := dto.UpdateRespondenRequest{
		UserID:             "user123",
		NoTelepon:          "+62812345678",
		Sektor:             "Logam",
		SektorLainnya:      "",
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

func TestIsValidSektor(t *testing.T) {
	if !isValidSektor("Logam") {
		t.Error("expected valid sector")
	}

	if isValidSektor("TidakAda") {
		t.Error("expected invalid sector")
	}
}