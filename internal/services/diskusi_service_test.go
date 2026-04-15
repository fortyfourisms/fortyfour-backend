package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repositories for Diskusi ────────────────────────────────────────────

type mockDiskusiRepo struct {
	CreateFn      func(diskusi *models.Diskusi) error
	FindByMateriFn func(idMateri string) ([]models.Diskusi, error)
	FindByIDFn    func(id string) (*models.Diskusi, error)
	UpdateFn      func(diskusi *models.Diskusi) error
	DeleteFn      func(id string) error
	FindRepliesFn func(idParent string) ([]models.Diskusi, error)
}

func (m *mockDiskusiRepo) Create(diskusi *models.Diskusi) error { return m.CreateFn(diskusi) }
func (m *mockDiskusiRepo) FindByMateri(idMateri string) ([]models.Diskusi, error) {
	return m.FindByMateriFn(idMateri)
}
func (m *mockDiskusiRepo) FindByID(id string) (*models.Diskusi, error) { return m.FindByIDFn(id) }
func (m *mockDiskusiRepo) Update(diskusi *models.Diskusi) error        { return m.UpdateFn(diskusi) }
func (m *mockDiskusiRepo) Delete(id string) error                     { return m.DeleteFn(id) }
func (m *mockDiskusiRepo) FindReplies(idParent string) ([]models.Diskusi, error) {
	if m.FindRepliesFn != nil {
		return m.FindRepliesFn(idParent)
	}
	return nil, nil
}

type mockUserRepoDiskusi struct {
	FindByIDFn func(id string) (*models.User, error)
}

func (m *mockUserRepoDiskusi) Create(user *models.User) error { return nil }
func (m *mockUserRepoDiskusi) FindByUsername(username string) (*models.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoDiskusi) FindByEmail(email string) (*models.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoDiskusi) FindByID(id string) (*models.User, error) { return m.FindByIDFn(id) }
func (m *mockUserRepoDiskusi) FindAll() ([]models.User, error)          { return nil, nil }
func (m *mockUserRepoDiskusi) Update(user *models.User) error           { return nil }
func (m *mockUserRepoDiskusi) UpdateWithPhoto(user *models.User) error  { return nil }
func (m *mockUserRepoDiskusi) UpdatePassword(id, hp string) error       { return nil }
func (m *mockUserRepoDiskusi) GetPasswordByID(id string) (string, error) {
	return "", errors.New("not found")
}
func (m *mockUserRepoDiskusi) Delete(id string) error                         { return nil }
func (m *mockUserRepoDiskusi) EmailExists(email string, ex *string) (bool, error) { return false, nil }
func (m *mockUserRepoDiskusi) UsernameExists(un string, ex *string) (bool, error) { return false, nil }
func (m *mockUserRepoDiskusi) SetMFA(uid string, s *string, e bool) error         { return nil }
func (m *mockUserRepoDiskusi) ExistsByPerusahaan(idP string) (bool, error)        { return false, nil }
func (m *mockUserRepoDiskusi) UpdateStatus(uid string, s models.UserStatus) error { return nil }
func (m *mockUserRepoDiskusi) IncrementLoginAttempts(uid string) (int, error)     { return 0, nil }
func (m *mockUserRepoDiskusi) ResetLoginAttempts(uid string) error                { return nil }
func (m *mockUserRepoDiskusi) UpdatePasswordChangedAt(uid string) error           { return nil }

/*
=====================================
 TEST CREATE DISKUSI
=====================================
*/

func TestCreateDiskusi_TopLevel(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		CreateFn: func(d *models.Diskusi) error {
			d.CreatedAt = now
			d.UpdatedAt = now
			return nil
		},
	}
	displayName := "John"
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "john", DisplayName: &displayName}, nil
		},
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	resp, err := svc.Create("m-1", "user-1", dto.CreateDiskusiRequest{
		Konten: "Halo semua",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Halo semua", resp.Konten)
	assert.Equal(t, "John", resp.NamaUser)
	assert.Nil(t, resp.IDParent)
}

func TestCreateDiskusi_Reply(t *testing.T) {
	now := time.Now()
	parentID := "d-parent"
	diskusiRepo := &mockDiskusiRepo{
		CreateFn: func(d *models.Diskusi) error {
			d.CreatedAt = now
			d.UpdatedAt = now
			return nil
		},
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "john"}, nil
		},
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	resp, err := svc.Create("m-1", "user-1", dto.CreateDiskusiRequest{
		IDParent: &parentID,
		Konten:   "Ini reply",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, &parentID, resp.IDParent)
}

func TestCreateDiskusi_RepoError(t *testing.T) {
	diskusiRepo := &mockDiskusiRepo{
		CreateFn: func(d *models.Diskusi) error { return errors.New("db error") },
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id, Username: "john"}, nil },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	resp, err := svc.Create("m-1", "user-1", dto.CreateDiskusiRequest{Konten: "Test"})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST GET BY MATERI (DISKUSI)
=====================================
*/

func TestGetDiskusiByMateri_Success(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		FindByMateriFn: func(idMateri string) ([]models.Diskusi, error) {
			return []models.Diskusi{
				{ID: "d1", IDMateri: idMateri, IDUser: "u1", Konten: "Hello", CreatedAt: now, UpdatedAt: now},
			}, nil
		},
		FindRepliesFn: func(idParent string) ([]models.Diskusi, error) {
			return []models.Diskusi{}, nil
		},
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "john"}, nil
		},
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	data, err := svc.GetByMateri("m-1")

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "Hello", data[0].Konten)
}

