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
	// Arrange - Menggunakan data master yang benar
	expectedSubSektor := []dto.SubSektorResponse{
		// ILMATE
		{
			ID:            "sub-1",
			NamaSubSektor: "Logam",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-2",
			NamaSubSektor: "Permesinan & alat pertanian",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-3",
			NamaSubSektor: "Transportasi, maritim & pertahanan",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-4",
			NamaSubSektor: "Elektronika & telematika",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		// Agro
		{
			ID:            "sub-5",
			NamaSubSektor: "Hasil hutan & perkebunan",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-6",
			NamaSubSektor: "Pangan & perikanan",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-7",
			NamaSubSektor: "Minuman, tembakau & bahan penyegar",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-8",
			NamaSubSektor: "Kemurgi, oleokimia & pakan",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		// IKFT
		{
			ID:            "sub-9",
			NamaSubSektor: "Kimia hulu",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-10",
			NamaSubSektor: "Kimia hilir & farmasi",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-11",
			NamaSubSektor: "Semen, keramik & nonlogam",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-12",
			NamaSubSektor: "Tekstil, kulit & alas kaki",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetAllFn: func() ([]dto.SubSektorResponse, error) {
			return expectedSubSektor, nil
		},
	}

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 12) // Total 12 sub sektor (4 ILMATE + 4 Agro + 4 IKFT)

	// Verify ILMATE sub-sektor
	assert.Equal(t, "Logam", result[0].NamaSubSektor)
	assert.Equal(t, "Permesinan & alat pertanian", result[1].NamaSubSektor)
	assert.Equal(t, "Transportasi, maritim & pertahanan", result[2].NamaSubSektor)
	assert.Equal(t, "Elektronika & telematika", result[3].NamaSubSektor)

	// Verify Agro sub-sektor
	assert.Equal(t, "Hasil hutan & perkebunan", result[4].NamaSubSektor)
	assert.Equal(t, "Pangan & perikanan", result[5].NamaSubSektor)
	assert.Equal(t, "Minuman, tembakau & bahan penyegar", result[6].NamaSubSektor)
	assert.Equal(t, "Kemurgi, oleokimia & pakan", result[7].NamaSubSektor)

	// Verify IKFT sub-sektor
	assert.Equal(t, "Kimia hulu", result[8].NamaSubSektor)
	assert.Equal(t, "Kimia hilir & farmasi", result[9].NamaSubSektor)
	assert.Equal(t, "Semen, keramik & nonlogam", result[10].NamaSubSektor)
	assert.Equal(t, "Tekstil, kulit & alas kaki", result[11].NamaSubSektor)
}

func TestGetAllSubSektor_EmptyResult(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetAllFn: func() ([]dto.SubSektorResponse, error) {
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo, nil)

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

	service := NewSubSektorService(repo, nil)

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

func TestGetSubSektorByID_Success_ILMATE(t *testing.T) {
	// Arrange - Test dengan sub-sektor dari ILMATE
	expectedSubSektor := &dto.SubSektorResponse{
		ID:            "sub-1",
		NamaSubSektor: "Logam",
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

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetByID("sub-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sub-1", result.ID)
	assert.Equal(t, "Logam", result.NamaSubSektor)
	assert.Equal(t, "sektor-1", result.IDSektor)
	assert.Equal(t, "ILMATE", result.NamaSektor)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestGetSubSektorByID_Success_Agro(t *testing.T) {
	// Arrange - Test dengan sub-sektor dari Agro
	expectedSubSektor := &dto.SubSektorResponse{
		ID:            "sub-5",
		NamaSubSektor: "Hasil hutan & perkebunan",
		IDSektor:      "sektor-2",
		NamaSektor:    "Agro",
		CreatedAt:     "2025-12-30 10:00:00",
		UpdatedAt:     "2025-12-30 10:00:00",
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			if id == "sub-5" {
				return expectedSubSektor, nil
			}
			return nil, errors.New("sub sektor not found")
		},
	}

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetByID("sub-5")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sub-5", result.ID)
	assert.Equal(t, "Hasil hutan & perkebunan", result.NamaSubSektor)
	assert.Equal(t, "sektor-2", result.IDSektor)
	assert.Equal(t, "Agro", result.NamaSektor)
}

func TestGetSubSektorByID_Success_IKFT(t *testing.T) {
	// Arrange - Test dengan sub-sektor dari IKFT
	expectedSubSektor := &dto.SubSektorResponse{
		ID:            "sub-9",
		NamaSubSektor: "Kimia hulu",
		IDSektor:      "sektor-3",
		NamaSektor:    "IKFT",
		CreatedAt:     "2025-12-30 10:00:00",
		UpdatedAt:     "2025-12-30 10:00:00",
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			if id == "sub-9" {
				return expectedSubSektor, nil
			}
			return nil, errors.New("sub sektor not found")
		},
	}

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetByID("sub-9")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sub-9", result.ID)
	assert.Equal(t, "Kimia hulu", result.NamaSubSektor)
	assert.Equal(t, "sektor-3", result.IDSektor)
	assert.Equal(t, "IKFT", result.NamaSektor)
}

func TestGetSubSektorByID_NotFound(t *testing.T) {
	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("sub sektor not found")
		},
	}

	service := NewSubSektorService(repo, nil)

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

	service := NewSubSektorService(repo, nil)

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

	service := NewSubSektorService(repo, nil)

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

func TestGetSubSektorBySektorID_Success_ILMATE(t *testing.T) {
	// Arrange - ILMATE dengan 4 sub-sektor
	expectedSubSektor := []dto.SubSektorResponse{
		{
			ID:            "sub-1",
			NamaSubSektor: "Logam",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-2",
			NamaSubSektor: "Permesinan & alat pertanian",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-3",
			NamaSubSektor: "Transportasi, maritim & pertahanan",
			IDSektor:      "sektor-1",
			NamaSektor:    "ILMATE",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-4",
			NamaSubSektor: "Elektronika & telematika",
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

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetBySektorID("sektor-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 4)
	assert.Equal(t, "Logam", result[0].NamaSubSektor)
	assert.Equal(t, "Permesinan & alat pertanian", result[1].NamaSubSektor)
	assert.Equal(t, "Transportasi, maritim & pertahanan", result[2].NamaSubSektor)
	assert.Equal(t, "Elektronika & telematika", result[3].NamaSubSektor)

	// Verify semua sub sektor punya IDSektor yang sama
	for _, sub := range result {
		assert.Equal(t, "sektor-1", sub.IDSektor)
		assert.Equal(t, "ILMATE", sub.NamaSektor)
	}
}

func TestGetSubSektorBySektorID_Success_Agro(t *testing.T) {
	// Arrange - Agro dengan 4 sub-sektor
	expectedSubSektor := []dto.SubSektorResponse{
		{
			ID:            "sub-5",
			NamaSubSektor: "Hasil hutan & perkebunan",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-6",
			NamaSubSektor: "Pangan & perikanan",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-7",
			NamaSubSektor: "Minuman, tembakau & bahan penyegar",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-8",
			NamaSubSektor: "Kemurgi, oleokimia & pakan",
			IDSektor:      "sektor-2",
			NamaSektor:    "Agro",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			if sektorID == "sektor-2" {
				return expectedSubSektor, nil
			}
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetBySektorID("sektor-2")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 4)
	assert.Equal(t, "Hasil hutan & perkebunan", result[0].NamaSubSektor)
	assert.Equal(t, "Pangan & perikanan", result[1].NamaSubSektor)

	// Verify semua sub sektor punya IDSektor yang sama
	for _, sub := range result {
		assert.Equal(t, "sektor-2", sub.IDSektor)
		assert.Equal(t, "Agro", sub.NamaSektor)
	}
}

func TestGetSubSektorBySektorID_Success_IKFT(t *testing.T) {
	// Arrange - IKFT dengan 4 sub-sektor
	expectedSubSektor := []dto.SubSektorResponse{
		{
			ID:            "sub-9",
			NamaSubSektor: "Kimia hulu",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-10",
			NamaSubSektor: "Kimia hilir & farmasi",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-11",
			NamaSubSektor: "Semen, keramik & nonlogam",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
		{
			ID:            "sub-12",
			NamaSubSektor: "Tekstil, kulit & alas kaki",
			IDSektor:      "sektor-3",
			NamaSektor:    "IKFT",
			CreatedAt:     "2025-12-30 10:00:00",
			UpdatedAt:     "2025-12-30 10:00:00",
		},
	}

	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			if sektorID == "sektor-3" {
				return expectedSubSektor, nil
			}
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetBySektorID("sektor-3")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 4)
	assert.Equal(t, "Kimia hulu", result[0].NamaSubSektor)
	assert.Equal(t, "Kimia hilir & farmasi", result[1].NamaSubSektor)
	assert.Equal(t, "Semen, keramik & nonlogam", result[2].NamaSubSektor)
	assert.Equal(t, "Tekstil, kulit & alas kaki", result[3].NamaSubSektor)
}

func TestGetSubSektorBySektorID_EmptyResult(t *testing.T) {
	// Arrange - Sektor tidak punya sub sektor
	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			return []dto.SubSektorResponse{}, nil
		},
	}

	service := NewSubSektorService(repo, nil)

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

	service := NewSubSektorService(repo, nil)

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

	service := NewSubSektorService(repo, nil)

	// Act
	result, err := service.GetBySektorID("sektor-1")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection failed", err.Error())
}

func TestGetSubSektorBySektorID_MultipleSektors(t *testing.T) {
	// Arrange - Test dengan berbagai sektor (menggunakan data master yang benar)
	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			switch sektorID {
			case "sektor-1": // ILMATE - 4 sub sektor
				return []dto.SubSektorResponse{
					{ID: "sub-1", NamaSubSektor: "Logam", IDSektor: "sektor-1", NamaSektor: "ILMATE"},
					{ID: "sub-2", NamaSubSektor: "Permesinan & alat pertanian", IDSektor: "sektor-1", NamaSektor: "ILMATE"},
					{ID: "sub-3", NamaSubSektor: "Transportasi, maritim & pertahanan", IDSektor: "sektor-1", NamaSektor: "ILMATE"},
					{ID: "sub-4", NamaSubSektor: "Elektronika & telematika", IDSektor: "sektor-1", NamaSektor: "ILMATE"},
				}, nil
			case "sektor-2": // Agro - 4 sub sektor
				return []dto.SubSektorResponse{
					{ID: "sub-5", NamaSubSektor: "Hasil hutan & perkebunan", IDSektor: "sektor-2", NamaSektor: "Agro"},
					{ID: "sub-6", NamaSubSektor: "Pangan & perikanan", IDSektor: "sektor-2", NamaSektor: "Agro"},
					{ID: "sub-7", NamaSubSektor: "Minuman, tembakau & bahan penyegar", IDSektor: "sektor-2", NamaSektor: "Agro"},
					{ID: "sub-8", NamaSubSektor: "Kemurgi, oleokimia & pakan", IDSektor: "sektor-2", NamaSektor: "Agro"},
				}, nil
			case "sektor-3": // IKFT - 4 sub sektor
				return []dto.SubSektorResponse{
					{ID: "sub-9", NamaSubSektor: "Kimia hulu", IDSektor: "sektor-3", NamaSektor: "IKFT"},
					{ID: "sub-10", NamaSubSektor: "Kimia hilir & farmasi", IDSektor: "sektor-3", NamaSektor: "IKFT"},
					{ID: "sub-11", NamaSubSektor: "Semen, keramik & nonlogam", IDSektor: "sektor-3", NamaSektor: "IKFT"},
					{ID: "sub-12", NamaSubSektor: "Tekstil, kulit & alas kaki", IDSektor: "sektor-3", NamaSektor: "IKFT"},
				}, nil
			default:
				return []dto.SubSektorResponse{}, nil
			}
		},
	}

	service := NewSubSektorService(repo, nil)

	// Act & Assert - ILMATE
	result1, err1 := service.GetBySektorID("sektor-1")
	assert.NoError(t, err1)
	assert.Len(t, result1, 4)

	// Act & Assert - Agro
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
			NamaSubSektor: "Logam",
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

	service := NewSubSektorService(repo, nil)

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

/*
=====================================
 TEST MASTER DATA COMPLIANCE
=====================================
*/

func TestSubSektor_VerifyMasterDataCompliance(t *testing.T) {
	// Test ini memastikan bahwa nama-nama sub sektor sesuai dengan master data
	expectedSubSektorNames := map[string][]string{
		"ILMATE": {
			"Logam",
			"Permesinan & alat pertanian",
			"Transportasi, maritim & pertahanan",
			"Elektronika & telematika",
		},
		"Agro": {
			"Hasil hutan & perkebunan",
			"Pangan & perikanan",
			"Minuman, tembakau & bahan penyegar",
			"Kemurgi, oleokimia & pakan",
		},
		"IKFT": {
			"Kimia hulu",
			"Kimia hilir & farmasi",
			"Semen, keramik & nonlogam",
			"Tekstil, kulit & alas kaki",
		},
	}

	// Arrange
	repo := &mockSubSektorRepositoryStandalone{
		GetBySektorIDFn: func(sektorID string) ([]dto.SubSektorResponse, error) {
			switch sektorID {
			case "sektor-1":
				return []dto.SubSektorResponse{
					{NamaSubSektor: "Logam", NamaSektor: "ILMATE"},
					{NamaSubSektor: "Permesinan & alat pertanian", NamaSektor: "ILMATE"},
					{NamaSubSektor: "Transportasi, maritim & pertahanan", NamaSektor: "ILMATE"},
					{NamaSubSektor: "Elektronika & telematika", NamaSektor: "ILMATE"},
				}, nil
			case "sektor-2":
				return []dto.SubSektorResponse{
					{NamaSubSektor: "Hasil hutan & perkebunan", NamaSektor: "Agro"},
					{NamaSubSektor: "Pangan & perikanan", NamaSektor: "Agro"},
					{NamaSubSektor: "Minuman, tembakau & bahan penyegar", NamaSektor: "Agro"},
					{NamaSubSektor: "Kemurgi, oleokimia & pakan", NamaSektor: "Agro"},
				}, nil
			case "sektor-3":
				return []dto.SubSektorResponse{
					{NamaSubSektor: "Kimia hulu", NamaSektor: "IKFT"},
					{NamaSubSektor: "Kimia hilir & farmasi", NamaSektor: "IKFT"},
					{NamaSubSektor: "Semen, keramik & nonlogam", NamaSektor: "IKFT"},
					{NamaSubSektor: "Tekstil, kulit & alas kaki", NamaSektor: "IKFT"},
				}, nil
			default:
				return []dto.SubSektorResponse{}, nil
			}
		},
	}

	service := NewSubSektorService(repo, nil)

	// Test untuk setiap sektor
	testCases := []struct {
		sektorID   string
		sektorName string
	}{
		{"sektor-1", "ILMATE"},
		{"sektor-2", "Agro"},
		{"sektor-3", "IKFT"},
	}

	for _, tc := range testCases {
		t.Run("Verify_"+tc.sektorName, func(t *testing.T) {
			result, err := service.GetBySektorID(tc.sektorID)

			assert.NoError(t, err)
			assert.Len(t, result, 4, "Setiap sektor harus punya 4 sub-sektor")

			// Verify nama-nama sub sektor sesuai master data
			expectedNames := expectedSubSektorNames[tc.sektorName]
			for i, subSektor := range result {
				assert.Equal(t, expectedNames[i], subSektor.NamaSubSektor,
					"Sub-sektor %s tidak sesuai master data", tc.sektorName)
			}
		})
	}
}
