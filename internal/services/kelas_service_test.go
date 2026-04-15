package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repositories for Kelas ──────────────────────────────────────────────

type mockKelasRepo struct {
	CreateFn  func(kelas *models.Kelas) error
	FindByIDFn func(id string) (*models.Kelas, error)
	FindAllFn func(onlyPublished bool) ([]models.Kelas, error)
	UpdateFn  func(kelas *models.Kelas) error
	DeleteFn  func(id string) error
}

func (m *mockKelasRepo) Create(kelas *models.Kelas) error          { return m.CreateFn(kelas) }
func (m *mockKelasRepo) FindByID(id string) (*models.Kelas, error) { return m.FindByIDFn(id) }
func (m *mockKelasRepo) FindAll(onlyPublished bool) ([]models.Kelas, error) {
	return m.FindAllFn(onlyPublished)
}
func (m *mockKelasRepo) Update(kelas *models.Kelas) error { return m.UpdateFn(kelas) }
func (m *mockKelasRepo) Delete(id string) error           { return m.DeleteFn(id) }

type mockMateriRepoKelas struct {
	FindByKelasFn func(idKelas string) ([]models.Materi, error)
	// stubs below
	CreateFn        func(m *models.Materi) error
	FindByIDFn      func(id string) (*models.Materi, error)
	UpdateFn        func(m *models.Materi) error
	DeleteFn        func(id string) error
	ReorderUrutanFn func(idKelas string) error
}

func (m *mockMateriRepoKelas) Create(materi *models.Materi) error { return nil }
func (m *mockMateriRepoKelas) FindByID(id string) (*models.Materi, error) {
	return nil, errors.New("not found")
}
func (m *mockMateriRepoKelas) FindByKelas(idKelas string) ([]models.Materi, error) {
	return m.FindByKelasFn(idKelas)
}
func (m *mockMateriRepoKelas) Update(materi *models.Materi) error { return nil }
func (m *mockMateriRepoKelas) Delete(id string) error             { return nil }
func (m *mockMateriRepoKelas) ReorderUrutan(idKelas string) error { return nil }

type mockProgressRepoKelas struct {
	FindByUserAndKelasFn    func(idUser, idKelas string) ([]models.UserMateriProgress, error)
	HasCompletedAllMateriFn func(idUser, idKelas string) (bool, error)
}

func (m *mockProgressRepoKelas) Upsert(p *models.UserMateriProgress) error { return nil }
func (m *mockProgressRepoKelas) FindByUserAndMateri(idUser, idMateri string) (*models.UserMateriProgress, error) {
	return nil, errors.New("not found")
}
func (m *mockProgressRepoKelas) FindByUserAndKelas(idUser, idKelas string) ([]models.UserMateriProgress, error) {
	return m.FindByUserAndKelasFn(idUser, idKelas)
}
func (m *mockProgressRepoKelas) HasCompletedAllMateri(idUser, idKelas string) (bool, error) {
	if m.HasCompletedAllMateriFn != nil {
		return m.HasCompletedAllMateriFn(idUser, idKelas)
	}
	return false, nil
}

type mockKuisRepoKelas struct {
	FindByKelasFn    func(idKelas string) ([]models.Kuis, error)
	FindByMateriFn   func(idMateri string) (*models.Kuis, error)
	FindFinalByKelasFn func(idKelas string) (*models.Kuis, error)
}

func (m *mockKuisRepoKelas) Create(kuis *models.Kuis) error              { return nil }
func (m *mockKuisRepoKelas) FindByID(id string) (*models.Kuis, error)    { return nil, errors.New("not found") }
func (m *mockKuisRepoKelas) FindByKelas(idKelas string) ([]models.Kuis, error) {
	return m.FindByKelasFn(idKelas)
}
func (m *mockKuisRepoKelas) FindByMateri(idMateri string) (*models.Kuis, error) {
	if m.FindByMateriFn != nil {
		return m.FindByMateriFn(idMateri)
	}
	return nil, errors.New("not found")
}
func (m *mockKuisRepoKelas) FindFinalByKelas(idKelas string) (*models.Kuis, error) {
	if m.FindFinalByKelasFn != nil {
		return m.FindFinalByKelasFn(idKelas)
	}
	return nil, errors.New("not found")
}
func (m *mockKuisRepoKelas) Update(kuis *models.Kuis) error { return nil }
func (m *mockKuisRepoKelas) Delete(id string) error         { return nil }