func TestGetDiskusiByMateri_WithReplies(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		FindByMateriFn: func(idMateri string) ([]models.Diskusi, error) {
			return []models.Diskusi{
				{ID: "d1", IDMateri: idMateri, IDUser: "u1", Konten: "Hello", CreatedAt: now, UpdatedAt: now},
			}, nil
		},
		FindRepliesFn: func(idParent string) ([]models.Diskusi, error) {
			return []models.Diskusi{
				{ID: "r1", IDMateri: "m-1", IDUser: "u2", Konten: "Reply!", CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "john"}, nil
		},
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	data, err := svc.GetByMateri("m-1")

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Len(t, data[0].Replies, 1)
	assert.Equal(t, "Reply!", data[0].Replies[0].Konten)
}

func TestGetDiskusiByMateri_Empty(t *testing.T) {
	diskusiRepo := &mockDiskusiRepo{
		FindByMateriFn: func(idMateri string) ([]models.Diskusi, error) {
			return []models.Diskusi{}, nil
		},
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return nil, errors.New("not found") },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	data, err := svc.GetByMateri("m-1")

	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

/*
=====================================
 TEST UPDATE DISKUSI
=====================================
*/

func TestUpdateDiskusi_Success(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		FindByIDFn: func(id string) (*models.Diskusi, error) {
			return &models.Diskusi{ID: id, IDMateri: "m-1", IDUser: "user-1", Konten: "Old", CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(d *models.Diskusi) error { return nil },
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "john"}, nil
		},
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	resp, err := svc.Update("d-1", "user-1", dto.UpdateDiskusiRequest{Konten: "Updated"})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.Konten)
}

func TestUpdateDiskusi_NotFound(t *testing.T) {
	diskusiRepo := &mockDiskusiRepo{
		FindByIDFn: func(id string) (*models.Diskusi, error) { return nil, errors.New("not found") },
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return nil, errors.New("not found") },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	resp, err := svc.Update("invalid", "user-1", dto.UpdateDiskusiRequest{Konten: "Test"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestUpdateDiskusi_WrongUser(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		FindByIDFn: func(id string) (*models.Diskusi, error) {
			return &models.Diskusi{ID: id, IDUser: "user-other", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id, Username: "john"}, nil },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	resp, err := svc.Update("d-1", "user-1", dto.UpdateDiskusiRequest{Konten: "Hack"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "milik sendiri")
}

/*
=====================================
 TEST DELETE DISKUSI
=====================================
*/

func TestDeleteDiskusi_ByAdmin(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		FindByIDFn: func(id string) (*models.Diskusi, error) {
			return &models.Diskusi{ID: id, IDUser: "user-other", CreatedAt: now, UpdatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return nil },
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id}, nil },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	err := svc.Delete("d-1", "admin-1", "admin") // admin bisa hapus semua
	assert.NoError(t, err)
}

func TestDeleteDiskusi_ByOwner(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		FindByIDFn: func(id string) (*models.Diskusi, error) {
			return &models.Diskusi{ID: id, IDUser: "user-1", CreatedAt: now, UpdatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return nil },
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id}, nil },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	err := svc.Delete("d-1", "user-1", "user") // owner bisa hapus miliknya
	assert.NoError(t, err)
}

func TestDeleteDiskusi_NotOwnerNotAdmin(t *testing.T) {
	now := time.Now()
	diskusiRepo := &mockDiskusiRepo{
		FindByIDFn: func(id string) (*models.Diskusi, error) {
			return &models.Diskusi{ID: id, IDUser: "user-other", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id}, nil },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	err := svc.Delete("d-1", "user-1", "user") // bukan owner, bukan admin
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "milik sendiri")
}

func TestDeleteDiskusi_NotFound(t *testing.T) {
	diskusiRepo := &mockDiskusiRepo{
		FindByIDFn: func(id string) (*models.Diskusi, error) { return nil, errors.New("not found") },
	}
	userRepo := &mockUserRepoDiskusi{
		FindByIDFn: func(id string) (*models.User, error) { return nil, errors.New("not found") },
	}
	svc := NewDiskusiService(diskusiRepo, userRepo)

	err := svc.Delete("invalid", "user-1", "user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}
