package dto

import "fortyfour-backend/internal/models"

// ── Kelas ────────────────────────────────────────────────────────────────────

type CreateKelasRequest struct {
	Judul     string  `json:"judul" validate:"required,min=3,max=255"`
	Deskripsi *string `json:"deskripsi,omitempty"`
	Thumbnail *string `json:"thumbnail,omitempty"`
}

type UpdateKelasRequest struct {
	Judul     *string `json:"judul,omitempty" validate:"omitempty,min=3,max=255"`
	Deskripsi *string `json:"deskripsi,omitempty"`
	Thumbnail *string `json:"thumbnail,omitempty"`
	Status    *string `json:"status,omitempty" validate:"omitempty,oneof=draft published"`
}

type KelasResponse struct {
	ID        string             `json:"id"`
	Judul     string             `json:"judul"`
	Deskripsi *string            `json:"deskripsi"`
	Thumbnail *string            `json:"thumbnail"`
	Status    models.KelasStatus `json:"status"`
	CreatedBy string             `json:"created_by"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`

	// Disertakan saat GET detail
	Materi     []MateriResponse    `json:"materi,omitempty"`
	KuisList   []KuisResponse      `json:"kuis_list,omitempty"`
	Progress   *KelasProgress      `json:"progress,omitempty"` // progress user saat ini
	Sertifikat *SertifikatResponse `json:"sertifikat,omitempty"`
}

// KelasProgress adalah ringkasan progress user dalam satu kelas
type KelasProgress struct {
	TotalMateri        int     `json:"total_materi"`
	MateriSelesai      int     `json:"materi_selesai"`
	TotalKuis          int     `json:"total_kuis"`
	KuisLulus          int     `json:"kuis_lulus"`
	KuisAkhirLulus     bool    `json:"kuis_akhir_lulus"`
	IsKelasSelesai     bool    `json:"is_kelas_selesai"`
	PersentaseProgress float64 `json:"persentase_progress"`
}

// ── Materi ───────────────────────────────────────────────────────────────────

type CreateMateriRequest struct {
	Judul            string  `json:"judul" validate:"required,min=3,max=255"`
	Tipe             string  `json:"tipe" validate:"required,oneof=video teks"`
	Urutan           int     `json:"urutan" validate:"required,min=1"`
	YoutubeID        *string `json:"youtube_id,omitempty"`
	DurasiDetik      *int    `json:"durasi_detik,omitempty"`
	KontenHTML       *string `json:"konten_html,omitempty"`
	DeskripsiSingkat *string `json:"deskripsi_singkat,omitempty"`
	Kategori         *string `json:"kategori,omitempty"`
}

type UpdateMateriRequest struct {
	Judul            *string `json:"judul,omitempty" validate:"omitempty,min=3,max=255"`
	Urutan           *int    `json:"urutan,omitempty" validate:"omitempty,min=1"`
	YoutubeID        *string `json:"youtube_id,omitempty"`
	DurasiDetik      *int    `json:"durasi_detik,omitempty"`
	KontenHTML       *string `json:"konten_html,omitempty"`
	DeskripsiSingkat *string `json:"deskripsi_singkat,omitempty"`
	Kategori         *string `json:"kategori,omitempty"`
}

type MateriResponse struct {
	ID               string            `json:"id"`
	IDKelas          string            `json:"id_kelas"`
	Judul            string            `json:"judul"`
	Tipe             models.MateriTipe `json:"tipe"`
	Urutan           int               `json:"urutan"`
	YoutubeID        *string           `json:"youtube_id,omitempty"`
	DurasiDetik      *int              `json:"durasi_detik,omitempty"`
	KontenHTML       *string           `json:"konten_html,omitempty"`
	DeskripsiSingkat *string           `json:"deskripsi_singkat,omitempty"`
	Kategori         *string           `json:"kategori,omitempty"`
	CreatedAt        string            `json:"created_at"`
	UpdatedAt        string            `json:"updated_at"`

	// Progress user untuk materi ini
	IsCompleted        bool `json:"is_completed,omitempty"`
	LastWatchedSeconds int  `json:"last_watched_seconds,omitempty"`

	// Nested data (opsional, disertakan saat GET detail)
	FilePendukung []FilePendukungResponse `json:"file_pendukung,omitempty"`
	Kuis          *KuisResponse           `json:"kuis,omitempty"` // kuis per-materi (opsional)
}

