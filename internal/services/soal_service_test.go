package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repositories for Soal ───────────────────────────────────────────────

type mockSoalRepo struct {
	CreateFn             func(soal *models.Soal, pilihan []models.PilihanJawaban) error
	FindByIDFn           func(id string) (*models.Soal, error)
	FindByKuisFn         func(idKuis string) ([]models.Soal, error)
	UpdateFn             func(soal *models.Soal, pilihan []models.PilihanJawaban) error
	DeleteFn             func(id string) error
	FindPilihanByIDFn    func(idPilihan string) (*models.PilihanJawaban, error)
	FindCorrectPilihanFn func(idSoal string) (*models.PilihanJawaban, error)
}

func (m *mockSoalRepo) Create(soal *models.Soal, pilihan []models.PilihanJawaban) error {
	return m.CreateFn(soal, pilihan)
}
func (m *mockSoalRepo) FindByID(id string) (*models.Soal, error) { return m.FindByIDFn(id) }
func (m *mockSoalRepo) FindByKuis(idKuis string) ([]models.Soal, error) {
	return m.FindByKuisFn(idKuis)
}
func (m *mockSoalRepo) Update(soal *models.Soal, pilihan []models.PilihanJawaban) error {
	return m.UpdateFn(soal, pilihan)
}
func (m *mockSoalRepo) Delete(id string) error { return m.DeleteFn(id) }
func (m *mockSoalRepo) FindPilihanByID(idPilihan string) (*models.PilihanJawaban, error) {
	if m.FindPilihanByIDFn != nil {
		return m.FindPilihanByIDFn(idPilihan)
	}
	return nil, errors.New("not found")
}
func (m *mockSoalRepo) FindCorrectPilihan(idSoal string) (*models.PilihanJawaban, error) {
	if m.FindCorrectPilihanFn != nil {
		return m.FindCorrectPilihanFn(idSoal)
	}
	return nil, errors.New("not found")
}

type mockKuisRepoSoal struct {
	FindByIDFn func(id string) (*models.Kuis, error)
}

func (m *mockKuisRepoSoal) Create(kuis *models.Kuis) error                    { return nil }
func (m *mockKuisRepoSoal) FindByID(id string) (*models.Kuis, error)          { return m.FindByIDFn(id) }
func (m *mockKuisRepoSoal) FindByKelas(idKelas string) ([]models.Kuis, error) { return nil, nil }
func (m *mockKuisRepoSoal) FindByMateri(idMateri string) (*models.Kuis, error) {
	return nil, errors.New("not found")
}
func (m *mockKuisRepoSoal) FindFinalByKelas(idKelas string) (*models.Kuis, error) {
	return nil, errors.New("not found")
}
func (m *mockKuisRepoSoal) Update(kuis *models.Kuis) error { return nil }
func (m *mockKuisRepoSoal) Delete(id string) error         { return nil }

func validPilihan() []dto.CreatePilihanRequest {
	return []dto.CreatePilihanRequest{
		{Teks: "Pilihan A", IsCorrect: true, Urutan: 1},
		{Teks: "Pilihan B", IsCorrect: false, Urutan: 2},
	}
}

/*
=====================================
 TEST CREATE SOAL
=====================================
*/

func TestCreateSoal_Success(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{
		CreateFn: func(soal *models.Soal, pilihan []models.PilihanJawaban) error { return nil },
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "Apa itu Go?",
		Urutan:     1,
		Pilihan:    validPilihan(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Apa itu Go?", resp.Pertanyaan)
	assert.Len(t, resp.Pilihan, 2)
}

func TestCreateSoal_KuisNotFound(t *testing.T) {
	soalRepo := &mockSoalRepo{}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, errors.New("not found") },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("invalid", dto.CreateSoalRequest{
		Pertanyaan: "Test?", Urutan: 1, Pilihan: validPilihan(),
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "kuis tidak ditemukan")
}

func TestCreateSoal_PertanyaanKosong(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "   ", Urutan: 1, Pilihan: validPilihan(),
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "pertanyaan wajib diisi")
}

func TestCreateSoal_PilihanKurangDari2(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "Test?", Urutan: 1,
		Pilihan: []dto.CreatePilihanRequest{{Teks: "A", IsCorrect: true, Urutan: 1}},
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "minimal 2 pilihan")
}

func TestCreateSoal_PilihanLebihDari5(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	pilihan := make([]dto.CreatePilihanRequest, 6)
	for i := range pilihan {
		pilihan[i] = dto.CreatePilihanRequest{Teks: "Option", IsCorrect: i == 0, Urutan: i + 1}
	}

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "Test?", Urutan: 1, Pilihan: pilihan,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "maksimal")
}

