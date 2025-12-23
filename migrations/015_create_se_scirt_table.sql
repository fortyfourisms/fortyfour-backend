CREATE TABLE se_csirt (
  id CHAR(36) PRIMARY KEY,
  id_csirt CHAR(36) NOT NULL,

  nama_se VARCHAR(255) NOT NULL,
  ip_se VARCHAR(100) NOT NULL,
  as_number_se VARCHAR(100) NOT NULL,
  pengelola_se VARCHAR(255) NOT NULL,
  fitur_se TEXT,
  kategori_se ENUM('rendah','tinggi','strategis') NOT NULL DEFAULT 'rendah',

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT fk_se_csirt_csirt
    FOREIGN KEY (id_csirt)
    REFERENCES csirt(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
);
