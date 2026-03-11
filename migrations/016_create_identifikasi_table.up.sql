CREATE TABLE IF NOT EXISTS identifikasi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    perusahaan_id CHAR(36) NULL,
    nilai_identifikasi FLOAT,
    nilai_subdomain1 FLOAT NOT NULL,
    nilai_subdomain2 FLOAT NOT NULL,
    nilai_subdomain3 FLOAT NOT NULL,
    nilai_subdomain4 FLOAT NOT NULL,
    nilai_subdomain5 FLOAT NOT NULL,

    CONSTRAINT uq_identifikasi_perusahaan UNIQUE (perusahaan_id),
    CONSTRAINT fk_identifikasi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);