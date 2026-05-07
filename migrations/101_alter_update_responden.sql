ALTER TABLE responden MODIFY id BIGINT AUTO_INCREMENT;

ALTER TABLE responden ADD COLUMN user_id CHAR(36) NOT NULL AFTER id;

ALTER TABLE responden 
MODIFY nama_lengkap VARCHAR(150) NOT NULL,
MODIFY jabatan VARCHAR(150) NOT NULL,
MODIFY email VARCHAR(150) NOT NULL,
MODIFY no_telepon VARCHAR(30) NOT NULL,
MODIFY id_perusahaan CHAR(36) NOT NULL,
MODIFY sertifikat_training TEXT NULL;

ALTER TABLE responden
ADD CONSTRAINT fk_survey_responden_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE responden
ADD CONSTRAINT fk_survey_responden_perusahaan
FOREIGN KEY (id_perusahaan) REFERENCES perusahaan(id)
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE responden
ADD CONSTRAINT uq_survey_responden_user UNIQUE (user_id);