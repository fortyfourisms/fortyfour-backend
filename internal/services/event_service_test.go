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

// MockEventRepository
type MockEventRepository struct {
	CreateFunc   func(event *models.Event) error
	FindAllFunc  func() ([]models.Event, error)
	FindByIDFunc func(id int64) (*models.Event, error)
	UpdateFunc   func(event *models.Event) error
	DeleteFunc   func(id int64) error
}

func (m *MockEventRepository) Create(event *models.Event) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(event)
	}
	return nil
}

func (m *MockEventRepository) FindAll() ([]models.Event, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil
}

func (m *MockEventRepository) FindByID(id int64) (*models.Event, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockEventRepository) Update(event *models.Event) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(event)
	}
	return nil
}

func (m *MockEventRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestEventService_Create(t *testing.T) {
	mockRepo := &MockEventRepository{
		CreateFunc: func(event *models.Event) error {
			event.ID = 1
			return nil
		},
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id, Judul: "Test"}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	req := dto.CreateEventRequest{
		Judul:     "Test",
		Deskripsi: "Desc",
		Lokasi:    "Loc",
		Tanggal:   "2026-12-31T10:00:00Z",
	}
	err := svc.Create(req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEventService_Create_WithProducer(t *testing.T) {
	mockRepo := &MockEventRepository{
		CreateFunc: func(event *models.Event) error {
			event.ID = 1
			return nil
		},
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id, Judul: "Test"}, nil
		},
	}
	pkgProducer := pkgRmq.NewProducer(nil)
	mockProducer := internalRmq.NewProducer(pkgProducer)
	svc := services.NewEventService(mockRepo, mockProducer)

	req := dto.CreateEventRequest{
		Judul:     "Test",
		Deskripsi: "Desc",
		Lokasi:    "Loc",
		Tanggal:   "2026-12-31T10:00:00Z",
	}
	err := svc.Create(req)

	if err == nil {
		t.Errorf("expected error from nil channel, got nil")
	}
}

func TestEventService_Create_InvalidDate(t *testing.T) {
	mockRepo := &MockEventRepository{}
	svc := services.NewEventService(mockRepo, nil)

	req := dto.CreateEventRequest{
		Judul:     "Test",
		Deskripsi: "Desc",
		Lokasi:    "Loc",
		Tanggal:   "invalid-date",
	}
	err := svc.Create(req)

	if err == nil {
		t.Error("expected error for invalid date")
	}
}

func TestEventService_GetAll(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindAllFunc: func() ([]models.Event, error) {
			now := time.Now()
			return []models.Event{
				{ID: 1, Judul: "Test 1", Tanggal: now, CreatedAt: now, UpdatedAt: now},
				{ID: 2, Judul: "Test 2", Tanggal: now, CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	res, err := svc.GetAll()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 2 {
		t.Errorf("expected length 2, got %d", len(res))
	}
}

func TestEventService_GetAll_Error(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindAllFunc: func() ([]models.Event, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	_, err := svc.GetAll()

	if err == nil {
		t.Error("expected error")
	}
}

func TestEventService_GetByID(t *testing.T) {
	now := time.Now()
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id, Judul: "Test", Tanggal: now, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	res, err := svc.GetByID(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res.ID != 1 {
		t.Errorf("expected ID 1, got %d", res.ID)
	}
}

func TestEventService_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return nil, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	_, err := svc.GetByID(1)

	if err == nil {
		t.Error("expected error")
	}
}

func TestEventService_GetByID_Error(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	_, err := svc.GetByID(1)

	if err == nil {
		t.Error("expected error")
	}
}

func TestEventService_Update(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id, Judul: "Old"}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	judul := "New"
	desc := "New Desc"
	lokasi := "New Loc"
	tgl := "2026-12-31T10:00:00Z"
	req := dto.UpdateEventRequest{Judul: &judul, Deskripsi: &desc, Lokasi: &lokasi, Tanggal: &tgl}
	err := svc.Update(1, req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEventService_Update_WithProducer(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id, Judul: "Old"}, nil
		},
	}
	pkgProducer := pkgRmq.NewProducer(nil)
	mockProducer := internalRmq.NewProducer(pkgProducer)
	svc := services.NewEventService(mockRepo, mockProducer)

	judul := "New"
	req := dto.UpdateEventRequest{Judul: &judul}
	err := svc.Update(1, req)

	if err == nil {
		t.Errorf("expected error from nil channel, got nil")
	}
}

func TestEventService_Update_InvalidDate(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id, Judul: "Old"}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	tgl := "invalid-date"
	req := dto.UpdateEventRequest{Tanggal: &tgl}
	err := svc.Update(1, req)

	if err == nil {
		t.Error("expected error")
	}
}

func TestEventService_Update_NotFound(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return nil, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	req := dto.UpdateEventRequest{}
	err := svc.Update(1, req)

	if err == nil {
		t.Error("expected error")
	}
}

func TestEventService_Update_ErrorFind(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	req := dto.UpdateEventRequest{}
	err := svc.Update(1, req)

	if err == nil {
		t.Error("expected error")
	}
}

func TestEventService_Delete(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	err := svc.Delete(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEventService_Delete_WithProducer(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return &models.Event{ID: id}, nil
		},
	}
	pkgProducer := pkgRmq.NewProducer(nil)
	mockProducer := internalRmq.NewProducer(pkgProducer)
	svc := services.NewEventService(mockRepo, mockProducer)

	err := svc.Delete(1)

	if err == nil {
		t.Errorf("expected error from nil channel, got nil")
	}
}

func TestEventService_Delete_NotFound(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return nil, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	err := svc.Delete(1)

	if err == nil {
		t.Error("expected error")
	}
}

func TestEventService_Delete_ErrorFind(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindByIDFunc: func(id int64) (*models.Event, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewEventService(mockRepo, nil)

	err := svc.Delete(1)

	if err == nil {
		t.Error("expected error")
	}
}
