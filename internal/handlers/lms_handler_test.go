package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"

	"github.com/stretchr/testify/assert"
)

// ── inline mock repos untuk handler tests ────────────────────────────────────

type lmsKelasRepo struct {
	CreateFn  func(k *models.Kelas) error
	FindByIDFn func(id string) (*models.Kelas, error)
	FindAllFn func(onlyPublished bool) ([]models.Kelas, error)
	UpdateFn  func(k *models.Kelas) error
	DeleteFn  func(id string) error
}
func (m *lmsKelasRepo) Create(k *models.Kelas) error { return m.CreateFn(k) }
func (m *lmsKelasRepo) FindByID(id string) (*models.Kelas, error) { return m.FindByIDFn(id) }
func (m *lmsKelasRepo) FindAll(onlyPublished bool) ([]models.Kelas, error) { return m.FindAllFn(onlyPublished) }
func (m *lmsKelasRepo) Update(k *models.Kelas) error { return m.UpdateFn(k) }
func (m *lmsKelasRepo) Delete(id string) error { return m.DeleteFn(id) }

type lmsMateriRepo struct {
	CreateFn   func(m *models.Materi) error
	FindByIDFn func(id string) (*models.Materi, error)
	FindByKelasFn func(idKelas string) ([]models.Materi, error)
	UpdateFn   func(m *models.Materi) error
	DeleteFn   func(id string) error
	ReorderUrutanFn func(idKelas string) error
}
func (m *lmsMateriRepo) Create(materi *models.Materi) error { return m.CreateFn(materi) }
func (m *lmsMateriRepo) FindByID(id string) (*models.Materi, error) { return m.FindByIDFn(id) }
func (m *lmsMateriRepo) FindByKelas(idKelas string) ([]models.Materi, error) {
	if m.FindByKelasFn != nil { return m.FindByKelasFn(idKelas) }
	return nil, nil
}
func (m *lmsMateriRepo) Update(materi *models.Materi) error { return m.UpdateFn(materi) }
func (m *lmsMateriRepo) Delete(id string) error { return m.DeleteFn(id) }
func (m *lmsMateriRepo) ReorderUrutan(idKelas string) error {
	if m.ReorderUrutanFn != nil { return m.ReorderUrutanFn(idKelas) }
	return nil
}

type lmsProgressRepo struct{}
func (m *lmsProgressRepo) Upsert(p *models.UserMateriProgress) error { return nil }
func (m *lmsProgressRepo) FindByUserAndMateri(idUser, idMateri string) (*models.UserMateriProgress, error) { return nil, errors.New("not found") }
func (m *lmsProgressRepo) FindByUserAndKelas(idUser, idKelas string) ([]models.UserMateriProgress, error) { return nil, nil }
func (m *lmsProgressRepo) HasCompletedAllMateri(idUser, idKelas string) (bool, error) { return false, nil }

type lmsKuisRepo struct {
	CreateFn   func(kuis *models.Kuis) error
	FindByIDFn func(id string) (*models.Kuis, error)
	FindByKelasFn func(idKelas string) ([]models.Kuis, error)
	FindByMateriFn func(idMateri string) (*models.Kuis, error)
	FindFinalByKelasFn func(idKelas string) (*models.Kuis, error)
	UpdateFn   func(kuis *models.Kuis) error
	DeleteFn   func(id string) error
}
func (m *lmsKuisRepo) Create(kuis *models.Kuis) error { return m.CreateFn(kuis) }
func (m *lmsKuisRepo) FindByID(id string) (*models.Kuis, error) { return m.FindByIDFn(id) }
func (m *lmsKuisRepo) FindByKelas(idKelas string) ([]models.Kuis, error) {
	if m.FindByKelasFn != nil { return m.FindByKelasFn(idKelas) }
	return nil, nil
}
func (m *lmsKuisRepo) FindByMateri(idMateri string) (*models.Kuis, error) { return nil, errors.New("not found") }
func (m *lmsKuisRepo) FindFinalByKelas(idKelas string) (*models.Kuis, error) {
	if m.FindFinalByKelasFn != nil { return m.FindFinalByKelasFn(idKelas) }
	return nil, errors.New("not found")
}
func (m *lmsKuisRepo) Update(kuis *models.Kuis) error { return m.UpdateFn(kuis) }
func (m *lmsKuisRepo) Delete(id string) error { return m.DeleteFn(id) }

