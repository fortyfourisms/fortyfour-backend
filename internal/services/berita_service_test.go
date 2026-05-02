package services_test

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	internalRmq "fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/services"
	pkgRmq "fortyfour-backend/pkg/rabbitmq"
	"testing"
	"time"
)

// MockBeritaRepository
type MockBeritaRepository struct {
	CreateFunc   func(berita *models.Berita) error
	FindAllFunc  func() ([]models.Berita, error)
	FindByIDFunc func(id int64) (*models.Berita, error)
	UpdateFunc   func(berita *models.Berita) error
	DeleteFunc   func(id int64) error
}

func (m *MockBeritaRepository) Create(berita *models.Berita) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(berita)
	}
	return nil
}

func (m *MockBeritaRepository) FindAll() ([]models.Berita, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil
}

func (m *MockBeritaRepository) FindByID(id int64) (*models.Berita, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockBeritaRepository) Update(berita *models.Berita) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(berita)
	}
	return nil
}

func (m *MockBeritaRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestBeritaService_Create(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		CreateFunc: func(berita *models.Berita) error {
			berita.ID = 1
			return nil
		},
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return &models.Berita{ID: id, Judul: "Test"}, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	req := dto.CreateBeritaRequest{Judul: "Test", Deskripsi: "Desc"}
	err := svc.Create("author1", req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBeritaService_Create_WithProducer(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		CreateFunc: func(berita *models.Berita) error {
			berita.ID = 1
			return nil
		},
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return &models.Berita{ID: id, Judul: "Test"}, nil
		},
	}
	pkgProducer := pkgRmq.NewProducer(nil) // nil rmq, will return error but not panic
	mockProducer := internalRmq.NewProducer(pkgProducer)
	svc := services.NewBeritaService(mockRepo, mockProducer, nil)

	req := dto.CreateBeritaRequest{Judul: "Test", Deskripsi: "Desc"}
	err := svc.Create("author1", req)

	// Since channel is nil, Publish returns error, so Create should return error
	if err == nil {
		t.Errorf("expected error from nil channel, got nil")
	}
}

func TestBeritaService_GetAll(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindAllFunc: func() ([]models.Berita, error) {
			now := time.Now()
			return []models.Berita{
				{ID: 1, Judul: "Test 1", CreatedAt: now, UpdatedAt: now},
				{ID: 2, Judul: "Test 2", CreatedAt: now, UpdatedAt: now, Author: &models.User{ID: "a", Username: "user"}},
			}, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	res, err := svc.GetAll()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 2 {
		t.Errorf("expected length 2, got %d", len(res))
	}
}

func TestBeritaService_GetAll_Error(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindAllFunc: func() ([]models.Berita, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	_, err := svc.GetAll()

	if err == nil {
		t.Error("expected error")
	}
}

func TestBeritaService_GetByID(t *testing.T) {
	now := time.Now()
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return &models.Berita{ID: id, Judul: "Test", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	res, err := svc.GetByID(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res.ID != 1 {
		t.Errorf("expected ID 1, got %d", res.ID)
	}
}

func TestBeritaService_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return nil, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	_, err := svc.GetByID(1)

	if err == nil {
		t.Error("expected error")
	}
}

func TestBeritaService_GetByID_Error(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	_, err := svc.GetByID(1)

	if err == nil {
		t.Error("expected error")
	}
}

func TestBeritaService_Update(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return &models.Berita{ID: id, Judul: "Old"}, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	judul := "New"
	desc := "New Desc"
	req := dto.UpdateBeritaRequest{Judul: &judul, Deskripsi: &desc}
	err := svc.Update(1, req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBeritaService_Update_WithProducer(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return &models.Berita{ID: id, Judul: "Old"}, nil
		},
	}
	pkgProducer := pkgRmq.NewProducer(nil)
	mockProducer := internalRmq.NewProducer(pkgProducer)
	svc := services.NewBeritaService(mockRepo, mockProducer, nil)

	judul := "New"
	req := dto.UpdateBeritaRequest{Judul: &judul}
	err := svc.Update(1, req)

	if err == nil {
		t.Errorf("expected error from nil channel, got nil")
	}
}

func TestBeritaService_Update_NotFound(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return nil, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	req := dto.UpdateBeritaRequest{}
	err := svc.Update(1, req)

	if err == nil {
		t.Error("expected error")
	}
}

func TestBeritaService_Update_ErrorFind(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	req := dto.UpdateBeritaRequest{}
	err := svc.Update(1, req)

	if err == nil {
		t.Error("expected error")
	}
}

func TestBeritaService_Delete(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return &models.Berita{ID: id}, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	err := svc.Delete(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBeritaService_Delete_WithProducer(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return &models.Berita{ID: id}, nil
		},
	}
	pkgProducer := pkgRmq.NewProducer(nil)
	mockProducer := internalRmq.NewProducer(pkgProducer)
	svc := services.NewBeritaService(mockRepo, mockProducer, nil)

	err := svc.Delete(1)

	if err == nil {
		t.Errorf("expected error from nil channel, got nil")
	}
}

func TestBeritaService_Delete_NotFound(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return nil, nil
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	err := svc.Delete(1)

	if err == nil {
		t.Error("expected error")
	}
}

func TestBeritaService_Delete_ErrorFind(t *testing.T) {
	mockRepo := &MockBeritaRepository{
		FindByIDFunc: func(id int64) (*models.Berita, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewBeritaService(mockRepo, nil, nil)

	err := svc.Delete(1)

	if err == nil {
		t.Error("expected error")
	}
}
