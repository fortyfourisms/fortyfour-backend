-- 1. Tambah kolom ikas_id setelah perusahaan_id
ALTER TABLE gulih 
ADD COLUMN ikas_id CHAR(36) AFTER perusahaan_id;

-- 2. Mapping data dari perusahaan_id → ikas_id
UPDATE gulih det
JOIN ikas i ON det.perusahaan_id = i.id_perusahaan
SET det.ikas_id = i.id;

-- 3. Drop foreign key lama
ALTER TABLE gulih 
DROP FOREIGN KEY fk_gulih_perusahaan;

-- 4. Drop unique/index lama
ALTER TABLE gulih 
DROP INDEX uq_gulih_perusahaan;

-- 5. Drop kolom lama
ALTER TABLE gulih 
DROP COLUMN perusahaan_id;

-- 6. Tambah unique baru
ALTER TABLE gulih 
ADD CONSTRAINT uq_gulih_ikas UNIQUE (ikas_id);

-- 7. Tambah foreign key baru ke ikas + CASCADE
ALTER TABLE gulih 
ADD CONSTRAINT fk_gulih_ikas
FOREIGN KEY (ikas_id)
REFERENCES ikas(id)
ON UPDATE CASCADE
ON DELETE CASCADE;