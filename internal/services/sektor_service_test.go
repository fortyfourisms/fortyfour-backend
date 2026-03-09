package services

import (
	"encoding/json"
	"errors"
	"fortyfour-backend/internal/dto"
	"testing"
	"time"

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
					NamaSubSektor: "Logam",
					IDSektor:      "sektor-1",
				},
				{
					ID:            "sub-2",
					NamaSubSektor: "Permesinan & alat pertanian",
					IDSektor:      "sektor-1",
				},
				{
					ID:            "sub-3",
					NamaSubSektor: "Transportasi, maritim & pertahanan",
					IDSektor:      "sektor-1",
				},
				{
					ID:            "sub-4",
					NamaSubSektor: "Elektronika & telematika",
					IDSektor:      "sektor-1",
				},
			},
			CreatedAt: "2025-12-30 10:00:00",
			UpdatedAt: "2025-12-30 10:00:00",
		},
		{
			ID:         "sektor-2",
			NamaSektor: "Agro",
			SubSektor: []dto.SubSektorResponse{
				{
					ID:            "sub-5",
					NamaSubSektor: "Hasil hutan & perkebunan",
					IDSektor:      "sektor-2",
				},
				{
					ID:            "sub-6",
					NamaSubSektor: "Pangan & perikanan",
					IDSektor:      "sektor-2",
				},
				{
					ID:            "sub-7",
					NamaSubSektor: "Minuman, tembakau & bahan penyegar",
					IDSektor:      "sektor-2",
				},
				{
					ID:            "sub-8",
					NamaSubSektor: "Kemurgi, oleokimia & pakan",
					IDSektor:      "sektor-2",
				},
			},
			CreatedAt: "2025-12-30 10:00:00",
			UpdatedAt: "2025-12-30 10:00:00",
		},
		{
			ID:         "sektor-3",
			NamaSektor: "IKFT",
			SubSektor: []dto.SubSektorResponse{
				{
					ID:            "sub-9",
					NamaSubSektor: "Kimia hulu",
					IDSektor:      "sektor-3",
				},
				{
					ID:            "sub-10",
					NamaSubSektor: "Kimia hilir & farmasi",
					IDSektor:      "sektor-3",
				},
				{
					ID:            "sub-11",
					NamaSubSektor: "Semen, keramik & nonlogam",
					IDSektor:      "sektor-3",
				},
				{
					ID:            "sub-12",
					NamaSubSektor: "Tekstil, kulit & alas kaki",
					IDSektor:      "sektor-3",
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

	service := NewSektorService(repo, nil)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)

	// Verify ILMATE
	assert.Equal(t, "ILMATE", result[0].NamaSektor)
	assert.Len(t, result[0].SubSektor, 4)
	assert.Equal(t, "Logam", result[0].SubSektor[0].NamaSubSektor)
	assert.Equal(t, "Permesinan & alat pertanian", result[0].SubSektor[1].NamaSubSektor)
	assert.Equal(t, "Transportasi, maritim & pertahanan", result[0].SubSektor[2].NamaSubSektor)
	assert.Equal(t, "Elektronika & telematika", result[0].SubSektor[3].NamaSubSektor)

	// Verify Agro
	assert.Equal(t, "Agro", result[1].NamaSektor)
	assert.Len(t, result[1].SubSektor, 4)
	assert.Equal(t, "Hasil hutan & perkebunan", result[1].SubSektor[0].NamaSubSektor)
	assert.Equal(t, "Pangan & perikanan", result[1].SubSektor[1].NamaSubSektor)
	assert.Equal(t, "Minuman, tembakau & bahan penyegar", result[1].SubSektor[2].NamaSubSektor)
	assert.Equal(t, "Kemurgi, oleokimia & pakan", result[1].SubSektor[3].NamaSubSektor)

	// Verify IKFT
	assert.Equal(t, "IKFT", result[2].NamaSektor)
	assert.Len(t, result[2].SubSektor, 4)
	assert.Equal(t, "Kimia hulu", result[2].SubSektor[0].NamaSubSektor)
	assert.Equal(t, "Kimia hilir & farmasi", result[2].SubSektor[1].NamaSubSektor)
	assert.Equal(t, "Semen, keramik & nonlogam", result[2].SubSektor[2].NamaSubSektor)
	assert.Equal(t, "Tekstil, kulit & alas kaki", result[2].SubSektor[3].NamaSubSektor)
}

func TestGetAllSektor_EmptyResult(t *testing.T) {
	// Arrange
	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return []dto.SektorResponse{}, nil
		},
	}

	service := NewSektorService(repo, nil)

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

	service := NewSektorService(repo, nil)

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

