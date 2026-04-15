package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repositories for Materi ─────────────────────────────────────────────

type mockMateriRepo struct {
	CreateFn        func(m *models.Materi) error
	FindByIDFn      func(id string) (*models.Materi, error)
	FindByKelasFn   func(idKelas string) ([]models.Materi, error)
	UpdateFn        func(m *models.Materi) error
	DeleteFn        func(id string) error
	ReorderUrutanFn func(idKelas string) error
}

func (m *mockMateriRepo) Create(materi *models.Materi) error         { return m.CreateFn(materi) }
func (m *mockMateriRepo) FindByID(id string) (*models.Materi, error) { return m.FindByIDFn(id) }
func (m *mockMateriRepo) FindByKelas(idKelas string) ([]models.Materi, error) {
	if m.FindByKelasFn != nil {
		return m.FindByKelasFn(idKelas)
	}
	return nil, nil
}
func (m *mockMateriRepo) Update(materi *models.Materi) error { return m.UpdateFn(materi) }
func (m *mockMateriRepo) Delete(id string) error             { return m.DeleteFn(id) }
func (m *mockMateriRepo) ReorderUrutan(idKelas string) error {
	if m.ReorderUrutanFn != nil {
		return m.ReorderUrutanFn(idKelas)
	}
	return nil
}

type mockKelasRepoMateri struct {
	FindByIDFn func(id string) (*models.Kelas, error)
}

func (m *mockKelasRepoMateri) Create(k *models.Kelas) error                       { return nil }
func (m *mockKelasRepoMateri) FindByID(id string) (*models.Kelas, error)          { return m.FindByIDFn(id) }
func (m *mockKelasRepoMateri) FindAll(onlyPublished bool) ([]models.Kelas, error) { return nil, nil }
func (m *mockKelasRepoMateri) Update(k *models.Kelas) error                       { return nil }
func (m *mockKelasRepoMateri) Delete(id string) error                             { return nil }

type mockProgressRepoMateri struct {
	FindByUserAndMateriFn func(idUser, idMateri string) (*models.UserMateriProgress, error)
	UpsertFn              func(p *models.UserMateriProgress) error
}

func (m *mockProgressRepoMateri) Upsert(p *models.UserMateriProgress) error {
	if m.UpsertFn != nil {
		return m.UpsertFn(p)
	}
	return nil
}
func (m *mockProgressRepoMateri) FindByUserAndMateri(idUser, idMateri string) (*models.UserMateriProgress, error) {
	return m.FindByUserAndMateriFn(idUser, idMateri)
}
func (m *mockProgressRepoMateri) FindByUserAndKelas(idUser, idKelas string) ([]models.UserMateriProgress, error) {
	return nil, nil
}
func (m *mockProgressRepoMateri) HasCompletedAllMateri(idUser, idKelas string) (bool, error) {
	return false, nil
}

func newTestMateriService(repo *mockMateriRepo, kelasRepo *mockKelasRepoMateri, progressRepo *mockProgressRepoMateri) *MateriService {
	if progressRepo == nil {
		progressRepo = &mockProgressRepoMateri{
			FindByUserAndMateriFn: func(string, string) (*models.UserMateriProgress, error) {
				return nil, errors.New("not found")
			},
		}
	}
	return NewMateriService(repo, kelasRepo, progressRepo, nil)
}

/*
=====================================
 TEST CREATE MATERI
=====================================
*/

