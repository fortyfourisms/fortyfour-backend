-- Tambah kembali kolom id_csirt
ALTER TABLE se ADD COLUMN id_csirt CHAR(36) DEFAULT NULL AFTER id_sub_sektor;

-- Tambah kembali FK
ALTER TABLE se ADD CONSTRAINT fk_se_csirt
    FOREIGN KEY (id_csirt) REFERENCES csirt(id)
    ON DELETE RESTRICT;
