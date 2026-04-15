package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// ── Mock Repositories for Sertifikat ─────────────────────────────────────────

type mockSertifikatRepo struct {
	CreateFn             func(s *models.Sertifikat) error
	FindByUserAndKelasFn func(idUser, idKelas string) (*models.Sertifikat, error)
	FindByIDFn           func(id string) (*models.Sertifikat, error)
	FindByUserFn         func(idUser string) ([]models.Sertifikat, error)
}

func (m *mockSertifikatRepo) Create(s *models.Sertifikat) error { return m.CreateFn(s) }
func (m *mockSertifikatRepo) FindByUserAndKelas(idUser, idKelas string) (*models.Sertifikat, error) {
	return m.FindByUserAndKelasFn(idUser, idKelas)
}
func (m *mockSertifikatRepo) FindByID(id string) (*models.Sertifikat, error) {
	return m.FindByIDFn(id)
}
func (m *mockSertifikatRepo) FindByUser(idUser string) ([]models.Sertifikat, error) {
	return m.FindByUserFn(idUser)
}

type mockKelasRepoSert struct {
	FindByIDFn func(id string) (*models.Kelas, error)
}

func (m *mockKelasRepoSert) Create(k *models.Kelas) error                      { return nil }
func (m *mockKelasRepoSert) FindByID(id string) (*models.Kelas, error)          { return m.FindByIDFn(id) }
func (m *mockKelasRepoSert) FindAll(onlyPublished bool) ([]models.Kelas, error) { return nil, nil }
func (m *mockKelasRepoSert) Update(k *models.Kelas) error                      { return nil }
func (m *mockKelasRepoSert) Delete(id string) error                            { return nil }

type mockProgressRepoSert struct {
	HasCompletedAllMateriFn func(idUser, idKelas string) (bool, error)
}

func (m *mockProgressRepoSert) Upsert(p *models.UserMateriProgress) error { return nil }
func (m *mockProgressRepoSert) FindByUserAndMateri(idUser, idMateri string) (*models.UserMateriProgress, error) {
	return nil, errors.New("not found")
}
func (m *mockProgressRepoSert) FindByUserAndKelas(idUser, idKelas string) ([]models.UserMateriProgress, error) {
	return nil, nil
}
func (m *mockProgressRepoSert) HasCompletedAllMateri(idUser, idKelas string) (bool, error) {
	return m.HasCompletedAllMateriFn(idUser, idKelas)
}

type mockAttemptRepoSert struct {
	HasPassedAllKuisInKelasFn func(idUser, idKelas string) (bool, error)
	FindByUserAndKuisFn       func(idUser, idKuis string) ([]models.KuisAttempt, error)
}

func (m *mockAttemptRepoSert) Create(a *models.KuisAttempt) error { return nil }
func (m *mockAttemptRepoSert) FindByID(id string) (*models.KuisAttempt, error) {
	return nil, errors.New("not found")
}
func (m *mockAttemptRepoSert) FindByUserAndKuis(idUser, idKuis string) ([]models.KuisAttempt, error) {
	if m.FindByUserAndKuisFn != nil {
		return m.FindByUserAndKuisFn(idUser, idKuis)
	}
	return nil, nil
}
func (m *mockAttemptRepoSert) FindLatestByUserAndKuis(idUser, idKuis string) (*models.KuisAttempt, error) {
	return nil, errors.New("not found")
}
func (m *mockAttemptRepoSert) Finish(id string, skor float64, totalBenar int, isPassed bool, jawaban []models.KuisJawaban) error {
	return nil
}
func (m *mockAttemptRepoSert) HasPassedAllKuisInKelas(idUser, idKelas string) (bool, error) {
	return m.HasPassedAllKuisInKelasFn(idUser, idKelas)
}
func (m *mockAttemptRepoSert) FindJawabanByAttempt(idAttempt string) ([]models.KuisJawaban, error) {
	return nil, nil
}

type mockKuisRepoSert struct {
	FindFinalByKelasFn func(idKelas string) (*models.Kuis, error)
}

