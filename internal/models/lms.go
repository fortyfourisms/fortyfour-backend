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
	MateriTipeTeks  MateriTipe = "teks"
)

type Materi struct {
	ID               string     `json:"id"`
	IDKelas          string     `json:"id_kelas"`
	Judul            string     `json:"judul"`
	Tipe             MateriTipe `json:"tipe"`
	Urutan           int        `json:"urutan"`
	YoutubeID        *string    `json:"youtube_id,omitempty"`   // hanya tipe video
	DurasiDetik      *int       `json:"durasi_detik,omitempty"` // hanya tipe video
	KontenHTML       *string    `json:"konten_html,omitempty"`  // rich content (blog-style)
	DeskripsiSingkat *string    `json:"deskripsi_singkat,omitempty"`
	Kategori         *string    `json:"kategori,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ── File Pendukung (PDF) ─────────────────────────────────────────────────────

type FilePendukung struct {
	ID        string    `json:"id"`
	IDMateri  string    `json:"id_materi"`
	NamaFile  string    `json:"nama_file"`
	FilePath  string    `json:"file_path"`
	Ukuran    int64     `json:"ukuran"` // bytes
	CreatedAt time.Time `json:"created_at"`
}

// ── Kuis ─────────────────────────────────────────────────────────────────────

type Kuis struct {
	ID           string    `json:"id"`
	IDKelas      string    `json:"id_kelas"`
	IDMateri     *string   `json:"id_materi,omitempty"` // NULL = kuis akhir
	Judul        string    `json:"judul"`
	Deskripsi    *string   `json:"deskripsi,omitempty"`
	DurasiMenit  *int      `json:"durasi_menit,omitempty"`
	PassingGrade float64   `json:"passing_grade"`
	IsFinal      bool      `json:"is_final"`
	Urutan       int       `json:"urutan"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ── Soal & Pilihan ───────────────────────────────────────────────────────────

type Soal struct {
	ID         string           `json:"id"`
	IDKuis     string           `json:"id_kuis"`
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

// ── Progress Materi (video & teks) ───────────────────────────────────────────

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
	IDKuis     string     `json:"id_kuis"`
	Skor       float64    `json:"skor"` // 0-100
	TotalSoal  int        `json:"total_soal"`
	TotalBenar int        `json:"total_benar"`
	IsPassed   bool       `json:"is_passed"`
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

// ── Diskusi ──────────────────────────────────────────────────────────────────

type Diskusi struct {
	ID        string    `json:"id"`
	IDMateri  string    `json:"id_materi"`
	IDUser    string    `json:"id_user"`
	IDParent  *string   `json:"id_parent,omitempty"` // NULL = top-level, NOT NULL = reply
	Konten    string    `json:"konten"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ── Catatan Pribadi ──────────────────────────────────────────────────────────

type CatatanPribadi struct {
	ID        string    `json:"id"`
	IDMateri  string    `json:"id_materi"`
	IDUser    string    `json:"id_user"`
	Konten    string    `json:"konten"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ── Sertifikat ───────────────────────────────────────────────────────────────

type Sertifikat struct {
	ID              string    `json:"id"`
	NomorSertifikat string    `json:"nomor_sertifikat"`
	IDKelas         string    `json:"id_kelas"`
	IDUser          string    `json:"id_user"`
	NamaPeserta     string    `json:"nama_peserta"`
	NamaKelas       string    `json:"nama_kelas"`
	TanggalTerbit   time.Time `json:"tanggal_terbit"`
	PDFPath         *string   `json:"pdf_path,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
