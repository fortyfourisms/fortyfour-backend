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
 MOCK PERUSAHAAN REPOSITORY
=====================================
*/

type mockPerusahaanRepository struct {
	CreateFn     func(req dto.CreatePerusahaanRequest, id string) error
	GetByIDFn    func(id string) (*dto.PerusahaanResponse, error)
	GetByNamaFn  func(nama string) (*dto.PerusahaanResponse, error)
	GetAllFn     func() ([]dto.PerusahaanResponse, error)
	UpdateFn     func(id string, perusahaan dto.PerusahaanResponse) error
	DeleteFn     func(id string) error
}

func (m *mockPerusahaanRepository) Create(req dto.CreatePerusahaanRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockPerusahaanRepository) GetByID(id string) (*dto.PerusahaanResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockPerusahaanRepository) GetAll() ([]dto.PerusahaanResponse, error) {
	return m.GetAllFn()
}

func (m *mockPerusahaanRepository) Update(id string, perusahaan dto.PerusahaanResponse) error {
	return m.UpdateFn(id, perusahaan)
}

func (m *mockPerusahaanRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

func (m *mockPerusahaanRepository) GetByNama(nama string) (*dto.PerusahaanResponse, error) {
	if m.GetByNamaFn != nil {
		return m.GetByNamaFn(nama)
	}
	return nil, errors.New("perusahaan not found")
}

/*
=====================================
 MOCK SUB SEKTOR REPOSITORY
=====================================
*/

type mockSubSektorRepository struct {
	GetByIDFn       func(id string) (*dto.SubSektorResponse, error)
	GetAllFn        func() ([]dto.SubSektorResponse, error)
	GetBySektorIDFn func(sektorID string) ([]dto.SubSektorResponse, error)
}

func (m *mockSubSektorRepository) GetByID(id string) (*dto.SubSektorResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockSubSektorRepository) GetAll() ([]dto.SubSektorResponse, error) {
	return m.GetAllFn()
}

func (m *mockSubSektorRepository) GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error) {
	return m.GetBySektorIDFn(sektorID)
}

/*
=====================================
 TEST CREATE PERUSAHAAN
=====================================
*/

func TestCreatePerusahaan_Success(t *testing.T) {
	nama := "PT Teknologi Maju"
	idSubSektor := "sub-sektor-id-123"

	perusahaanRepo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: nama,
				SubSektor: &dto.SubSektorResponse{
					ID:            idSubSektor,
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-id-123",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return &dto.SubSektorResponse{
				ID:            idSubSektor,
				NamaSubSektor: "Elektronik",
				IDSektor:      "sektor-id-123",
				NamaSektor:    "ILMATE",
			}, nil
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.NamaPerusahaan)
	assert.NotNil(t, result.SubSektor)
	assert.Equal(t, "Elektronik", result.SubSektor.NamaSubSektor)
}

func TestCreatePerusahaan_ValidationFailed_NamaKosong(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{}
	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	idSubSektor := "sub-sektor-id-123"
	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: nil,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "nama_perusahaan wajib diisi", err.Error())
}

func TestCreatePerusahaan_ValidationFailed_SubSektorKosong(t *testing.T) {
    nama := "PT ABC"
    
    perusahaanRepo := &mockPerusahaanRepository{
        CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
            return nil
        },
        GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
            return &dto.PerusahaanResponse{
                ID:             id,
                NamaPerusahaan: nama,
                SubSektor:      nil,
            }, nil
        },
    }
    
    subSektorRepo := &mockSubSektorRepository{}
    service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)
    
    req := dto.CreatePerusahaanRequest{
        NamaPerusahaan: &nama,
        IDSubSektor:    nil,
    }
    
    result, err := service.Create(req)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, nama, result.NamaPerusahaan)
    assert.Nil(t, result.SubSektor)
}

func TestCreatePerusahaan_SubSektorNotFound(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	nama := "PT ABC"
	idSubSektor := "invalid-sub-sektor-id"

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sub sektor tidak ditemukan", err.Error())
}

