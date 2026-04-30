-- 1. Hapus foreign key ke users
ALTER TABLE responden 
DROP FOREIGN KEY fk_responden_user;

-- 2. Hapus kolom yang tidak dipakai
ALTER TABLE responden 
DROP COLUMN user_id,
DROP COLUMN perusahaan,
DROP COLUMN sektor,
DROP COLUMN sektor_lainnya;

-- 3. Tambah kolom baru
ALTER TABLE responden 
ADD COLUMN id_perusahaan CHAR(36) NOT NULL AFTER id,
ADD COLUMN nama_lengkap VARCHAR(150) NOT NULL AFTER id_perusahaan,
ADD COLUMN jabatan VARCHAR(100) NOT NULL AFTER nama_lengkap,
ADD COLUMN email VARCHAR(150) NOT NULL AFTER jabatan;

-- 4. Ubah panjang kolom sertifikat
ALTER TABLE responden 
MODIFY sertifikat_training VARCHAR(255);

-- 5. Tambahkan foreign key ke perusahaan
ALTER TABLE responden 
ADD CONSTRAINT fk_responden_perusahaan
FOREIGN KEY (id_perusahaan)
REFERENCES perusahaan(id)
ON DELETE CASCADE
ON UPDATE CASCADE;