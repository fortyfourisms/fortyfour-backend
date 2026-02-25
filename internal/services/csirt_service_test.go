package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
=====================================
 MOCK CSIRT REPOSITORY
=====================================
*/

type mockCsirtRepo struct {
	CreateFn                func(req dto.CreateCsirtRequest, id string) error
	GetByIDFn               func(id string) (*models.Csirt, error)
	GetAllWithPerusahaanFn  func() ([]dto.CsirtResponse, error)
	GetByIDWithPerusahaanFn func(id string) (*dto.CsirtResponse, error)
	UpdateFn                func(id string, csirt models.Csirt) error
	DeleteFn                func(id string) error
}

func (m *mockCsirtRepo) Create(req dto.CreateCsirtRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockCsirtRepo) GetByID(id string) (*models.Csirt, error) {
	return m.GetByIDFn(id)
}

func (m *mockCsirtRepo) GetAllWithPerusahaan() ([]dto.CsirtResponse, error) {
	return m.GetAllWithPerusahaanFn()
}

func (m *mockCsirtRepo) GetByIDWithPerusahaan(id string) (*dto.CsirtResponse, error) {
	return m.GetByIDWithPerusahaanFn(id)
}

func (m *mockCsirtRepo) Update(id string, csirt models.Csirt) error {
	return m.UpdateFn(id, csirt)
}

func (m *mockCsirtRepo) Delete(id string) error {
	return m.DeleteFn(id)
}

/*
=====================================
 TEST CREATE CSIRT - SUCCESS CASES
=====================================
*/

func TestCsirtService_Create_Success(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "CSIRT Test",
			}, nil
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT Test",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "CSIRT Test", res.NamaCsirt)
	assert.NotEmpty(t, res.ID)
}

func TestCsirtService_Create_SuccessWithAllFields(t *testing.T) {
	telepon := "081234567890"
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:           id,
				NamaCsirt:    "CSIRT Advanced",
				TeleponCsirt: &telepon,
				WebCsirt:     "https://csirt.example.com",
				IdPerusahaan: "perusahaan-123",
			}, nil
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt:    "CSIRT Advanced",
		TeleponCsirt: "081234567890",
		WebCsirt:     "https://csirt.example.com",
		IdPerusahaan: "perusahaan-123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "CSIRT Advanced", res.NamaCsirt)
	assert.NotNil(t, res.TeleponCsirt)
	assert.Equal(t, "081234567890", *res.TeleponCsirt)
}

/*
=====================================
 TEST CREATE CSIRT - VALIDATION ERRORS
 (Validation is handled at handler/repository layer)
=====================================
*/

func TestCsirtService_Create_EmptyName_RepositoryRejects(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return errors.New("nama_csirt tidak boleh kosong")
		},
	}
	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "nama_csirt")
}

func TestCsirtService_Create_InvalidPhoneNumber_RepositoryRejects(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return errors.New("nomor telepon tidak valid")
		},
	}
	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt:    "CSIRT Test",
		TeleponCsirt: "12345", // Too short
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCsirtService_Create_InvalidWebsite_RepositoryRejects(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return errors.New("format website tidak valid")
		},
	}
	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT Test",
		WebCsirt:  "not-a-url",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

/*
=====================================
 TEST CREATE CSIRT - REPOSITORY ERRORS
=====================================
*/

func TestCsirtService_Create_RepositoryError_CreateFailed(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return errors.New("database connection error")
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT Test",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "database connection error", err.Error())
}

func TestCsirtService_Create_RepositoryError_GetByIDFailed(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return nil, errors.New("failed to fetch created csirt")
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT Test",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCsirtService_Create_DuplicateCSIRTName(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return errors.New("nama csirt sudah digunakan")
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT Existing",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "nama csirt sudah digunakan", err.Error())
}

/*
=====================================
 TEST GET ALL CSIRT
=====================================
*/