func TestCreateSoal_TidakAdaJawabanBenar(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "Test?", Urutan: 1,
		Pilihan: []dto.CreatePilihanRequest{
			{Teks: "A", IsCorrect: false, Urutan: 1},
			{Teks: "B", IsCorrect: false, Urutan: 2},
		},
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tepat 1 pilihan jawaban yang benar")
}

func TestCreateSoal_LebihDari1JawabanBenar(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "Test?", Urutan: 1,
		Pilihan: []dto.CreatePilihanRequest{
			{Teks: "A", IsCorrect: true, Urutan: 1},
			{Teks: "B", IsCorrect: true, Urutan: 2},
		},
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "hanya boleh ada 1")
}

func TestCreateSoal_TeksPilihanKosong(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "Test?", Urutan: 1,
		Pilihan: []dto.CreatePilihanRequest{
			{Teks: "", IsCorrect: true, Urutan: 1},
			{Teks: "B", IsCorrect: false, Urutan: 2},
		},
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "teks pilihan jawaban tidak boleh kosong")
}

func TestCreateSoal_RepoError(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{
		CreateFn: func(soal *models.Soal, pilihan []models.PilihanJawaban) error {
			return errors.New("db error")
		},
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) {
			return &models.Kuis{ID: id, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	resp, err := svc.Create("kuis-1", dto.CreateSoalRequest{
		Pertanyaan: "Test?", Urutan: 1, Pilihan: validPilihan(),
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST UPDATE SOAL
=====================================
*/

func TestUpdateSoal_Success(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{
		FindByIDFn: func(id string) (*models.Soal, error) {
			return &models.Soal{ID: id, IDKuis: "k1", Pertanyaan: "Old?", Urutan: 1, CreatedAt: now}, nil
		},
		UpdateFn: func(soal *models.Soal, pilihan []models.PilihanJawaban) error { return nil },
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return &models.Kuis{ID: id}, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)
	newQ := "Updated?"

	resp, err := svc.Update("soal-1", dto.UpdateSoalRequest{Pertanyaan: &newQ})

	assert.NoError(t, err)
	assert.Equal(t, "Updated?", resp.Pertanyaan)
}

func TestUpdateSoal_NotFound(t *testing.T) {
	soalRepo := &mockSoalRepo{
		FindByIDFn: func(id string) (*models.Soal, error) { return nil, errors.New("not found") },
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)
	q := "New?"

	resp, err := svc.Update("invalid", dto.UpdateSoalRequest{Pertanyaan: &q})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestUpdateSoal_PertanyaanKosong(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{
		FindByIDFn: func(id string) (*models.Soal, error) {
			return &models.Soal{ID: id, IDKuis: "k1", Pertanyaan: "Old?", CreatedAt: now}, nil
		},
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return &models.Kuis{ID: id}, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)
	empty := "   "

	resp, err := svc.Update("soal-1", dto.UpdateSoalRequest{Pertanyaan: &empty})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

/*
=====================================
 TEST DELETE SOAL
=====================================
*/

func TestDeleteSoal_Success(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{
		FindByIDFn: func(id string) (*models.Soal, error) {
			return &models.Soal{ID: id, CreatedAt: now}, nil
		},
		DeleteFn: func(id string) error { return nil },
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	err := svc.Delete("soal-1")
	assert.NoError(t, err)
}

func TestDeleteSoal_NotFound(t *testing.T) {
	soalRepo := &mockSoalRepo{
		FindByIDFn: func(id string) (*models.Soal, error) { return nil, errors.New("not found") },
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	err := svc.Delete("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

/*
=====================================
 TEST GET BY KUIS
=====================================
*/

func TestGetByKuis_Success(t *testing.T) {
	now := time.Now()
	soalRepo := &mockSoalRepo{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) {
			return []models.Soal{
				{ID: "s1", IDKuis: idKuis, Pertanyaan: "Q1", Urutan: 1, CreatedAt: now,
					Pilihan: []models.PilihanJawaban{{ID: "p1", Teks: "A", IsCorrect: true, Urutan: 1}}},
			}, nil
		},
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	data, err := svc.GetByKuis("kuis-1")

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "Q1", data[0].Pertanyaan)
	assert.Len(t, data[0].Pilihan, 1)
}

func TestGetByKuis_Empty(t *testing.T) {
	soalRepo := &mockSoalRepo{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) { return []models.Soal{}, nil },
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	data, err := svc.GetByKuis("kuis-1")

	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestGetByKuis_RepoError(t *testing.T) {
	soalRepo := &mockSoalRepo{
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) { return nil, errors.New("db error") },
	}
	kuisRepo := &mockKuisRepoSoal{
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, nil },
	}
	svc := NewSoalService(soalRepo, kuisRepo, nil)

	data, err := svc.GetByKuis("kuis-1")

	assert.Error(t, err)
	assert.Nil(t, data)
}
