CREATE TABLE IF NOT EXISTS jawaban_identifikasi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pertanyaan_identifikasi_id INT NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_identifikasi DECIMAL(3, 2) NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_jawaban_identifikasi_pertanyaan
        FOREIGN KEY (pertanyaan_identifikasi_id)
        REFERENCES pertanyaan_identifikasi(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_identifikasi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_jawaban_identifikasi_perusahaan
        UNIQUE (perusahaan_id, pertanyaan_identifikasi_id)
);