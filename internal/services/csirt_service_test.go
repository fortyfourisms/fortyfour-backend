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
	ExistsByPerusahaanFn    func(idPerusahaan string) (bool, error)
	GetByIDFn               func(id string) (*models.Csirt, error)
	GetAllWithPerusahaanFn  func() ([]dto.CsirtResponse, error)
	GetByIDWithPerusahaanFn func(id string) (*dto.CsirtResponse, error)
	UpdateFn                func(id string, csirt models.Csirt) error
	DeleteFn                func(id string) error
}

func (m *mockCsirtRepo) Create(req dto.CreateCsirtRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockCsirtRepo) ExistsByPerusahaan(idPerusahaan string) (bool, error) {
	if m.ExistsByPerusahaanFn != nil {
		return m.ExistsByPerusahaanFn(idPerusahaan)
	}
	return false, nil
}

func (m *mockCsirtRepo) GetByID(id string) (*models.Csirt, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
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

func (m *mockCsirtRepo) GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error) {
	return []dto.CsirtResponse{}, nil
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

	service := NewCsirtService(repo, nil, nil)

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
	tglReg := "2024-01-15"
	tglKad := "2025-01-15"
	tglRegU := "2025-01-20"
	fileStr := "uploads/str_csirt/str.pdf"

	var capturedReq dto.CreateCsirtRequest
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			capturedReq = req
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:                     id,
				NamaCsirt:              "CSIRT Advanced",
				TeleponCsirt:           &telepon,
				WebCsirt:               "https://csirt.example.com",
				IdPerusahaan:           "perusahaan-123",
				FileStr:                &fileStr,
				TanggalRegistrasi:      &tglReg,
				TanggalKadaluarsa:      &tglKad,
				TanggalRegistrasiUlang: &tglRegU,
			}, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt:              "CSIRT Advanced",
		TeleponCsirt:           "081234567890",
		WebCsirt:               "https://csirt.example.com",
		IdPerusahaan:           "perusahaan-123",
		FileStr:                fileStr,
		TanggalRegistrasi:      tglReg,
		TanggalKadaluarsa:      tglKad,
		TanggalRegistrasiUlang: tglRegU,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "CSIRT Advanced", res.NamaCsirt)
	assert.NotNil(t, res.TeleponCsirt)
	assert.Equal(t, "081234567890", *res.TeleponCsirt)
	// Verifikasi field baru diteruskan ke repo
	assert.Equal(t, fileStr, capturedReq.FileStr)
	assert.Equal(t, tglReg, capturedReq.TanggalRegistrasi)
	assert.Equal(t, tglKad, capturedReq.TanggalKadaluarsa)
	assert.Equal(t, tglRegU, capturedReq.TanggalRegistrasiUlang)
	// Verifikasi repo mengembalikan field baru
	assert.NotNil(t, res.FileStr)
	assert.Equal(t, fileStr, *res.FileStr)
	assert.NotNil(t, res.TanggalRegistrasi)
	assert.Equal(t, tglReg, *res.TanggalRegistrasi)
	assert.NotNil(t, res.TanggalKadaluarsa)
	assert.Equal(t, tglKad, *res.TanggalKadaluarsa)
	assert.NotNil(t, res.TanggalRegistrasiUlang)
	assert.Equal(t, tglRegU, *res.TanggalRegistrasiUlang)
}

// ─────────────────────────────────────────────────────────
// CREATE — field baru nullable (kosong → NULL)
// ─────────────────────────────────────────────────────────

func TestCsirtService_Create_NewFieldsNullable_WhenEmpty(t *testing.T) {
	var capturedReq dto.CreateCsirtRequest
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			capturedReq = req
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "CSIRT Minimal",
				// FileStr dan Tanggal* semua nil
			}, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt:    "CSIRT Minimal",
		IdPerusahaan: "perusahaan-1",
		// FileStr dan Tanggal* sengaja tidak diisi
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	// Field kosong → diteruskan sebagai empty string ke repo (nullableStr di repo yg ubah ke nil)
	assert.Equal(t, "", capturedReq.FileStr)
	assert.Equal(t, "", capturedReq.TanggalRegistrasi)
	assert.Equal(t, "", capturedReq.TanggalKadaluarsa)
	assert.Equal(t, "", capturedReq.TanggalRegistrasiUlang)
	// Hasil dari GetByID: field nullable nil
	assert.Nil(t, res.FileStr)
	assert.Nil(t, res.TanggalRegistrasi)
	assert.Nil(t, res.TanggalKadaluarsa)
	assert.Nil(t, res.TanggalRegistrasiUlang)
}