// ── File Pendukung ───────────────────────────────────────────────────────────

type FilePendukungResponse struct {
	ID        string `json:"id"`
	IDMateri  string `json:"id_materi"`
	NamaFile  string `json:"nama_file"`
	FilePath  string `json:"file_path"`
	Ukuran    int64  `json:"ukuran"`
	CreatedAt string `json:"created_at"`
}

// ── Kuis ─────────────────────────────────────────────────────────────────────

type CreateKuisRequest struct {
	IDMateri     *string `json:"id_materi,omitempty"` // NULL = kuis akhir
	Judul        string  `json:"judul" validate:"required,min=3,max=255"`
	Deskripsi    *string `json:"deskripsi,omitempty"`
	DurasiMenit  *int    `json:"durasi_menit,omitempty"`
	PassingGrade float64 `json:"passing_grade" validate:"required,min=0,max=100"`
	IsFinal      bool    `json:"is_final"`
	Urutan       int     `json:"urutan" validate:"required,min=1"`
}

type UpdateKuisRequest struct {
	Judul        *string  `json:"judul,omitempty" validate:"omitempty,min=3,max=255"`
	Deskripsi    *string  `json:"deskripsi,omitempty"`
	DurasiMenit  *int     `json:"durasi_menit,omitempty"`
	PassingGrade *float64 `json:"passing_grade,omitempty" validate:"omitempty,min=0,max=100"`
	IsFinal      *bool    `json:"is_final,omitempty"`
	Urutan       *int     `json:"urutan,omitempty"`
}

type KuisResponse struct {
	ID           string  `json:"id"`
	IDKelas      string  `json:"id_kelas"`
	IDMateri     *string `json:"id_materi,omitempty"`
	Judul        string  `json:"judul"`
	Deskripsi    *string `json:"deskripsi,omitempty"`
	DurasiMenit  *int    `json:"durasi_menit,omitempty"`
	PassingGrade float64 `json:"passing_grade"`
	IsFinal      bool    `json:"is_final"`
	Urutan       int     `json:"urutan"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`

	// Nested — hanya admin GET detail
	Soal []SoalResponse `json:"soal,omitempty"`
}

// ── Soal ─────────────────────────────────────────────────────────────────────

type CreateSoalRequest struct {
	Pertanyaan string                 `json:"pertanyaan" validate:"required"`
	Urutan     int                    `json:"urutan" validate:"required,min=1"`
	Pilihan    []CreatePilihanRequest `json:"pilihan" validate:"required,min=2,max=5,dive"`
}

type CreatePilihanRequest struct {
	Teks      string `json:"teks" validate:"required"`
	IsCorrect bool   `json:"is_correct"`
	Urutan    int    `json:"urutan" validate:"required,min=1"`
}

type UpdateSoalRequest struct {
	Pertanyaan *string                `json:"pertanyaan,omitempty"`
	Urutan     *int                   `json:"urutan,omitempty"`
	Pilihan    []CreatePilihanRequest `json:"pilihan,omitempty"`
}

// SoalResponse untuk admin (tampilkan is_correct)
type SoalResponse struct {
	ID         string            `json:"id"`
	IDKuis     string            `json:"id_kuis"`
	Pertanyaan string            `json:"pertanyaan"`
	Urutan     int               `json:"urutan"`
	Pilihan    []PilihanResponse `json:"pilihan"`
}

// SoalUserResponse untuk user (sembunyikan is_correct)
type SoalUserResponse struct {
	ID         string                `json:"id"`
	Pertanyaan string                `json:"pertanyaan"`
	Urutan     int                   `json:"urutan"`
	Pilihan    []PilihanUserResponse `json:"pilihan"`
}

type PilihanResponse struct {
	ID        string `json:"id"`
	Teks      string `json:"teks"`
	IsCorrect bool   `json:"is_correct"`
	Urutan    int    `json:"urutan"`
}

