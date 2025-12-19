CREATE TABLE se_csirt (
  id VARCHAR(36) PRIMARY KEY,
  id_csirt VARCHAR(36),
  nama_se VARCHAR(255),
  ip_se VARCHAR(100),
  as_number_se VARCHAR(100),
  pengelola_se VARCHAR(255),
  fitur_se TEXT,
  kategori_se ENUM('rendah','tinggi','strategis') DEFAULT 'rendah',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (id_csirt) REFERENCES csirt(id) ON DELETE SET NULL
);