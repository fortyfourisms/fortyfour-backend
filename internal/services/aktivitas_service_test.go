package services_test

import (
	"database/sql"
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
)

// ═══════════════════════════════════════════════════════════════════════════
// MOCK: AktivitasRepository
// ═══════════════════════════════════════════════════════════════════════════

type MockAktivitasRepository struct {
	CreateFunc          func(req dto.CreateAktivitasRequest) (int64, error)
	GetAllFunc          func() ([]dto.AktivitasResponse, error)
	GetByIDFunc         func(id int) (*dto.AktivitasResponse, error)
	GetByPerusahaanFunc func(perusahaanID string) ([]dto.AktivitasResponse, error)
	UpdateFunc          func(id int, req dto.UpdateAktivitasRequest) error
	DeleteFunc          func(id int) error
}

func (m *MockAktivitasRepository) Create(req dto.CreateAktivitasRequest) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(req)
	}
	return 1, nil
}

func (m *MockAktivitasRepository) GetAll() ([]dto.AktivitasResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockAktivitasRepository) GetByID(id int) (*dto.AktivitasResponse, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockAktivitasRepository) GetByPerusahaanID(perusahaanID string) ([]dto.AktivitasResponse, error) {
	if m.GetByPerusahaanFunc != nil {
		return m.GetByPerusahaanFunc(perusahaanID)
	}
	return nil, nil
}

func (m *MockAktivitasRepository) Update(id int, req dto.UpdateAktivitasRequest) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, req)
	}
	return nil
}

func (m *MockAktivitasRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// ═══════════════════════════════════════════════════════════════════════════
// MOCK: PerusahaanRepository (hanya GetByID yang dipakai validasi)
// ═══════════════════════════════════════════════════════════════════════════

type MockPerusahaanRepoForAktivitas struct {
	GetByIDFunc func(id string) (*dto.PerusahaanResponse, error)
}

func (m *MockPerusahaanRepoForAktivitas) Create(req dto.CreatePerusahaanRequest, id string) error {
	return nil
}
func (m *MockPerusahaanRepoForAktivitas) GetByID(id string) (*dto.PerusahaanResponse, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return &dto.PerusahaanResponse{}, nil
}
func (m *MockPerusahaanRepoForAktivitas) GetByNama(nama string) (*dto.PerusahaanResponse, error) {
	return nil, nil
}
func (m *MockPerusahaanRepoForAktivitas) GetAll() ([]dto.PerusahaanResponse, error) {
	return nil, nil
}
func (m *MockPerusahaanRepoForAktivitas) Update(id string, p dto.PerusahaanResponse) error {
	return nil
}
func (m *MockPerusahaanRepoForAktivitas) Delete(id string) error { return nil }

// ═══════════════════════════════════════════════════════════════════════════
// TESTS
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasService_GetAllowedJenis(t *testing.T) {
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		&MockPerusahaanRepoForAktivitas{},
		nil, nil,
	)

	jenis := svc.GetAllowedJenis()
	if len(jenis) == 0 {
		t.Error("expected non-empty allowed jenis list")
	}
}

func TestAktivitasService_GetAll_Success(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetAllFunc: func() ([]dto.AktivitasResponse, error) {
			return []dto.AktivitasResponse{
				{ID: 1, Judul: "Aktivitas 1"},
				{ID: 2, Judul: "Aktivitas 2"},
			}, nil
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	result, err := svc.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestAktivitasService_GetAll_Error(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetAllFunc: func() ([]dto.AktivitasResponse, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	_, err := svc.GetAll()
	if err == nil {
		t.Error("expected error")
	}
}

func TestAktivitasService_GetByID_Success(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return &dto.AktivitasResponse{ID: id, Judul: "Test"}, nil
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	result, err := svc.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestAktivitasService_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	_, err := svc.GetByID(999)
	if err == nil {
		t.Error("expected error for not found")
	}
	if err.Error() != "data tidak ditemukan" {
		t.Errorf("expected 'data tidak ditemukan', got '%s'", err.Error())
	}
}

func TestAktivitasService_GetByID_DBError(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	_, err := svc.GetByID(1)
	if err == nil {
		t.Error("expected error")
	}
}

func TestAktivitasService_GetByPerusahaanID_Success(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByPerusahaanFunc: func(perusahaanID string) ([]dto.AktivitasResponse, error) {
			return []dto.AktivitasResponse{
				{ID: 1, PerusahaanID: perusahaanID, Judul: "Act 1"},
			}, nil
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	result, err := svc.GetByPerusahaanID("uuid-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}

func TestAktivitasService_GetByPerusahaanID_Error(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByPerusahaanFunc: func(perusahaanID string) ([]dto.AktivitasResponse, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	_, err := svc.GetByPerusahaanID("uuid-123")
	if err == nil {
		t.Error("expected error")
	}
}

func TestAktivitasService_Create_ValidationEmptyPerusahaanID(t *testing.T) {
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		&MockPerusahaanRepoForAktivitas{},
		nil, nil,
	)

	req := dto.CreateAktivitasRequest{
		Judul:          "Test",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas"},
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected validation error for empty perusahaan_id")
	}
}

func TestAktivitasService_Create_ValidationEmptyJudul(t *testing.T) {
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		&MockPerusahaanRepoForAktivitas{},
		nil, nil,
	)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-123",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas"},
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected validation error for empty judul")
	}
}

func TestAktivitasService_Create_ValidationInvalidJenis(t *testing.T) {
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		&MockPerusahaanRepoForAktivitas{},
		nil, nil,
	)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-123",
		Judul:          "Test",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"invalid_type"},
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected validation error for invalid jenis_aktivitas")
	}
}

