package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
=====================================
 MOCK SEKTOR REPOSITORY
=====================================
*/

type mockSektorRepository struct {
	GetAllFn  func() ([]dto.SektorResponse, error)
	GetByIDFn func(id string) (*dto.SektorResponse, error)
}

func (m *mockSektorRepository) GetAll() ([]dto.SektorResponse, error) {
	return m.GetAllFn()
}

func (m *mockSektorRepository) GetByID(id string) (*dto.SektorResponse, error) {
	return m.GetByIDFn(id)
}

/*
=====================================
 TEST GET ALL SEKTOR
=====================================
*/

func TestGetAllSektor_Success(t *testing.T) {
	// Arrange
	expectedSektor := []dto.SektorResponse{
		{
			ID:         "sektor-1",
			NamaSektor: "ILMATE",
			SubSektor: []dto.SubSektorResponse{
				{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
				},
				{
					ID:            "sub-2",
					NamaSubSektor: "Otomotif",
					IDSektor:      "sektor-1",
				},
			},
			CreatedAt: "2025-12-30 10:00:00",
			UpdatedAt: "2025-12-30 10:00:00",
		},
		{
			ID:         "sektor-2",
			NamaSektor: "Industri Agro",
			SubSektor: []dto.SubSektorResponse{
				{
					ID:            "sub-3",
					NamaSubSektor: "Agro Bisnis",
					IDSektor:      "sektor-2",
				},
			},
			CreatedAt: "2025-12-30 10:00:00",
			UpdatedAt: "2025-12-30 10:00:00",
		},
	}

	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return expectedSektor, nil
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "ILMATE", result[0].NamaSektor)
	assert.Equal(t, "Industri Agro", result[1].NamaSektor)
	assert.Len(t, result[0].SubSektor, 2)
	assert.Len(t, result[1].SubSektor, 1)
}

func TestGetAllSektor_EmptyResult(t *testing.T) {
	// Arrange
	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return []dto.SektorResponse{}, nil
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestGetAllSektor_RepositoryError(t *testing.T) {
	// Arrange
	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection error", err.Error())
}

/*
=====================================
 TEST GET SEKTOR BY ID
=====================================
*/

func TestGetSektorByID_Success(t *testing.T) {
	// Arrange
	expectedSektor := &dto.SektorResponse{
		ID:         "sektor-1",
		NamaSektor: "ILMATE",
		SubSektor: []dto.SubSektorResponse{
			{
				ID:            "sub-1",
				NamaSubSektor: "Elektronik",
				IDSektor:      "sektor-1",
			},
			{
				ID:            "sub-2",
				NamaSubSektor: "Otomotif",
				IDSektor:      "sektor-1",
			},
			{
				ID:            "sub-3",
				NamaSubSektor: "Keamanan Siber",
				IDSektor:      "sektor-1",
			},
		},
		CreatedAt: "2025-12-30 10:00:00",
		UpdatedAt: "2025-12-30 10:00:00",
	}

	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			if id == "sektor-1" {
				return expectedSektor, nil
			}
			return nil, errors.New("sektor not found")
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetByID("sektor-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sektor-1", result.ID)
	assert.Equal(t, "ILMATE", result.NamaSektor)
	assert.Len(t, result.SubSektor, 3)
	assert.Equal(t, "Elektronik", result.SubSektor[0].NamaSubSektor)
	assert.Equal(t, "Otomotif", result.SubSektor[1].NamaSubSektor)
	assert.Equal(t, "Keamanan Siber", result.SubSektor[2].NamaSubSektor)
}

func TestGetSektorByID_NotFound(t *testing.T) {
	// Arrange
	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return nil, errors.New("sektor not found")
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetByID("invalid-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sektor not found", err.Error())
}

func TestGetSektorByID_WithoutSubSektor(t *testing.T) {
	// Arrange - Sektor tanpa sub sektor
	expectedSektor := &dto.SektorResponse{
		ID:         "sektor-empty",
		NamaSektor: "Sektor Kosong",
		SubSektor:  []dto.SubSektorResponse{}, // Empty sub sektor
		CreatedAt:  "2025-12-30 10:00:00",
		UpdatedAt:  "2025-12-30 10:00:00",
	}

	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return expectedSektor, nil
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetByID("sektor-empty")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Sektor Kosong", result.NamaSektor)
	assert.Len(t, result.SubSektor, 0)
}

func TestGetSektorByID_RepositoryError(t *testing.T) {
	// Arrange
	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return nil, errors.New("database timeout")
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetByID("sektor-1")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database timeout", err.Error())
}

/*
=====================================
 TEST EDGE CASES
=====================================
*/

func TestGetSektorByID_EmptyID(t *testing.T) {
	// Arrange
	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			if id == "" {
				return nil, errors.New("id cannot be empty")
			}
			return nil, errors.New("sektor not found")
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetByID("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAllSektor_VerifySubSektorData(t *testing.T) {
	// Arrange - Verify all sub sektor memiliki data lengkap
	expectedSektor := []dto.SektorResponse{
		{
			ID:         "sektor-1",
			NamaSektor: "ILMATE",
			SubSektor: []dto.SubSektorResponse{
				{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
					NamaSektor:    "ILMATE",
					CreatedAt:     "2025-12-30 10:00:00",
					UpdatedAt:     "2025-12-30 10:00:00",
				},
			},
		},
	}

	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return expectedSektor, nil
		},
	}

	service := NewSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify sub sektor data lengkap
	subSektor := result[0].SubSektor[0]
	assert.Equal(t, "sub-1", subSektor.ID)
	assert.Equal(t, "Elektronik", subSektor.NamaSubSektor)
	assert.Equal(t, "sektor-1", subSektor.IDSektor)
	assert.Equal(t, "ILMATE", subSektor.NamaSektor)
	assert.NotEmpty(t, subSektor.CreatedAt)
	assert.NotEmpty(t, subSektor.UpdatedAt)
}