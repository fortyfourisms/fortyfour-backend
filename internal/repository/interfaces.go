package repository

import (
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

// UserRepositoryInterface defines methods for user data access
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindAllAdmins() ([]models.User, error)
	FindUsersByPerusahaan(idPerusahaan string) ([]models.User, error)
	Update(user *models.User) error
	UpdateWithPhoto(user *models.User) error
	UpdatePassword(id, hashedPassword string) error
	GetPasswordByID(id string) (string, error)
	Delete(id string) error
	EmailExists(email string, excludeID *string) (bool, error)
	UsernameExists(username string, excludeID *string) (bool, error)
	SetMFA(userID string, secret *string, enabled bool) error
	ExistsByPerusahaan(idPerusahaan string) (bool, error)

	// Security fields
	UpdateStatus(userID string, status models.UserStatus) error
	IncrementLoginAttempts(userID string) (int, error)
	ResetLoginAttempts(userID string) error
	UpdatePasswordChangedAt(userID string) error
}

type TokenRepositoryInterface interface {
	GenerateTokenPair(userID, username, role string) (*models.TokenPair, error)
	RevokeRefreshToken(refreshToken string) error
}

// PostRepositoryInterface defines methods for post data access
type PostRepositoryInterface interface {
	Create(post *models.Post) error
	FindAll() ([]*models.Post, error)
	FindByID(id int) (*models.Post, error)
	FindByAuthorID(authorID string) ([]*models.Post, error)
	Update(post *models.Post) error
	Delete(id int) error
}

// PerusahaanRepositoryInterface
type PerusahaanRepositoryInterface interface {
	Create(req dto.CreatePerusahaanRequest, id string) error
	GetByID(id string) (*dto.PerusahaanResponse, error)
	GetByNama(nama string) (*dto.PerusahaanResponse, error)
	GetAll() ([]dto.PerusahaanResponse, error)
	Update(id string, perusahaan dto.PerusahaanResponse) error
	Delete(id string) error
}

// PICPerusahaanRepositoryInterface
type PICRepositoryInterface interface {
	Create(req dto.CreatePICRequest, id string) error
	GetByID(id string) (*dto.PICResponse, error)
	GetAll() ([]dto.PICResponse, error)
	GetByPerusahaan(idPerusahaan string) ([]dto.PICResponse, error)
	Update(id string, req dto.UpdatePICRequest) error
	Delete(id string) error
}

// CsirtRepositoryInterface
type CsirtRepositoryInterface interface {
	Create(req dto.CreateCsirtRequest, id string) error
	ExistsByPerusahaan(idPerusahaan string) (bool, error)
	GetByID(id string) (*models.Csirt, error)
	GetAllWithPerusahaan() ([]dto.CsirtResponse, error)
	GetByIDWithPerusahaan(id string) (*dto.CsirtResponse, error)
	GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error)
	Update(id string, csirt models.Csirt) error
	Delete(id string) error
	GetByPerusahaanModel(idPerusahaan string) (*models.Csirt, error)
}

// SdmCsirtRepositoryInterface
type SdmCsirtRepositoryInterface interface {
	Create(req dto.CreateSdmCsirtRequest, id string) error
	GetAll() ([]dto.SdmCsirtResponse, error)
	GetByID(id string) (*dto.SdmCsirtResponse, error)
	GetByCsirt(idCsirt string) ([]dto.SdmCsirtResponse, error)
	Update(id string, req dto.SdmCsirtResponse) error
	Delete(id string) error
}

// SektorRepositoryInterface
type SektorRepositoryInterface interface {
	GetAll() ([]dto.SektorResponse, error)
	GetByID(id string) (*dto.SektorResponse, error)
	Create(req dto.SektorRequest) (*dto.SektorResponse, error)
	Update(id string, req dto.SektorRequest) (*dto.SektorResponse, error)
	Delete(id string) error
}

