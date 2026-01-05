package services

import (
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSeCsirtRepo struct {
	CreateFn  func(req dto.CreateSeCsirtRequest, id string) error
	GetAllFn  func() ([]dto.SeCsirtResponse, error)
	GetByIDFn func(id string) (*dto.SeCsirtResponse, error)
	UpdateFn  func(id string, req dto.SeCsirtResponse) error
	DeleteFn  func(id string) error
}

func (m *mockSeCsirtRepo) Create(req dto.CreateSeCsirtRequest, id string) error {
	return m.CreateFn(req, id)
}
func (m *mockSeCsirtRepo) GetAll() ([]dto.SeCsirtResponse, error) {
	return m.GetAllFn()
}
func (m *mockSeCsirtRepo) GetByID(id string) (*dto.SeCsirtResponse, error) {
	return m.GetByIDFn(id)
}
func (m *mockSeCsirtRepo) Update(id string, req dto.SeCsirtResponse) error {
	return m.UpdateFn(id, req)
}
func (m *mockSeCsirtRepo) Delete(id string) error {
	return m.DeleteFn(id)
}

func TestSeCsirtService_Create(t *testing.T) {
	repo := &mockSeCsirtRepo{
		CreateFn: func(req dto.CreateSeCsirtRequest, id string) error {
			assert.Equal(t, "Security Engine", *req.NamaSe)
			assert.Equal(t, "192.168.1.1", *req.IpSe)
			return nil
		},
	}

	service := NewSeCsirtService(repo)

	id, err := service.Create(dto.CreateSeCsirtRequest{
		NamaSe:      strPtr("Security Engine"),
		IpSe:        strPtr("192.168.1.1"),
		AsNumberSe:  strPtr("AS12345"),
		PengelolaSe: strPtr("BSSN"),
		FiturSe:     strPtr("Monitoring"),
		KategoriSe:  strPtr("Critical"),
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestSeCsirtService_GetAll(t *testing.T) {
	repo := &mockSeCsirtRepo{
		GetAllFn: func() ([]dto.SeCsirtResponse, error) {
			return []dto.SeCsirtResponse{
				{
					ID:     "1",
					NamaSe: "SE-1",
				},
			}, nil
		},
	}

	service := NewSeCsirtService(repo)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "SE-1", res[0].NamaSe)
}

func TestSeCsirtService_GetByID(t *testing.T) {
	repo := &mockSeCsirtRepo{
		GetByIDFn: func(id string) (*dto.SeCsirtResponse, error) {
			return &dto.SeCsirtResponse{
				ID:     id,
				NamaSe: "SE-Detail",
			}, nil
		},
	}

	service := NewSeCsirtService(repo)

	res, err := service.GetByID("1")

	assert.NoError(t, err)
	assert.Equal(t, "SE-Detail", res.NamaSe)
}

func TestSeCsirtService_Update(t *testing.T) {
	repo := &mockSeCsirtRepo{
		GetByIDFn: func(id string) (*dto.SeCsirtResponse, error) {
			return &dto.SeCsirtResponse{
				ID:     id,
				NamaSe: "Old SE",
				IpSe:   "10.0.0.1",
			}, nil
		},
		UpdateFn: func(id string, req dto.SeCsirtResponse) error {
			assert.Equal(t, "New SE", req.NamaSe)
			assert.Equal(t, "192.168.0.1", req.IpSe)
			return nil
		},
	}

	service := NewSeCsirtService(repo)

	err := service.Update("1", dto.UpdateSeCsirtRequest{
		NamaSe: strPtr("New SE"),
		IpSe:   strPtr("192.168.0.1"),
	})

	assert.NoError(t, err)
}

func TestSeCsirtService_Delete(t *testing.T) {
	repo := &mockSeCsirtRepo{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewSeCsirtService(repo)

	err := service.Delete("1")

	assert.NoError(t, err)
}
