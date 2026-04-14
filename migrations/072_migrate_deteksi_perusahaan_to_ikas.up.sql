-- 1. Tambah kolom ikas_id setelah perusahaan_id
ALTER TABLE deteksi 
ADD COLUMN ikas_id CHAR(36) AFTER perusahaan_id;

-- 2. Mapping data dari perusahaan_id → ikas_id
UPDATE deteksi det
JOIN ikas i ON det.perusahaan_id = i.id_perusahaan
SET det.ikas_id = i.id;

-- 3. Drop foreign key lama
ALTER TABLE deteksi 
DROP FOREIGN KEY fk_deteksi_perusahaan;

-- 4. Drop unique/index lama
ALTER TABLE deteksi 
DROP INDEX uq_deteksi_perusahaan;

-- 5. Drop kolom lama
ALTER TABLE deteksi 
DROP COLUMN perusahaan_id;

-- 6. Tambah unique baru
ALTER TABLE deteksi 
ADD CONSTRAINT uq_deteksi_ikas UNIQUE (ikas_id);

-- 7. Tambah foreign key baru ke ikas + CASCADE
ALTER TABLE deteksi 
ADD CONSTRAINT fk_deteksi_ikas
FOREIGN KEY (ikas_id)
REFERENCES ikas(id)
ON UPDATE CASCADE
ON DELETE CASCADE;