// SubSektorRepositoryInterface
type SubSektorRepositoryInterface interface {
	GetAll() ([]dto.SubSektorResponse, error)
	GetByID(id string) (*dto.SubSektorResponse, error)
	GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error)
	Create(req dto.SubSektorRequest) (*dto.SubSektorResponse, error)
	Update(id string, req dto.SubSektorRequest) (*dto.SubSektorResponse, error)
	Delete(id string) error
}

// SERepositoryInterface
type SERepositoryInterface interface {
	Create(req dto.CreateSERequest, id string, totalBobot int, kategori string) error
	GetAll() ([]dto.SEResponse, error)
	GetByID(id string) (*dto.SEResponse, error)
	GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error)
	Update(id string, req dto.UpdateSERequest, totalBobot int, kategori string) error
	Delete(id string) error
}

// SEEditRequestRepositoryInterface
type SEEditRequestRepositoryInterface interface {
	Create(req *models.SEEditRequest) error
	FindByID(id string) (*models.SEEditRequest, error)
	FindPendingBySE(idSE string) ([]models.SEEditRequest, error)
	FindAllPending() ([]models.SEEditRequest, error)
	FindByUser(idUser string) ([]models.SEEditRequest, error)
	UpdateStatus(id string, status models.SEEditRequestStatus, catatan *string) error
}

// ── Kelas ────────────────────────────────────────────────────────────────────

type KelasRepositoryInterface interface {
	Create(kelas *models.Kelas) error
	FindByID(id string) (*models.Kelas, error)
	FindAll(onlyPublished bool) ([]models.Kelas, error)
	Update(kelas *models.Kelas) error
	Delete(id string) error
}

// ── Materi ───────────────────────────────────────────────────────────────────

type MateriRepositoryInterface interface {
	Create(materi *models.Materi) error
	FindByID(id string) (*models.Materi, error)
	FindByKelas(idKelas string) ([]models.Materi, error)
	Update(materi *models.Materi) error
	Delete(id string) error
	// ReorderUrutan dipakai saat materi dihapus agar urutan tetap rapi
	ReorderUrutan(idKelas string) error
}

// ── File Pendukung ───────────────────────────────────────────────────────────

type FilePendukungRepositoryInterface interface {
	Create(fp *models.FilePendukung) error
	FindByMateri(idMateri string) ([]models.FilePendukung, error)
	FindByID(id string) (*models.FilePendukung, error)
	Delete(id string) error
}

// ── Kuis ─────────────────────────────────────────────────────────────────────

type KuisRepositoryInterface interface {
	Create(kuis *models.Kuis) error
	FindByID(id string) (*models.Kuis, error)
	FindByKelas(idKelas string) ([]models.Kuis, error)
	FindByMateri(idMateri string) (*models.Kuis, error)
	FindFinalByKelas(idKelas string) (*models.Kuis, error)
	Update(kuis *models.Kuis) error
	Delete(id string) error
}

// ── Soal ─────────────────────────────────────────────────────────────────────

type SoalRepositoryInterface interface {
	Create(soal *models.Soal, pilihan []models.PilihanJawaban) error
	FindByID(id string) (*models.Soal, error)
	FindByKuis(idKuis string) ([]models.Soal, error)
	Update(soal *models.Soal, pilihan []models.PilihanJawaban) error
	Delete(id string) error

	// FindPilihanByID dipakai saat validasi submit kuis
	FindPilihanByID(idPilihan string) (*models.PilihanJawaban, error)
	// FindCorrectPilihan dipakai saat hitung skor dan tampilkan hasil
	FindCorrectPilihan(idSoal string) (*models.PilihanJawaban, error)
}

// ── Progress ─────────────────────────────────────────────────────────────────

type ProgressRepositoryInterface interface {
	// Upsert: insert jika belum ada, update jika sudah ada
	Upsert(progress *models.UserMateriProgress) error
	FindByUserAndMateri(idUser, idMateri string) (*models.UserMateriProgress, error)
	FindByUserAndKelas(idUser, idKelas string) ([]models.UserMateriProgress, error)
	// HasCompletedAllMateri cek apakah user sudah selesai semua materi dalam kelas
	HasCompletedAllMateri(idUser, idKelas string) (bool, error)
}

