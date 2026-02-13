package services

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
)

/*
=====================================
 MOCK JABATAN REPOSITORY
=====================================
*/

type mockJabatanRepository struct {
	CreateFn  func(req dto.CreateJabatanRequest, id string) error
	GetByIDFn func(id string) (*dto.JabatanResponse, error)
	GetAllFn  func() ([]dto.JabatanResponse, error)
	UpdateFn  func(id string, jabatan dto.JabatanResponse) error
	DeleteFn  func(id string) error
}

func (m *mockJabatanRepository) Create(req dto.CreateJabatanRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockJabatanRepository) GetByID(id string) (*dto.JabatanResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockJabatanRepository) GetAll() ([]dto.JabatanResponse, error) {
	return m.GetAllFn()
}

func (m *mockJabatanRepository) Update(id string, jabatan dto.JabatanResponse) error {
	return m.UpdateFn(id, jabatan)
}

func (m *mockJabatanRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

/*
=====================================
 TEST CREATE - SUCCESS CASES
=====================================
*/

func TestCreateJabatan_Success(t *testing.T) {
	nama := "Manager IT"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: nama,
				CreatedAt:   "2024-01-01 10:00:00",
				UpdatedAt:   "2024-01-01 10:00:00",
			}, nil
		},
	}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: &nama,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.NamaJabatan)
	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestCreateJabatan_SuccessWithVariousNames(t *testing.T) {
	testCases := []struct {
		name         string
		namaJabatan  string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Simple name",
			namaJabatan: "Manager",
			expectError: false,
		},
		{
			name:        "Name with space",
			namaJabatan: "Manager IT",
			expectError: false,
		},
		{
			name:        "Name with numbers",
			namaJabatan: "Manager Level 1",
			expectError: false,
		},
		{
			name:        "Long name",
			namaJabatan: "Senior Executive Manager IT Infrastructure and Development",
			expectError: false,
		},
		{
			name:        "Name with special characters",
			namaJabatan: "Manager IT & Development",
			expectError: false,
		},
		{
			name:        "Indonesian name",
			namaJabatan: "Kepala Divisi Teknologi Informasi",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockJabatanRepository{
				CreateFn: func(req dto.CreateJabatanRequest, id string) error {
					return nil
				},
				GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
					return &dto.JabatanResponse{
						ID:          id,
						NamaJabatan: tc.namaJabatan,
					}, nil
				},
			}

			service := NewJabatanService(repo)

			req := dto.CreateJabatanRequest{
				NamaJabatan: &tc.namaJabatan,
			}

			result, err := service.Create(req)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.namaJabatan, result.NamaJabatan)
			}
		})
	}
}

/*
=====================================
 TEST CREATE - VALIDATION ERRORS
=====================================
*/

func TestCreateJabatan_ValidationFailed_NilName(t *testing.T) {
	repo := &mockJabatanRepository{}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: nil,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "nama_jabatan")
}

func TestCreateJabatan_ValidationFailed_EmptyName(t *testing.T) {
	emptyName := ""
	repo := &mockJabatanRepository{}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: &emptyName,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "nama_jabatan")
}

func TestCreateJabatan_ValidationFailed_WhitespaceName(t *testing.T) {
	testCases := []struct {
		name      string
		namaInput string
	}{
		{"Single space", " "},
		{"Multiple spaces", "   "},
		{"Tabs", "\t\t"},
		{"Newlines", "\n\n"},
		{"Mixed whitespace", " \t \n "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockJabatanRepository{}
			service := NewJabatanService(repo)

			req := dto.CreateJabatanRequest{
				NamaJabatan: &tc.namaInput,
			}

			result, err := service.Create(req)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "nama_jabatan")
		})
	}
}

/*
=====================================
 TEST CREATE - REPOSITORY ERRORS
=====================================
*/

func TestCreateJabatan_RepoFailed_CreateError(t *testing.T) {
	nama := "Manager"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			return errors.New("database connection error")
		},
	}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: &nama,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection error", err.Error())
}

