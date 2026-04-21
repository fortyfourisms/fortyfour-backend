-- Hapus FK constraint csirt dari tabel se
ALTER TABLE se DROP FOREIGN KEY fk_se_csirt;

-- Hapus kolom id_csirt
ALTER TABLE se DROP COLUMN id_csirt;