// ── Kuis Attempt ─────────────────────────────────────────────────────────────

type KuisAttemptRepositoryInterface interface {
	Create(attempt *models.KuisAttempt) error
	FindByID(id string) (*models.KuisAttempt, error)
	FindByUserAndKuis(idUser, idKuis string) ([]models.KuisAttempt, error)
	// FindLatestByUserAndKuis untuk cek apakah ada attempt yang belum selesai
	FindLatestByUserAndKuis(idUser, idKuis string) (*models.KuisAttempt, error)
	Finish(id string, skor float64, totalBenar int, isPassed bool, jawaban []models.KuisJawaban) error
	// HasPassedAllKuisInKelas cek apakah user sudah lulus semua kuis (non-final) dalam kelas
	HasPassedAllKuisInKelas(idUser, idKelas string) (bool, error)

	// FindJawabanByAttempt untuk tampilkan detail hasil
	FindJawabanByAttempt(idAttempt string) ([]models.KuisJawaban, error)
}

// ── Feedback ─────────────────────────────────────────────────────────────────

type FeedbackRepositoryInterface interface {
	Upsert(feedback *models.Feedback) error
	FindByUserAndMateri(idUser, idMateri string) (*models.Feedback, error)
	FindByMateri(idMateri string) ([]dto.FeedbackListItem, error)
	Delete(id string) error
}

// ── Sertifikat ───────────────────────────────────────────────────────────────

type SertifikatRepositoryInterface interface {
	Create(sertifikat *models.Sertifikat) error
	FindByUserAndKelas(idUser, idKelas string) (*models.Sertifikat, error)
	FindByID(id string) (*models.Sertifikat, error)
	FindByUser(idUser string) ([]models.Sertifikat, error)
}

// ── Notifications ────────────────────────────────────────────────────────────

type NotificationRepositoryInterface interface {
	Create(notif *models.Notification) error
	FindAll() ([]models.Notification, error)
	FindAllByUserID(userID string) ([]models.Notification, error)
	MarkRead(userID string, notifID int64) error
	MarkAllRead(userID string) error
	MarkAllReadGlobal() error // admin: tandai semua notif semua user
	Delete(userID string, notifID int64) error
	DeleteAllByUserID(userID string) error
	DeleteAll() error // admin: hapus semua notif semua user
}

// ── Berita ───────────────────────────────────────────────────────────────────

type BeritaRepositoryInterface interface {
	Create(berita *models.Berita) error
	FindAll() ([]models.Berita, error)
	FindByID(id int64) (*models.Berita, error)
	Update(berita *models.Berita) error
	Delete(id int64) error
}

type EventRepositoryInterface interface {
	Create(event *models.Event) error
	FindAll() ([]models.Event, error)
	FindByID(id string) (*models.Event, error)
	FindBySlug(slug string) (*models.Event, error)
	Update(event *models.Event) error
	Delete(id string) error
	CreateRegistration(reg *models.EventRegistration) error
	FindRegistrationByID(id string) (*models.EventRegistration, error)
	ExistsRegistrationByEventAndEmail(eventID string, email string) (bool, error)
	UpdateRegistrationPayload(id string, payload string) error
}

type AktivitasRepositoryInterface interface {
	Create(req dto.CreateAktivitasRequest) (int64, error)
	GetAll() ([]dto.AktivitasResponse, error)
	GetByID(id int) (*dto.AktivitasResponse, error)
	GetByPerusahaanID(perusahaanID string) ([]dto.AktivitasResponse, error)
	Update(id int, req dto.UpdateAktivitasRequest) error
	Delete(id int) error
}

// DTOnya tidak dipakai langsung di interface ini, tapi diimport
// agar tetap terkompilasi jika ada helper yang butuh dto.
var _ = dto.KelasResponse{}