func TestCsirtService_Create_OnlyTanggalRegistrasi(t *testing.T) {
	tglReg := "2025-03-01"
	var capturedReq dto.CreateCsirtRequest
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			capturedReq = req
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:                id,
				NamaCsirt:         "CSIRT Reg Only",
				TanggalRegistrasi: &tglReg,
			}, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt:         "CSIRT Reg Only",
		IdPerusahaan:      "perusahaan-1",
		TanggalRegistrasi: tglReg,
		// Kadaluarsa dan RegistrasiUlang sengaja kosong
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, tglReg, capturedReq.TanggalRegistrasi)
	assert.Equal(t, "", capturedReq.TanggalKadaluarsa)
	assert.Equal(t, "", capturedReq.TanggalRegistrasiUlang)
	assert.NotNil(t, res.TanggalRegistrasi)
	assert.Equal(t, tglReg, *res.TanggalRegistrasi)
}

func TestCsirtService_Create_WithFileStrOnly(t *testing.T) {
	fileStr := "uploads/str_csirt/abc123.pdf"
	var capturedReq dto.CreateCsirtRequest
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			capturedReq = req
			return nil
		},
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "CSIRT STR Only",
				FileStr:   &fileStr,
			}, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt:    "CSIRT STR Only",
		IdPerusahaan: "perusahaan-1",
		FileStr:      fileStr,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, fileStr, capturedReq.FileStr)
	assert.NotNil(t, res.FileStr)
	assert.Equal(t, fileStr, *res.FileStr)
	// Tanggal lain harus kosong
	assert.Equal(t, "", capturedReq.TanggalRegistrasi)
	assert.Nil(t, res.TanggalRegistrasi)
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
	service := NewCsirtService(repo, nil, nil)

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
	service := NewCsirtService(repo, nil, nil)

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
	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT Test",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCsirtService_Create_DuplicatePerusahaan(t *testing.T) {
	repo := &mockCsirtRepo{
		ExistsByPerusahaanFn: func(idPerusahaan string) (bool, error) {
			return true, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		IdPerusahaan: "perusahaan-123",
		NamaCsirt:    "CSIRT Test",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "perusahaan ini sudah memiliki data CSIRT")
}

