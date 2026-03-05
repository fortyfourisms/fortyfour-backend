ALTER TABLE identifikasi
    DROP FOREIGN KEY fk_identifikasi_perusahaan,
    DROP INDEX uq_identifikasi_perusahaan,
    DROP COLUMN perusahaan_id;
