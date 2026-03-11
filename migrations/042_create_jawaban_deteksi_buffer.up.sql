CREATE TABLE IF NOT EXISTS jawaban_deteksi_buffer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pertanyaan_deteksi_id INT NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_deteksi DECIMAL(3, 2) NOT NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uq_jawaban_deteksi_buffer (perusahaan_id, pertanyaan_deteksi_id)
);