func TestCreatePerusahaan_RepoFailed(t *testing.T) {
	nama := "PT ABC"
	idSubSektor := "sub-sektor-id-123"

	perusahaanRepo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error {
			return errors.New("db error")
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return &dto.SubSektorResponse{
				ID:            idSubSektor,
				NamaSubSektor: "Elektronik",
			}, nil
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	req := dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    &idSubSektor,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST GET
=====================================
*/

func TestGetAllPerusahaan_Success(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetAllFn: func() ([]dto.PerusahaanResponse, error) {
			return []dto.PerusahaanResponse{
				{
					ID:             "1",
					NamaPerusahaan: "PT A",
					SubSektor: &dto.SubSektorResponse{
						ID:            "sub-1",
						NamaSubSektor: "Elektronik",
						NamaSektor:    "ILMATE",
					},
				},
				{
					ID:             "2",
					NamaPerusahaan: "PT B",
					SubSektor: &dto.SubSektorResponse{
						ID:            "sub-2",
						NamaSubSektor: "Tekstil",
						NamaSektor:    "IKFT",
					},
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, "Elektronik", data[0].SubSektor.NamaSubSektor)
	assert.Equal(t, "Tekstil", data[1].SubSektor.NamaSubSektor)
}

func TestGetPerusahaanByID_Success(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	result, err := service.GetByID("123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PT A", result.NamaPerusahaan)
	assert.Equal(t, "Elektronik", result.SubSektor.NamaSubSektor)
}

/*
=====================================
 TEST UPDATE PERUSAHAAN
=====================================
*/

func TestUpdatePerusahaan_Success(t *testing.T) {
	oldNama := "PT Lama"
	newNama := "PT Baru"

	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: oldNama,
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
		UpdateFn: func(id string, perusahaan dto.PerusahaanResponse) error {
			return nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	req := dto.UpdatePerusahaanRequest{
		NamaPerusahaan: &newNama,
	}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.Equal(t, newNama, result.NamaPerusahaan)
	assert.NotNil(t, result.SubSektor)
	assert.Equal(t, "Elektronik", result.SubSektor.NamaSubSektor)
}

func TestUpdatePerusahaan_UpdateSubSektor_Success(t *testing.T) {
	oldSubSektorID := "sub-1"
	newSubSektorID := "sub-2"

	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				SubSektor: &dto.SubSektorResponse{
					ID:            oldSubSektorID,
					NamaSubSektor: "Elektronik",
					IDSektor:      "sektor-1",
					NamaSektor:    "ILMATE",
				},
			}, nil
		},
		UpdateFn: func(id string, perusahaan dto.PerusahaanResponse) error {
			return nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return &dto.SubSektorResponse{
				ID:            newSubSektorID,
				NamaSubSektor: "Otomotif",
				IDSektor:      "sektor-1",
				NamaSektor:    "ILMATE",
			}, nil
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	req := dto.UpdatePerusahaanRequest{
		IDSubSektor: &newSubSektorID,
	}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.NotNil(t, result.SubSektor)
	assert.Equal(t, newSubSektorID, result.SubSektor.ID)
	assert.Equal(t, "Otomotif", result.SubSektor.NamaSubSektor)
}

func TestUpdatePerusahaan_InvalidSubSektor(t *testing.T) {
	invalidSubSektorID := "invalid-sub-sektor"

	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{
				ID:             id,
				NamaPerusahaan: "PT A",
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-1",
					NamaSubSektor: "Elektronik",
				},
			}, nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{
		GetByIDFn: func(id string) (*dto.SubSektorResponse, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	req := dto.UpdatePerusahaanRequest{
		IDSubSektor: &invalidSubSektorID,
	}

	result, err := service.Update("123", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sub sektor tidak ditemukan", err.Error())
}

/*
=====================================
 TEST DELETE PERUSAHAAN
=====================================
*/

func TestDeletePerusahaan_Success(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	err := service.Delete("123")

	assert.NoError(t, err)
}

func TestDeletePerusahaan_Failed(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		DeleteFn: func(id string) error {
			return errors.New("db error")
		},
	}

	subSektorRepo := &mockSubSektorRepository{}
	service := NewPerusahaanService(perusahaanRepo, subSektorRepo, nil)

	err := service.Delete("123")

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
/*
=====================================
 TEST CREATE — tambahan
=====================================
*/

func TestCreatePerusahaan_NamaKosong_Error(t *testing.T) {
	nama := "   "
	idSubSektor := "sub-1"
	service := NewPerusahaanService(&mockPerusahaanRepository{}, &mockSubSektorRepository{}, nil)

	result, err := service.Create(dto.CreatePerusahaanRequest{
		NamaPerusahaan: &nama,
		IDSubSektor:    &idSubSektor,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "nama_perusahaan wajib diisi")
	assert.Nil(t, result)
}

func TestCreatePerusahaan_GetByIDAfterCreate_Error(t *testing.T) {
	nama := "PT Test"
	perusahaanRepo := &mockPerusahaanRepository{
		CreateFn: func(req dto.CreatePerusahaanRequest, id string) error { return nil },
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return nil, errors.New("tidak ditemukan setelah create")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	result, err := service.Create(dto.CreatePerusahaanRequest{NamaPerusahaan: &nama})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreatePerusahaan_InvalidatesListCache(t *testing.T) {
	rc := newPerusahaanTestRedis()
	setPerusahaanCache(rc, keyList("perusahaan"), []dto.PerusahaanResponse{{ID: "lama"}})

	nama := "PT Baru"
	perusahaanRepo := &mockPerusahaanRepository{
		CreateFn:  func(req dto.CreatePerusahaanRequest, id string) error { return nil },
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) { return &dto.PerusahaanResponse{ID: id, NamaPerusahaan: nama}, nil },
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, rc)

	_, err := service.Create(dto.CreatePerusahaanRequest{NamaPerusahaan: &nama})

	assert.NoError(t, err)
	exists, _ := rc.Exists(keyList("perusahaan"))
	assert.False(t, exists, "cache list harus dihapus setelah create")
}

/*
=====================================
 TEST GET ALL — cache
=====================================
*/

func TestGetAllPerusahaan_CacheHit_SkipRepo(t *testing.T) {
	rc := newPerusahaanTestRedis()
	cached := []dto.PerusahaanResponse{{ID: "cache-1", NamaPerusahaan: "PT Cache"}}
	setPerusahaanCache(rc, keyList("perusahaan"), cached)

	repoCalled := false
	perusahaanRepo := &mockPerusahaanRepository{
		GetAllFn: func() ([]dto.PerusahaanResponse, error) {
			repoCalled = true
			return nil, errors.New("tidak boleh dipanggil")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, rc)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.False(t, repoCalled)
	assert.Len(t, result, 1)
	assert.Equal(t, "PT Cache", result[0].NamaPerusahaan)
}

func TestGetAllPerusahaan_CacheMiss_SetsCache(t *testing.T) {
	rc := newPerusahaanTestRedis()
	perusahaanRepo := &mockPerusahaanRepository{
		GetAllFn: func() ([]dto.PerusahaanResponse, error) {
			return []dto.PerusahaanResponse{{ID: "db-1", NamaPerusahaan: "PT DB"}}, nil
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, rc)

	_, err := service.GetAll()

	assert.NoError(t, err)
	exists, _ := rc.Exists(keyList("perusahaan"))
	assert.True(t, exists, "data harus di-cache setelah GetAll")
}

func TestGetAllPerusahaan_RepoError(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetAllFn: func() ([]dto.PerusahaanResponse, error) {
			return nil, errors.New("db timeout")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	result, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST GET BY ID — cache & error
=====================================
*/

func TestGetPerusahaanByID_CacheHit_SkipRepo(t *testing.T) {
	rc := newPerusahaanTestRedis()
	setPerusahaanCache(rc, keyDetail("perusahaan", "p-1"), dto.PerusahaanResponse{ID: "p-1", NamaPerusahaan: "PT Cache"})

	repoCalled := false
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			repoCalled = true
			return nil, errors.New("tidak boleh dipanggil")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, rc)

	result, err := service.GetByID("p-1")

	assert.NoError(t, err)
	assert.False(t, repoCalled)
	assert.Equal(t, "PT Cache", result.NamaPerusahaan)
}

func TestGetPerusahaanByID_CacheMiss_SetsCache(t *testing.T) {
	rc := newPerusahaanTestRedis()
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: id, NamaPerusahaan: "PT DB"}, nil
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, rc)

	_, err := service.GetByID("p-1")

	assert.NoError(t, err)
	exists, _ := rc.Exists(keyDetail("perusahaan", "p-1"))
	assert.True(t, exists, "data harus di-cache setelah GetByID")
}

func TestGetPerusahaanByID_RepoError(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return nil, errors.New("not found")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	result, err := service.GetByID("invalid")

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST GET BY NAMA
=====================================
*/

func TestGetByNama_Success(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByNamaFn: func(nama string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: "p-1", NamaPerusahaan: nama}, nil
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	result, err := service.GetByNama("PT ABC")

	assert.NoError(t, err)
	assert.Equal(t, "PT ABC", result.NamaPerusahaan)
}

func TestGetByNama_NotFound(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByNamaFn: func(nama string) (*dto.PerusahaanResponse, error) {
			return nil, errors.New("perusahaan not found")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	result, err := service.GetByNama("PT Tidak Ada")

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST UPDATE — tambahan
=====================================
*/

func TestUpdatePerusahaan_GetByIDError(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return nil, errors.New("perusahaan not found")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	nama := "PT Baru"
	result, err := service.Update("invalid", dto.UpdatePerusahaanRequest{NamaPerusahaan: &nama})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdatePerusahaan_RepoUpdateError(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: id, NamaPerusahaan: "PT Lama"}, nil
		},
		UpdateFn: func(id string, p dto.PerusahaanResponse) error {
			return errors.New("update failed")
		},
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	nama := "PT Baru"
	result, err := service.Update("p-1", dto.UpdatePerusahaanRequest{NamaPerusahaan: &nama})

	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
	assert.Nil(t, result)
}

func TestUpdatePerusahaan_UpdateAllOptionalFields(t *testing.T) {
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: id, NamaPerusahaan: "PT Lama"}, nil
		},
		UpdateFn: func(id string, p dto.PerusahaanResponse) error { return nil },
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, nil)

	alamat := "Jl. Test 1"
	telepon := "08123456789"
	email := "test@test.com"
	website := "https://test.com"
	photo := "photo.jpg"

	result, err := service.Update("p-1", dto.UpdatePerusahaanRequest{
		Alamat:  &alamat,
		Telepon: &telepon,
		Email:   &email,
		Website: &website,
		Photo:   &photo,
	})

	assert.NoError(t, err)
	assert.Equal(t, alamat, result.Alamat)
	assert.Equal(t, telepon, result.Telepon)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, website, result.Website)
	assert.Equal(t, photo, result.Photo)
}

func TestUpdatePerusahaan_InvalidatesCache(t *testing.T) {
	rc := newPerusahaanTestRedis()
	setPerusahaanCache(rc, keyDetail("perusahaan", "p-1"), dto.PerusahaanResponse{ID: "p-1"})
	setPerusahaanCache(rc, keyList("perusahaan"), []dto.PerusahaanResponse{{ID: "p-1"}})

	nama := "PT Baru"
	perusahaanRepo := &mockPerusahaanRepository{
		GetByIDFn: func(id string) (*dto.PerusahaanResponse, error) {
			return &dto.PerusahaanResponse{ID: id, NamaPerusahaan: "PT Lama"}, nil
		},
		UpdateFn: func(id string, p dto.PerusahaanResponse) error { return nil },
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, rc)

	_, err := service.Update("p-1", dto.UpdatePerusahaanRequest{NamaPerusahaan: &nama})

	assert.NoError(t, err)
	existsDetail, _ := rc.Exists(keyDetail("perusahaan", "p-1"))
	existsList, _ := rc.Exists(keyList("perusahaan"))
	assert.False(t, existsDetail, "cache detail harus dihapus setelah update")
	assert.False(t, existsList, "cache list harus dihapus setelah update")
}

/*
=====================================
 TEST DELETE — cache
=====================================
*/

func TestDeletePerusahaan_InvalidatesCache(t *testing.T) {
	rc := newPerusahaanTestRedis()
	setPerusahaanCache(rc, keyDetail("perusahaan", "p-1"), dto.PerusahaanResponse{ID: "p-1"})
	setPerusahaanCache(rc, keyList("perusahaan"), []dto.PerusahaanResponse{{ID: "p-1"}})

	perusahaanRepo := &mockPerusahaanRepository{
		DeleteFn: func(id string) error { return nil },
	}
	service := NewPerusahaanService(perusahaanRepo, &mockSubSektorRepository{}, rc)

	err := service.Delete("p-1")

	assert.NoError(t, err)
	existsDetail, _ := rc.Exists(keyDetail("perusahaan", "p-1"))
	existsList, _ := rc.Exists(keyList("perusahaan"))
	assert.False(t, existsDetail, "cache detail harus dihapus setelah delete")
	assert.False(t, existsList, "cache list harus dihapus setelah delete")
}

/*
=====================================
 HELPERS REDIS UNTUK PERUSAHAAN TEST
=====================================
*/

func newPerusahaanTestRedis() *perusahaanTestRedis {
	return &perusahaanTestRedis{data: make(map[string]string)}
}

func setPerusahaanCache(rc *perusahaanTestRedis, key string, value interface{}) {
	b, _ := json.Marshal(value)
	rc.data[key] = string(b)
}

type perusahaanTestRedis struct {
	data map[string]string
}

func (r *perusahaanTestRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}
func (r *perusahaanTestRedis) Get(key string) (string, error) {
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (r *perusahaanTestRedis) Delete(key string) error { delete(r.data, key); return nil }
func (r *perusahaanTestRedis) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}
func (r *perusahaanTestRedis) Scan(pattern string) ([]string, error) { return nil, nil }
func (r *perusahaanTestRedis) Close() error                          { return nil }