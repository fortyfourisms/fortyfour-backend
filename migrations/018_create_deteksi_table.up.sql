CREATE TABLE IF NOT EXISTS deteksi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    perusahaan_id CHAR(36) NULL,
    nilai_deteksi FLOAT,
    nilai_subdomain1 FLOAT NOT NULL,
    nilai_subdomain2 FLOAT NOT NULL,
    nilai_subdomain3 FLOAT NOT NULL,

    CONSTRAINT uq_deteksi_perusahaan UNIQUE (perusahaan_id),
    CONSTRAINT fk_deteksi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);