type lmsSoalRepo struct {
	CreateFn   func(soal *models.Soal, pilihan []models.PilihanJawaban) error
	FindByIDFn func(id string) (*models.Soal, error)
	FindByKuisFn func(idKuis string) ([]models.Soal, error)
	UpdateFn   func(soal *models.Soal, pilihan []models.PilihanJawaban) error
	DeleteFn   func(id string) error
	FindPilihanByIDFn func(idPilihan string) (*models.PilihanJawaban, error)
	FindCorrectPilihanFn func(idSoal string) (*models.PilihanJawaban, error)
}
func (m *lmsSoalRepo) Create(soal *models.Soal, pilihan []models.PilihanJawaban) error { return m.CreateFn(soal, pilihan) }
func (m *lmsSoalRepo) FindByID(id string) (*models.Soal, error) { return m.FindByIDFn(id) }
func (m *lmsSoalRepo) FindByKuis(idKuis string) ([]models.Soal, error) { return m.FindByKuisFn(idKuis) }
func (m *lmsSoalRepo) Update(soal *models.Soal, pilihan []models.PilihanJawaban) error { return m.UpdateFn(soal, pilihan) }
func (m *lmsSoalRepo) Delete(id string) error { return m.DeleteFn(id) }
func (m *lmsSoalRepo) FindPilihanByID(idPilihan string) (*models.PilihanJawaban, error) { return nil, errors.New("not found") }
func (m *lmsSoalRepo) FindCorrectPilihan(idSoal string) (*models.PilihanJawaban, error) { return nil, errors.New("not found") }

type lmsAttemptRepo struct{}
func (m *lmsAttemptRepo) Create(a *models.KuisAttempt) error { return nil }
func (m *lmsAttemptRepo) FindByID(id string) (*models.KuisAttempt, error) { return nil, errors.New("not found") }
func (m *lmsAttemptRepo) FindByUserAndKuis(idUser, idKuis string) ([]models.KuisAttempt, error) { return nil, nil }
func (m *lmsAttemptRepo) FindLatestByUserAndKuis(idUser, idKuis string) (*models.KuisAttempt, error) { return nil, errors.New("not found") }
func (m *lmsAttemptRepo) Finish(id string, skor float64, totalBenar int, isPassed bool, jawaban []models.KuisJawaban) error { return nil }
func (m *lmsAttemptRepo) HasPassedAllKuisInKelas(idUser, idKelas string) (bool, error) { return false, nil }
func (m *lmsAttemptRepo) FindJawabanByAttempt(idAttempt string) ([]models.KuisJawaban, error) { return nil, nil }

type lmsFPRepo struct {
	CreateFn func(fp *models.FilePendukung) error
	FindByMateriFn func(idMateri string) ([]models.FilePendukung, error)
	FindByIDFn func(id string) (*models.FilePendukung, error)
	DeleteFn func(id string) error
}
func (m *lmsFPRepo) Create(fp *models.FilePendukung) error { return m.CreateFn(fp) }
func (m *lmsFPRepo) FindByMateri(idMateri string) ([]models.FilePendukung, error) { return m.FindByMateriFn(idMateri) }
func (m *lmsFPRepo) FindByID(id string) (*models.FilePendukung, error) { return m.FindByIDFn(id) }
func (m *lmsFPRepo) Delete(id string) error { return m.DeleteFn(id) }

