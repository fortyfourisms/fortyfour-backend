-- Tambah kolom jabatan (string) ke users
ALTER TABLE users ADD COLUMN jabatan VARCHAR(255) DEFAULT NULL AFTER role_id;

-- Migrasi data: copy nama_jabatan dari tabel jabatan ke kolom baru
UPDATE users u
LEFT JOIN jabatan j ON u.id_jabatan = j.id
SET u.jabatan = j.nama_jabatan
WHERE u.id_jabatan IS NOT NULL;

-- Hapus FK dan kolom id_jabatan
ALTER TABLE users DROP FOREIGN KEY fk_users_jabatan;
ALTER TABLE users DROP INDEX idx_jabatan;
ALTER TABLE users DROP COLUMN id_jabatan;

-- Hapus tabel jabatan
DROP TABLE IF EXISTS jabatan;
