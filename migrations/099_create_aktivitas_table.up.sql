CREATE TABLE IF NOT EXISTS aktivitas (
    id INT AUTO_INCREMENT PRIMARY KEY,
    perusahaan_id CHAR(36) NOT NULL,
    judul VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    tanggal_mulai DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    jenis_aktivitas JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_aktivitas_perusahaan FOREIGN KEY (perusahaan_id) REFERENCES perusahaan(id) ON DELETE CASCADE
);