func TestCsirtService_GetAll_Success(t *testing.T) {
	repo := &mockCsirtRepo{
		GetAllWithPerusahaanFn: func() ([]dto.CsirtResponse, error) {
			return []dto.CsirtResponse{
				{
					ID:        "csirt-1",
					NamaCsirt: "CSIRT A",
					WebCsirt:  "https://csirta.example.com",
				},
				{
					ID:        "csirt-2",
					NamaCsirt: "CSIRT B",
					WebCsirt:  "https://csirtb.example.com",
				},
			}, nil
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "CSIRT A", res[0].NamaCsirt)
	assert.Equal(t, "CSIRT B", res[1].NamaCsirt)
}

func TestCsirtService_GetAll_EmptyResult(t *testing.T) {
	repo := &mockCsirtRepo{
		GetAllWithPerusahaanFn: func() ([]dto.CsirtResponse, error) {
			return []dto.CsirtResponse{}, nil
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)
}

func TestCsirtService_GetAll_RepositoryError(t *testing.T) {
	repo := &mockCsirtRepo{
		GetAllWithPerusahaanFn: func() ([]dto.CsirtResponse, error) {
			return nil, errors.New("database timeout")
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "database timeout", err.Error())
}

/*
=====================================
 TEST GET CSIRT BY ID
=====================================
*/

func TestCsirtService_GetByID_Success(t *testing.T) {
	telepon := "081234567890"
	repo := &mockCsirtRepo{
		GetByIDWithPerusahaanFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID:           id,
				NamaCsirt:    "CSIRT Test",
				WebCsirt:     "https://test.csirt.com",
				TeleponCsirt: &telepon,
			}, nil
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetByID("csirt-123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "csirt-123", res.ID)
	assert.Equal(t, "CSIRT Test", res.NamaCsirt)
	assert.Equal(t, "https://test.csirt.com", res.WebCsirt)
}

func TestCsirtService_GetByID_NotFound(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDWithPerusahaanFn: func(id string) (*dto.CsirtResponse, error) {
			return nil, errors.New("csirt tidak ditemukan")
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "csirt tidak ditemukan", err.Error())
}

func TestCsirtService_GetByID_EmptyID(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDWithPerusahaanFn: func(id string) (*dto.CsirtResponse, error) {
			if id == "" {
				return nil, errors.New("id cannot be empty")
			}
			return nil, errors.New("csirt tidak ditemukan")
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetByID("")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCsirtService_GetByID_RepositoryError(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDWithPerusahaanFn: func(id string) (*dto.CsirtResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetByID("csirt-123")

	assert.Error(t, err)
	assert.Nil(t, res)
}

/*
=====================================
 TEST UPDATE CSIRT - SUCCESS CASES
=====================================
*/

func TestCsirtService_Update_Success(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "Old CSIRT Name",
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return nil
		},
	}

	service := NewCsirtService(repo, nil)

	newName := "New CSIRT Name"
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		NamaCsirt: &newName,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "New CSIRT Name", res.NamaCsirt)
}

func TestCsirtService_Update_PartialUpdate(t *testing.T) {
	telepon := "081111111111"
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:           id,
				NamaCsirt:    "CSIRT A",
				WebCsirt:     "https://old.csirt.com",
				TeleponCsirt: &telepon,
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return nil
		},
	}

	service := NewCsirtService(repo, nil)

	// Hanya update web
	newWeb := "https://new.csirt.com"
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		WebCsirt: &newWeb,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "CSIRT A", res.NamaCsirt)              // Nama tetap sama
	assert.Equal(t, "https://new.csirt.com", res.WebCsirt) // Web berubah
}

func TestCsirtService_Update_AllFields(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "Old Name",
				WebCsirt:  "https://old.com",
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return nil
		},
	}

	service := NewCsirtService(repo, nil)

	newName := "New Name"
	newWeb := "https://newcsirt.com"
	newPhone := "082222222222"

	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		NamaCsirt:    &newName,
		WebCsirt:     &newWeb,
		TeleponCsirt: &newPhone,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "New Name", res.NamaCsirt)
	assert.Equal(t, "https://newcsirt.com", res.WebCsirt)
	assert.NotNil(t, res.TeleponCsirt)
	assert.Equal(t, "082222222222", *res.TeleponCsirt)
}

