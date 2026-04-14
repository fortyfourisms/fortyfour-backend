package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repositories for Kuis ───────────────────────────────────────────────

type mockKuisAttemptRepo struct {
	CreateFn                 func(a *models.KuisAttempt) error
	FindByIDFn               func(id string) (*models.KuisAttempt, error)
	FindByUserAndKuisFn      func(idUser, idKuis string) ([]models.KuisAttempt, error)
	FindLatestByUserAndKuisFn func(idUser, idKuis string) (*models.KuisAttempt, error)
	FinishFn                 func(id string, skor float64, totalBenar int, isPassed bool, jawaban []models.KuisJawaban) error
	HasPassedAllKuisInKelasFn func(idUser, idKelas string) (bool, error)
	FindJawabanByAttemptFn   func(idAttempt string) ([]models.KuisJawaban, error)
}

func (m *mockKuisAttemptRepo) Create(a *models.KuisAttempt) error { return m.CreateFn(a) }
func (m *mockKuisAttemptRepo) FindByID(id string) (*models.KuisAttempt, error) {
	return m.FindByIDFn(id)
}
func (m *mockKuisAttemptRepo) FindByUserAndKuis(idUser, idKuis string) ([]models.KuisAttempt, error) {
	if m.FindByUserAndKuisFn != nil {
		return m.FindByUserAndKuisFn(idUser, idKuis)
	}
	return nil, nil
}
func (m *mockKuisAttemptRepo) FindLatestByUserAndKuis(idUser, idKuis string) (*models.KuisAttempt, error) {
	return m.FindLatestByUserAndKuisFn(idUser, idKuis)
}
func (m *mockKuisAttemptRepo) Finish(id string, skor float64, totalBenar int, isPassed bool, jawaban []models.KuisJawaban) error {
	return m.FinishFn(id, skor, totalBenar, isPassed, jawaban)
}
func (m *mockKuisAttemptRepo) HasPassedAllKuisInKelas(idUser, idKelas string) (bool, error) {
	if m.HasPassedAllKuisInKelasFn != nil {
		return m.HasPassedAllKuisInKelasFn(idUser, idKelas)
	}
	return true, nil
}
func (m *mockKuisAttemptRepo) FindJawabanByAttempt(idAttempt string) ([]models.KuisJawaban, error) {
	if m.FindJawabanByAttemptFn != nil {
		return m.FindJawabanByAttemptFn(idAttempt)
	}
	return nil, nil
}

type mockSoalRepoKuis struct {
	FindByKuisFn       func(idKuis string) ([]models.Soal, error)
	FindPilihanByIDFn  func(idPilihan string) (*models.PilihanJawaban, error)
	FindCorrectPilihanFn func(idSoal string) (*models.PilihanJawaban, error)
}

func (m *mockSoalRepoKuis) Create(soal *models.Soal, pilihan []models.PilihanJawaban) error {
	return nil
}
func (m *mockSoalRepoKuis) FindByID(id string) (*models.Soal, error) {
	return nil, errors.New("not found")
}
func (m *mockSoalRepoKuis) FindByKuis(idKuis string) ([]models.Soal, error) {
	return m.FindByKuisFn(idKuis)
}
func (m *mockSoalRepoKuis) Update(soal *models.Soal, pilihan []models.PilihanJawaban) error {
	return nil
}
func (m *mockSoalRepoKuis) Delete(id string) error { return nil }
func (m *mockSoalRepoKuis) FindPilihanByID(idPilihan string) (*models.PilihanJawaban, error) {
	return m.FindPilihanByIDFn(idPilihan)
}
func (m *mockSoalRepoKuis) FindCorrectPilihan(idSoal string) (*models.PilihanJawaban, error) {
	return m.FindCorrectPilihanFn(idSoal)
}

type mockKuisRepoForKuis struct {
	CreateFn           func(kuis *models.Kuis) error
	FindByIDFn         func(id string) (*models.Kuis, error)
	FindByKelasFn      func(idKelas string) ([]models.Kuis, error)
	FindByMateriFn     func(idMateri string) (*models.Kuis, error)
	FindFinalByKelasFn func(idKelas string) (*models.Kuis, error)
	UpdateFn           func(kuis *models.Kuis) error
	DeleteFn           func(id string) error
}