type lmsDiskusiRepo struct {
	CreateFn func(d *models.Diskusi) error
	FindByMateriFn func(idMateri string) ([]models.Diskusi, error)
	FindByIDFn func(id string) (*models.Diskusi, error)
	UpdateFn func(d *models.Diskusi) error
	DeleteFn func(id string) error
	FindRepliesFn func(idParent string) ([]models.Diskusi, error)
}
func (m *lmsDiskusiRepo) Create(d *models.Diskusi) error { return m.CreateFn(d) }
func (m *lmsDiskusiRepo) FindByMateri(idMateri string) ([]models.Diskusi, error) { return m.FindByMateriFn(idMateri) }
func (m *lmsDiskusiRepo) FindByID(id string) (*models.Diskusi, error) { return m.FindByIDFn(id) }
func (m *lmsDiskusiRepo) Update(d *models.Diskusi) error { return m.UpdateFn(d) }
func (m *lmsDiskusiRepo) Delete(id string) error { return m.DeleteFn(id) }
func (m *lmsDiskusiRepo) FindReplies(idParent string) ([]models.Diskusi, error) {
	if m.FindRepliesFn != nil { return m.FindRepliesFn(idParent) }
	return nil, nil
}

type lmsCatatanRepo struct {
	UpsertFn func(c *models.CatatanPribadi) error
	FindByUserAndMateriFn func(idUser, idMateri string) (*models.CatatanPribadi, error)
	DeleteFn func(id string) error
}
func (m *lmsCatatanRepo) Upsert(c *models.CatatanPribadi) error { return m.UpsertFn(c) }
func (m *lmsCatatanRepo) FindByUserAndMateri(idUser, idMateri string) (*models.CatatanPribadi, error) { return m.FindByUserAndMateriFn(idUser, idMateri) }
func (m *lmsCatatanRepo) Delete(id string) error { if m.DeleteFn != nil { return m.DeleteFn(id) }; return nil }

type lmsSertifikatRepo struct {
	CreateFn func(s *models.Sertifikat) error
	FindByUserAndKelasFn func(idUser, idKelas string) (*models.Sertifikat, error)
	FindByIDFn func(id string) (*models.Sertifikat, error)
	FindByUserFn func(idUser string) ([]models.Sertifikat, error)
}
func (m *lmsSertifikatRepo) Create(s *models.Sertifikat) error { return m.CreateFn(s) }
func (m *lmsSertifikatRepo) FindByUserAndKelas(idUser, idKelas string) (*models.Sertifikat, error) { return m.FindByUserAndKelasFn(idUser, idKelas) }
func (m *lmsSertifikatRepo) FindByID(id string) (*models.Sertifikat, error) { return m.FindByIDFn(id) }
func (m *lmsSertifikatRepo) FindByUser(idUser string) ([]models.Sertifikat, error) { return m.FindByUserFn(idUser) }

type lmsUserRepo struct {
	FindByIDFn func(id string) (*models.User, error)
}
func (m *lmsUserRepo) Create(user *models.User) error { return nil }
func (m *lmsUserRepo) FindByUsername(username string) (*models.User, error) { return nil, errors.New("not found") }
func (m *lmsUserRepo) FindByEmail(email string) (*models.User, error) { return nil, errors.New("not found") }
func (m *lmsUserRepo) FindByID(id string) (*models.User, error) { return m.FindByIDFn(id) }
func (m *lmsUserRepo) FindAll() ([]models.User, error) { return nil, nil }
func (m *lmsUserRepo) Update(user *models.User) error { return nil }
func (m *lmsUserRepo) UpdateWithPhoto(user *models.User) error { return nil }
func (m *lmsUserRepo) UpdatePassword(id, hp string) error { return nil }
func (m *lmsUserRepo) GetPasswordByID(id string) (string, error) { return "", errors.New("not found") }
func (m *lmsUserRepo) Delete(id string) error { return nil }
func (m *lmsUserRepo) EmailExists(email string, ex *string) (bool, error) { return false, nil }
func (m *lmsUserRepo) UsernameExists(un string, ex *string) (bool, error) { return false, nil }
func (m *lmsUserRepo) SetMFA(uid string, s *string, e bool) error { return nil }
func (m *lmsUserRepo) ExistsByPerusahaan(idP string) (bool, error) { return false, nil }
func (m *lmsUserRepo) UpdateStatus(uid string, s models.UserStatus) error { return nil }
func (m *lmsUserRepo) IncrementLoginAttempts(uid string) (int, error) { return 0, nil }
func (m *lmsUserRepo) ResetLoginAttempts(uid string) error { return nil }
func (m *lmsUserRepo) UpdatePasswordChangedAt(uid string) error { return nil }

