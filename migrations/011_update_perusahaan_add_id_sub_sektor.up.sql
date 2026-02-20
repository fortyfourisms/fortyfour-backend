ALTER TABLE perusahaan
    DROP COLUMN sektor;

ALTER TABLE perusahaan
    ADD COLUMN id_sub_sektor CHAR(36) NOT NULL AFTER nama_perusahaan;

ALTER TABLE perusahaan
    ADD CONSTRAINT fk_perusahaan_sub_sektor
        FOREIGN KEY (id_sub_sektor)
        REFERENCES sub_sektor(id)
        ON DELETE RESTRICT;