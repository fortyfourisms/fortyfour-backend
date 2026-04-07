ALTER TABLE perusahaan
    ADD CONSTRAINT uq_perusahaan_nama UNIQUE (nama_perusahaan);