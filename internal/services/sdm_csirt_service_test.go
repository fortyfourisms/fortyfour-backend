package services

import (
	"errors"
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
=====================================
 MOCK SDM CSIRT REPOSITORY
=====================================
*/

type mockSdmCsirtRepo struct {
	CreateFn     func(req dto.CreateSdmCsirtRequest, id string) error
	GetAllFn     func() ([]dto.SdmCsirtResponse, error)
	GetByIDFn    func(id string) (*dto.SdmCsirtResponse, error)
	GetByCsirtFn func(idCsirt string) ([]dto.SdmCsirtResponse, error)
	UpdateFn     func(id string, req dto.SdmCsirtResponse) error
	DeleteFn     func(id string) error
}

func (m *mockSdmCsirtRepo) Create(req dto.CreateSdmCsirtRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockSdmCsirtRepo) GetAll() ([]dto.SdmCsirtResponse, error) {
	return m.GetAllFn()
}

func (m *mockSdmCsirtRepo) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}

func (m *mockSdmCsirtRepo) GetByCsirt(idCsirt string) ([]dto.SdmCsirtResponse, error) {
	if m.GetByCsirtFn != nil {
		return m.GetByCsirtFn(idCsirt)
	}
	return []dto.SdmCsirtResponse{}, nil
}

func (m *mockSdmCsirtRepo) Update(id string, req dto.SdmCsirtResponse) error {
	return m.UpdateFn(id, req)
}

func (m *mockSdmCsirtRepo) Delete(id string) error {
	return m.DeleteFn(id)
}

/*
=====================================
 HELPER FUNCTIONS
=====================================
*/

func strPtr(s string) *string {
	return &s
}

/*
=====================================
 TEST CREATE SDM CSIRT - SUCCESS CASES
=====================================
*/

func TestSdmCsirtService_Create_Success(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Andi Pratama",
			}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: strPtr("Andi Pratama"),
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestSdmCsirtService_Create_SuccessWithAllFields(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id, NamaPersonel: "Budi Santoso"}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel:      strPtr("Budi Santoso"),
		JabatanCsirt:      strPtr("Security Analyst"),
		JabatanPerusahaan: strPtr("IT Manager"),
		Skill:             strPtr("Penetration Testing, SIEM"),
		Sertifikasi:       strPtr("CISSP, CEH"),
		IdCsirt:           strPtr("csirt-123"),
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

/*
=====================================
 TEST CREATE SDM CSIRT - VALIDATION ERRORS
 (Validation is handled at handler/repository layer)
=====================================
*/

func TestSdmCsirtService_Create_EmptyName_RepositoryRejects(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return errors.New("nama_personel tidak boleh kosong")
		},
	}
	service := NewSdmCsirtService(repo, nil)

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: strPtr(""),
	})

	assert.Error(t, err)
	assert.Empty(t, id)
	assert.Contains(t, err.Error(), "nama_personel")
}

func TestSdmCsirtService_Create_NilName_RepositoryRejects(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return errors.New("nama_personel tidak boleh kosong")
		},
	}
	service := NewSdmCsirtService(repo, nil)

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: nil,
	})

	assert.Error(t, err)
	assert.Empty(t, id)
}

/*
=====================================
 TEST CREATE SDM CSIRT - REPOSITORY ERRORS
=====================================
*/

func TestSdmCsirtService_Create_RepositoryError(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return errors.New("database connection error")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: strPtr("Andi"),
	})

	assert.Error(t, err)
	assert.Empty(t, id)
	assert.Equal(t, "database connection error", err.Error())
}

func TestSdmCsirtService_Create_InvalidCSIRTReference(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return errors.New("csirt tidak ditemukan")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: strPtr("Andi"),
		IdCsirt:      strPtr("invalid-csirt-id"),
	})

	assert.Error(t, err)
	assert.Empty(t, id)
}

/*
=====================================
 TEST GET ALL SDM CSIRT
=====================================
*/

