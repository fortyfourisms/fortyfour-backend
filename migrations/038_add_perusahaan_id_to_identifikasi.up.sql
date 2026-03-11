ALTER TABLE identifikasi
    ADD COLUMN perusahaan_id CHAR(36) NULL AFTER id,
    ADD CONSTRAINT fk_identifikasi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    ADD CONSTRAINT uq_identifikasi_perusahaan
        UNIQUE (perusahaan_id);
