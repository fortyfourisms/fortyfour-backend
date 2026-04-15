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

// JabatanRepositoryInterface defines methods for jabatan data access
type JabatanRepositoryInterface interface {
	Create(req dto.CreateJabatanRequest, id string) error
	GetAll() ([]dto.JabatanResponse, error)
	GetByID(id string) (*dto.JabatanResponse, error)
	Update(id string, jabatan dto.JabatanResponse) error
	Delete(id string) error
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
}

// SubSektorRepositoryInterface
type SubSektorRepositoryInterface interface {
	GetAll() ([]dto.SubSektorResponse, error)
	GetByID(id string) (*dto.SubSektorResponse, error)
	GetBySektorID(sektorID string) ([]dto.SubSektorResponse, error)
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

// ── Diskusi ──────────────────────────────────────────────────────────────────

type DiskusiRepositoryInterface interface {
	Create(diskusi *models.Diskusi) error
	FindByMateri(idMateri string) ([]models.Diskusi, error)
	FindByID(id string) (*models.Diskusi, error)
	Update(diskusi *models.Diskusi) error
	Delete(id string) error
	// FindReplies untuk memuat replies secara nested
	FindReplies(idParent string) ([]models.Diskusi, error)
}

// ── Catatan Pribadi ──────────────────────────────────────────────────────────

type CatatanRepositoryInterface interface {
	Upsert(catatan *models.CatatanPribadi) error
	FindByUserAndMateri(idUser, idMateri string) (*models.CatatanPribadi, error)
	Delete(id string) error
}

// ── Sertifikat ───────────────────────────────────────────────────────────────

type SertifikatRepositoryInterface interface {
	Create(sertifikat *models.Sertifikat) error
	FindByUserAndKelas(idUser, idKelas string) (*models.Sertifikat, error)
	FindByID(id string) (*models.Sertifikat, error)
	FindByUser(idUser string) ([]models.Sertifikat, error)
}

// DTOnya tidak dipakai langsung di interface ini, tapi diimport
// agar tetap terkompilasi jika ada helper yang butuh dto.
var _ = dto.KelasResponse{}