func TestGetSektorByID_Success_ILMATE(t *testing.T) {
	// Arrange - Test dengan sektor ILMATE lengkap
	expectedSektor := &dto.SektorResponse{
		ID:         "sektor-1",
		NamaSektor: "ILMATE",
		SubSektor: []dto.SubSektorResponse{
			{
				ID:            "sub-1",
				NamaSubSektor: "Logam",
				IDSektor:      "sektor-1",
			},
			{
				ID:            "sub-2",
				NamaSubSektor: "Permesinan & alat pertanian",
				IDSektor:      "sektor-1",
			},
			{
				ID:            "sub-3",
				NamaSubSektor: "Transportasi, maritim & pertahanan",
				IDSektor:      "sektor-1",
			},
			{
				ID:            "sub-4",
				NamaSubSektor: "Elektronika & telematika",
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

	service := NewSektorService(repo, nil)

	// Act
	result, err := service.GetByID("sektor-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sektor-1", result.ID)
	assert.Equal(t, "ILMATE", result.NamaSektor)
	assert.Len(t, result.SubSektor, 4)
	assert.Equal(t, "Logam", result.SubSektor[0].NamaSubSektor)
	assert.Equal(t, "Permesinan & alat pertanian", result.SubSektor[1].NamaSubSektor)
	assert.Equal(t, "Transportasi, maritim & pertahanan", result.SubSektor[2].NamaSubSektor)
	assert.Equal(t, "Elektronika & telematika", result.SubSektor[3].NamaSubSektor)
}

func TestGetSektorByID_Success_Agro(t *testing.T) {
	// Arrange - Test dengan sektor Agro
	expectedSektor := &dto.SektorResponse{
		ID:         "sektor-2",
		NamaSektor: "Agro",
		SubSektor: []dto.SubSektorResponse{
			{
				ID:            "sub-5",
				NamaSubSektor: "Hasil hutan & perkebunan",
				IDSektor:      "sektor-2",
			},
			{
				ID:            "sub-6",
				NamaSubSektor: "Pangan & perikanan",
				IDSektor:      "sektor-2",
			},
			{
				ID:            "sub-7",
				NamaSubSektor: "Minuman, tembakau & bahan penyegar",
				IDSektor:      "sektor-2",
			},
			{
				ID:            "sub-8",
				NamaSubSektor: "Kemurgi, oleokimia & pakan",
				IDSektor:      "sektor-2",
			},
		},
		CreatedAt: "2025-12-30 10:00:00",
		UpdatedAt: "2025-12-30 10:00:00",
	}

	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			if id == "sektor-2" {
				return expectedSektor, nil
			}
			return nil, errors.New("sektor not found")
		},
	}

	service := NewSektorService(repo, nil)

	// Act
	result, err := service.GetByID("sektor-2")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sektor-2", result.ID)
	assert.Equal(t, "Agro", result.NamaSektor)
	assert.Len(t, result.SubSektor, 4)
	assert.Equal(t, "Hasil hutan & perkebunan", result.SubSektor[0].NamaSubSektor)
	assert.Equal(t, "Pangan & perikanan", result.SubSektor[1].NamaSubSektor)
}

func TestGetSektorByID_Success_IKFT(t *testing.T) {
	// Arrange - Test dengan sektor IKFT
	expectedSektor := &dto.SektorResponse{
		ID:         "sektor-3",
		NamaSektor: "IKFT",
		SubSektor: []dto.SubSektorResponse{
			{
				ID:            "sub-9",
				NamaSubSektor: "Kimia hulu",
				IDSektor:      "sektor-3",
			},
			{
				ID:            "sub-10",
				NamaSubSektor: "Kimia hilir & farmasi",
				IDSektor:      "sektor-3",
			},
			{
				ID:            "sub-11",
				NamaSubSektor: "Semen, keramik & nonlogam",
				IDSektor:      "sektor-3",
			},
			{
				ID:            "sub-12",
				NamaSubSektor: "Tekstil, kulit & alas kaki",
				IDSektor:      "sektor-3",
			},
		},
		CreatedAt: "2025-12-30 10:00:00",
		UpdatedAt: "2025-12-30 10:00:00",
	}

	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			if id == "sektor-3" {
				return expectedSektor, nil
			}
			return nil, errors.New("sektor not found")
		},
	}

	service := NewSektorService(repo, nil)

	// Act
	result, err := service.GetByID("sektor-3")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sektor-3", result.ID)
	assert.Equal(t, "IKFT", result.NamaSektor)
	assert.Len(t, result.SubSektor, 4)
	assert.Equal(t, "Kimia hulu", result.SubSektor[0].NamaSubSektor)
	assert.Equal(t, "Kimia hilir & farmasi", result.SubSektor[1].NamaSubSektor)
	assert.Equal(t, "Semen, keramik & nonlogam", result.SubSektor[2].NamaSubSektor)
	assert.Equal(t, "Tekstil, kulit & alas kaki", result.SubSektor[3].NamaSubSektor)
}