func TestCreateJabatan_RepoFailed_GetByIDError(t *testing.T) {
	nama := "Manager"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return nil, errors.New("failed to retrieve created jabatan")
		},
	}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: &nama,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateJabatan_DuplicateName(t *testing.T) {
	nama := "Manager IT"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			return errors.New("nama jabatan sudah ada")
		},
	}

	service := NewJabatanService(repo)

	req := dto.CreateJabatanRequest{
		NamaJabatan: &nama,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST GET ALL - SUCCESS CASES
=====================================
*/

func TestGetAllJabatan_Success(t *testing.T) {
	repo := &mockJabatanRepository{
		GetAllFn: func() ([]dto.JabatanResponse, error) {
			return []dto.JabatanResponse{
				{
					ID:          "1",
					NamaJabatan: "Manager",
					CreatedAt:   "2024-01-01 10:00:00",
					UpdatedAt:   "2024-01-01 10:00:00",
				},
				{
					ID:          "2",
					NamaJabatan: "Staff",
					CreatedAt:   "2024-01-02 10:00:00",
					UpdatedAt:   "2024-01-02 10:00:00",
				},
			}, nil
		},
	}

	service := NewJabatanService(repo)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, "Manager", data[0].NamaJabatan)
	assert.Equal(t, "Staff", data[1].NamaJabatan)
}

func TestGetAllJabatan_EmptyResult(t *testing.T) {
	repo := &mockJabatanRepository{
		GetAllFn: func() ([]dto.JabatanResponse, error) {
			return []dto.JabatanResponse{}, nil
		},
	}

	service := NewJabatanService(repo)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 0)
}

func TestGetAllJabatan_MultipleRecords(t *testing.T) {
	repo := &mockJabatanRepository{
		GetAllFn: func() ([]dto.JabatanResponse, error) {
			result := make([]dto.JabatanResponse, 100)
			for i := 0; i < 100; i++ {
				result[i] = dto.JabatanResponse{
					ID:          string(rune(i)),
					NamaJabatan: "Jabatan " + string(rune(i)),
				}
			}
			return result, nil
		},
	}

	service := NewJabatanService(repo)

	data, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, data, 100)
}

/*
=====================================
 TEST GET ALL - ERROR CASES
=====================================
*/

func TestGetAllJabatan_RepositoryError(t *testing.T) {
	repo := &mockJabatanRepository{
		GetAllFn: func() ([]dto.JabatanResponse, error) {
			return nil, errors.New("database timeout")
		},
	}

	service := NewJabatanService(repo)

	data, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Equal(t, "database timeout", err.Error())
}

func TestGetAllJabatan_ConnectionError(t *testing.T) {
	repo := &mockJabatanRepository{
		GetAllFn: func() ([]dto.JabatanResponse, error) {
			return nil, errors.New("failed to connect to database")
		},
	}

	service := NewJabatanService(repo)

	data, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, data)
}

/*
=====================================
 TEST GET BY ID - SUCCESS CASES
=====================================
*/

func TestGetJabatanByID_Success(t *testing.T) {
	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: "Manager IT",
				CreatedAt:   "2024-01-01 10:00:00",
				UpdatedAt:   "2024-01-01 10:00:00",
			}, nil
		},
	}

	service := NewJabatanService(repo)

	result, err := service.GetByID("123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123", result.ID)
	assert.Equal(t, "Manager IT", result.NamaJabatan)
}

func TestGetJabatanByID_DifferentIDs(t *testing.T) {
	testCases := []string{
		"uuid-1",
		"uuid-abc-123",
		"12345",
		"very-long-uuid-string-12345678901234567890",
	}

	for _, id := range testCases {
		t.Run("ID: "+id, func(t *testing.T) {
			repo := &mockJabatanRepository{
				GetByIDFn: func(receivedID string) (*dto.JabatanResponse, error) {
					return &dto.JabatanResponse{
						ID:          receivedID,
						NamaJabatan: "Test Jabatan",
					}, nil
				},
			}

			service := NewJabatanService(repo)

			result, err := service.GetByID(id)

			assert.NoError(t, err)
			assert.Equal(t, id, result.ID)
		})
	}
}