func (m *mockKuisRepoSert) Create(kuis *models.Kuis) error           { return nil }
func (m *mockKuisRepoSert) FindByID(id string) (*models.Kuis, error) { return nil, errors.New("not found") }
func (m *mockKuisRepoSert) FindByKelas(idKelas string) ([]models.Kuis, error) { return nil, nil }
func (m *mockKuisRepoSert) FindByMateri(idMateri string) (*models.Kuis, error) {
	return nil, errors.New("not found")
}
func (m *mockKuisRepoSert) FindFinalByKelas(idKelas string) (*models.Kuis, error) {
	return m.FindFinalByKelasFn(idKelas)
}
func (m *mockKuisRepoSert) Update(kuis *models.Kuis) error { return nil }
func (m *mockKuisRepoSert) Delete(id string) error         { return nil }

type mockUserRepoSert struct {
	FindByIDFn func(id string) (*models.User, error)
}

func (m *mockUserRepoSert) Create(user *models.User) error { return nil }
func (m *mockUserRepoSert) FindByUsername(username string) (*models.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoSert) FindByEmail(email string) (*models.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepoSert) FindByID(id string) (*models.User, error) { return m.FindByIDFn(id) }
func (m *mockUserRepoSert) FindAll() ([]models.User, error)          { return nil, nil }
func (m *mockUserRepoSert) Update(user *models.User) error           { return nil }
func (m *mockUserRepoSert) UpdateWithPhoto(user *models.User) error  { return nil }
func (m *mockUserRepoSert) UpdatePassword(id, hp string) error       { return nil }
func (m *mockUserRepoSert) GetPasswordByID(id string) (string, error) {
	return "", errors.New("not found")
}
func (m *mockUserRepoSert) Delete(id string) error                           { return nil }
func (m *mockUserRepoSert) EmailExists(email string, ex *string) (bool, error)   { return false, nil }
func (m *mockUserRepoSert) UsernameExists(un string, ex *string) (bool, error)   { return false, nil }
func (m *mockUserRepoSert) SetMFA(uid string, s *string, e bool) error           { return nil }
func (m *mockUserRepoSert) ExistsByPerusahaan(idP string) (bool, error)          { return false, nil }
func (m *mockUserRepoSert) UpdateStatus(uid string, s models.UserStatus) error   { return nil }
func (m *mockUserRepoSert) IncrementLoginAttempts(uid string) (int, error)       { return 0, nil }
func (m *mockUserRepoSert) ResetLoginAttempts(uid string) error                  { return nil }
func (m *mockUserRepoSert) UpdatePasswordChangedAt(uid string) error             { return nil }

/*
=====================================
 TEST GENERATE SERTIFIKAT
=====================================
*/

func TestGenerateSertifikat_AlreadyExists(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return &models.Sertifikat{
				ID: "cert-1", IDKelas: idKelas, IDUser: idUser,
				NomorSertifikat: "CERT/20260101/abc", NamaPeserta: "John",
				NamaKelas: "Go Class", TanggalTerbit: now, CreatedAt: now,
			}, nil
		},
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "cert-1", resp.ID) // returns existing
}

func TestGenerateSertifikat_KelasNotFound(t *testing.T) {
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, errors.New("not found") },
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, nil, nil, nil, nil)

	resp, err := svc.Generate("user-1", "invalid")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "kelas tidak ditemukan")
}

func TestGenerateSertifikat_UserNotFound(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Go Class", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) { return nil, errors.New("not found") },
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, nil, nil, nil, userRepo)

	resp, err := svc.Generate("invalid", "kelas-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "user tidak ditemukan")
}

func TestGenerateSertifikat_MateriBelumSelesai(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Go", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id, Username: "john"}, nil },
	}
	progressRepo := &mockProgressRepoSert{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return false, nil },
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, progressRepo, nil, nil, userRepo)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "selesaikan semua materi")
}

func TestGenerateSertifikat_KuisNonFinalBelumLulus(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Go", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id, Username: "john"}, nil },
	}
	progressRepo := &mockProgressRepoSert{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	attemptRepo := &mockAttemptRepoSert{
		HasPassedAllKuisInKelasFn: func(idUser, idKelas string) (bool, error) { return false, nil },
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, progressRepo, attemptRepo, nil, userRepo)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "lulus semua kuis")
}