func TestCreateMateri_SuccessVideo(t *testing.T) {
	ytID := "abc123"
	durasi := 600
	repo := &mockMateriRepo{
		CreateFn: func(m *models.Materi) error { return nil },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id}, nil
		},
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("kelas-1", dto.CreateMateriRequest{
		Judul:       "Intro Go",
		Tipe:        "video",
		Urutan:      1,
		YoutubeID:   &ytID,
		DurasiDetik: &durasi,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Intro Go", resp.Judul)
	assert.Equal(t, models.MateriTipeVideo, resp.Tipe)
	assert.Equal(t, &ytID, resp.YoutubeID)
}

func TestCreateMateri_SuccessTeks(t *testing.T) {
	konten := "<p>Hello World</p>"
	repo := &mockMateriRepo{
		CreateFn: func(m *models.Materi) error { return nil },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("kelas-1", dto.CreateMateriRequest{
		Judul:      "Artikel Go",
		Tipe:       "teks",
		Urutan:     1,
		KontenHTML: &konten,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, models.MateriTipeTeks, resp.Tipe)
	assert.Equal(t, &konten, resp.KontenHTML)
}

func TestCreateMateri_KelasNotFound(t *testing.T) {
	repo := &mockMateriRepo{}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, errors.New("not found") },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("invalid", dto.CreateMateriRequest{Judul: "Test", Tipe: "video", Urutan: 1})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "kelas tidak ditemukan")
}

func TestCreateMateri_JudulKosong(t *testing.T) {
	repo := &mockMateriRepo{}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("kelas-1", dto.CreateMateriRequest{Judul: "", Tipe: "video", Urutan: 1})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "judul wajib diisi")
}

func TestCreateMateri_VideoMissingYoutubeID(t *testing.T) {
	repo := &mockMateriRepo{}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("kelas-1", dto.CreateMateriRequest{Judul: "Test", Tipe: "video", Urutan: 1})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "youtube_id wajib")
}

func TestCreateMateri_TeksMissingKontenHTML(t *testing.T) {
	repo := &mockMateriRepo{}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("kelas-1", dto.CreateMateriRequest{Judul: "Test", Tipe: "teks", Urutan: 1})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "konten_html wajib")
}

func TestCreateMateri_TipeTidakValid(t *testing.T) {
	repo := &mockMateriRepo{}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("kelas-1", dto.CreateMateriRequest{Judul: "Test", Tipe: "audio", Urutan: 1})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tipe tidak valid")
}

func TestCreateMateri_RepoError(t *testing.T) {
	ytID := "abc123"
	repo := &mockMateriRepo{
		CreateFn: func(m *models.Materi) error { return errors.New("db error") },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.Create("kelas-1", dto.CreateMateriRequest{
		Judul: "Test", Tipe: "video", Urutan: 1, YoutubeID: &ytID,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST UPDATE MATERI
=====================================
*/

func TestUpdateMateri_Success(t *testing.T) {
	now := time.Now()
	ytID := "abc123"
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, IDKelas: "k-1", Judul: "Old", Tipe: models.MateriTipeVideo, YoutubeID: &ytID, CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(m *models.Materi) error { return nil },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)
	newJudul := "Updated"

	resp, err := svc.Update("m-1", dto.UpdateMateriRequest{Judul: &newJudul})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.Judul)
}

func TestUpdateMateri_NotFound(t *testing.T) {
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) { return nil, errors.New("not found") },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)
	judul := "New"

	resp, err := svc.Update("invalid", dto.UpdateMateriRequest{Judul: &judul})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestUpdateMateri_JudulKosong(t *testing.T) {
	now := time.Now()
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, IDKelas: "k-1", Judul: "Old", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)
	empty := "   "

	resp, err := svc.Update("m-1", dto.UpdateMateriRequest{Judul: &empty})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateMateri_RepoError(t *testing.T) {
	now := time.Now()
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, IDKelas: "k-1", Judul: "Old", CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(m *models.Materi) error { return errors.New("db error") },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)
	judul := "New"

	resp, err := svc.Update("m-1", dto.UpdateMateriRequest{Judul: &judul})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST DELETE MATERI
=====================================
*/

func TestDeleteMateri_Success(t *testing.T) {
	now := time.Now()
	reorderCalled := false
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, IDKelas: "k-1", CreatedAt: now, UpdatedAt: now}, nil
		},
		DeleteFn:        func(id string) error { return nil },
		ReorderUrutanFn: func(idKelas string) error { reorderCalled = true; return nil },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	err := svc.Delete("m-1")

	assert.NoError(t, err)
	assert.True(t, reorderCalled)
}

func TestDeleteMateri_NotFound(t *testing.T) {
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) { return nil, errors.New("not found") },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	err := svc.Delete("invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestDeleteMateri_RepoError(t *testing.T) {
	now := time.Now()
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, IDKelas: "k-1", CreatedAt: now, UpdatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return errors.New("db error") },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	err := svc.Delete("m-1")

	assert.Error(t, err)
}

