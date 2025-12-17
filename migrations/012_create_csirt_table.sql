CREATE TABLE csirt (
    id CHAR(36) PRIMARY KEY,
    id_perusahaan CHAR(36) NOT NULL,
    nama_csirt VARCHAR(255) NOT NULL,
    web_csirt VARCHAR(255),
    photo_csirt varchar(255) DEFAULT NULL,
    file_rfc2350 VARCHAR(255),
    file_public_key_pgp VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_perusahaan) REFERENCES perusahaan(id)
);
