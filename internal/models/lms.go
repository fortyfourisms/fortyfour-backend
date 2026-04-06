package models

import "time"

// ── Kelas ────────────────────────────────────────────────────────────────────

type KelasStatus string

const (
	KelasStatusDraft     KelasStatus = "draft"
	KelasStatusPublished KelasStatus = "published"
)

type Kelas struct {
	ID        string      `json:"id"`
	Judul     string      `json:"judul"`
	Deskripsi *string     `json:"deskripsi"`
	Thumbnail *string     `json:"thumbnail"`
	Status    KelasStatus `json:"status"`
	CreatedBy string      `json:"created_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ── Materi ───────────────────────────────────────────────────────────────────

type MateriTipe string

const (
	MateriTipeVideo MateriTipe = "video"
	MateriTipePDF   MateriTipe = "pdf"
	MateriTipeKuis  MateriTipe = "kuis"
)

type Materi struct {
	ID          string     `json:"id"`
	IDKelas     string     `json:"id_kelas"`
	Judul       string     `json:"judul"`
	Tipe        MateriTipe `json:"tipe"`
	Urutan      int        `json:"urutan"`
	YoutubeID   *string    `json:"youtube_id,omitempty"`   // hanya tipe video
	PDFPath     *string    `json:"pdf_path,omitempty"`     // hanya tipe pdf
	DurasiDetik *int       `json:"durasi_detik,omitempty"` // hanya tipe video
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ── Soal & Pilihan ───────────────────────────────────────────────────────────

type Soal struct {
	ID         string           `json:"id"`
	IDMateri   string           `json:"id_materi"` // materi bertipe kuis
	Pertanyaan string           `json:"pertanyaan"`
	Urutan     int              `json:"urutan"`
	Pilihan    []PilihanJawaban `json:"pilihan,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
}

type PilihanJawaban struct {
	ID        string `json:"id"`
	IDSoal    string `json:"id_soal"`
	Teks      string `json:"teks"`
	IsCorrect bool   `json:"is_correct,omitempty"` // disembunyikan dari response user
	Urutan    int    `json:"urutan"`
}

// ── Progress Materi (video & pdf) ────────────────────────────────────────────

type UserMateriProgress struct {
	ID                 string     `json:"id"`
	IDUser             string     `json:"id_user"`
	IDMateri           string     `json:"id_materi"`
	IsCompleted        bool       `json:"is_completed"`
	LastWatchedSeconds int        `json:"last_watched_seconds"` // hanya video
	CompletedAt        *time.Time `json:"completed_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ── Kuis Attempt ─────────────────────────────────────────────────────────────

type KuisAttempt struct {
	ID         string     `json:"id"`
	IDUser     string     `json:"id_user"`
	IDMateri   string     `json:"id_materi"`
	Skor       float64    `json:"skor"` // 0-100
	TotalSoal  int        `json:"total_soal"`
	TotalBenar int        `json:"total_benar"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

type KuisJawaban struct {
	ID        string `json:"id"`
	IDAttempt string `json:"id_attempt"`
	IDSoal    string `json:"id_soal"`
	IDPilihan string `json:"id_pilihan"` // jawaban yang dipilih user
	IsCorrect bool   `json:"is_correct"`
}
