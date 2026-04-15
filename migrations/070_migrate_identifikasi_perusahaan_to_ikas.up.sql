-- 1. Tambah kolom ikas_id setelah perusahaan_id
ALTER TABLE identifikasi 
ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER perusahaan_id;

-- 2. Mapping data dari perusahaan_id → ikas_id
UPDATE identifikasi det
JOIN ikas i ON det.perusahaan_id = i.id_perusahaan
SET det.ikas_id = i.id;

-- 3. Drop foreign key lama
ALTER TABLE identifikasi 
DROP FOREIGN KEY fk_identifikasi_perusahaan;

-- 4. Drop unique/index lama
ALTER TABLE identifikasi 
DROP INDEX uq_identifikasi_perusahaan;

-- 5. Drop kolom lama
ALTER TABLE identifikasi 
DROP COLUMN perusahaan_id;

-- 6. Tambah unique baru
ALTER TABLE identifikasi 
ADD CONSTRAINT uq_identifikasi_ikas UNIQUE (ikas_id);

-- 7. Tambah foreign key baru ke ikas + CASCADE
ALTER TABLE identifikasi 
ADD CONSTRAINT fk_identifikasi_ikas
FOREIGN KEY (ikas_id)
REFERENCES ikas(id)
ON UPDATE CASCADE
ON DELETE CASCADE;