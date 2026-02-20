CREATE TABLE jawaban_deteksi (
    id CHAR(36) PRIMARY KEY,
    pertanyaan_deteksi_id CHAR(36) NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_deteksi TEXT NOT NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,

    CONSTRAINT chk_validasi_evidence
        CHECK (evidence IS NOT NULL OR validasi IS NULL),

    CONSTRAINT fk_jawaban_deteksi_pertanyaan
        FOREIGN KEY (pertanyaan_deteksi_id)
        REFERENCES pertanyaan_deteksi(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_deteksi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_jawaban_deteksi_perusahaan
        UNIQUE (perusahaan_id, pertanyaan_deteksi_id)
);