func (m *mockKuisRepoForKuis) Create(kuis *models.Kuis) error { return m.CreateFn(kuis) }
func (m *mockKuisRepoForKuis) FindByID(id string) (*models.Kuis, error) { return m.FindByIDFn(id) }
func (m *mockKuisRepoForKuis) FindByKelas(idKelas string) ([]models.Kuis, error) {
	if m.FindByKelasFn != nil {
		return m.FindByKelasFn(idKelas)
	}
	return nil, nil
}
func (m *mockKuisRepoForKuis) FindByMateri(idMateri string) (*models.Kuis, error) {
	if m.FindByMateriFn != nil {
		return m.FindByMateriFn(idMateri)
	}
	return nil, errors.New("not found")
}
func (m *mockKuisRepoForKuis) FindFinalByKelas(idKelas string) (*models.Kuis, error) {
	if m.FindFinalByKelasFn != nil {
		return m.FindFinalByKelasFn(idKelas)
	}
	return nil, errors.New("not found")
}
func (m *mockKuisRepoForKuis) Update(kuis *models.Kuis) error { return m.UpdateFn(kuis) }
func (m *mockKuisRepoForKuis) Delete(id string) error         { return m.DeleteFn(id) }

type mockProgressRepoKuis struct {
	HasCompletedAllMateriFn func(idUser, idKelas string) (bool, error)
}

func (m *mockProgressRepoKuis) Upsert(p *models.UserMateriProgress) error { return nil }
func (m *mockProgressRepoKuis) FindByUserAndMateri(idUser, idMateri string) (*models.UserMateriProgress, error) {
	return nil, errors.New("not found")
}
func (m *mockProgressRepoKuis) FindByUserAndKelas(idUser, idKelas string) ([]models.UserMateriProgress, error) {
	return nil, nil
}
func (m *mockProgressRepoKuis) HasCompletedAllMateri(idUser, idKelas string) (bool, error) {
	return m.HasCompletedAllMateriFn(idUser, idKelas)
}

/*
=====================================
 TEST CREATE KUIS
=====================================
*/