type PilihanUserResponse struct {
	ID     string `json:"id"`
	Teks   string `json:"teks"`
	Urutan int    `json:"urutan"`
}

// ── Progress Materi ───────────────────────────────────────────────────────────

// UpdateProgressRequest dipakai untuk update progress video maupun tandai teks selesai
type UpdateProgressRequest struct {
	LastWatchedSeconds *int `json:"last_watched_seconds,omitempty"` // khusus video
	IsCompleted        bool `json:"is_completed"`
}

type ProgressResponse struct {
	IDMateri           string  `json:"id_materi"`
	IsCompleted        bool    `json:"is_completed"`
	LastWatchedSeconds int     `json:"last_watched_seconds"`
	CompletedAt        *string `json:"completed_at"`
}

// ── Kuis Flow ────────────────────────────────────────────────────────────────

// StartKuisResponse berisi soal-soal yang harus dijawab (tanpa is_correct)
type StartKuisResponse struct {
	AttemptID string             `json:"attempt_id"`
	IDKuis    string             `json:"id_kuis"`
	Soal      []SoalUserResponse `json:"soal"`
}

// SubmitKuisRequest berisi jawaban user
type SubmitKuisRequest struct {
	Jawaban []JawabanItem `json:"jawaban" validate:"required,dive"`
}

type JawabanItem struct {
	IDSoal    string `json:"id_soal" validate:"required"`
	IDPilihan string `json:"id_pilihan" validate:"required"`
}

// KuisResultResponse berisi hasil kuis beserta pembahasan
type KuisResultResponse struct {
	AttemptID  string              `json:"attempt_id"`
	Skor       float64             `json:"skor"`
	TotalSoal  int                 `json:"total_soal"`
	TotalBenar int                 `json:"total_benar"`
	IsPassed   bool                `json:"is_passed"`
	FinishedAt string              `json:"finished_at"`
	Detail     []HasilSoalResponse `json:"detail"`
}

type HasilSoalResponse struct {
	IDSoal         string `json:"id_soal"`
	Pertanyaan     string `json:"pertanyaan"`
	IDPilihanUser  string `json:"id_pilihan_user"`  // jawaban yang dipilih user
	IDPilihanBenar string `json:"id_pilihan_benar"` // jawaban benar
	IsCorrect      bool   `json:"is_correct"`
}

// ── Diskusi ──────────────────────────────────────────────────────────────────

type CreateDiskusiRequest struct {
	IDParent *string `json:"id_parent,omitempty"` // NULL = top-level, NOT NULL = reply
	Konten   string  `json:"konten" validate:"required"`
}

type UpdateDiskusiRequest struct {
	Konten string `json:"konten" validate:"required"`
}

type DiskusiResponse struct {
	ID          string            `json:"id"`
	IDMateri    string            `json:"id_materi"`
	IDUser      string            `json:"id_user"`
	NamaUser    string            `json:"nama_user"`
	FotoProfile *string           `json:"foto_profile,omitempty"`
	IDParent    *string           `json:"id_parent,omitempty"`
	Konten      string            `json:"konten"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Replies     []DiskusiResponse `json:"replies,omitempty"`
}

// ── Catatan Pribadi ──────────────────────────────────────────────────────────

type UpsertCatatanRequest struct {
	Konten string `json:"konten" validate:"required"`
}

type CatatanPribadiResponse struct {
	ID        string `json:"id"`
	IDMateri  string `json:"id_materi"`
	Konten    string `json:"konten"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ── Sertifikat ───────────────────────────────────────────────────────────────

type SertifikatResponse struct {
	ID              string  `json:"id"`
	NomorSertifikat string  `json:"nomor_sertifikat"`
	IDKelas         string  `json:"id_kelas"`
	IDUser          string  `json:"id_user"`
	NamaPeserta     string  `json:"nama_peserta"`
	NamaKelas       string  `json:"nama_kelas"`
	TanggalTerbit   string  `json:"tanggal_terbit"`
	PDFPath         *string `json:"pdf_path,omitempty"`
	CreatedAt       string  `json:"created_at"`
}
