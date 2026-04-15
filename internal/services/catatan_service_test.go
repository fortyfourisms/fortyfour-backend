package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repository for Catatan ──────────────────────────────────────────────

type mockCatatanRepo struct {
	UpsertFn              func(catatan *models.CatatanPribadi) error
	FindByUserAndMateriFn func(idUser, idMateri string) (*models.CatatanPribadi, error)
	DeleteFn              func(id string) error
}

func (m *mockCatatanRepo) Upsert(catatan *models.CatatanPribadi) error {
	return m.UpsertFn(catatan)
}
func (m *mockCatatanRepo) FindByUserAndMateri(idUser, idMateri string) (*models.CatatanPribadi, error) {
	return m.FindByUserAndMateriFn(idUser, idMateri)
}
func (m *mockCatatanRepo) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

/*
=====================================
 TEST UPSERT CATATAN
=====================================
*/

func TestUpsertCatatan_Success(t *testing.T) {
	now := time.Now()
	repo := &mockCatatanRepo{
		UpsertFn: func(catatan *models.CatatanPribadi) error { return nil },
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.CatatanPribadi, error) {
			return &models.CatatanPribadi{
				ID: "c-1", IDMateri: idMateri, IDUser: idUser, Konten: "Catatan saya",
				CreatedAt: now, UpdatedAt: now,
			}, nil
		},
	}
	svc := NewCatatanService(repo)

	resp, err := svc.Upsert("m-1", "user-1", dto.UpsertCatatanRequest{Konten: "Catatan saya"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Catatan saya", resp.Konten)
	assert.Equal(t, "m-1", resp.IDMateri)
}

func TestUpsertCatatan_KontenKosong(t *testing.T) {
	repo := &mockCatatanRepo{}
	svc := NewCatatanService(repo)

	resp, err := svc.Upsert("m-1", "user-1", dto.UpsertCatatanRequest{Konten: ""})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "konten tidak boleh kosong")
}

func TestUpsertCatatan_UpsertError(t *testing.T) {
	repo := &mockCatatanRepo{
		UpsertFn: func(catatan *models.CatatanPribadi) error { return errors.New("db error") },
	}
	svc := NewCatatanService(repo)

	resp, err := svc.Upsert("m-1", "user-1", dto.UpsertCatatanRequest{Konten: "Test"})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpsertCatatan_FindAfterUpsertError(t *testing.T) {
	repo := &mockCatatanRepo{
		UpsertFn: func(catatan *models.CatatanPribadi) error { return nil },
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.CatatanPribadi, error) {
			return nil, errors.New("not found after upsert")
		},
	}
	svc := NewCatatanService(repo)

	resp, err := svc.Upsert("m-1", "user-1", dto.UpsertCatatanRequest{Konten: "Test"})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST GET BY USER AND MATERI
=====================================
*/

func TestGetCatatanByUserAndMateri_Success(t *testing.T) {
	now := time.Now()
	repo := &mockCatatanRepo{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.CatatanPribadi, error) {
			return &models.CatatanPribadi{
				ID: "c-1", IDMateri: idMateri, IDUser: idUser, Konten: "Catatan penting",
				CreatedAt: now, UpdatedAt: now,
			}, nil
		},
	}
	svc := NewCatatanService(repo)

	resp, err := svc.GetByUserAndMateri("user-1", "m-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Catatan penting", resp.Konten)
}

func TestGetCatatanByUserAndMateri_NotFound(t *testing.T) {
	repo := &mockCatatanRepo{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.CatatanPribadi, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewCatatanService(repo)

	resp, err := svc.GetByUserAndMateri("user-1", "m-invalid")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}