func TestGenerateSertifikat_KuisAkhirBelumLulus(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Go", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) { return &models.User{ID: id, Username: "john"}, nil },
	}
	progressRepo := &mockProgressRepoSert{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	attemptRepo := &mockAttemptRepoSert{
		HasPassedAllKuisInKelasFn: func(idUser, idKelas string) (bool, error) { return true, nil },
		FindByUserAndKuisFn: func(idUser, idKuis string) ([]models.KuisAttempt, error) {
			return []models.KuisAttempt{{ID: "a-1", IsPassed: false}}, nil // belum lulus
		},
	}
	kuisRepo := &mockKuisRepoSert{
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return &models.Kuis{ID: "kuis-final", IDKelas: idKelas, IsFinal: true, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, progressRepo, attemptRepo, kuisRepo, userRepo)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "belum lulus kuis akhir")
}

func TestGenerateSertifikat_Success_TanpaKuisAkhir(t *testing.T) {
	now := time.Now()
	displayName := "John Doe"
	var createdSertifikat *models.Sertifikat

	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(s *models.Sertifikat) error {
			createdSertifikat = s
			return nil
		},
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Kelas Golang", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "johndoe", DisplayName: &displayName}, nil
		},
	}
	progressRepo := &mockProgressRepoSert{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	attemptRepo := &mockAttemptRepoSert{
		HasPassedAllKuisInKelasFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	kuisRepo := &mockKuisRepoSert{
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return nil, errors.New("not found") // tidak ada kuis akhir
		},
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, progressRepo, attemptRepo, kuisRepo, userRepo)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "John Doe", resp.NamaPeserta)
	assert.Equal(t, "Kelas Golang", resp.NamaKelas)
	assert.NotEmpty(t, resp.NomorSertifikat)
	assert.Contains(t, resp.NomorSertifikat, "CERT/")
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.TanggalTerbit)
	// Verifikasi create dipanggil
	assert.NotNil(t, createdSertifikat)
	assert.Equal(t, "kelas-1", createdSertifikat.IDKelas)
	assert.Equal(t, "user-1", createdSertifikat.IDUser)
}

func TestGenerateSertifikat_Success_DenganKuisAkhirLulus(t *testing.T) {
	now := time.Now()

	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(s *models.Sertifikat) error { return nil },
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Kelas Security", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "janedoe"}, nil
		},
	}
	progressRepo := &mockProgressRepoSert{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	attemptRepo := &mockAttemptRepoSert{
		HasPassedAllKuisInKelasFn: func(idUser, idKelas string) (bool, error) { return true, nil },
		FindByUserAndKuisFn: func(idUser, idKuis string) ([]models.KuisAttempt, error) {
			return []models.KuisAttempt{
				{ID: "att-1", IsPassed: false},
				{ID: "att-2", IsPassed: true}, // lulus di percobaan kedua
			}, nil
		},
	}
	kuisRepo := &mockKuisRepoSert{
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return &models.Kuis{ID: "final-kuis-1", IDKelas: idKelas, IsFinal: true, CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, progressRepo, attemptRepo, kuisRepo, userRepo)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Tanpa DisplayName → pakai Username
	assert.Equal(t, "janedoe", resp.NamaPeserta)
	assert.Equal(t, "Kelas Security", resp.NamaKelas)
	assert.Contains(t, resp.NomorSertifikat, "CERT/")
}

func TestGenerateSertifikat_Success_UsesDisplayNameOverUsername(t *testing.T) {
	now := time.Now()
	displayName := "Admin Budi"

	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(s *models.Sertifikat) error { return nil },
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Go Advanced", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "budi123", DisplayName: &displayName}, nil
		},
	}
	progressRepo := &mockProgressRepoSert{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	attemptRepo := &mockAttemptRepoSert{
		HasPassedAllKuisInKelasFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	kuisRepo := &mockKuisRepoSert{
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, progressRepo, attemptRepo, kuisRepo, userRepo)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Admin Budi", resp.NamaPeserta, "harus pakai DisplayName, bukan Username")
}

func TestGenerateSertifikat_CreateRepoError(t *testing.T) {
	now := time.Now()

	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(s *models.Sertifikat) error {
			return errors.New("database full")
		},
	}
	kelasRepo := &mockKelasRepoSert{
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Go", CreatedAt: now, UpdatedAt: now}, nil
		},
	}
	userRepo := &mockUserRepoSert{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "test"}, nil
		},
	}
	progressRepo := &mockProgressRepoSert{
		HasCompletedAllMateriFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	attemptRepo := &mockAttemptRepoSert{
		HasPassedAllKuisInKelasFn: func(idUser, idKelas string) (bool, error) { return true, nil },
	}
	kuisRepo := &mockKuisRepoSert{
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewSertifikatService(sertRepo, kelasRepo, progressRepo, attemptRepo, kuisRepo, userRepo)

	resp, err := svc.Generate("user-1", "kelas-1")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "database full")
}

