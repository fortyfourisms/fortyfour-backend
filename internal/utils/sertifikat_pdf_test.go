package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fortyfour-backend/internal/models"
)

func TestGenerateSertifikatPDF_CreatesPDFFile(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	sertifikat := &models.Sertifikat{
		ID:              "sertifikat-1",
		NomorSertifikat: "CERT-001",
		IDKelas:         "kelas-1",
		IDUser:          "user-1",
		NamaPeserta:     "Afiif Najmi",
		NamaKelas:       "Fundamental Keamanan Siber",
		TanggalTerbit:   time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC),
	}

	path, err := GenerateSertifikatPDF(sertifikat)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedPath := filepath.Join("uploads", "sertifikat", "sertifikat-1.pdf")
	if path != expectedPath {
		t.Fatalf("expected path %q, got %q", expectedPath, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected generated PDF file to exist: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected generated PDF file to be non-empty")
	}
}

func TestGenerateSertifikatPDF_HandlesUnsafeCharacters(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	sertifikat := &models.Sertifikat{
		ID:              "sertifikat-unsafe",
		NomorSertifikat: "CERT-UTF-001",
		IDKelas:         "kelas-1",
		IDUser:          "user-1",
		NamaPeserta:     "Peserta – Utama",
		NamaKelas:       "Kelas Keamanan Siber → Lanjutan",
		TanggalTerbit:   time.Now(),
	}

	path, err := GenerateSertifikatPDF(sertifikat)
	if err != nil {
		t.Fatalf("expected no error for unsafe characters, got %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read generated PDF: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected generated PDF bytes to be non-empty")
	}
	if !strings.HasSuffix(path, ".pdf") {
		t.Fatalf("expected generated file to have .pdf suffix, got %q", path)
	}
}