func TestAktivitasService_Create_ValidationInvalidDateFormat(t *testing.T) {
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		&MockPerusahaanRepoForAktivitas{},
		nil, nil,
	)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-123",
		Judul:          "Test",
		TanggalMulai:   "invalid-date",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas"},
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected validation error for invalid date format")
	}
}

func TestAktivitasService_Create_ValidationChronologicalOrder(t *testing.T) {
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		&MockPerusahaanRepoForAktivitas{},
		nil, nil,
	)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-123",
		Judul:          "Test",
		TanggalMulai:   "2024-01-02",
		TanggalSelesai: "2024-01-01",
		JenisAktivitas: []string{"dinas"},
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected validation error for chronological order")
	}
}

func TestAktivitasService_Create_PerusahaanNotFound(t *testing.T) {
	perusahaanRepo := &MockPerusahaanRepoForAktivitas{
		GetByIDFunc: func(id string) (*dto.PerusahaanResponse, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		perusahaanRepo,
		nil, nil,
	)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "nonexistent",
		Judul:          "Test",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas"},
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected error for nonexistent perusahaan")
	}
}

func TestAktivitasService_Create_Success_NilProducer(t *testing.T) {
	perusahaanRepo := &MockPerusahaanRepoForAktivitas{
		GetByIDFunc: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: id}, nil
		},
	}
	svc := services.NewAktivitasService(
		&MockAktivitasRepository{},
		perusahaanRepo,
		nil, nil,
	)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-123",
		Judul:          "Test",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas"},
	}
	// With nil producer, Create returns (nil, nil) — no publish, just validates
	_, err := svc.Create(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAktivitasService_Update_NotFound(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	req := dto.UpdateAktivitasRequest{}
	_, err := svc.Update(999, req)
	if err == nil {
		t.Error("expected error for not found")
	}
}

func TestAktivitasService_Update_InvalidJenis(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return &dto.AktivitasResponse{ID: id, PerusahaanID: "uuid-123"}, nil
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	invalid := []string{"invalid_type"}
	req := dto.UpdateAktivitasRequest{JenisAktivitas: &invalid}
	_, err := svc.Update(1, req)
	if err == nil {
		t.Error("expected validation error for invalid jenis")
	}
}

func TestAktivitasService_Update_Success_NilProducer(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return &dto.AktivitasResponse{ID: id, PerusahaanID: "uuid-123"}, nil
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	judul := "Updated"
	req := dto.UpdateAktivitasRequest{Judul: &judul}
	_, err := svc.Update(1, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAktivitasService_Delete_NotFound(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	err := svc.Delete(999)
	if err == nil {
		t.Error("expected error for not found")
	}
}

func TestAktivitasService_Delete_Success_NilProducer(t *testing.T) {
	mockRepo := &MockAktivitasRepository{
		GetByIDFunc: func(id int) (*dto.AktivitasResponse, error) {
			return &dto.AktivitasResponse{ID: id, PerusahaanID: "uuid-123"}, nil
		},
	}
	svc := services.NewAktivitasService(mockRepo, &MockPerusahaanRepoForAktivitas{}, nil, nil)

	err := svc.Delete(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
