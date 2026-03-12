CREATE TABLE IF NOT EXISTS jawaban_gulih_buffer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pertanyaan_gulih_id INT NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_gulih DECIMAL(3, 2) NOT NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uq_jawaban_gulih_buffer (perusahaan_id, pertanyaan_gulih_id)
);
