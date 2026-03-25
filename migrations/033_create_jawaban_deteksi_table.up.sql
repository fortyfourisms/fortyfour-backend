CREATE TABLE IF NOT EXISTS jawaban_deteksi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pertanyaan_deteksi_id INT NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_deteksi DECIMAL(3, 2) NOT NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

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