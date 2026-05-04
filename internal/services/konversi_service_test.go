package services_test

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
)

// ═══════════════════════════════════════════════════════════════════════════
// MOCK: KonversiRepository
// ═══════════════════════════════════════════════════════════════════════════

type MockKonversiRepository struct {
	GetAllKonversiFn func(perusahaanID string) ([]dto.KonversiResponse, error)
}

func (m *MockKonversiRepository) GetAllKonversi(perusahaanID string) ([]dto.KonversiResponse, error) {
	if m.GetAllKonversiFn != nil {
		return m.GetAllKonversiFn(perusahaanID)
	}
	return nil, nil
}

// ═══════════════════════════════════════════════════════════════════════════
// TESTS
// ═══════════════════════════════════════════════════════════════════════════

func TestKonversiService_GetKonversi_AllPerusahaan(t *testing.T) {
	mockRepo := &MockKonversiRepository{
		GetAllKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			return []dto.KonversiResponse{
				{PerusahaanID: "uuid-1", NamaPerusahaan: "PT A", PoinIkas: 1, PoinKse: 1, PoinSurvey: 0, PoinCsirt: 1},
				{PerusahaanID: "uuid-2", NamaPerusahaan: "PT B", PoinIkas: 0, PoinKse: 0, PoinSurvey: 1, PoinCsirt: 0},
			}, nil
		},
	}
	svc := services.NewKonversiService(mockRepo, nil)

	results, err := svc.GetKonversi("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	// Check computation: PT A has 1+1+0+1 = 3, percentage = (3/4)*100 = 75.0
	if results[0].TotalPoin != 3 {
		t.Errorf("expected TotalPoin 3, got %d", results[0].TotalPoin)
	}
	if results[0].Persentase != 75.0 {
		t.Errorf("expected Persentase 75.0, got %f", results[0].Persentase)
	}

	// Check computation: PT B has 0+0+1+0 = 1, percentage = (1/4)*100 = 25.0
	if results[1].TotalPoin != 1 {
		t.Errorf("expected TotalPoin 1, got %d", results[1].TotalPoin)
	}
	if results[1].Persentase != 25.0 {
		t.Errorf("expected Persentase 25.0, got %f", results[1].Persentase)
	}
}

func TestKonversiService_GetKonversi_ByPerusahaanID(t *testing.T) {
	mockRepo := &MockKonversiRepository{
		GetAllKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			if perusahaanID != "uuid-1" {
				return nil, nil
			}
			return []dto.KonversiResponse{
				{PerusahaanID: "uuid-1", NamaPerusahaan: "PT A", PoinIkas: 1, PoinKse: 1, PoinSurvey: 1, PoinCsirt: 1},
			}, nil
		},
	}
	svc := services.NewKonversiService(mockRepo, nil)

	results, err := svc.GetKonversi("uuid-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	// Full participation: 4/4 * 100 = 100.0%
	if results[0].TotalPoin != 4 {
		t.Errorf("expected TotalPoin 4, got %d", results[0].TotalPoin)
	}
	if results[0].Persentase != 100.0 {
		t.Errorf("expected Persentase 100.0, got %f", results[0].Persentase)
	}
}

func TestKonversiService_GetKonversi_ZeroPoin(t *testing.T) {
	mockRepo := &MockKonversiRepository{
		GetAllKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			return []dto.KonversiResponse{
				{PerusahaanID: "uuid-3", NamaPerusahaan: "PT C", PoinIkas: 0, PoinKse: 0, PoinSurvey: 0, PoinCsirt: 0},
			}, nil
		},
	}
	svc := services.NewKonversiService(mockRepo, nil)

	results, err := svc.GetKonversi("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if results[0].TotalPoin != 0 {
		t.Errorf("expected TotalPoin 0, got %d", results[0].TotalPoin)
	}
	if results[0].Persentase != 0.0 {
		t.Errorf("expected Persentase 0.0, got %f", results[0].Persentase)
	}
}

func TestKonversiService_GetKonversi_RepoError(t *testing.T) {
	mockRepo := &MockKonversiRepository{
		GetAllKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewKonversiService(mockRepo, nil)

	_, err := svc.GetKonversi("")
	if err == nil {
		t.Error("expected error")
	}
}

func TestKonversiService_GetKonversi_EmptyResult(t *testing.T) {
	mockRepo := &MockKonversiRepository{
		GetAllKonversiFn: func(perusahaanID string) ([]dto.KonversiResponse, error) {
			return []dto.KonversiResponse{}, nil
		},
	}
	svc := services.NewKonversiService(mockRepo, nil)

	results, err := svc.GetKonversi("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
