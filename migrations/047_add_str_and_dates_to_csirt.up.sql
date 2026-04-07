ALTER TABLE csirt
    ADD COLUMN file_str          VARCHAR(255) NULL,
    ADD COLUMN tanggal_registrasi     DATE     NULL,
    ADD COLUMN tanggal_kadaluarsa     DATE     NULL,
    ADD COLUMN tanggal_registrasi_ulang DATE   NULL;