/*
=====================================
 TEST UPDATE PROGRESS
=====================================
*/

func TestUpdateProgress_VideoAutoComplete(t *testing.T) {
	now := time.Now()
	durasi := 100
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, Tipe: models.MateriTipeVideo, DurasiDetik: &durasi, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	progressRepo := &mockProgressRepoMateri{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.UserMateriProgress, error) {
			return nil, errors.New("not found") // buat baru
		},
		UpsertFn: func(p *models.UserMateriProgress) error { return nil },
	}
	svc := newTestMateriService(repo, kelasRepo, progressRepo)
	watched := 80 // >= 80%

	resp, err := svc.UpdateProgress("user-1", "m-1", dto.UpdateProgressRequest{
		LastWatchedSeconds: &watched,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsCompleted) // auto-completed at 80%
}

func TestUpdateProgress_VideoNotYetComplete(t *testing.T) {
	now := time.Now()
	durasi := 100
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, Tipe: models.MateriTipeVideo, DurasiDetik: &durasi, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	progressRepo := &mockProgressRepoMateri{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.UserMateriProgress, error) {
			return nil, errors.New("not found")
		},
		UpsertFn: func(p *models.UserMateriProgress) error { return nil },
	}
	svc := newTestMateriService(repo, kelasRepo, progressRepo)
	watched := 50 // < 80%

	resp, err := svc.UpdateProgress("user-1", "m-1", dto.UpdateProgressRequest{
		LastWatchedSeconds: &watched,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsCompleted)
	assert.Equal(t, 50, resp.LastWatchedSeconds)
}

func TestUpdateProgress_TeksMarkCompleted(t *testing.T) {
	now := time.Now()
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, Tipe: models.MateriTipeTeks, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	progressRepo := &mockProgressRepoMateri{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.UserMateriProgress, error) {
			return nil, errors.New("not found")
		},
		UpsertFn: func(p *models.UserMateriProgress) error { return nil },
	}
	svc := newTestMateriService(repo, kelasRepo, progressRepo)

	resp, err := svc.UpdateProgress("user-1", "m-1", dto.UpdateProgressRequest{IsCompleted: true})

	assert.NoError(t, err)
	assert.True(t, resp.IsCompleted)
	assert.NotNil(t, resp.CompletedAt)
}

func TestUpdateProgress_AlreadyCompletedNoRegress(t *testing.T) {
	now := time.Now()
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, Tipe: models.MateriTipeVideo, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	completedAt := now
	progressRepo := &mockProgressRepoMateri{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.UserMateriProgress, error) {
			return &models.UserMateriProgress{
				ID: "p-1", IDUser: idUser, IDMateri: idMateri,
				IsCompleted: true, CompletedAt: &completedAt, LastWatchedSeconds: 100,
			}, nil
		},
	}
	svc := newTestMateriService(repo, kelasRepo, progressRepo)

	resp, err := svc.UpdateProgress("user-1", "m-1", dto.UpdateProgressRequest{IsCompleted: false})

	assert.NoError(t, err)
	assert.True(t, resp.IsCompleted) // stays completed
}

func TestUpdateProgress_MateriNotFound(t *testing.T) {
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) { return nil, errors.New("not found") },
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, nil },
	}
	svc := newTestMateriService(repo, kelasRepo, nil)

	resp, err := svc.UpdateProgress("user-1", "invalid", dto.UpdateProgressRequest{IsCompleted: true})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestUpdateProgress_UpsertError(t *testing.T) {
	now := time.Now()
	repo := &mockMateriRepo{
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, Tipe: models.MateriTipeTeks, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	kelasRepo := &mockKelasRepoMateri{
		FindByIDFn: func(id string) (*models.Kelas, error) { return &models.Kelas{ID: id}, nil },
	}
	progressRepo := &mockProgressRepoMateri{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.UserMateriProgress, error) {
			return nil, errors.New("not found")
		},
		UpsertFn: func(p *models.UserMateriProgress) error { return errors.New("db error") },
	}
	svc := newTestMateriService(repo, kelasRepo, progressRepo)

	resp, err := svc.UpdateProgress("user-1", "m-1", dto.UpdateProgressRequest{IsCompleted: true})

	assert.Error(t, err)
	assert.Nil(t, resp)
}
