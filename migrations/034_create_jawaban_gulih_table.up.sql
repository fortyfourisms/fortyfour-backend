CREATE TABLE IF NOT EXISTS jawaban_gulih (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pertanyaan_gulih_id INT NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_gulih DECIMAL(3, 2) NOT NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT chk_validasi_evidence
        CHECK (evidence IS NOT NULL OR validasi IS NULL),

    CONSTRAINT fk_jawaban_gulih_pertanyaan
        FOREIGN KEY (pertanyaan_gulih_id)
        REFERENCES pertanyaan_gulih(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_gulih_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_jawaban_gulih_perusahaan
        UNIQUE (perusahaan_id, pertanyaan_gulih_id)
);