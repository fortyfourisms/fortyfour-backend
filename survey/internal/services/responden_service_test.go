package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
)

type mockRepo struct {
	createFn        func(m models.Responden) (int64, error)
	getAllFn        func() ([]models.RespondenDetail, error)
	getDetailByIDFn func(id int) (*models.RespondenDetail, error)
	getByIDFn       func(id int) (*models.Responden, error)
	updateFn        func(id int, m models.Responden) error
}

func (m *mockRepo) Create(r models.Responden) (int64, error) {
	return m.createFn(r)
}
func (m *mockRepo) GetAllDetail() ([]models.RespondenDetail, error) {
	return m.getAllFn()
}
func (m *mockRepo) GetDetailByID(id int) (*models.RespondenDetail, error) {
	return m.getDetailByIDFn(id)
}
func (m *mockRepo) GetByID(id int) (*models.Responden, error) {
	return m.getByIDFn(id)
}
func (m *mockRepo) Update(id int, r models.Responden) error {
	return m.updateFn(id, r)
}

type mockValidator struct{}

func (m mockValidator) ValidateCreate(req dto.CreateRespondenRequest) error {
	if req.IdPerusahaan == "" {
		return errors.New("id_perusahaan wajib diisi")
	}
	return nil
}

func (m mockValidator) ValidateUpdate(req dto.UpdateRespondenRequest) error {
	if req.NoTelepon == "" {
		return errors.New("nomor telepon tidak boleh kosong")
	}
	return nil
}

// mockCacheForService is a no-op cache for testing
type mockCacheForService struct{}

func (m *mockCacheForService) Get(ctx context.Context, key string) (string, bool, error) {
	return "", false, nil
}
func (m *mockCacheForService) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return nil
}
func (m *mockCacheForService) Del(ctx context.Context, key string) error {
	return nil
}

func TestCreate_Success(t *testing.T) {
	mock := &mockRepo{
		createFn: func(m models.Responden) (int64, error) { return 1, nil },
		getDetailByIDFn: func(id int) (*models.RespondenDetail, error) {
			return &models.RespondenDetail{
				ID:           id,
				IdPerusahaan: "perusahaan1",
				NamaLengkap:  "Nama Lengkap",
				Jabatan:      "Manager",
				Email:        "email@mail.com",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
	}

	svc := NewRespondenService(mock, mockValidator{}, &mockCacheForService{})

	req := dto.CreateRespondenRequest{
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "08123456789",
		SertifikatTraining: "yes",
	}

	res, err := svc.Create(req)

	if err != nil {
		t.Fatal(err)
	}

	if res.IdPerusahaan != "perusahaan1" {
		t.Error("invalid response")
	}
}

func TestGetAll_Success(t *testing.T) {
	mock := &mockRepo{
		getAllFn: func() ([]models.RespondenDetail, error) {
			return []models.RespondenDetail{
				{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, nil
		},
	}

	svc := NewRespondenService(mock, mockValidator{}, &mockCacheForService{})

	res, err := svc.GetAll()

	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Error("expected 1 data")
	}
}
