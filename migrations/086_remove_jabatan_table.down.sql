-- Buat ulang tabel jabatan
CREATE TABLE IF NOT EXISTS jabatan (
    id CHAR(36) PRIMARY KEY,
    nama_jabatan VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tambah kolom id_jabatan kembali ke users
ALTER TABLE users ADD COLUMN id_jabatan CHAR(36) DEFAULT NULL AFTER role_id;
ALTER TABLE users ADD INDEX idx_jabatan (id_jabatan);
ALTER TABLE users ADD CONSTRAINT fk_users_jabatan
    FOREIGN KEY (id_jabatan) REFERENCES jabatan(id)
    ON DELETE SET NULL ON UPDATE CASCADE;

-- Hapus kolom jabatan string
ALTER TABLE users DROP COLUMN jabatan;
