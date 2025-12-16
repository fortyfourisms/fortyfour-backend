CREATE TABLE sdm_csirt (
  id VARCHAR(36) PRIMARY KEY,
  id_csirt VARCHAR(36),
  nama_personel VARCHAR(255),
  jabatan_csirt VARCHAR(255),
  jabatan_perusahaan VARCHAR(255),
  skill TEXT,
  sertifikasi TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (id_csirt) REFERENCES csirt(id) ON DELETE SET NULL
);