func TestCreateKuis_Success(t *testing.T) {
	kuisRepo := &mockKuisRepoForKuis{
		CreateFn: func(kuis *models.Kuis) error { return nil },
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	resp, err := svc.CreateKuis("kelas-1", dto.CreateKuisRequest{
		Judul:        "Kuis 1",
		PassingGrade: 70,
		Urutan:       1,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Kuis 1", resp.Judul)
	assert.Equal(t, float64(70), resp.PassingGrade)
}

func TestCreateKuis_FinalSuccess(t *testing.T) {
	kuisRepo := &mockKuisRepoForKuis{
		CreateFn: func(kuis *models.Kuis) error { return nil },
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return nil, errors.New("not found") // belum ada final
		},
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	resp, err := svc.CreateKuis("kelas-1", dto.CreateKuisRequest{
		Judul: "Kuis Akhir", PassingGrade: 80, IsFinal: true, Urutan: 1,
	})

	assert.NoError(t, err)
	assert.True(t, resp.IsFinal)
	assert.Nil(t, resp.IDMateri) // final tidak terikat materi
}

func TestCreateKuis_DuplicateFinal(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return &models.Kuis{ID: "existing", IsFinal: true, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	resp, err := svc.CreateKuis("kelas-1", dto.CreateKuisRequest{
		Judul: "Kuis Akhir 2", PassingGrade: 80, IsFinal: true, Urutan: 2,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "sudah ada kuis akhir")
}

func TestCreateKuis_RepoError(t *testing.T) {
	kuisRepo := &mockKuisRepoForKuis{
		CreateFn:           func(kuis *models.Kuis) error { return errors.New("db error") },
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) { return nil, errors.New("not found") },
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	resp, err := svc.CreateKuis("kelas-1", dto.CreateKuisRequest{
		Judul: "Kuis 1", PassingGrade: 70, Urutan: 1,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST UPDATE KUIS
=====================================
*/

func TestUpdateKuis_Success(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, IDKelas: "k-1", Judul: "Old", PassingGrade: 70, CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(kuis *models.Kuis) error { return nil },
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)
	newJudul := "Updated"

	resp, err := svc.UpdateKuis("kuis-1", dto.UpdateKuisRequest{Judul: &newJudul})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.Judul)
}

func TestUpdateKuis_NotFound(t *testing.T) {
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, errors.New("not found") },
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)
	judul := "New"

	resp, err := svc.UpdateKuis("invalid", dto.UpdateKuisRequest{Judul: &judul})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST DELETE KUIS
=====================================
*/

func TestDeleteKuis_Success(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return nil },
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	err := svc.DeleteKuis("kuis-1")
	assert.NoError(t, err)
}

func TestDeleteKuis_NotFound(t *testing.T) {
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, errors.New("not found") },
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	err := svc.DeleteKuis("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

/*
=====================================
 TEST GET KUIS BY KELAS
=====================================
*/

func TestGetKuisByKelas_Success(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByKelasFn: func(idKelas string) ([]models.Kuis, error) {
			return []models.Kuis{
				{ID: "k1", IDKelas: idKelas, Judul: "Kuis 1", CreatedAt: now, UpdatedAt: now},
				{ID: "k2", IDKelas: idKelas, Judul: "Kuis 2", CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	data, err := svc.GetKuisByKelas("kelas-1")
	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestGetKuisByKelas_Empty(t *testing.T) {
	kuisRepo := &mockKuisRepoForKuis{
		FindByKelasFn: func(idKelas string) ([]models.Kuis, error) {
			return []models.Kuis{}, nil
		},
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	data, err := svc.GetKuisByKelas("kelas-1")
	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

/*
=====================================
 TEST START KUIS
=====================================
*/

func TestStartKuis_Success(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, IDKelas: "k-1", IsFinal: false, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	attemptRepo := &mockKuisAttemptRepo{
		FindLatestByUserAndKuisFn: func(idUser, idKuis string) (*models.KuisAttempt, error) {
			return nil, errors.New("not found") // no existing attempt
		},
		CreateFn: func(a *models.KuisAttempt) error { return nil },
	}
	soalRepo := &mockSoalRepoKuis{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) {
			return []models.Soal{
				{ID: "s1", Pertanyaan: "Q1", Pilihan: []models.PilihanJawaban{{ID: "p1", Teks: "A", Urutan: 1}}},
			}, nil
		},
	}
	svc := NewKuisService(attemptRepo, soalRepo, kuisRepo, nil, nil)

	resp, err := svc.Start("user-1", "kuis-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AttemptID)
	assert.Len(t, resp.Soal, 1)
}

func TestStartKuis_ResumeUnfinished(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, IDKelas: "k-1", IsFinal: false, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	attemptRepo := &mockKuisAttemptRepo{
		FindLatestByUserAndKuisFn: func(idUser, idKuis string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: "existing-attempt", IDUser: idUser, IDKuis: idKuis, FinishedAt: nil}, nil
		},
	}
	soalRepo := &mockSoalRepoKuis{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) {
			return []models.Soal{
				{ID: "s1", Pertanyaan: "Q1", Pilihan: []models.PilihanJawaban{{ID: "p1", Teks: "A", Urutan: 1}}},
			}, nil
		},
	}
	svc := NewKuisService(attemptRepo, soalRepo, kuisRepo, nil, nil)

	resp, err := svc.Start("user-1", "kuis-1")

	assert.NoError(t, err)
	assert.Equal(t, "existing-attempt", resp.AttemptID)
}

func TestStartKuis_NotFound(t *testing.T) {
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, errors.New("not found") },
	}
	svc := NewKuisService(nil, nil, kuisRepo, nil, nil)

	resp, err := svc.Start("user-1", "invalid")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestStartKuis_FinalMateriBelumSelesai(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, IDKelas: "k-1", IsFinal: true, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	progressRepo := &mockProgressRepoKuis{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return false, nil },
	}
	svc := NewKuisService(nil, nil, kuisRepo, progressRepo, nil)

	resp, err := svc.Start("user-1", "kuis-final")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "selesaikan semua materi")
}

func TestStartKuis_FinalKuisNonFinalBelumLulus(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, IDKelas: "k-1", IsFinal: true, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	progressRepo := &mockProgressRepoKuis{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	attemptRepo := &mockKuisAttemptRepo{
		HasPassedAllKuisInKelasFn: func(idUser, idKelas string) (bool, error) { return false, nil },
		FindLatestByUserAndKuisFn: func(idUser, idKuis string) (*models.KuisAttempt, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(a *models.KuisAttempt) error { return nil },
	}
	svc := NewKuisService(attemptRepo, nil, kuisRepo, progressRepo, nil)

	resp, err := svc.Start("user-1", "kuis-final")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "lulus semua kuis")
}

func TestStartKuis_NoSoal(t *testing.T) {
	now := time.Now()
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, IDKelas: "k-1", IsFinal: false, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	attemptRepo := &mockKuisAttemptRepo{
		FindLatestByUserAndKuisFn: func(idUser, idKuis string) (*models.KuisAttempt, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(a *models.KuisAttempt) error { return nil },
	}
	soalRepo := &mockSoalRepoKuis{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) { return []models.Soal{}, nil },
	}
	svc := NewKuisService(attemptRepo, soalRepo, kuisRepo, nil, nil)

	resp, err := svc.Start("user-1", "kuis-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "belum memiliki soal")
}

/*
=====================================
 TEST SUBMIT KUIS
=====================================
*/

func TestSubmitKuis_SuccessLulus(t *testing.T) {
	now := time.Now()
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: id, IDUser: "user-1", IDKuis: "kuis-1", FinishedAt: nil}, nil
		},
		FinishFn: func(id string, skor float64, totalBenar int, isPassed bool, jawaban []models.KuisJawaban) error {
			return nil
		},
	}
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, PassingGrade: 70, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	soalRepo := &mockSoalRepoKuis{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) {
			return []models.Soal{
				{ID: "s1", Pertanyaan: "Q1"},
			}, nil
		},
		FindPilihanByIDFn: func(idPilihan string) (*models.PilihanJawaban, error) {
			return &models.PilihanJawaban{ID: idPilihan, IDSoal: "s1", IsCorrect: true}, nil
		},
		FindCorrectPilihanFn: func(idSoal string) (*models.PilihanJawaban, error) {
			return &models.PilihanJawaban{ID: "p-correct", IDSoal: idSoal}, nil
		},
	}
	svc := NewKuisService(attemptRepo, soalRepo, kuisRepo, nil, nil)

	resp, err := svc.Submit("user-1", "attempt-1", dto.SubmitKuisRequest{
		Jawaban: []dto.JawabanItem{{IDSoal: "s1", IDPilihan: "p1"}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, float64(100), resp.Skor)
	assert.True(t, resp.IsPassed)
}

func TestSubmitKuis_AttemptNotFound(t *testing.T) {
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) { return nil, errors.New("not found") },
	}
	svc := NewKuisService(attemptRepo, nil, nil, nil, nil)

	resp, err := svc.Submit("user-1", "invalid", dto.SubmitKuisRequest{})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestSubmitKuis_WrongUser(t *testing.T) {
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: id, IDUser: "user-other", IDKuis: "kuis-1"}, nil
		},
	}
	svc := NewKuisService(attemptRepo, nil, nil, nil, nil)

	resp, err := svc.Submit("user-1", "attempt-1", dto.SubmitKuisRequest{})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "bukan milik user")
}

func TestSubmitKuis_AlreadyFinished(t *testing.T) {
	finished := time.Now()
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: id, IDUser: "user-1", IDKuis: "kuis-1", FinishedAt: &finished}, nil
		},
	}
	svc := NewKuisService(attemptRepo, nil, nil, nil, nil)

	resp, err := svc.Submit("user-1", "attempt-1", dto.SubmitKuisRequest{})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "sudah pernah dikerjakan")
}