type mockAttemptRepoKelas struct {
	FindByUserAndKuisFn func(idUser, idKuis string) ([]models.KuisAttempt, error)
}

func (m *mockAttemptRepoKelas) Create(a *models.KuisAttempt) error { return nil }
func (m *mockAttemptRepoKelas) FindByID(id string) (*models.KuisAttempt, error) {
	return nil, errors.New("not found")
}
func (m *mockAttemptRepoKelas) FindByUserAndKuis(idUser, idKuis string) ([]models.KuisAttempt, error) {
	if m.FindByUserAndKuisFn != nil {
		return m.FindByUserAndKuisFn(idUser, idKuis)
	}
	return nil, nil
}
func (m *mockAttemptRepoKelas) FindLatestByUserAndKuis(idUser, idKuis string) (*models.KuisAttempt, error) {
	return nil, errors.New("not found")
}
func (m *mockAttemptRepoKelas) Finish(id string, skor float64, totalBenar int, isPassed bool, jawaban []models.KuisJawaban) error {
	return nil
}
func (m *mockAttemptRepoKelas) HasPassedAllKuisInKelas(idUser, idKelas string) (bool, error) {
	return false, nil
}
func (m *mockAttemptRepoKelas) FindJawabanByAttempt(idAttempt string) ([]models.KuisJawaban, error) {
	return nil, nil
}

type mockSertifikatRepoKelas struct {
	FindByUserAndKelasFn func(idUser, idKelas string) (*models.Sertifikat, error)
}

func (m *mockSertifikatRepoKelas) Create(s *models.Sertifikat) error { return nil }
func (m *mockSertifikatRepoKelas) FindByUserAndKelas(idUser, idKelas string) (*models.Sertifikat, error) {
	if m.FindByUserAndKelasFn != nil {
		return m.FindByUserAndKelasFn(idUser, idKelas)
	}
	return nil, errors.New("not found")
}
func (m *mockSertifikatRepoKelas) FindByID(id string) (*models.Sertifikat, error) {
	return nil, errors.New("not found")
}
func (m *mockSertifikatRepoKelas) FindByUser(idUser string) ([]models.Sertifikat, error) {
	return nil, nil
}

type mockFPRepoKelas struct {
	FindByMateriFn func(idMateri string) ([]models.FilePendukung, error)
}

func (m *mockFPRepoKelas) Create(fp *models.FilePendukung) error { return nil }
func (m *mockFPRepoKelas) FindByMateri(idMateri string) ([]models.FilePendukung, error) {
	if m.FindByMateriFn != nil {
		return m.FindByMateriFn(idMateri)
	}
	return nil, nil
}
func (m *mockFPRepoKelas) FindByID(id string) (*models.FilePendukung, error) {
	return nil, errors.New("not found")
}
func (m *mockFPRepoKelas) Delete(id string) error { return nil }

// ── helper: build KelasService with defaults ─────────────────────────────────

func newTestKelasService(
	repo *mockKelasRepo,
	materiRepo *mockMateriRepoKelas,
	progressRepo *mockProgressRepoKelas,
	kuisRepo *mockKuisRepoKelas,
	attemptRepo *mockAttemptRepoKelas,
	sertifikatRepo *mockSertifikatRepoKelas,
	fpRepo *mockFPRepoKelas,
) *KelasService {
	if materiRepo == nil {
		materiRepo = &mockMateriRepoKelas{FindByKelasFn: func(string) ([]models.Materi, error) { return nil, nil }}
	}
	if progressRepo == nil {
		progressRepo = &mockProgressRepoKelas{
			FindByUserAndKelasFn: func(string, string) ([]models.UserMateriProgress, error) { return nil, nil },
		}
	}
	if kuisRepo == nil {
		kuisRepo = &mockKuisRepoKelas{FindByKelasFn: func(string) ([]models.Kuis, error) { return nil, nil }}
	}
	if attemptRepo == nil {
		attemptRepo = &mockAttemptRepoKelas{}
	}
	if sertifikatRepo == nil {
		sertifikatRepo = &mockSertifikatRepoKelas{}
	}
	if fpRepo == nil {
		fpRepo = &mockFPRepoKelas{}
	}
	return NewKelasService(repo, materiRepo, progressRepo, kuisRepo, attemptRepo, sertifikatRepo, fpRepo, nil)
}

