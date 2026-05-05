CREATE TABLE responden (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,

  user_id CHAR(36) NOT NULL,
  id_perusahaan CHAR(36) NOT NULL,

  nama_lengkap VARCHAR(150) NOT NULL,
  jabatan VARCHAR(150) NOT NULL,
  email VARCHAR(150) NOT NULL,
  no_telepon VARCHAR(30) NOT NULL,
  sertifikat_training TEXT NULL,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT fk_survey_responden_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

  CONSTRAINT fk_survey_responden_perusahaan
    FOREIGN KEY (id_perusahaan) REFERENCES perusahaan(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

  CONSTRAINT uq_survey_responden_user UNIQUE (user_id)
);