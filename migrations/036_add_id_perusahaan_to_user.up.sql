ALTER TABLE users
    ADD COLUMN id_perusahaan CHAR(36) NULL AFTER id_jabatan,
    ADD CONSTRAINT fk_users_perusahaan
        FOREIGN KEY (id_perusahaan)
        REFERENCES perusahaan(id)
        ON DELETE SET NULL;

CREATE INDEX idx_users_perusahaan ON users(id_perusahaan);