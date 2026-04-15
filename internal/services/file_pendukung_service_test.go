package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repositories for FilePendukung ──────────────────────────────────────

type mockFilePendukungRepo struct {
	CreateFn       func(fp *models.FilePendukung) error
	FindByMateriFn func(idMateri string) ([]models.FilePendukung, error)
	FindByIDFn     func(id string) (*models.FilePendukung, error)
	DeleteFn       func(id string) error
}

func (m *mockFilePendukungRepo) Create(fp *models.FilePendukung) error { return m.CreateFn(fp) }
func (m *mockFilePendukungRepo) FindByMateri(idMateri string) ([]models.FilePendukung, error) {
	return m.FindByMateriFn(idMateri)
}
func (m *mockFilePendukungRepo) FindByID(id string) (*models.FilePendukung, error) {
	return m.FindByIDFn(id)
}
func (m *mockFilePendukungRepo) Delete(id string) error { return m.DeleteFn(id) }

type mockMateriRepoFP struct {
	FindByIDFn func(id string) (*models.Materi, error)
}

func (m *mockMateriRepoFP) Create(materi *models.Materi) error                  { return nil }
func (m *mockMateriRepoFP) FindByID(id string) (*models.Materi, error)          { return m.FindByIDFn(id) }
func (m *mockMateriRepoFP) FindByKelas(idKelas string) ([]models.Materi, error) { return nil, nil }
func (m *mockMateriRepoFP) Update(materi *models.Materi) error                  { return nil }
func (m *mockMateriRepoFP) Delete(id string) error                              { return nil }
func (m *mockMateriRepoFP) ReorderUrutan(idKelas string) error                  { return nil }

/*
=====================================
 TEST CREATE FILE PENDUKUNG
=====================================
*/

func TestCreateFilePendukung_Success(t *testing.T) {
	fpRepo := &mockFilePendukungRepo{
		CreateFn: func(fp *models.FilePendukung) error { return nil },
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id}, nil
		},
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	resp, err := svc.Create("m-1", "document.pdf", "/uploads/doc.pdf", 1024)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "document.pdf", resp.NamaFile)
	assert.Equal(t, "/uploads/doc.pdf", resp.FilePath)
	assert.Equal(t, int64(1024), resp.Ukuran)
	assert.NotEmpty(t, resp.ID)
}

func TestCreateFilePendukung_MateriNotFound(t *testing.T) {
	fpRepo := &mockFilePendukungRepo{}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return nil, errors.New("not found") },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	resp, err := svc.Create("invalid", "doc.pdf", "/path", 100)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "materi tidak ditemukan")
}

func TestCreateFilePendukung_RepoError(t *testing.T) {
	fpRepo := &mockFilePendukungRepo{
		CreateFn: func(fp *models.FilePendukung) error { return errors.New("db error") },
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return &models.Materi{ID: id}, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	resp, err := svc.Create("m-1", "doc.pdf", "/path", 100)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST GET BY MATERI
=====================================
*/

func TestGetByMateri_Success(t *testing.T) {
	now := time.Now()
	fpRepo := &mockFilePendukungRepo{
		FindByMateriFn: func(idMateri string) ([]models.FilePendukung, error) {
			return []models.FilePendukung{
				{ID: "fp-1", IDMateri: idMateri, NamaFile: "doc1.pdf", FilePath: "/p1", Ukuran: 100, CreatedAt: now},
				{ID: "fp-2", IDMateri: idMateri, NamaFile: "doc2.pdf", FilePath: "/p2", Ukuran: 200, CreatedAt: now},
			}, nil
		},
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return &models.Materi{ID: id}, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	data, err := svc.GetByMateri("m-1")

	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestGetByMateri_Empty(t *testing.T) {
	fpRepo := &mockFilePendukungRepo{
		FindByMateriFn: func(idMateri string) ([]models.FilePendukung, error) {
			return []models.FilePendukung{}, nil
		},
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return &models.Materi{ID: id}, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	data, err := svc.GetByMateri("m-1")

	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestGetByMateri_RepoError(t *testing.T) {
	fpRepo := &mockFilePendukungRepo{
		FindByMateriFn: func(idMateri string) ([]models.FilePendukung, error) {
			return nil, errors.New("db error")
		},
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return &models.Materi{ID: id}, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	data, err := svc.GetByMateri("m-1")

	assert.Error(t, err)
	assert.Nil(t, data)
}

/*
=====================================
 TEST DELETE FILE PENDUKUNG
=====================================
*/

func TestDeleteFilePendukung_Success(t *testing.T) {
	now := time.Now()
	fpRepo := &mockFilePendukungRepo{
		FindByIDFn: func(id string) (*models.FilePendukung, error) {
			return &models.FilePendukung{ID: id, CreatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return nil },
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return &models.Materi{ID: id}, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	err := svc.Delete("fp-1")
	assert.NoError(t, err)
}

func TestDeleteFilePendukung_NotFound(t *testing.T) {
	fpRepo := &mockFilePendukungRepo{
		FindByIDFn: func(id string) (*models.FilePendukung, error) { return nil, errors.New("not found") },
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return nil, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	err := svc.Delete("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

/*
=====================================
 TEST FIND BY ID
=====================================
*/

func TestFindByID_FilePendukung_Success(t *testing.T) {
	now := time.Now()
	fpRepo := &mockFilePendukungRepo{
		FindByIDFn: func(id string) (*models.FilePendukung, error) {
			return &models.FilePendukung{ID: id, NamaFile: "doc.pdf", CreatedAt: now}, nil
		},
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return nil, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	fp, err := svc.FindByID("fp-1")

	assert.NoError(t, err)
	assert.Equal(t, "doc.pdf", fp.NamaFile)
}

func TestFindByID_FilePendukung_NotFound(t *testing.T) {
	fpRepo := &mockFilePendukungRepo{
		FindByIDFn: func(id string) (*models.FilePendukung, error) { return nil, errors.New("not found") },
	}
	materiRepo := &mockMateriRepoFP{
		FindByIDFn: func(id string) (*models.Materi, error) { return nil, nil },
	}
	svc := NewFilePendukungService(fpRepo, materiRepo, nil)

	fp, err := svc.FindByID("invalid")

	assert.Error(t, err)
	assert.Nil(t, fp)
}
