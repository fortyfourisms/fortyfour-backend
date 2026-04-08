-- ── Kelas ─────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS kelas (
    id          VARCHAR(36)  NOT NULL PRIMARY KEY,
    judul       VARCHAR(255) NOT NULL,
    deskripsi   TEXT,
    thumbnail   VARCHAR(255),
    status      ENUM('draft', 'published') NOT NULL DEFAULT 'draft',
    created_by  VARCHAR(36)  NOT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_kelas_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

-- ── Materi ────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS materi (
    id           VARCHAR(36)  NOT NULL PRIMARY KEY,
    id_kelas     VARCHAR(36)  NOT NULL,
    judul        VARCHAR(255) NOT NULL,
    tipe         ENUM('video', 'pdf', 'kuis') NOT NULL,
    urutan       INT          NOT NULL,
    youtube_id   VARCHAR(50),
    pdf_path     VARCHAR(255),
    durasi_detik INT,
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_materi_kelas FOREIGN KEY (id_kelas) REFERENCES kelas(id) ON DELETE CASCADE,
    UNIQUE KEY uq_materi_urutan (id_kelas, urutan)
);

-- ── Soal Kuis ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS soal (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    id_materi   VARCHAR(36) NOT NULL,
    pertanyaan  TEXT        NOT NULL,
    urutan      INT         NOT NULL,
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_soal_materi FOREIGN KEY (id_materi) REFERENCES materi(id) ON DELETE CASCADE,
    UNIQUE KEY uq_soal_urutan (id_materi, urutan)
);

-- ── Pilihan Jawaban ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS pilihan_jawaban (
    id         VARCHAR(36) NOT NULL PRIMARY KEY,
    id_soal    VARCHAR(36) NOT NULL,
    teks       TEXT        NOT NULL,
    is_correct TINYINT(1)  NOT NULL DEFAULT 0,
    urutan     INT         NOT NULL,

    CONSTRAINT fk_pilihan_soal FOREIGN KEY (id_soal) REFERENCES soal(id) ON DELETE CASCADE,
    UNIQUE KEY uq_pilihan_urutan (id_soal, urutan)
);

-- ── Progress User (video & pdf) ───────────────────────────────────────────────
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

-- ── Kuis Attempt ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS kuis_attempt (
    id          VARCHAR(36)   NOT NULL PRIMARY KEY,
    id_user     VARCHAR(36)   NOT NULL,
    id_materi   VARCHAR(36)   NOT NULL,
    skor        DECIMAL(5, 2) NOT NULL DEFAULT 0.00,
    total_soal  INT           NOT NULL DEFAULT 0,
    total_benar INT           NOT NULL DEFAULT 0,
    started_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,

    CONSTRAINT fk_attempt_user   FOREIGN KEY (id_user)   REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_attempt_materi FOREIGN KEY (id_materi) REFERENCES materi(id) ON DELETE CASCADE,
    INDEX idx_attempt_user_materi (id_user, id_materi)
);

-- ── Jawaban Per Attempt ───────────────────────────────────────────────────────
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