/*
=====================================
 TEST UPDATE CSIRT - ERROR CASES
=====================================
*/

func TestCsirtService_Update_NotFound(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return nil, errors.New("csirt tidak ditemukan")
		},
	}

	service := NewCsirtService(repo, nil)

	newName := "New CSIRT"
	res, err := service.Update("invalid-id", dto.UpdateCsirtRequest{
		NamaCsirt: &newName,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "csirt tidak ditemukan", err.Error())
}

func TestCsirtService_Update_InvalidWebsite_RepositoryRejects(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{ID: id, NamaCsirt: "CSIRT A"}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return errors.New("format website tidak valid")
		},
	}

	service := NewCsirtService(repo, nil)

	invalidWeb := "not-a-url"
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		WebCsirt: &invalidWeb,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCsirtService_Update_EmptyName_RepositoryRejects(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{ID: id, NamaCsirt: "CSIRT A"}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return errors.New("nama_csirt tidak boleh kosong")
		},
	}

	service := NewCsirtService(repo, nil)

	emptyName := ""
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		NamaCsirt: &emptyName,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCsirtService_Update_RepositoryError_UpdateFailed(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{ID: id, NamaCsirt: "CSIRT A"}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return errors.New("database update failed")
		},
	}

	service := NewCsirtService(repo, nil)

	newName := "New CSIRT"
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		NamaCsirt: &newName,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCsirtService_Update_DuplicateName(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{ID: id, NamaCsirt: "CSIRT A"}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return errors.New("nama csirt sudah digunakan")
		},
	}

	service := NewCsirtService(repo, nil)

	duplicateName := "CSIRT B" // Assuming already exists
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		NamaCsirt: &duplicateName,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

/*
=====================================
 TEST DELETE CSIRT
=====================================
*/

func TestCsirtService_Delete_Success(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewCsirtService(repo, nil)

	err := service.Delete("csirt-1")

	assert.NoError(t, err)
}

func TestCsirtService_Delete_NotFound(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			return errors.New("csirt tidak ditemukan")
		},
	}

	service := NewCsirtService(repo, nil)

	err := service.Delete("invalid-id")

	assert.Error(t, err)
	assert.Equal(t, "csirt tidak ditemukan", err.Error())
}

func TestCsirtService_Delete_EmptyID(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			if id == "" {
				return errors.New("id cannot be empty")
			}
			return nil
		},
	}

	service := NewCsirtService(repo, nil)

	err := service.Delete("")

	assert.Error(t, err)
}

func TestCsirtService_Delete_RepositoryError(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			return errors.New("database error")
		},
	}

	service := NewCsirtService(repo, nil)

	err := service.Delete("csirt-1")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestCsirtService_Delete_HasRelatedData(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			return errors.New("cannot delete csirt with related data")
		},
	}

	service := NewCsirtService(repo, nil)

	err := service.Delete("csirt-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "related data")
}

/*
=====================================
 TEST EDGE CASES
=====================================
*/

func TestCsirtService_Create_SpecialCharactersInName(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "CSIRT @#$%",
			}, nil
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT @#$%",
	})

	// Depending on validation rules, this might be valid or invalid
	// Adjust assertion based on business rules
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCsirtService_Update_NoFieldsToUpdate(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "CSIRT A",
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			return nil
		},
	}

	service := NewCsirtService(repo, nil)

	// Update tanpa field apapun
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{})

	// Seharusnya tetap success tapi data tidak berubah
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "CSIRT A", res.NamaCsirt)
}

func TestCsirtService_GetAll_WithPerusahaanRelation(t *testing.T) {
	repo := &mockCsirtRepo{
		GetAllWithPerusahaanFn: func() ([]dto.CsirtResponse, error) {
			return []dto.CsirtResponse{
				{
					ID:        "csirt-1",
					NamaCsirt: "CSIRT A",
					Perusahaan: dto.PerusahaanResponse{
						ID:             "perusahaan-1",
						NamaPerusahaan: "PT Test",
					},
				},
			}, nil
		},
	}

	service := NewCsirtService(repo, nil)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "PT Test", res[0].Perusahaan.NamaPerusahaan)
}