func TestGetSektorByID_NotFound(t *testing.T) {
	// Arrange
	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return nil, errors.New("sektor not found")
		},
	}

	service := NewSektorService(repo, nil)

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

	service := NewSektorService(repo, nil)

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

	service := NewSektorService(repo, nil)

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

	service := NewSektorService(repo, nil)

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
					NamaSubSektor: "Logam",
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

	service := NewSektorService(repo, nil)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify sub sektor data lengkap
	subSektor := result[0].SubSektor[0]
	assert.Equal(t, "sub-1", subSektor.ID)
	assert.Equal(t, "Logam", subSektor.NamaSubSektor)
	assert.Equal(t, "sektor-1", subSektor.IDSektor)
	assert.Equal(t, "ILMATE", subSektor.NamaSektor)
	assert.NotEmpty(t, subSektor.CreatedAt)
	assert.NotEmpty(t, subSektor.UpdatedAt)
}

/*
=====================================
 TEST ADDITIONAL VALIDATION
=====================================
*/

func TestGetAllSektor_VerifyAllThreeSektors(t *testing.T) {
	// Arrange - Pastikan semua 3 sektor ada
	expectedSektor := []dto.SektorResponse{
		{ID: "sektor-1", NamaSektor: "ILMATE"},
		{ID: "sektor-2", NamaSektor: "Agro"},
		{ID: "sektor-3", NamaSektor: "IKFT"},
	}

	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return expectedSektor, nil
		},
	}

	service := NewSektorService(repo, nil)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// Verify sektor names match master data
	sektorNames := []string{}
	for _, sektor := range result {
		sektorNames = append(sektorNames, sektor.NamaSektor)
	}

	assert.Contains(t, sektorNames, "ILMATE")
	assert.Contains(t, sektorNames, "Agro")
	assert.Contains(t, sektorNames, "IKFT")
}

/*
=====================================
 TEST CACHE — GetAll
=====================================
*/