func TestSubmitKuis_JumlahJawabanTidakSesuai(t *testing.T) {
	now := time.Now()
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: id, IDUser: "user-1", IDKuis: "kuis-1", FinishedAt: nil}, nil
		},
	}
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, PassingGrade: 70, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	soalRepo := &mockSoalRepoKuis{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) {
			return []models.Soal{
				{ID: "s1", Pertanyaan: "Q1"},
				{ID: "s2", Pertanyaan: "Q2"},
			}, nil
		},
	}
	svc := NewKuisService(attemptRepo, soalRepo, kuisRepo, nil, nil)

	resp, err := svc.Submit("user-1", "attempt-1", dto.SubmitKuisRequest{
		Jawaban: []dto.JawabanItem{{IDSoal: "s1", IDPilihan: "p1"}}, // hanya 1, harusnya 2
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "jumlah jawaban tidak sesuai")
}

func TestSubmitKuis_DuplicateAnswer(t *testing.T) {
	now := time.Now()
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: id, IDUser: "user-1", IDKuis: "kuis-1", FinishedAt: nil}, nil
		},
	}
	kuisRepo := &mockKuisRepoForKuis{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, PassingGrade: 70, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	soalRepo := &mockSoalRepoKuis{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) {
			return []models.Soal{
				{ID: "s1", Pertanyaan: "Q1"},
				{ID: "s2", Pertanyaan: "Q2"},
			}, nil
		},
	}
	svc := NewKuisService(attemptRepo, soalRepo, kuisRepo, nil, nil)

	resp, err := svc.Submit("user-1", "attempt-1", dto.SubmitKuisRequest{
		Jawaban: []dto.JawabanItem{
			{IDSoal: "s1", IDPilihan: "p1"},
			{IDSoal: "s1", IDPilihan: "p2"}, // duplicate soal
		},
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "dijawab lebih dari satu kali")
}