/*
=====================================
 TEST CREATE KELAS
=====================================
*/

func TestCreateKelas_Success(t *testing.T) {
	repo := &mockKelasRepo{
		CreateFn: func(kelas *models.Kelas) error { return nil },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	resp, err := svc.Create(dto.CreateKelasRequest{Judul: "Kelas Go"}, "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Kelas Go", resp.Judul)
	assert.Equal(t, models.KelasStatusDraft, resp.Status)
	assert.Equal(t, "user-1", resp.CreatedBy)
	assert.NotEmpty(t, resp.ID)
}

func TestCreateKelas_JudulKosong(t *testing.T) {
	repo := &mockKelasRepo{}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	resp, err := svc.Create(dto.CreateKelasRequest{Judul: ""}, "user-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "judul wajib diisi")
}

func TestCreateKelas_JudulWhitespace(t *testing.T) {
	repo := &mockKelasRepo{}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	resp, err := svc.Create(dto.CreateKelasRequest{Judul: "   "}, "user-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateKelas_RepoError(t *testing.T) {
	repo := &mockKelasRepo{
		CreateFn: func(kelas *models.Kelas) error { return errors.New("db error") },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	resp, err := svc.Create(dto.CreateKelasRequest{Judul: "Kelas Go"}, "user-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateKelas_WithDeskripsiAndThumbnail(t *testing.T) {
	repo := &mockKelasRepo{
		CreateFn: func(kelas *models.Kelas) error { return nil },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)
	desc := "Deskripsi kelas"
	thumb := "/img/thumb.jpg"

	resp, err := svc.Create(dto.CreateKelasRequest{
		Judul:     "Kelas Go",
		Deskripsi: &desc,
		Thumbnail: &thumb,
	}, "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, &desc, resp.Deskripsi)
	assert.Equal(t, &thumb, resp.Thumbnail)
}

/*
=====================================
 TEST UPDATE KELAS
=====================================
*/

func TestUpdateKelas_Success(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Old", Status: models.KelasStatusDraft, CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(kelas *models.Kelas) error { return nil },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)
	newJudul := "Updated"

	resp, err := svc.Update("kelas-1", dto.UpdateKelasRequest{Judul: &newJudul})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.Judul)
}

func TestUpdateKelas_NotFound(t *testing.T) {
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, errors.New("not found") },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)
	judul := "New"

	resp, err := svc.Update("invalid", dto.UpdateKelasRequest{Judul: &judul})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestUpdateKelas_JudulKosong(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Old", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)
	empty := "   "

	resp, err := svc.Update("kelas-1", dto.UpdateKelasRequest{Judul: &empty})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "judul tidak boleh kosong")
}

func TestUpdateKelas_PartialUpdate(t *testing.T) {
	now := time.Now()
	desc := "Old Desc"
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Old", Deskripsi: &desc, Status: models.KelasStatusDraft, CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(kelas *models.Kelas) error { return nil },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)
	newDesc := "New Desc"

	resp, err := svc.Update("kelas-1", dto.UpdateKelasRequest{Deskripsi: &newDesc})

	assert.NoError(t, err)
	assert.Equal(t, "Old", resp.Judul) // unchanged
	assert.Equal(t, &newDesc, resp.Deskripsi)
}

func TestUpdateKelas_ChangeStatus(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Kelas", Status: models.KelasStatusDraft, CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(kelas *models.Kelas) error { return nil },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)
	status := "published"

	resp, err := svc.Update("kelas-1", dto.UpdateKelasRequest{Status: &status})

	assert.NoError(t, err)
	assert.Equal(t, models.KelasStatusPublished, resp.Status)
}

func TestUpdateKelas_RepoError(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Old", CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(kelas *models.Kelas) error { return errors.New("db error") },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)
	judul := "New"

	resp, err := svc.Update("kelas-1", dto.UpdateKelasRequest{Judul: &judul})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST DELETE KELAS
=====================================
*/

func TestDeleteKelas_Success(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return nil },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	err := svc.Delete("kelas-1")

	assert.NoError(t, err)
}

func TestDeleteKelas_NotFound(t *testing.T) {
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, errors.New("not found") },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	err := svc.Delete("invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestDeleteKelas_RepoError(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return errors.New("db error") },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	err := svc.Delete("kelas-1")

	assert.Error(t, err)
}

/*
=====================================
 TEST GET ALL KELAS
=====================================
*/

func TestGetAllKelas_Success(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindAllFn: func(onlyPublished bool) ([]models.Kelas, error) {
			return []models.Kelas{
				{ID: "1", Judul: "Kelas A", Status: models.KelasStatusPublished, CreatedAt: now, UpdatedAt: now},
				{ID: "2", Judul: "Kelas B", Status: models.KelasStatusDraft, CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	data, err := svc.GetAll(false)

	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestGetAllKelas_OnlyPublished(t *testing.T) {
	now := time.Now()
	calledWithPublished := false
	repo := &mockKelasRepo{
		FindAllFn: func(onlyPublished bool) ([]models.Kelas, error) {
			calledWithPublished = onlyPublished
			return []models.Kelas{
				{ID: "1", Judul: "Kelas A", Status: models.KelasStatusPublished, CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	data, err := svc.GetAll(true)

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.True(t, calledWithPublished)
}

func TestGetAllKelas_Empty(t *testing.T) {
	repo := &mockKelasRepo{
		FindAllFn: func(onlyPublished bool) ([]models.Kelas, error) {
			return []models.Kelas{}, nil
		},
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	data, err := svc.GetAll(false)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 0)
}

func TestGetAllKelas_RepoError(t *testing.T) {
	repo := &mockKelasRepo{
		FindAllFn: func(onlyPublished bool) ([]models.Kelas, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	data, err := svc.GetAll(false)

	assert.Error(t, err)
	assert.Nil(t, data)
}

/*
=====================================
 TEST GET DETAIL KELAS
=====================================
*/

func TestGetDetailKelas_Success(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Kelas Go", Status: models.KelasStatusPublished, CreatedBy: "admin-1", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	materiRepo := &mockMateriRepoKelas{
		FindByKelasFn: func(idKelas string) ([]models.Materi, error) {
			return []models.Materi{
				{ID: "m-1", IDKelas: idKelas, Judul: "Materi 1", Tipe: models.MateriTipeVideo, Urutan: 1, CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	progressRepo := &mockProgressRepoKelas{
		FindByUserAndKelasFn: func(idUser, idKelas string) ([]models.UserMateriProgress, error) {
			return []models.UserMateriProgress{
				{IDMateri: "m-1", IsCompleted: true},
			}, nil
		},
	}
	kuisRepo := &mockKuisRepoKelas{
		FindByKelasFn: func(idKelas string) ([]models.Kuis, error) {
			return []models.Kuis{}, nil
		},
		FindByMateriFn: func(idMateri string) (*models.Kuis, error) {
			return nil, errors.New("not found")
		},
	}
	fpRepo := &mockFPRepoKelas{
		FindByMateriFn: func(idMateri string) ([]models.FilePendukung, error) {
			return nil, nil
		},
	}

	svc := newTestKelasService(repo, materiRepo, progressRepo, kuisRepo, nil, nil, fpRepo)

	resp, err := svc.GetDetail("kelas-1", "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Kelas Go", resp.Judul)
	assert.Len(t, resp.Materi, 1)
	assert.True(t, resp.Materi[0].IsCompleted)
}

func TestGetDetailKelas_NotFound(t *testing.T) {
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, errors.New("not found") },
	}
	svc := newTestKelasService(repo, nil, nil, nil, nil, nil, nil)

	resp, err := svc.GetDetail("invalid", "user-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestGetDetailKelas_WithoutUserID(t *testing.T) {
	now := time.Now()
	repo := &mockKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Kelas Go", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	materiRepo := &mockMateriRepoKelas{
		FindByKelasFn: func(idKelas string) ([]models.Materi, error) { return nil, nil },
	}
	kuisRepo := &mockKuisRepoKelas{
		FindByKelasFn: func(idKelas string) ([]models.Kuis, error) { return nil, nil },
	}

	svc := newTestKelasService(repo, materiRepo, nil, kuisRepo, nil, nil, nil)

	resp, err := svc.GetDetail("kelas-1", "")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.Progress) // no progress without userID
}