func TestCsirtService_Create_DuplicateCSIRTName(t *testing.T) {
	repo := &mockCsirtRepo{
		CreateFn: func(req dto.CreateCsirtRequest, id string) error {
			return errors.New("nama csirt sudah digunakan")
		},
	}

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

	res, err := service.GetByID("csirt-123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "csirt-123", res.ID)
	assert.Equal(t, "CSIRT Test", res.NamaCsirt)
	assert.Equal(t, "https://test.csirt.com", res.WebCsirt)
}

// ─────────────────────────────────────────────────────────
// GET BY ID — memastikan field baru ikut dikembalikan
// ─────────────────────────────────────────────────────────

func TestCsirtService_GetByID_ReturnsNewFields(t *testing.T) {
	tglReg := "2024-01-15"
	tglKad := "2025-01-15"
	tglRegU := "2025-01-20"
	fileStr := "uploads/str_csirt/abc.pdf"
	telepon := "081234567890"

	repo := &mockCsirtRepo{
		GetByIDWithPerusahaanFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID:                     id,
				NamaCsirt:              "CSIRT Lengkap",
				TeleponCsirt:           &telepon,
				FileStr:                fileStr,
				TanggalRegistrasi:      tglReg,
				TanggalKadaluarsa:      tglKad,
				TanggalRegistrasiUlang: tglRegU,
			}, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.GetByID("csirt-123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, fileStr, res.FileStr)
	assert.Equal(t, tglReg, res.TanggalRegistrasi)
	assert.Equal(t, tglKad, res.TanggalKadaluarsa)
	assert.Equal(t, tglRegU, res.TanggalRegistrasiUlang)
}

func TestCsirtService_GetByID_NewFieldsNullWhenNotSet(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDWithPerusahaanFn: func(id string) (*dto.CsirtResponse, error) {
			return &dto.CsirtResponse{
				ID:        id,
				NamaCsirt: "CSIRT Tanpa Tanggal",
				// FileStr dan Tanggal* sengaja kosong (zero value string)
			}, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.GetByID("csirt-123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "", res.FileStr)
	assert.Equal(t, "", res.TanggalRegistrasi)
	assert.Equal(t, "", res.TanggalKadaluarsa)
	assert.Equal(t, "", res.TanggalRegistrasiUlang)
}

func TestCsirtService_GetByID_NotFound(t *testing.T) {
	repo := &mockCsirtRepo{
		GetByIDWithPerusahaanFn: func(id string) (*dto.CsirtResponse, error) {
			return nil, errors.New("csirt tidak ditemukan")
		},
	}

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

// ─────────────────────────────────────────────────────────
// UPDATE — field baru: TanggalRegistrasi, Kadaluarsa, RegistrasiUlang, FileStr
// ─────────────────────────────────────────────────────────

func TestCsirtService_Update_AllFields(t *testing.T) {
	tglReg := "2024-06-01"
	tglKad := "2025-06-01"
	tglRegU := "2025-06-10"
	fileStr := "uploads/str_csirt/new.pdf"

	var capturedCsirt models.Csirt
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "Old Name",
				WebCsirt:  "https://old.com",
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			capturedCsirt = csirt
			return nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	newName := "New Name"
	newWeb := "https://newcsirt.com"
	newPhone := "082222222222"

	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		NamaCsirt:              &newName,
		WebCsirt:               &newWeb,
		TeleponCsirt:           &newPhone,
		TanggalRegistrasi:      &tglReg,
		TanggalKadaluarsa:      &tglKad,
		TanggalRegistrasiUlang: &tglRegU,
		FileStr:                &fileStr,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "New Name", res.NamaCsirt)
	assert.Equal(t, "https://newcsirt.com", res.WebCsirt)
	assert.NotNil(t, res.TeleponCsirt)
	assert.Equal(t, "082222222222", *res.TeleponCsirt)
	// Verifikasi field baru masuk ke repo
	assert.NotNil(t, capturedCsirt.TanggalRegistrasi)
	assert.Equal(t, tglReg, *capturedCsirt.TanggalRegistrasi)
	assert.NotNil(t, capturedCsirt.TanggalKadaluarsa)
	assert.Equal(t, tglKad, *capturedCsirt.TanggalKadaluarsa)
	assert.NotNil(t, capturedCsirt.TanggalRegistrasiUlang)
	assert.Equal(t, tglRegU, *capturedCsirt.TanggalRegistrasiUlang)
	assert.NotNil(t, capturedCsirt.FileStr)
	assert.Equal(t, fileStr, *capturedCsirt.FileStr)
	// Verifikasi result juga membawa field baru
	assert.NotNil(t, res.TanggalRegistrasi)
	assert.Equal(t, tglReg, *res.TanggalRegistrasi)
	assert.NotNil(t, res.FileStr)
	assert.Equal(t, fileStr, *res.FileStr)
}

func TestCsirtService_Update_OnlyTanggalKadaluarsa(t *testing.T) {
	tglKad := "2026-12-31"
	oldTglReg := "2024-01-01"

	var capturedCsirt models.Csirt
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:                id,
				NamaCsirt:         "CSIRT Existing",
				TanggalRegistrasi: &oldTglReg,
				// TanggalKadaluarsa belum diset
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			capturedCsirt = csirt
			return nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		TanggalKadaluarsa: &tglKad,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	// TanggalRegistrasi lama harus tetap terjaga
	assert.NotNil(t, capturedCsirt.TanggalRegistrasi)
	assert.Equal(t, oldTglReg, *capturedCsirt.TanggalRegistrasi)
	// TanggalKadaluarsa harus ter-update
	assert.NotNil(t, capturedCsirt.TanggalKadaluarsa)
	assert.Equal(t, tglKad, *capturedCsirt.TanggalKadaluarsa)
	// TanggalRegistrasiUlang tidak di-update → tetap nil
	assert.Nil(t, capturedCsirt.TanggalRegistrasiUlang)
}

func TestCsirtService_Update_FileStr_ReplacesExisting(t *testing.T) {
	oldFileStr := "uploads/str_csirt/old.pdf"
	newFileStr := "uploads/str_csirt/new.pdf"

	var capturedCsirt models.Csirt
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "CSIRT Existing",
				FileStr:   &oldFileStr,
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			capturedCsirt = csirt
			return nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		FileStr: &newFileStr,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, capturedCsirt.FileStr)
	assert.Equal(t, newFileStr, *capturedCsirt.FileStr)
	assert.NotNil(t, res.FileStr)
	assert.Equal(t, newFileStr, *res.FileStr)
}

func TestCsirtService_Update_NewFieldsNotProvided_ExistingValuesKept(t *testing.T) {
	oldTglReg := "2024-01-01"
	oldTglKad := "2025-01-01"
	oldFileStr := "uploads/str_csirt/existing.pdf"

	var capturedCsirt models.Csirt
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:                id,
				NamaCsirt:         "CSIRT Lama",
				TanggalRegistrasi: &oldTglReg,
				TanggalKadaluarsa: &oldTglKad,
				FileStr:           &oldFileStr,
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			capturedCsirt = csirt
			return nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	// Hanya update nama, field baru tidak disertakan
	newName := "CSIRT Diperbarui"
	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		NamaCsirt: &newName,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "CSIRT Diperbarui", res.NamaCsirt)
	// Nilai lama harus tetap terjaga karena tidak di-update
	assert.NotNil(t, capturedCsirt.TanggalRegistrasi)
	assert.Equal(t, oldTglReg, *capturedCsirt.TanggalRegistrasi)
	assert.NotNil(t, capturedCsirt.TanggalKadaluarsa)
	assert.Equal(t, oldTglKad, *capturedCsirt.TanggalKadaluarsa)
	assert.NotNil(t, capturedCsirt.FileStr)
	assert.Equal(t, oldFileStr, *capturedCsirt.FileStr)
}

func TestCsirtService_Update_TanggalRegistrasiUlang_Independent(t *testing.T) {
	tglRegU := "2026-02-01"

	var capturedCsirt models.Csirt
	repo := &mockCsirtRepo{
		GetByIDFn: func(id string) (*models.Csirt, error) {
			return &models.Csirt{
				ID:        id,
				NamaCsirt: "CSIRT A",
			}, nil
		},
		UpdateFn: func(id string, csirt models.Csirt) error {
			capturedCsirt = csirt
			return nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Update("csirt-1", dto.UpdateCsirtRequest{
		TanggalRegistrasiUlang: &tglRegU,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, capturedCsirt.TanggalRegistrasiUlang)
	assert.Equal(t, tglRegU, *capturedCsirt.TanggalRegistrasiUlang)
	// Tanggal lain tidak di-update → tetap nil
	assert.Nil(t, capturedCsirt.TanggalRegistrasi)
	assert.Nil(t, capturedCsirt.TanggalKadaluarsa)
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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

	err := service.Delete("csirt-1")

	assert.NoError(t, err)
}

func TestCsirtService_Delete_NotFound(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			return errors.New("csirt tidak ditemukan")
		},
	}

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

	err := service.Delete("")

	assert.Error(t, err)
}

func TestCsirtService_Delete_RepositoryError(t *testing.T) {
	repo := &mockCsirtRepo{
		DeleteFn: func(id string) error {
			return errors.New("database error")
		},
	}

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

	res, err := service.Create(dto.CreateCsirtRequest{
		NamaCsirt: "CSIRT @#$%",
	})

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

	service := NewCsirtService(repo, nil, nil)

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

	service := NewCsirtService(repo, nil, nil)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "PT Test", res[0].Perusahaan.NamaPerusahaan)
}

// ─────────────────────────────────────────────────────────
// GetAll — memastikan field baru muncul di response list
// ─────────────────────────────────────────────────────────

func TestCsirtService_GetAll_IncludesNewFields(t *testing.T) {
	tglReg := "2024-01-15"
	fileStr := "uploads/str_csirt/abc.pdf"

	repo := &mockCsirtRepo{
		GetAllWithPerusahaanFn: func() ([]dto.CsirtResponse, error) {
			return []dto.CsirtResponse{
				{
					ID:                "csirt-1",
					NamaCsirt:         "CSIRT A",
					TanggalRegistrasi: tglReg,
					FileStr:           fileStr,
				},
				{
					ID:        "csirt-2",
					NamaCsirt: "CSIRT B",
					// field baru kosong
				},
			}, nil
		},
	}

	service := NewCsirtService(repo, nil, nil)

	res, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, tglReg, res[0].TanggalRegistrasi)
	assert.Equal(t, fileStr, res[0].FileStr)
	assert.Equal(t, "", res[1].TanggalRegistrasi)
	assert.Equal(t, "", res[1].FileStr)
}