/*
=====================================
 TEST GET BY ID - ERROR CASES
=====================================
*/

func TestGetJabatanByID_NotFound(t *testing.T) {
	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return nil, errors.New("jabatan tidak ditemukan")
		},
	}

	service := NewJabatanService(repo)

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "jabatan tidak ditemukan", err.Error())
}

func TestGetJabatanByID_EmptyID(t *testing.T) {
	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			if id == "" {
				return nil, errors.New("id cannot be empty")
			}
			return nil, errors.New("jabatan tidak ditemukan")
		},
	}

	service := NewJabatanService(repo)

	result, err := service.GetByID("")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetJabatanByID_RepositoryError(t *testing.T) {
	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := NewJabatanService(repo)

	result, err := service.GetByID("123")

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST UPDATE - SUCCESS CASES
=====================================
*/

func TestUpdateJabatan_Success(t *testing.T) {
	newName := "Updated Jabatan"

	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: "Old Jabatan",
				CreatedAt:   "2024-01-01 10:00:00",
				UpdatedAt:   "2024-01-01 10:00:00",
			}, nil
		},
		UpdateFn: func(id string, jabatan dto.JabatanResponse) error {
			return nil
		},
	}

	service := NewJabatanService(repo)

	req := dto.UpdateJabatanRequest{
		NamaJabatan: &newName,
	}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.Equal(t, newName, result.NamaJabatan)
	assert.Equal(t, "123", result.ID)
}

func TestUpdateJabatan_NoFieldsToUpdate(t *testing.T) {
	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: "Original Name",
			}, nil
		},
		UpdateFn: func(id string, jabatan dto.JabatanResponse) error {
			return nil
		},
	}

	service := NewJabatanService(repo)

	// Empty update request
	req := dto.UpdateJabatanRequest{}

	result, err := service.Update("123", req)

	assert.NoError(t, err)
	assert.Equal(t, "Original Name", result.NamaJabatan) // Unchanged
}

func TestUpdateJabatan_DifferentNames(t *testing.T) {
	testCases := []struct {
		name        string
		newName     string
		expectError bool
	}{
		{"Simple name", "Manager", false},
		{"Name with spaces", "Senior Manager", false},
		{"Long name", "Chief Executive Officer of Technology Department", false},
		{"Indonesian name", "Kepala Bagian IT", false},
		{"Name with numbers", "Manager Level 5", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockJabatanRepository{
				GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
					return &dto.JabatanResponse{
						ID:          id,
						NamaJabatan: "Old Name",
					}, nil
				},
				UpdateFn: func(id string, jabatan dto.JabatanResponse) error {
					return nil
				},
			}

			service := NewJabatanService(repo)

			req := dto.UpdateJabatanRequest{
				NamaJabatan: &tc.newName,
			}

			result, err := service.Update("123", req)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.newName, result.NamaJabatan)
			}
		})
	}
}

/*
=====================================
 TEST UPDATE - ERROR CASES
=====================================
*/

func TestUpdateJabatan_NotFound(t *testing.T) {
	newName := "New Name"

	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return nil, errors.New("jabatan tidak ditemukan")
		},
	}

	service := NewJabatanService(repo)

	req := dto.UpdateJabatanRequest{
		NamaJabatan: &newName,
	}

	result, err := service.Update("invalid-id", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "jabatan tidak ditemukan", err.Error())
}

