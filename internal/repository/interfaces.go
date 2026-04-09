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
	// FindByKelasBeforeUrutan dipakai untuk cek prerequisite kuis:
	// ambil semua materi dalam kelas dengan urutan < urutanKuis
	FindByKelasBeforeUrutan(idKelas string, urutan int) ([]models.Materi, error)
	Update(materi *models.Materi) error
	Delete(id string) error
	// ReorderUrutan dipakai saat materi dihapus agar urutan tetap rapi
	ReorderUrutan(idKelas string) error
}

// ── Soal ─────────────────────────────────────────────────────────────────────

type SoalRepositoryInterface interface {
	Create(soal *models.Soal, pilihan []models.PilihanJawaban) error
	FindByID(id string) (*models.Soal, error)
	FindByMateri(idMateri string) ([]models.Soal, error)
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
	// HasCompletedAnyMedia cek apakah user sudah selesai minimal 1 video/pdf dalam kelas
	HasCompletedAnyMedia(idUser, idKelas string) (bool, error)
}

// ── Kuis Attempt ─────────────────────────────────────────────────────────────

type KuisAttemptRepositoryInterface interface {
	Create(attempt *models.KuisAttempt) error
	FindByID(id string) (*models.KuisAttempt, error)
	FindByUserAndMateri(idUser, idMateri string) ([]models.KuisAttempt, error)
	// FindLatestByUserAndMateri untuk cek apakah ada attempt yang belum selesai
	FindLatestByUserAndMateri(idUser, idMateri string) (*models.KuisAttempt, error)
	Finish(id string, skor float64, totalBenar int, jawaban []models.KuisJawaban) error

	// FindJawabanByAttempt untuk tampilkan detail hasil
	FindJawabanByAttempt(idAttempt string) ([]models.KuisJawaban, error)
}

// ── Interface gabungan untuk LMS service ─────────────────────────────────────

// Pastikan semua interface ini diimplementasikan di masing-masing repository file.
// Contoh:
//   var _ KelasRepositoryInterface   = (*KelasRepository)(nil)
//   var _ MateriRepositoryInterface  = (*MateriRepository)(nil)
//   var _ SoalRepositoryInterface    = (*SoalRepository)(nil)
//   var _ ProgressRepositoryInterface = (*ProgressRepository)(nil)
//   var _ KuisAttemptRepositoryInterface = (*KuisAttemptRepository)(nil)

// DTOnya tidak dipakai langsung di interface ini, tapi diimport
// agar tetap terkompilasi jika ada helper yang butuh dto.
var _ = dto.KelasResponse{}