func TestSdmCsirtService_GetAll_Success(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetAllFn: func() ([]dto.SdmCsirtResponse, error) {
			return []dto.SdmCsirtResponse{
				{
					ID:           "sdm-1",
					NamaPersonel: "Budi",
					JabatanCsirt: "Security Analyst",
				},
				{
					ID:           "sdm-2",
					NamaPersonel: "Citra",
					JabatanCsirt: "SOC Manager",
				},
			}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Budi", res[0].NamaPersonel)
	assert.Equal(t, "Citra", res[1].NamaPersonel)
}

func TestSdmCsirtService_GetAll_EmptyResult(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetAllFn: func() ([]dto.SdmCsirtResponse, error) {
			return []dto.SdmCsirtResponse{}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)
}

func TestSdmCsirtService_GetAll_RepositoryError(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetAllFn: func() ([]dto.SdmCsirtResponse, error) {
			return nil, errors.New("database timeout")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "database timeout", err.Error())
}

/*
=====================================
 TEST GET SDM CSIRT BY ID
=====================================
*/

func TestSdmCsirtService_GetByID_Success(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:                id,
				NamaPersonel:      "Andi Pratama",
				JabatanCsirt:      "Security Analyst",
				JabatanPerusahaan: "IT Manager",
				Skill:             "Penetration Testing, SIEM",
				Sertifikasi:       "CISSP, CEH",
			}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetByID("sdm-123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "sdm-123", res.ID)
	assert.Equal(t, "Andi Pratama", res.NamaPersonel)
	assert.Equal(t, "Security Analyst", res.JabatanCsirt)
}

func TestSdmCsirtService_GetByID_NotFound(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return nil, errors.New("sdm csirt tidak ditemukan")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "sdm csirt tidak ditemukan", err.Error())
}

func TestSdmCsirtService_GetByID_EmptyID(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			if id == "" {
				return nil, errors.New("id cannot be empty")
			}
			return nil, errors.New("sdm csirt tidak ditemukan")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetByID("")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSdmCsirtService_GetByID_RepositoryError(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetByID("sdm-123")

	assert.Error(t, err)
	assert.Nil(t, res)
}

/*
=====================================
 TEST UPDATE SDM CSIRT - SUCCESS CASES
=====================================
*/

func TestSdmCsirtService_Update_Success(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Old Name",
			}, nil
		},
		UpdateFn: func(id string, req dto.SdmCsirtResponse) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	newName := "New Name"
	err := service.Update("sdm-1", dto.UpdateSdmCsirtRequest{
		NamaPersonel: &newName,
	})

	assert.NoError(t, err)
}

func TestSdmCsirtService_Update_PartialUpdate(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Andi",
				JabatanCsirt: "Analyst",
			}, nil
		},
		UpdateFn: func(id string, req dto.SdmCsirtResponse) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	// Hanya update jabatan csirt
	newJabatan := "Senior Analyst"
	err := service.Update("sdm-1", dto.UpdateSdmCsirtRequest{
		JabatanCsirt: &newJabatan,
	})

	assert.NoError(t, err)
}

func TestSdmCsirtService_Update_AllFields(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Old",
			}, nil
		},
		UpdateFn: func(id string, req dto.SdmCsirtResponse) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	newName := "New Name"
	newJabatanCsirt := "Senior Analyst"
	newJabatanPerusahaan := "IT Director"
	newSkill := "Incident Response, Forensics"
	newCert := "CISSP, CISM"

	err := service.Update("sdm-1", dto.UpdateSdmCsirtRequest{
		NamaPersonel:      &newName,
		JabatanCsirt:      &newJabatanCsirt,
		JabatanPerusahaan: &newJabatanPerusahaan,
		Skill:             &newSkill,
		Sertifikasi:       &newCert,
	})

	assert.NoError(t, err)
}

/*
=====================================
 TEST UPDATE SDM CSIRT - ERROR CASES
=====================================
*/

func TestSdmCsirtService_Update_NotFound(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return nil, errors.New("sdm csirt tidak ditemukan")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	newName := "New Name"
	err := service.Update("invalid-id", dto.UpdateSdmCsirtRequest{
		NamaPersonel: &newName,
	})

	assert.Error(t, err)
	assert.Equal(t, "sdm csirt tidak ditemukan", err.Error())
}