func TestGetAllSektor_CacheHit_SkipRepo(t *testing.T) {
	rc := newSektorTestRedis()
	cached := []dto.SektorResponse{
		{ID: "sektor-cache", NamaSektor: "Dari Cache"},
	}
	setSektorCache(rc, keyList("sektor"), cached)

	repoCalled := false
	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			repoCalled = true
			return nil, errors.New("seharusnya tidak dipanggil")
		},
	}
	service := NewSektorService(repo, rc)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.False(t, repoCalled, "repo tidak boleh dipanggil saat cache hit")
	assert.Len(t, result, 1)
	assert.Equal(t, "Dari Cache", result[0].NamaSektor)
}

func TestGetAllSektor_CacheMiss_SetsCache(t *testing.T) {
	rc := newSektorTestRedis()
	repo := &mockSektorRepository{
		GetAllFn: func() ([]dto.SektorResponse, error) {
			return []dto.SektorResponse{
				{ID: "sektor-1", NamaSektor: "ILMATE"},
			}, nil
		},
	}
	service := NewSektorService(repo, rc)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	exists, _ := rc.Exists(keyList("sektor"))
	assert.True(t, exists, "data harus di-cache setelah GetAll")
}

/*
=====================================
 TEST CACHE — GetByID
=====================================
*/

func TestGetSektorByID_CacheHit_SkipRepo(t *testing.T) {
	rc := newSektorTestRedis()
	cached := dto.SektorResponse{ID: "sektor-1", NamaSektor: "ILMATE dari Cache"}
	setSektorCache(rc, keyDetail("sektor", "sektor-1"), cached)

	repoCalled := false
	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			repoCalled = true
			return nil, errors.New("seharusnya tidak dipanggil")
		},
	}
	service := NewSektorService(repo, rc)

	result, err := service.GetByID("sektor-1")

	assert.NoError(t, err)
	assert.False(t, repoCalled, "repo tidak boleh dipanggil saat cache hit")
	assert.Equal(t, "ILMATE dari Cache", result.NamaSektor)
}

func TestGetSektorByID_CacheMiss_SetsCache(t *testing.T) {
	rc := newSektorTestRedis()
	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return &dto.SektorResponse{ID: id, NamaSektor: "ILMATE"}, nil
		},
	}
	service := NewSektorService(repo, rc)

	_, err := service.GetByID("sektor-1")

	assert.NoError(t, err)

	exists, _ := rc.Exists(keyDetail("sektor", "sektor-1"))
	assert.True(t, exists, "data harus di-cache setelah GetByID")
}

func TestGetSektorByID_CacheMiss_DifferentIDs_CachedSeparately(t *testing.T) {
	rc := newSektorTestRedis()
	repo := &mockSektorRepository{
		GetByIDFn: func(id string) (*dto.SektorResponse, error) {
			return &dto.SektorResponse{ID: id, NamaSektor: "Sektor " + id}, nil
		},
	}
	service := NewSektorService(repo, rc)

	_, _ = service.GetByID("sektor-1")
	_, _ = service.GetByID("sektor-2")

	exists1, _ := rc.Exists(keyDetail("sektor", "sektor-1"))
	exists2, _ := rc.Exists(keyDetail("sektor", "sektor-2"))
	assert.True(t, exists1, "sektor-1 harus di-cache")
	assert.True(t, exists2, "sektor-2 harus di-cache terpisah")
}

/*
=====================================
 HELPERS REDIS UNTUK SEKTOR TEST
=====================================
*/

func newSektorTestRedis() *sektorTestRedis {
	return &sektorTestRedis{data: make(map[string]string)}
}

func setSektorCache(rc *sektorTestRedis, key string, value interface{}) {
	b, _ := json.Marshal(value)
	rc.data[key] = string(b)
}

type sektorTestRedis struct {
	data map[string]string
}

func (r *sektorTestRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}
func (r *sektorTestRedis) Get(key string) (string, error) {
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (r *sektorTestRedis) Delete(key string) error { delete(r.data, key); return nil }
func (r *sektorTestRedis) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}
func (r *sektorTestRedis) Scan(pattern string) ([]string, error) { return nil, nil }
func (r *sektorTestRedis) Close() error                          { return nil }
