CREATE TABLE IF NOT EXISTS jawaban_proteksi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pertanyaan_proteksi_id INT NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_proteksi DECIMAL(3, 2) NOT NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT chk_validasi_evidence
        CHECK (evidence IS NOT NULL OR validasi IS NULL),

    CONSTRAINT fk_jawaban_proteksi_pertanyaan
        FOREIGN KEY (pertanyaan_proteksi_id)
        REFERENCES pertanyaan_proteksi(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_proteksi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_jawaban_proteksi_perusahaan
        UNIQUE (perusahaan_id, pertanyaan_proteksi_id)
);