/*
=====================================
 TEST GET RESULT
=====================================
*/

func TestGetResult_Success(t *testing.T) {
	now := time.Now()
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{
				ID: id, IDUser: "user-1", IDKuis: "kuis-1",
				Skor: 100, TotalSoal: 1, TotalBenar: 1, IsPassed: true, FinishedAt: &now,
			}, nil
		},
		FindJawabanByAttemptFn: func(idAttempt string) ([]models.KuisJawaban, error) {
			return []models.KuisJawaban{
				{ID: "j1", IDAttempt: idAttempt, IDSoal: "s1", IDPilihan: "p1", IsCorrect: true},
			}, nil
		},
	}
	soalRepo := &mockSoalRepoKuis{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) {
			return []models.Soal{{ID: "s1", Pertanyaan: "Q1"}}, nil
		},
		FindCorrectPilihanFn: func(idSoal string) (*models.PilihanJawaban, error) {
			return &models.PilihanJawaban{ID: "p1", IDSoal: idSoal}, nil
		},
	}
	svc := NewKuisService(attemptRepo, soalRepo, nil, nil, nil)

	resp, err := svc.GetResult("user-1", "attempt-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, float64(100), resp.Skor)
	assert.True(t, resp.IsPassed)
	assert.Len(t, resp.Detail, 1)
}

func TestGetResult_AttemptNotFound(t *testing.T) {
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) { return nil, errors.New("not found") },
	}
	svc := NewKuisService(attemptRepo, nil, nil, nil, nil)

	resp, err := svc.GetResult("user-1", "invalid")

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetResult_WrongUser(t *testing.T) {
	now := time.Now()
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: id, IDUser: "user-other", FinishedAt: &now}, nil
		},
	}
	svc := NewKuisService(attemptRepo, nil, nil, nil, nil)

	resp, err := svc.GetResult("user-1", "attempt-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "bukan milik user")
}

func TestGetResult_NotFinished(t *testing.T) {
	attemptRepo := &mockKuisAttemptRepo{
		FindByIDFn: func(id string) (*models.KuisAttempt, error) {
			return &models.KuisAttempt{ID: id, IDUser: "user-1", FinishedAt: nil}, nil
		},
	}
	svc := NewKuisService(attemptRepo, nil, nil, nil, nil)

	resp, err := svc.GetResult("user-1", "attempt-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "belum selesai")
}
