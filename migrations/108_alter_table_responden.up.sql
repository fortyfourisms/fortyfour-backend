ALTER TABLE responden
MODIFY COLUMN user_id CHAR(36) NOT NULL,
MODIFY COLUMN id_perusahaan CHAR(36) NOT NULL,
MODIFY COLUMN nama_lengkap VARCHAR(150) NOT NULL,
MODIFY COLUMN jabatan VARCHAR(150) NOT NULL,
MODIFY COLUMN email VARCHAR(150) NOT NULL,
MODIFY COLUMN no_telepon VARCHAR(30) NOT NULL,
MODIFY COLUMN sertifikat_training TEXT NULL;

-- Tambahkan unique constraint
ALTER TABLE responden
ADD CONSTRAINT uq_survey_responden_user UNIQUE (user_id);

-- Tambahkan foreign key ke users
ALTER TABLE responden
ADD CONSTRAINT fk_survey_responden_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE
ON UPDATE CASCADE;

-- Tambahkan foreign key ke perusahaan
ALTER TABLE responden
ADD CONSTRAINT fk_survey_responden_perusahaan
FOREIGN KEY (id_perusahaan)
REFERENCES perusahaan(id)
ON DELETE CASCADE
ON UPDATE CASCADE;