CREATE TABLE IF NOT EXISTS gulih (
    id INT AUTO_INCREMENT PRIMARY KEY,
    perusahaan_id CHAR(36) NULL,
    nilai_gulih FLOAT,
    nilai_subdomain1 FLOAT NOT NULL,
    nilai_subdomain2 FLOAT NOT NULL,
    nilai_subdomain3 FLOAT NOT NULL,
    nilai_subdomain4 FLOAT NOT NULL,

    CONSTRAINT uq_gulih_perusahaan UNIQUE (perusahaan_id),
    CONSTRAINT fk_gulih_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);