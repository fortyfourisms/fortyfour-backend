CREATE TABLE sdm_csirt (
  id CHAR(36) PRIMARY KEY,
  id_csirt CHAR(36) NOT NULL,

  nama_personel VARCHAR(255) NOT NULL,
  jabatan_csirt VARCHAR(255) NOT NULL,
  jabatan_perusahaan VARCHAR(255) NOT NULL,
  skill TEXT,
  sertifikasi TEXT,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT fk_sdm_csirt_csirt
    FOREIGN KEY (id_csirt)
    REFERENCES csirt(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
);
