package services_test

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockEventRepository implements repository.EventRepositoryInterface for tests
type MockEventRepository struct {
	CreateFunc                    func(event *models.Event) error
	FindAllFunc                   func() ([]models.Event, error)
	FindByIDFunc                  func(id string) (*models.Event, error)
	FindBySlugFunc                func(slug string) (*models.Event, error)
	UpdateFunc                    func(event *models.Event) error
	DeleteFunc                    func(id string) error
	CreateRegistrationFunc        func(reg *models.EventRegistration) error
	FindRegistrationByIDFunc      func(id string) (*models.EventRegistration, error)
	ExistsRegistrationFunc        func(eventID string, email string) (bool, error)
	UpdateRegistrationPayloadFunc func(id string, payload string) error

	// Helper fields for assertions
	updatedPayload string
}

type MockRedis struct{}

func (m *MockRedis) Set(key string, value interface{}, expiration time.Duration) error {
	return nil
}
func (m *MockRedis) Get(key string) (string, error) {
	return "", nil
}
func (m *MockRedis) Delete(key string) error {
	return nil
}
func (m *MockRedis) Exists(key string) (bool, error) {
	return false, nil
}
func (m *MockRedis) Scan(pattern string) ([]string, error) {
	return nil, nil
}
func (m *MockRedis) Close() error {
	return nil
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

func (m *MockEventRepository) FindByID(id string) (*models.Event, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockEventRepository) FindBySlug(slug string) (*models.Event, error) {
	if m.FindBySlugFunc != nil {
		return m.FindBySlugFunc(slug)
	}
	return nil, nil
}

func (m *MockEventRepository) Update(event *models.Event) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(event)
	}
	return nil
}

func (m *MockEventRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockEventRepository) CreateRegistration(reg *models.EventRegistration) error {
	if m.CreateRegistrationFunc != nil {
		return m.CreateRegistrationFunc(reg)
	}
	reg.ID = "101"
	reg.CreatedAt = time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC)
	reg.UpdatedAt = reg.CreatedAt
	return nil
}

func (m *MockEventRepository) FindRegistrationByID(id string) (*models.EventRegistration, error) {
	if m.FindRegistrationByIDFunc != nil {
		return m.FindRegistrationByIDFunc(id)
	}
	return nil, nil
}

func (m *MockEventRepository) ExistsRegistrationByEventAndEmail(eventID string, email string) (bool, error) {
	if m.ExistsRegistrationFunc != nil {
		return m.ExistsRegistrationFunc(eventID, email)
	}
	return false, nil
}

func (m *MockEventRepository) UpdateRegistrationPayload(id string, payload string) error {
	if m.UpdateRegistrationPayloadFunc != nil {
		return m.UpdateRegistrationPayloadFunc(id, payload)
	}
	m.updatedPayload = payload
	return nil
}

// ── CRUD Tests ───────────────────────────────────────────────────────────────

func TestEventService_Create_Success(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindBySlugFunc: func(slug string) (*models.Event, error) {
			return nil, nil // Not exist
		},
	}
	svc := services.NewEventService(mockRepo, nil, &MockRedis{})

	req := dto.CreateEventRequest{
		Judul:     "Test Event",
		Deskripsi: "Deskripsi",
		Lokasi:    "Lokasi",
		Tanggal:   "2026-12-31T10:00:00Z",
	}

	err := svc.Create(req)
	assert.NoError(t, err)
}

func TestEventService_Create_DuplicateSlug(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindBySlugFunc: func(slug string) (*models.Event, error) {
			return &models.Event{ID: "existing-id", Slug: slug}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil, &MockRedis{})

	req := dto.CreateEventRequest{
		Judul:   "Test Event",
		Tanggal: "2026-12-31T10:00:00Z",
	}

	err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sudah ada")
}

func TestEventService_Create_InvalidDate(t *testing.T) {
	mockRepo := &MockEventRepository{
		FindAllFunc: func() ([]models.Event, error) {
			now := time.Now()
			return []models.Event{
				{ID: "1", Judul: "Test 1", Tanggal: now, CreatedAt: now, UpdatedAt: now},
				{ID: "2", Judul: "Test 2", Tanggal: now, CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	svc := services.NewEventService(mockRepo, nil, &MockRedis{})

	req := dto.CreateEventRequest{
		Judul:     "Test Event",
		Deskripsi: "Deskripsi",
		Lokasi:    "Lokasi",
		Tanggal:   "invalid-date",
	}

	err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "format tanggal tidak valid")
}