// ── setup helpers ────────────────────────────────────────────────────────────

func newDefaultKelasRepo() *lmsKelasRepo {
	now := time.Now()
	return &lmsKelasRepo{
		CreateFn:  func(k *models.Kelas) error { return nil },
		FindByIDFn: func(id string) (*models.Kelas, error) {
			return &models.Kelas{ID: id, Judul: "Test Class", Status: models.KelasStatusPublished, CreatedAt: now, UpdatedAt: now}, nil
		},
		FindAllFn: func(onlyPublished bool) ([]models.Kelas, error) {
			return []models.Kelas{{ID: "k-1", Judul: "Class 1", Status: models.KelasStatusPublished, CreatedAt: now, UpdatedAt: now}}, nil
		},
		UpdateFn: func(k *models.Kelas) error { return nil },
		DeleteFn: func(id string) error { return nil },
	}
}

func newDefaultMateriRepo() *lmsMateriRepo {
	now := time.Now()
	return &lmsMateriRepo{
		CreateFn: func(m *models.Materi) error { return nil },
		FindByIDFn: func(id string) (*models.Materi, error) {
			return &models.Materi{ID: id, IDKelas: "k-1", Judul: "Test Materi", Tipe: models.MateriTipeVideo, CreatedAt: now, UpdatedAt: now}, nil
		},
		FindByKelasFn: func(idKelas string) ([]models.Materi, error) { return nil, nil },
		UpdateFn: func(m *models.Materi) error { return nil },
		DeleteFn: func(id string) error { return nil },
		ReorderUrutanFn: func(idKelas string) error { return nil },
	}
}