func TestUpdateJabatan_RepositoryUpdateFailed(t *testing.T) {
	newName := "New Name"

	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: "Old Name",
			}, nil
		},
		UpdateFn: func(id string, jabatan dto.JabatanResponse) error {
			return errors.New("database update failed")
		},
	}

	service := NewJabatanService(repo)

	req := dto.UpdateJabatanRequest{
		NamaJabatan: &newName,
	}

	result, err := service.Update("123", req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateJabatan_EmptyID(t *testing.T) {
	newName := "New Name"

	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			if id == "" {
				return nil, errors.New("id cannot be empty")
			}
			return nil, errors.New("jabatan tidak ditemukan")
		},
	}

	service := NewJabatanService(repo)

	req := dto.UpdateJabatanRequest{
		NamaJabatan: &newName,
	}

	result, err := service.Update("", req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateJabatan_DuplicateName(t *testing.T) {
	duplicateName := "Existing Manager"

	repo := &mockJabatanRepository{
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: "Old Name",
			}, nil
		},
		UpdateFn: func(id string, jabatan dto.JabatanResponse) error {
			return errors.New("nama jabatan sudah ada")
		},
	}

	service := NewJabatanService(repo)

	req := dto.UpdateJabatanRequest{
		NamaJabatan: &duplicateName,
	}

	result, err := service.Update("123", req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

/*
=====================================
 TEST DELETE - SUCCESS CASES
=====================================
*/

func TestDeleteJabatan_Success(t *testing.T) {
	repo := &mockJabatanRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewJabatanService(repo)

	err := service.Delete("123")

	assert.NoError(t, err)
}

func TestDeleteJabatan_MultipleIDs(t *testing.T) {
	ids := []string{"id-1", "id-2", "id-3", "uuid-abc-123"}

	for _, id := range ids {
		t.Run("Delete ID: "+id, func(t *testing.T) {
			deletedID := ""

			repo := &mockJabatanRepository{
				DeleteFn: func(receivedID string) error {
					deletedID = receivedID
					return nil
				},
			}

			service := NewJabatanService(repo)

			err := service.Delete(id)

			assert.NoError(t, err)
			assert.Equal(t, id, deletedID)
		})
	}
}

/*
=====================================
 TEST DELETE - ERROR CASES
=====================================
*/

func TestDeleteJabatan_NotFound(t *testing.T) {
	repo := &mockJabatanRepository{
		DeleteFn: func(id string) error {
			return errors.New("jabatan tidak ditemukan")
		},
	}

	service := NewJabatanService(repo)

	err := service.Delete("invalid-id")

	assert.Error(t, err)
	assert.Equal(t, "jabatan tidak ditemukan", err.Error())
}

func TestDeleteJabatan_EmptyID(t *testing.T) {
	repo := &mockJabatanRepository{
		DeleteFn: func(id string) error {
			if id == "" {
				return errors.New("id cannot be empty")
			}
			return nil
		},
	}

	service := NewJabatanService(repo)

	err := service.Delete("")

	assert.Error(t, err)
}

func TestDeleteJabatan_RepositoryError(t *testing.T) {
	repo := &mockJabatanRepository{
		DeleteFn: func(id string) error {
			return errors.New("database error")
		},
	}

	service := NewJabatanService(repo)

	err := service.Delete("123")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestDeleteJabatan_HasDependencies(t *testing.T) {
	repo := &mockJabatanRepository{
		DeleteFn: func(id string) error {
			return errors.New("cannot delete jabatan with associated users")
		},
	}

	service := NewJabatanService(repo)

	err := service.Delete("123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "associated users")
}

func TestDeleteJabatan_ForeignKeyConstraint(t *testing.T) {
	repo := &mockJabatanRepository{
		DeleteFn: func(id string) error {
			return errors.New("foreign key constraint violation")
		},
	}

	service := NewJabatanService(repo)

	err := service.Delete("123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "constraint")
}

/*
=====================================
 TEST EDGE CASES
=====================================
*/

func TestJabatanService_CreateThenGetByID(t *testing.T) {
	var createdID string
	nama := "Test Manager"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			createdID = id
			return nil
		},
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			if id == createdID {
				return &dto.JabatanResponse{
					ID:          id,
					NamaJabatan: nama,
				}, nil
			}
			return nil, errors.New("not found")
		},
	}

	service := NewJabatanService(repo)

	// Create
	created, err := service.Create(dto.CreateJabatanRequest{
		NamaJabatan: &nama,
	})

	assert.NoError(t, err)
	assert.NotNil(t, created)

	// GetByID
	retrieved, err := service.GetByID(created.ID)

	assert.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.NamaJabatan, retrieved.NamaJabatan)
}

