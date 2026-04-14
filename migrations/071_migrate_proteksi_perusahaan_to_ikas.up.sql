-- 1. Tambah kolom ikas_id setelah perusahaan_id
ALTER TABLE proteksi 
ADD COLUMN ikas_id CHAR(36) AFTER perusahaan_id;

-- 2. Mapping data dari perusahaan_id → ikas_id
UPDATE proteksi det
JOIN ikas i ON det.perusahaan_id = i.id_perusahaan
SET det.ikas_id = i.id;

-- 3. Drop foreign key lama
ALTER TABLE proteksi 
DROP FOREIGN KEY fk_proteksi_perusahaan;

-- 4. Drop unique/index lama
ALTER TABLE proteksi 
DROP INDEX uq_proteksi_perusahaan;

-- 5. Drop kolom lama
ALTER TABLE proteksi 
DROP COLUMN perusahaan_id;

-- 6. Tambah unique baru
ALTER TABLE proteksi 
ADD CONSTRAINT uq_proteksi_ikas UNIQUE (ikas_id);

-- 7. Tambah foreign key baru ke ikas + CASCADE
ALTER TABLE proteksi 
ADD CONSTRAINT fk_proteksi_ikas
FOREIGN KEY (ikas_id)
REFERENCES ikas(id)
ON UPDATE CASCADE
ON DELETE CASCADE;