func setupLMSHandler() *LMSHandler {
	kelasRepo := newDefaultKelasRepo()
	materiRepo := newDefaultMateriRepo()
	kuisRepo := &lmsKuisRepo{
		CreateFn: func(kuis *models.Kuis) error { return nil },
		FindByIDFn: func(id string) (*models.Kuis, error) { return nil, errors.New("not found") },
		FindByKelasFn: func(idKelas string) ([]models.Kuis, error) { return nil, nil },
		FindFinalByKelasFn: func(idKelas string) (*models.Kuis, error) { return nil, errors.New("not found") },
		UpdateFn: func(kuis *models.Kuis) error { return nil },
		DeleteFn: func(id string) error { return nil },
	}
	soalRepo := &lmsSoalRepo{
		CreateFn: func(soal *models.Soal, pilihan []models.PilihanJawaban) error { return nil },
		FindByIDFn: func(id string) (*models.Soal, error) { return nil, errors.New("not found") },
		FindByKuisFn: func(idKuis string) ([]models.Soal, error) { return nil, nil },
		UpdateFn: func(soal *models.Soal, pilihan []models.PilihanJawaban) error { return nil },
		DeleteFn: func(id string) error { return nil },
	}
	fpRepo := &lmsFPRepo{
		CreateFn: func(fp *models.FilePendukung) error { return nil },
		FindByMateriFn: func(idMateri string) ([]models.FilePendukung, error) { return nil, nil },
		FindByIDFn: func(id string) (*models.FilePendukung, error) { return nil, errors.New("not found") },
		DeleteFn: func(id string) error { return nil },
	}
	diskusiRepo := &lmsDiskusiRepo{
		CreateFn: func(d *models.Diskusi) error { return nil },
		FindByMateriFn: func(idMateri string) ([]models.Diskusi, error) { return []models.Diskusi{}, nil },
		FindByIDFn: func(id string) (*models.Diskusi, error) { return nil, errors.New("not found") },
		UpdateFn: func(d *models.Diskusi) error { return nil },
		DeleteFn: func(id string) error { return nil },
	}
	catatanRepo := &lmsCatatanRepo{
		UpsertFn: func(c *models.CatatanPribadi) error { return nil },
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.CatatanPribadi, error) {
			return nil, errors.New("not found")
		},
	}
	sertifikatRepo := &lmsSertifikatRepo{
		CreateFn: func(s *models.Sertifikat) error { return nil },
		FindByUserAndKelasFn: func(idUser, idKelas string) (*models.Sertifikat, error) { return nil, errors.New("not found") },
		FindByIDFn: func(id string) (*models.Sertifikat, error) { return nil, errors.New("not found") },
		FindByUserFn: func(idUser string) ([]models.Sertifikat, error) { return nil, nil },
	}
	userRepo := &lmsUserRepo{
		FindByIDFn: func(id string) (*models.User, error) {
			return &models.User{ID: id, Username: "testuser"}, nil
		},
	}

	kelasSvc := services.NewKelasService(kelasRepo, materiRepo, &lmsProgressRepo{}, kuisRepo, &lmsAttemptRepo{}, sertifikatRepo, fpRepo, nil)
	materiSvc := services.NewMateriService(materiRepo, kelasRepo, &lmsProgressRepo{}, nil)
	soalSvc := services.NewSoalService(soalRepo, kuisRepo, nil)
	kuisSvc := services.NewKuisService(&lmsAttemptRepo{}, soalRepo, kuisRepo, &lmsProgressRepo{}, nil)
	fpSvc := services.NewFilePendukungService(fpRepo, materiRepo, nil)
	diskusiSvc := services.NewDiskusiService(diskusiRepo, userRepo)
	catatanSvc := services.NewCatatanService(catatanRepo)
	sertifikatSvc := services.NewSertifikatService(sertifikatRepo, kelasRepo, &lmsProgressRepo{}, &lmsAttemptRepo{}, kuisRepo, userRepo)
	sseSvc := services.NewSSEService()

	return NewLMSHandler(kelasSvc, materiSvc, soalSvc, kuisSvc, fpSvc, diskusiSvc, catatanSvc, sertifikatSvc, sseSvc)
}

// withUserCtx is defined in user_handler_test.go

// ════════════════════════════════════════════════════════════════════════════
// TEST KELAS HANDLER
// ════════════════════════════════════════════════════════════════════════════

