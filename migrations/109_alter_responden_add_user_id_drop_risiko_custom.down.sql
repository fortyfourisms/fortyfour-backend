-- Buat kembali tabel risiko_custom jika di-rollback
CREATE TABLE IF NOT EXISTS risiko_custom (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    nama_risiko VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Hapus foreign key dan constraint dari tabel responden
ALTER TABLE responden
DROP FOREIGN KEY fk_survey_responden_user;

ALTER TABLE responden
DROP INDEX uq_survey_responden_user;

-- Hapus kolom user_id dari tabel responden
ALTER TABLE responden
DROP COLUMN user_id;
