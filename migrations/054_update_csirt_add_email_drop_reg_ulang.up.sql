ALTER TABLE csirt
    ADD COLUMN email_csirt VARCHAR(255) NULL AFTER web_csirt,
    DROP COLUMN tanggal_registrasi_ulang;