/*
=====================================
 TEST GET BY USER AND KELAS
=====================================
*/

func TestGetSertifikatByUserAndKelas_Success(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return &models.Sertifikat{
				ID: "cert-1", NomorSertifikat: "CERT/123", IDKelas: idKelas, IDUser: idUser,
				NamaPeserta: "John", NamaKelas: "Go", TanggalTerbit: now, CreatedAt: now,
			}, nil
		},
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	resp, err := svc.GetByUserAndKelas("user-1", "kelas-1")

	assert.NoError(t, err)
	assert.Equal(t, "cert-1", resp.ID)
}

func TestGetSertifikatByUserAndKelas_NotFound(t *testing.T) {
	sertRepo := &mockSertifikatRepo{
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	resp, err := svc.GetByUserAndKelas("user-1", "invalid")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "belum tersedia")
}

/*
=====================================
 TEST GET BY ID
=====================================
*/

func TestGetSertifikatByID_Success(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByIDFn: func(id string) (*models.Sertifikat, error) {
			return &models.Sertifikat{
				ID: id, NomorSertifikat: "CERT/123",
				NamaPeserta: "John", NamaKelas: "Go", TanggalTerbit: now, CreatedAt: now,
			}, nil
		},
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	resp, err := svc.GetByID("cert-1")

	assert.NoError(t, err)
	assert.Equal(t, "CERT/123", resp.NomorSertifikat)
}

func TestGetSertifikatByID_NotFound(t *testing.T) {
	sertRepo := &mockSertifikatRepo{
		FindByIDFn: func(id string) (*models.Sertifikat, error) { return nil, errors.New("not found") },
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	resp, err := svc.GetByID("invalid")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

/*
=====================================
 TEST GET BY USER
=====================================
*/

func TestGetSertifikatByUser_Success(t *testing.T) {
	now := time.Now()
	sertRepo := &mockSertifikatRepo{
		FindByUserFn: func(idUser string) ([]models.Sertifikat, error) {
			return []models.Sertifikat{
				{ID: "c1", NomorSertifikat: "CERT/1", IDUser: idUser, NamaPeserta: "John", NamaKelas: "Go1", TanggalTerbit: now, CreatedAt: now},
				{ID: "c2", NomorSertifikat: "CERT/2", IDUser: idUser, NamaPeserta: "John", NamaKelas: "Go2", TanggalTerbit: now, CreatedAt: now},
			}, nil
		},
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	data, err := svc.GetByUser("user-1")

	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestGetSertifikatByUser_Empty(t *testing.T) {
	sertRepo := &mockSertifikatRepo{
		FindByUserFn: func(idUser string) ([]models.Sertifikat, error) {
			return []models.Sertifikat{}, nil
		},
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	data, err := svc.GetByUser("user-1")

	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestGetSertifikatByUser_RepoError(t *testing.T) {
	sertRepo := &mockSertifikatRepo{
		FindByUserFn: func(idUser string) ([]models.Sertifikat, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewSertifikatService(sertRepo, nil, nil, nil, nil, nil)

	data, err := svc.GetByUser("user-1")

	assert.Error(t, err)
	assert.Nil(t, data)
}
