package utils

import (
	"ikas/internal/dto"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateIkasPDF(t *testing.T) {
	// Change working directory to project root so assets can be found
	originalWd, _ := os.Getwd()
	// Navigate up from ikas/internal/utils to root
	root := filepath.Join(originalWd, "..", "..", "..")
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Failed to change directory to root: %v", err)
	}
	defer os.Chdir(originalWd)

	// Sample data for coverage
	levels := []string{"Level 1", "Level 2", "Level 3", "Level 4", "Level 5", "Unknown"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			data := &dto.IkasResponse{
				Perusahaan: &dto.PerusahaanInIkas{
					NamaPerusahaan: "Test Corp",
					Alamat:         "Test Street",
					Email:          "test@corp.com",
					Sektor:         "Technology",
				},
				Telepon:                         "08123456789",
				Responden:                       "John Doe",
				Tanggal:                         "2026-04-30T20:42:23+07:00",
				Jabatan:                         "Manager",
				NilaiKematangan:                 3.5,
				KategoriKematanganKeamananSiber: level,
				Identifikasi: &dto.IdentifikasiInIkas{
					NilaiIdentifikasi: 4.0,
				},
				Proteksi: &dto.ProteksiInIkas{
					NilaiProteksi: 3.5,
				},
				Deteksi: &dto.DeteksiInIkas{
					NilaiDeteksi: 3.0,
				},
				Gulih: &dto.GulihInIkas{
					NilaiGulih: 3.2,
				},
			}

			pdfBytes, err := GenerateIkasPDF(data)
			if err != nil {
				t.Fatalf("GenerateIkasPDF failed for %s: %v", level, err)
			}
			if len(pdfBytes) == 0 {
				t.Errorf("Generated PDF is empty for %s", level)
			}
		})
	}

	// Test with missing domain data
	t.Run("MissingDomains", func(t *testing.T) {
		data := &dto.IkasResponse{
			Perusahaan: &dto.PerusahaanInIkas{
				NamaPerusahaan: "Test",
			},
			Tanggal: "2026-04-30",
		}
		_, err := GenerateIkasPDF(data)
		if err != nil {
			t.Fatalf("GenerateIkasPDF failed for MissingDomains: %v", err)
		}
	})
}

func TestToSafe(t *testing.T) {
	input := "\u2013 \u201Cquote\u201D \u00A0 \u2026"
	want := "- \"quote\"   ..."
	got := toSafe(input)
	if got != want {
		t.Errorf("toSafe() = %q, want %q", got, want)
	}
}
