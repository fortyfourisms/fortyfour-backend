-- ════════════════════════════════════════════════════════════════════════════
-- 071: Restructure LMS
-- ════════════════════════════════════════════════════════════════════════════

-- 1. Drop old FK constraints and tables that depend on old structure
--    (data belum ada / dummy, aman di-reset)
DROP TABLE IF EXISTS kuis_jawaban;
DROP TABLE IF EXISTS kuis_attempt;
DROP TABLE IF EXISTS user_materi_progress;
DROP TABLE IF EXISTS pilihan_jawaban;
DROP TABLE IF EXISTS soal;
DROP TABLE IF EXISTS materi;

-- ── Materi (restructured) ─────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS materi (
    id               VARCHAR(36)  NOT NULL PRIMARY KEY,
    id_kelas         VARCHAR(36)  NOT NULL,
    judul            VARCHAR(255) NOT NULL,
    tipe             ENUM('video', 'teks') NOT NULL,
    urutan           INT          NOT NULL,
    youtube_id       VARCHAR(50),
    durasi_detik     INT,
    konten_html      LONGTEXT,                          -- rich content (blog-style, gambar embedded)
    deskripsi_singkat VARCHAR(500),
    kategori         VARCHAR(100),
    created_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_materi_kelas FOREIGN KEY (id_kelas) REFERENCES kelas(id) ON DELETE CASCADE,
    UNIQUE KEY uq_materi_urutan (id_kelas, urutan)
);

-- ── File Pendukung (PDF only) ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS file_pendukung (
    id          VARCHAR(36)  NOT NULL PRIMARY KEY,
    id_materi   VARCHAR(36)  NOT NULL,
    nama_file   VARCHAR(255) NOT NULL,
    file_path   VARCHAR(500) NOT NULL,
    ukuran      BIGINT       NOT NULL DEFAULT 0,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_file_pendukung_materi FOREIGN KEY (id_materi) REFERENCES materi(id) ON DELETE CASCADE
);

-- ── Kuis (entitas terpisah dari materi) ───────────────────────────────────
CREATE TABLE IF NOT EXISTS kuis (
    id            VARCHAR(36)   NOT NULL PRIMARY KEY,
    id_kelas      VARCHAR(36)   NOT NULL,
    id_materi     VARCHAR(36),                          -- NULL = kuis akhir
    judul         VARCHAR(255)  NOT NULL,
    deskripsi     TEXT,
    durasi_menit  INT,
    passing_grade DECIMAL(5,2)  NOT NULL DEFAULT 70.00,
    is_final      TINYINT(1)    NOT NULL DEFAULT 0,
    urutan        INT           NOT NULL DEFAULT 0,
    created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_kuis_kelas  FOREIGN KEY (id_kelas)  REFERENCES kelas(id)  ON DELETE CASCADE,
    CONSTRAINT fk_kuis_materi FOREIGN KEY (id_materi) REFERENCES materi(id) ON DELETE SET NULL
);

-- ── Soal Kuis (FK ke kuis, bukan materi) ──────────────────────────────────
CREATE TABLE IF NOT EXISTS soal (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    id_kuis     VARCHAR(36) NOT NULL,
    pertanyaan  TEXT        NOT NULL,
    urutan      INT         NOT NULL,
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_soal_kuis FOREIGN KEY (id_kuis) REFERENCES kuis(id) ON DELETE CASCADE,
    UNIQUE KEY uq_soal_urutan (id_kuis, urutan)
);

-- ── Pilihan Jawaban ───────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS pilihan_jawaban (
    id         VARCHAR(36) NOT NULL PRIMARY KEY,
    id_soal    VARCHAR(36) NOT NULL,
    teks       TEXT        NOT NULL,
    is_correct TINYINT(1)  NOT NULL DEFAULT 0,
    urutan     INT         NOT NULL,

    CONSTRAINT fk_pilihan_soal FOREIGN KEY (id_soal) REFERENCES soal(id) ON DELETE CASCADE,
    UNIQUE KEY uq_pilihan_urutan (id_soal, urutan)
);

-- ── Progress User (video & teks) ──────────────────────────────────────────
CREATE TABLE IF NOT EXISTS user_materi_progress (
    id                   VARCHAR(36) NOT NULL PRIMARY KEY,
    id_user              VARCHAR(36) NOT NULL,
    id_materi            VARCHAR(36) NOT NULL,
    is_completed         TINYINT(1)  NOT NULL DEFAULT 0,
    last_watched_seconds INT         NOT NULL DEFAULT 0,
    completed_at         DATETIME,
    created_at           DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_progress_user   FOREIGN KEY (id_user)   REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_progress_materi FOREIGN KEY (id_materi) REFERENCES materi(id) ON DELETE CASCADE,
    UNIQUE KEY uq_progress_user_materi (id_user, id_materi),
    INDEX idx_progress_user_kelas (id_user, is_completed)
);