func TestLMSHandler_KelasGetAll(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/kelas", nil)
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLMSHandler_KelasGetAll_ContentType(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/kelas", nil)
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestLMSHandler_KelasGetDetail(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/kelas/k-1", nil)
	req = withUserCtx(req, "user-1", "user")
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLMSHandler_KelasGetDetail_NotFound(t *testing.T) {
	// Override the FindByID to return not found
	// Since we can't easily override after setup, create a fresh handler
	kelasRepo := &lmsKelasRepo{
		FindByIDFn: func(id string) (*models.Kelas, error) { return nil, errors.New("not found") },
		FindAllFn: func(onlyPublished bool) ([]models.Kelas, error) { return nil, nil },
		CreateFn: func(k *models.Kelas) error { return nil },
		UpdateFn: func(k *models.Kelas) error { return nil },
		DeleteFn: func(id string) error { return nil },
	}
	materiRepo := newDefaultMateriRepo()
	kelasSvc := services.NewKelasService(kelasRepo, materiRepo, &lmsProgressRepo{}, &lmsKuisRepo{
		CreateFn: func(k *models.Kuis) error { return nil }, FindByIDFn: func(id string) (*models.Kuis, error) { return nil, errors.New("not found") },
		UpdateFn: func(k *models.Kuis) error { return nil }, DeleteFn: func(id string) error { return nil },
	}, &lmsAttemptRepo{}, &lmsSertifikatRepo{
		CreateFn: func(s *models.Sertifikat) error { return nil },
		FindByUserAndKelasFn: func(a, b string) (*models.Sertifikat, error) { return nil, errors.New("not found") },
		FindByIDFn: func(id string) (*models.Sertifikat, error) { return nil, errors.New("not found") },
		FindByUserFn: func(id string) ([]models.Sertifikat, error) { return nil, nil },
	}, &lmsFPRepo{
		CreateFn: func(fp *models.FilePendukung) error { return nil },
		FindByMateriFn: func(id string) ([]models.FilePendukung, error) { return nil, nil },
		FindByIDFn: func(id string) (*models.FilePendukung, error) { return nil, errors.New("not found") },
		DeleteFn: func(id string) error { return nil },
	}, nil)
	sseSvc := services.NewSSEService()
	h := NewLMSHandler(kelasSvc, nil, nil, nil, nil, nil, nil, nil, sseSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/kelas/invalid", nil)
	w := httptest.NewRecorder()
	h.ServeKelas(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLMSHandler_KelasCreate(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.CreateKelasRequest{Judul: "New Class"})
	req := httptest.NewRequest(http.MethodPost, "/api/kelas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserCtx(req, "admin-1", "admin")
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLMSHandler_KelasCreate_InvalidBody(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/kelas", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_KelasCreate_WithID(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.CreateKelasRequest{Judul: "New"})
	req := httptest.NewRequest(http.MethodPost, "/api/kelas/some-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_KelasUpdate(t *testing.T) {
	handler := setupLMSHandler()
	judul := "Updated"
	body, _ := json.Marshal(dto.UpdateKelasRequest{Judul: &judul})
	req := httptest.NewRequest(http.MethodPut, "/api/kelas/k-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserCtx(req, "admin-1", "admin")
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLMSHandler_KelasUpdate_NoID(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.UpdateKelasRequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/kelas", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_KelasDelete(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/kelas/k-1", nil)
	req = withUserCtx(req, "admin-1", "admin")
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLMSHandler_KelasDelete_NoID(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/kelas", nil)
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_KelasMethodNotAllowed(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodPatch, "/api/kelas", nil)
	w := httptest.NewRecorder()
	handler.ServeKelas(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// TEST MATERI HANDLER
// ════════════════════════════════════════════════════════════════════════════

func TestLMSHandler_MateriUpdate(t *testing.T) {
	handler := setupLMSHandler()
	judul := "Updated Materi"
	body, _ := json.Marshal(dto.UpdateMateriRequest{Judul: &judul})
	req := httptest.NewRequest(http.MethodPut, "/api/materi/m-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserCtx(req, "admin-1", "admin")
	w := httptest.NewRecorder()
	handler.ServeMateri(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLMSHandler_MateriUpdate_NoID(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.UpdateMateriRequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/materi", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.ServeMateri(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_MateriDelete(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/materi/m-1", nil)
	req = withUserCtx(req, "admin-1", "admin")
	w := httptest.NewRecorder()
	handler.ServeMateri(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLMSHandler_MateriProgress_Unauthorized(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.UpdateProgressRequest{IsCompleted: true})
	req := httptest.NewRequest(http.MethodPost, "/api/materi/m-1/progress", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// no user context
	w := httptest.NewRecorder()
	handler.ServeMateri(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLMSHandler_MateriProgress_MethodNotAllowed(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/materi/m-1/progress", nil)
	w := httptest.NewRecorder()
	handler.ServeMateri(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// TEST SOAL HANDLER
// ════════════════════════════════════════════════════════════════════════════

func TestLMSHandler_SoalDelete_NotFound(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/soal/s-1", nil)
	w := httptest.NewRecorder()
	handler.ServeSoal(w, req)
	// Mock soal repo returns not found → service returns error → handler returns 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_SoalDelete_NoID(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/soal", nil)
	w := httptest.NewRecorder()
	handler.ServeSoal(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_SoalMethodNotAllowed(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/soal", nil)
	w := httptest.NewRecorder()
	handler.ServeSoal(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// TEST DISKUSI HANDLER
// ════════════════════════════════════════════════════════════════════════════

func TestLMSHandler_DiskusiUpdate_Unauthorized(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.UpdateDiskusiRequest{Konten: "New"})
	req := httptest.NewRequest(http.MethodPut, "/api/diskusi/d-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// no user
	w := httptest.NewRecorder()
	handler.ServeDiskusi(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLMSHandler_DiskusiDelete_NoID(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/diskusi", nil)
	w := httptest.NewRecorder()
	handler.ServeDiskusi(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_DiskusiMethodNotAllowed(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/diskusi", nil)
	w := httptest.NewRecorder()
	handler.ServeDiskusi(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// TEST FILE PENDUKUNG HANDLER
// ════════════════════════════════════════════════════════════════════════════

func TestLMSHandler_FilePendukungDelete_NotFound(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/file-pendukung/fp-1", nil)
	w := httptest.NewRecorder()
	handler.ServeFilePendukung(w, req)
	// Mock FP repo's FindByID returns not found → service validates → handler returns 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_FilePendukungDelete_NoID(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/file-pendukung", nil)
	w := httptest.NewRecorder()
	handler.ServeFilePendukung(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_FilePendukungMethodNotAllowed(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/file-pendukung", nil)
	w := httptest.NewRecorder()
	handler.ServeFilePendukung(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// TEST SERTIFIKAT HANDLER
// ════════════════════════════════════════════════════════════════════════════

func TestLMSHandler_SertifikatMe_Unauthorized(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/sertifikat/me", nil)
	// no user
	w := httptest.NewRecorder()
	handler.ServeSertifikat(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLMSHandler_SertifikatMe(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/sertifikat/me", nil)
	req = withUserCtx(req, "user-1", "user")
	w := httptest.NewRecorder()
	handler.ServeSertifikat(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLMSHandler_SertifikatByID_NotFound(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/sertifikat/invalid", nil)
	req = withUserCtx(req, "user-1", "user")
	w := httptest.NewRecorder()
	handler.ServeSertifikat(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLMSHandler_SertifikatNotFoundPath(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/sertifikat", nil)
	req = withUserCtx(req, "user-1", "user")
	w := httptest.NewRecorder()
	handler.ServeSertifikat(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// TEST KUIS HANDLER
// ════════════════════════════════════════════════════════════════════════════

func TestLMSHandler_KuisUpdate_NoID(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.UpdateKuisRequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/kuis/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.ServeKuis(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_KuisDelete_NoID(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/kuis/", nil)
	w := httptest.NewRecorder()
	handler.ServeKuis(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLMSHandler_KuisStart_Unauthorized(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/kuis/kuis-1/start", nil)
	// no user
	w := httptest.NewRecorder()
	handler.ServeKuis(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLMSHandler_KuisSubmit_Unauthorized(t *testing.T) {
	handler := setupLMSHandler()
	body, _ := json.Marshal(dto.SubmitKuisRequest{})
	req := httptest.NewRequest(http.MethodPost, "/api/kuis/attempt/att-1/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// no user
	w := httptest.NewRecorder()
	handler.ServeKuis(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLMSHandler_KuisResult_Unauthorized(t *testing.T) {
	handler := setupLMSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/kuis/attempt/att-1/result", nil)
	// no user
	w := httptest.NewRecorder()
	handler.ServeKuis(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ════════════════════════════════════════════════════════════════════════════
// TEST ROUTING HELPER
// ════════════════════════════════════════════════════════════════════════════

func TestTrimID(t *testing.T) {
	assert.Equal(t, "abc", trimID("/api/kelas/abc", "/api/kelas"))
	assert.Equal(t, "", trimID("/api/kelas", "/api/kelas"))
	assert.Equal(t, "xyz/materi", trimID("/api/kelas/xyz/materi", "/api/kelas"))
}

func TestGetUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	assert.Equal(t, "", getUserID(req))

	req = withUserCtx(req, "user-123", "user")
	assert.Equal(t, "user-123", getUserID(req))
}
