package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
=====================================
 MOCK SUB SEKTOR REPOSITORY (STANDALONE)
=====================================
*/

type mockSubSektorRepositoryStandalone struct {
	GetAllFn        func() ([]dto.SubSektorResponse, error)
	GetByIDFn       func(id string) (*dto.SubSektorResponse, error)
	GetBySektorIDFn func(sektorID string) ([]dto.SubSektorResponse, error)
}

func (m *mockSubSektorRepositoryStandalone) GetAll() ([]dto.SubSektorResponse, error) {
	return m.GetAllFn()
}

func (m *mockSubSektorRepositoryStandalone) GetByID(id string) (*dto.SubSektorResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockSubSektorRepositoryStandalone) GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error) {
	return m.GetBySektorIDFn(sektorID)
}

/*
=====================================
 TEST GET ALL SUB SEKTOR
=====================================
*/

func TestGetAllSubSektor_Success(t *testing.T) {
	// Arrange
	expectedSubSektor := []dto.SubSektorResponse{
		{
			ID:            "sub-1",
			NamaSubSektor: "Elektronik",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-2",
			NamaSubSektor: "Otomotif",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-3",
			NamaSubSektor: "Agro Bisnis",
			IDSektor:      "sektor-2",
			NamaSektor:    "Industri Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetAllFn: func() ([]dto.SubSektorResponse, error) {
			return expectedSubSektor, nil
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Equal(t, "Elektronik", result[0].NamaSubSektor)
	assert.Equal(t, "Otomotif", result[1].NamaSubSektor)
	assert.Equal(t, "Agro Bisnis", result[2].NamaSubSektor)
}

func TestGetAllSubSektor_EmptyResult(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetAllFn: func() ([]dto.SubSektorResponse, error) {
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestGetAllSubSektor_RepositoryError(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetAllFn: func() ([]dto.SubSektorResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection error", err.Error())
}

/*
=====================================
 TEST GET SUB SEKTOR BY ID
=====================================
*/

func TestGetSubSektorByID_Success(t *testing.T) {
	// Arrange
	expectedSubSektor := &dto.SubSektorResponse{
		ID:            "sub-1",
		NamaSubSektor: "Elektronik",
		IDSektor:      "sektor-1",
		NamaSektor:    "ILMATE",
		CreatedAt:     "2025-12-30 10:00:00",
		UpdatedAt:     "2025-12-30 10:00:00",
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			if id == "sub-1" {
				return expectedSubSektor, nil
			}
			return nil, errors.New("sub sektor not found")
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetByID("sub-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sub-1", result.ID)
	assert.Equal(t, "Elektronik", result.NamaSubSektor)
	assert.Equal(t, "sektor-1", result.IDSektor)
	assert.Equal(t, "ILMATE", result.NamaSektor)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestGetSubSektorByID_NotFound(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("sub sektor not found")
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetByID("invalid-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sub sektor not found", err.Error())
}

func TestGetSubSektorByID_EmptyID(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			if id == "" {
				return nil, errors.New("id cannot be empty")
			}
			return nil, errors.New("sub sektor not found")
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetByID("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetSubSektorByID_RepositoryError(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("database timeout")
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetByID("sub-1")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database timeout", err.Error())
}

/*
=====================================
 TEST GET SUB SEKTOR BY SEKTOR ID
=====================================
*/

func TestGetSubSektorBySektorID_Success(t *testing.T) {
	// Arrange
	expectedSubSektor := []dto.SubSektorResponse{
		{
			ID:            "sub-1",
			NamaSubSektor: "Elektronik",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-2",
			NamaSubSektor: "Otomotif",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-3",
			NamaSubSektor: "Keamanan Siber",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			if sektorID == "sektor-1" {
				return expectedSubSektor, nil
			}
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetBySektorID("sektor-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Equal(t, "Elektronik", result[0].NamaSubSektor)
	assert.Equal(t, "Otomotif", result[1].NamaSubSektor)
	assert.Equal(t, "Keamanan Siber", result[2].NamaSubSektor)
	
	// Verify semua sub sektor punya IDSektor yang sama
	for _, sub := range result {
		assert.Equal(t, "sektor-1", sub.IDSektor)
		assert.Equal(t, "ILMATE", sub.NamaSektor)
	}
}

func TestGetSubSektorBySektorID_EmptyResult(t *testing.T) {
	// Arrange - Sektor tidak punya sub sektor
	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetBySektorID("sektor-empty")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestGetSubSektorBySektorID_InvalidSektorID(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetBySektorID("invalid-sektor-id")

	// Assert
	assert.NoError(t, err) // Tidak error, tapi result kosong
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestGetSubSektorBySektorID_RepositoryError(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetBySektorID("sektor-1")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection failed", err.Error())
}

func TestGetSubSektorBySektorID_MultipleSektors(t *testing.T) {
	// Arrange - Test dengan berbagai sektor
	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			switch sektorID {
			case "sektor-1": // ILMATE - 3 sub sektor
				return []dto.SubSektorResponse{
					{ID: "sub-1", NamaSubSektor: "Elektronik", IDSektor: "sektor-1", NamaSektor: "ILMATE"},
					{ID: "sub-2", NamaSubSektor: "Otomotif", IDSektor: "sektor-1", NamaSektor: "ILMATE"},
					{ID: "sub-3", NamaSubSektor: "Keamanan Siber", IDSektor: "sektor-1", NamaSektor: "ILMATE"},
				}, nil
			case "sektor-2": // Industri Agro - 4 sub sektor
				return []dto.SubSektorResponse{
					{ID: "sub-4", NamaSubSektor: "Agro Bisnis", IDSektor: "sektor-2", NamaSektor: "Industri Agro"},
					{ID: "sub-5", NamaSubSektor: "Konstruksi", IDSektor: "sektor-2", NamaSektor: "Industri Agro"},
					{ID: "sub-6", NamaSubSektor: "Jasa", IDSektor: "sektor-2", NamaSektor: "Industri Agro"},
					{ID: "sub-7", NamaSubSektor: "Surveyor", IDSektor: "sektor-2", NamaSektor: "Industri Agro"},
				}, nil
			case "sektor-3": // IKFT - 4 sub sektor
				return []dto.SubSektorResponse{
					{ID: "sub-8", NamaSubSektor: "Tekstil", IDSektor: "sektor-3", NamaSektor: "IKFT"},
					{ID: "sub-9", NamaSubSektor: "Kimia", IDSektor: "sektor-3", NamaSektor: "IKFT"},
					{ID: "sub-10", NamaSubSektor: "Kawasan Industri", IDSektor: "sektor-3", NamaSektor: "IKFT"},
					{ID: "sub-11", NamaSubSektor: "Farmasi", IDSektor: "sektor-3", NamaSektor: "IKFT"},
				}, nil
			default:
				return []dto.SubSektorResponse{}, nil
			}
		},
	}

	service := NewSubSektorService(repo)

	// Act & Assert - ILMATE
	result1, err1 := service.GetBySektorID("sektor-1")
	assert.NoError(t, err1)
	assert.Len(t, result1, 3)

	// Act & Assert - Industri Agro
	result2, err2 := service.GetBySektorID("sektor-2")
	assert.NoError(t, err2)
	assert.Len(t, result2, 4)

	// Act & Assert - IKFT
	result3, err3 := service.GetBySektorID("sektor-3")
	assert.NoError(t, err3)
	assert.Len(t, result3, 4)
}

/*
=====================================
 TEST DATA INTEGRITY
=====================================
*/

func TestSubSektor_DataIntegrity(t *testing.T) {
	// Arrange - Verify sub sektor selalu punya parent sektor info
	expectedSubSektor := []dto.SubSektorResponse{
		{
			ID:            "sub-1",
			NamaSubSektor: "Elektronik",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetAllFn: func() ([]dto.SubSektorResponse, error) {
			return expectedSubSektor, nil
		},
	}

	service := NewSubSektorService(repo)

	// Act
	result, err := service.GetAll()

	// Assert - Verify semua field terisi
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	
	subSektor := result[0]
	assert.NotEmpty(t, subSektor.ID, "ID should not be empty")
	assert.NotEmpty(t, subSektor.NamaSubSektor, "NamaSubSektor should not be empty")
	assert.NotEmpty(t, subSektor.IDSektor, "IDSektor should not be empty")
	assert.NotEmpty(t, subSektor.NamaSektor, "NamaSektor should not be empty")
	assert.NotEmpty(t, subSektor.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, subSektor.UpdatedAt, "UpdatedAt should not be empty")
}