func TestSdmCsirtService_Update_EmptyName_RepositoryRejects(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Andi",
			}, nil
		},
		UpdateFn: func(id string, req dto.SdmCsirtResponse) error {
			return errors.New("nama_personel tidak boleh kosong")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	emptyName := ""
	err := service.Update("sdm-1", dto.UpdateSdmCsirtRequest{
		NamaPersonel: &emptyName,
	})

	assert.Error(t, err)
}

func TestSdmCsirtService_Update_RepositoryError(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Andi",
			}, nil
		},
		UpdateFn: func(id string, req dto.SdmCsirtResponse) error {
			return errors.New("database update failed")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	newName := "New Name"
	err := service.Update("sdm-1", dto.UpdateSdmCsirtRequest{
		NamaPersonel: &newName,
	})

	assert.Error(t, err)
}

/*
=====================================
 TEST DELETE SDM CSIRT
=====================================
*/

func TestSdmCsirtService_Delete_Success(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	err := service.Delete("sdm-1")

	assert.NoError(t, err)
}

func TestSdmCsirtService_Delete_NotFound(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		DeleteFn: func(id string) error {
			return errors.New("sdm csirt tidak ditemukan")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	err := service.Delete("invalid-id")

	assert.Error(t, err)
	assert.Equal(t, "sdm csirt tidak ditemukan", err.Error())
}

func TestSdmCsirtService_Delete_EmptyID(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		DeleteFn: func(id string) error {
			if id == "" {
				return errors.New("id cannot be empty")
			}
			return nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	err := service.Delete("")

	assert.Error(t, err)
}

func TestSdmCsirtService_Delete_RepositoryError(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		DeleteFn: func(id string) error {
			return errors.New("database error")
		},
	}

	service := NewSdmCsirtService(repo, nil)

	err := service.Delete("sdm-1")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

/*
=====================================
 TEST EDGE CASES
=====================================
*/

func TestSdmCsirtService_Create_VeryLongName(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	longName := "Andi Pratama Budi Santoso Citra Dewi Rahman Wijaya Kusuma Putra Mahardika" // Very long name

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: &longName,
	})

	// Depending on validation rules
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestSdmCsirtService_Create_SpecialCharactersInName(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	nameWithSpecialChars := "O'Brien"

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: &nameWithSpecialChars,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestSdmCsirtService_Update_NoFieldsToUpdate(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Andi",
			}, nil
		},
		UpdateFn: func(id string, req dto.SdmCsirtResponse) error {
			return nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	// Update without any fields
	err := service.Update("sdm-1", dto.UpdateSdmCsirtRequest{})

	// Should be successful but no changes
	assert.NoError(t, err)
}

func TestSdmCsirtService_GetByID_WithCSIRTRelation(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{
				ID:           id,
				NamaPersonel: "Andi",
				Csirt: &dto.CsirtMiniResponse{
					ID:        "csirt-123",
					NamaCsirt: "CSIRT Test",
				},
			}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	res, err := service.GetByID("sdm-1")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Csirt)
	assert.Equal(t, "CSIRT Test", res.Csirt.NamaCsirt)
}

func TestSdmCsirtService_Create_MultipleCertifications(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	multipleCerts := "CISSP, CEH, OSCP, Security+, GIAC"

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: strPtr("Senior Expert"),
		Sertifikasi:  &multipleCerts,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestSdmCsirtService_Create_WithSkillSet(t *testing.T) {
	repo := &mockSdmCsirtRepo{
		CreateFn: func(req dto.CreateSdmCsirtRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.SdmCsirtResponse, error) {
			return &dto.SdmCsirtResponse{ID: id}, nil
		},
	}

	service := NewSdmCsirtService(repo, nil)

	skills := "Incident Response, Malware Analysis, Network Security, SIEM, Forensics"

	id, err := service.Create(dto.CreateSdmCsirtRequest{
		NamaPersonel: strPtr("Security Expert"),
		Skill:        &skills,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}