func TestJabatanService_CreateUpdateDelete(t *testing.T) {
	storage := make(map[string]*dto.JabatanResponse)
	nama := "Initial Name"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			storage[id] = &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: *req.NamaJabatan,
			}
			return nil
		},
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			if val, ok := storage[id]; ok {
				return val, nil
			}
			return nil, errors.New("not found")
		},
		UpdateFn: func(id string, jabatan dto.JabatanResponse) error {
			storage[id] = &jabatan
			return nil
		},
		DeleteFn: func(id string) error {
			delete(storage, id)
			return nil
		},
	}

	service := NewJabatanService(repo)

	// 1. Create
	created, err := service.Create(dto.CreateJabatanRequest{NamaJabatan: &nama})
	assert.NoError(t, err)

	// 2. Update
	newName := "Updated Name"
	updated, err := service.Update(created.ID, dto.UpdateJabatanRequest{NamaJabatan: &newName})
	assert.NoError(t, err)
	assert.Equal(t, newName, updated.NamaJabatan)

	// 3. Delete
	err = service.Delete(created.ID)
	assert.NoError(t, err)

	// 4. Verify deleted
	_, err = service.GetByID(created.ID)
	assert.Error(t, err)
}

func TestJabatanService_SpecialCharactersInName(t *testing.T) {
	testCases := []string{
		"Manager & Supervisor",
		"IT (Information Technology)",
		"Level 1/2/3",
		"Director - Operations",
		"VP, Finance & Accounting",
		"Staf IT: Junior",
	}

	for _, name := range testCases {
		t.Run("Name: "+name, func(t *testing.T) {
			repo := &mockJabatanRepository{
				CreateFn: func(req dto.CreateJabatanRequest, id string) error {
					return nil
				},
				GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
					return &dto.JabatanResponse{
						ID:          id,
						NamaJabatan: name,
					}, nil
				},
			}

			service := NewJabatanService(repo)

			result, err := service.Create(dto.CreateJabatanRequest{NamaJabatan: &name})

			assert.NoError(t, err)
			assert.Equal(t, name, result.NamaJabatan)
		})
	}
}

func TestJabatanService_VeryLongName(t *testing.T) {
	longName := "Chief Executive Officer and Director of Technology, Innovation, Digital Transformation, and Information Systems Management"

	repo := &mockJabatanRepository{
		CreateFn: func(req dto.CreateJabatanRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
			return &dto.JabatanResponse{
				ID:          id,
				NamaJabatan: longName,
			}, nil
		},
	}

	service := NewJabatanService(repo)

	result, err := service.Create(dto.CreateJabatanRequest{NamaJabatan: &longName})

	assert.NoError(t, err)
	assert.Equal(t, longName, result.NamaJabatan)
}

func TestJabatanService_UnicodeCharacters(t *testing.T) {
	testCases := []string{
		"المدير التنفيذي", // Arabic
		"主管经理",          // Chinese
		"マネージャー",        // Japanese
		"Менеджер",        // Russian
	}

	for _, name := range testCases {
		t.Run("Unicode: "+name, func(t *testing.T) {
			repo := &mockJabatanRepository{
				CreateFn: func(req dto.CreateJabatanRequest, id string) error {
					return nil
				},
				GetByIDFn: func(id string) (*dto.JabatanResponse, error) {
					return &dto.JabatanResponse{
						ID:          id,
						NamaJabatan: name,
					}, nil
				},
			}

			service := NewJabatanService(repo)

			result, err := service.Create(dto.CreateJabatanRequest{NamaJabatan: &name})

			assert.NoError(t, err)
			assert.Equal(t, name, result.NamaJabatan)
		})
	}
}