-- ── Kuis Attempt (FK ke kuis, bukan materi) ───────────────────────────────
CREATE TABLE IF NOT EXISTS kuis_attempt (
    id          VARCHAR(36)   NOT NULL PRIMARY KEY,
    id_user     VARCHAR(36)   NOT NULL,
    id_kuis     VARCHAR(36)   NOT NULL,
    skor        DECIMAL(5, 2) NOT NULL DEFAULT 0.00,
    total_soal  INT           NOT NULL DEFAULT 0,
    total_benar INT           NOT NULL DEFAULT 0,
    is_passed   TINYINT(1)    NOT NULL DEFAULT 0,
    started_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,

    CONSTRAINT fk_attempt_user FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_attempt_kuis FOREIGN KEY (id_kuis) REFERENCES kuis(id)  ON DELETE CASCADE,
    INDEX idx_attempt_user_kuis (id_user, id_kuis)
);

-- ── Jawaban Per Attempt ───────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS kuis_jawaban (
    id         VARCHAR(36) NOT NULL PRIMARY KEY,
    id_attempt VARCHAR(36) NOT NULL,
    id_soal    VARCHAR(36) NOT NULL,
    id_pilihan VARCHAR(36) NOT NULL,
    is_correct TINYINT(1)  NOT NULL DEFAULT 0,

    CONSTRAINT fk_jawaban_attempt FOREIGN KEY (id_attempt) REFERENCES kuis_attempt(id) ON DELETE CASCADE,
    CONSTRAINT fk_jawaban_soal    FOREIGN KEY (id_soal)    REFERENCES soal(id),
    CONSTRAINT fk_jawaban_pilihan FOREIGN KEY (id_pilihan) REFERENCES pilihan_jawaban(id),
    UNIQUE KEY uq_jawaban_attempt_soal (id_attempt, id_soal)
);

-- ── Diskusi ───────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS diskusi (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    id_materi   VARCHAR(36) NOT NULL,
    id_user     VARCHAR(36) NOT NULL,
    id_parent   VARCHAR(36),
    konten      TEXT        NOT NULL,
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_diskusi_materi FOREIGN KEY (id_materi) REFERENCES materi(id) ON DELETE CASCADE,
    CONSTRAINT fk_diskusi_user   FOREIGN KEY (id_user)   REFERENCES users(id)  ON DELETE CASCADE,
    CONSTRAINT fk_diskusi_parent FOREIGN KEY (id_parent) REFERENCES diskusi(id) ON DELETE CASCADE,
    INDEX idx_diskusi_materi (id_materi)
);

-- ── Catatan Pribadi ───────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS catatan_pribadi (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    id_materi   VARCHAR(36) NOT NULL,
    id_user     VARCHAR(36) NOT NULL,
    konten      TEXT        NOT NULL,
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_catatan_materi FOREIGN KEY (id_materi) REFERENCES materi(id) ON DELETE CASCADE,
    CONSTRAINT fk_catatan_user   FOREIGN KEY (id_user)   REFERENCES users(id)  ON DELETE CASCADE,
    UNIQUE KEY uq_catatan_user_materi (id_user, id_materi)
);

-- ── Sertifikat ────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS sertifikat (
    id                VARCHAR(36)  NOT NULL PRIMARY KEY,
    nomor_sertifikat  VARCHAR(50)  NOT NULL UNIQUE,
    id_kelas          VARCHAR(36)  NOT NULL,
    id_user           VARCHAR(36)  NOT NULL,
    nama_peserta      VARCHAR(255) NOT NULL,
    nama_kelas        VARCHAR(255) NOT NULL,
    tanggal_terbit    DATE         NOT NULL,
    pdf_path          VARCHAR(500),
    created_at        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_sertifikat_kelas FOREIGN KEY (id_kelas) REFERENCES kelas(id) ON DELETE CASCADE,
    CONSTRAINT fk_sertifikat_user  FOREIGN KEY (id_user)  REFERENCES users(id)  ON DELETE CASCADE,
    UNIQUE KEY uq_sertifikat_user_kelas (id_user, id_kelas)
);
