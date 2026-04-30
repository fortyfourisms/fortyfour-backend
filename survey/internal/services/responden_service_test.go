package services

import (
	"errors"
	"testing"
	"time"

	"survey/internal/dto"
	"survey/internal/models"
)

type mockRepo struct {
	createFn          func(m models.Responden) error
	getDetailByUserFn func(userID string) (*models.RespondenDetail, error)
	getAllFn          func() ([]models.RespondenDetail, error)
	getDetailByIDFn   func(id int) (*models.RespondenDetail, error)
	getByIDFn         func(id int) (*models.Responden, error)
	updateFn          func(id int, m models.Responden) error
}

func (m *mockRepo) Create(r models.Responden) error {
	return m.createFn(r)
}
func (m *mockRepo) GetDetailByUserID(userID string) (*models.RespondenDetail, error) {
	return m.getDetailByUserFn(userID)
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
	if req.UserID == "" {
		return errors.New("user_id wajib diisi")
	}
	return nil
}

func (m mockValidator) ValidateUpdate(req dto.UpdateRespondenRequest) error {
	if req.NoTelepon == "" {
		return errors.New("nomor telepon tidak boleh kosong")
	}
	return nil
}

func TestCreate_Success(t *testing.T) {
	mock := &mockRepo{
		createFn: func(m models.Responden) error { return nil },
		getDetailByUserFn: func(userID string) (*models.RespondenDetail, error) {
			return &models.RespondenDetail{
				ID:        1,
				UserID:    userID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	svc := NewRespondenService(mock, mockValidator{})

	req := dto.CreateRespondenRequest{
		UserID:             "user1",
		NoTelepon:          "08123456789",
		Sektor:             "IT",
		SertifikatTraining: "yes",
	}

	res, err := svc.Create(req)

	if err != nil {
		t.Fatal(err)
	}

	if res.UserID != "user1" {
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

	svc := NewRespondenService(mock, mockValidator{})

	res, err := svc.GetAll()

	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Error("expected 1 data")
	}
}