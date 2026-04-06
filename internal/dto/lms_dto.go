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
	Materi   []MateriResponse `json:"materi,omitempty"`
	Progress *KelasProgress   `json:"progress,omitempty"` // progress user saat ini
}

// KelasProgress adalah ringkasan progress user dalam satu kelas
type KelasProgress struct {
	TotalMateri    int  `json:"total_materi"`
	MateriSelesai  int  `json:"materi_selesai"`
	KuisSelesai    bool `json:"kuis_selesai"`
	IsKelasSelesai bool `json:"is_kelas_selesai"`
}

// ── Materi ───────────────────────────────────────────────────────────────────

type CreateMateriRequest struct {
	Judul       string  `json:"judul" validate:"required,min=3,max=255"`
	Tipe        string  `json:"tipe" validate:"required,oneof=video pdf kuis"`
	Urutan      int     `json:"urutan" validate:"required,min=1"`
	YoutubeID   *string `json:"youtube_id,omitempty"`
	PDFPath     *string `json:"pdf_path,omitempty"`
	DurasiDetik *int    `json:"durasi_detik,omitempty"`
}

type UpdateMateriRequest struct {
	Judul       *string `json:"judul,omitempty" validate:"omitempty,min=3,max=255"`
	Urutan      *int    `json:"urutan,omitempty" validate:"omitempty,min=1"`
	YoutubeID   *string `json:"youtube_id,omitempty"`
	PDFPath     *string `json:"pdf_path,omitempty"`
	DurasiDetik *int    `json:"durasi_detik,omitempty"`
}

type MateriResponse struct {
	ID          string            `json:"id"`
	IDKelas     string            `json:"id_kelas"`
	Judul       string            `json:"judul"`
	Tipe        models.MateriTipe `json:"tipe"`
	Urutan      int               `json:"urutan"`
	YoutubeID   *string           `json:"youtube_id,omitempty"`
	PDFPath     *string           `json:"pdf_path,omitempty"`
	DurasiDetik *int              `json:"durasi_detik,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`

	// Progress user untuk materi ini (disertakan saat GET detail kelas)
	IsCompleted        bool `json:"is_completed,omitempty"`
	LastWatchedSeconds int  `json:"last_watched_seconds,omitempty"`

	// Soal-soal disertakan hanya untuk tipe kuis
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
	IDMateri   string            `json:"id_materi"`
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

// UpdateProgressRequest dipakai untuk update progress video maupun tandai pdf selesai
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

// ── Kuis ─────────────────────────────────────────────────────────────────────

// StartKuisResponse berisi soal-soal yang harus dijawab (tanpa is_correct)
type StartKuisResponse struct {
	AttemptID string             `json:"attempt_id"`
	IDMateri  string             `json:"id_materi"`
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
