-- +migrate Up
CREATE TABLE IF NOT EXISTS jawaban_identifikasi_buffer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pertanyaan_identifikasi_id INT NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_identifikasi DECIMAL(3, 2) NULL,
    evidence TEXT NULL,
    validasi ENUM('yes', 'no') NULL,
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE (perusahaan_id, pertanyaan_identifikasi_id)
);

-- +migrate Down
DROP TABLE IF EXISTS jawaban_identifikasi_buffer;
