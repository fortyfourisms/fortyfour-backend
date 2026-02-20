CREATE TABLE IF NOT EXISTS pertanyaan_deteksi (
    id CHAR(36) PRIMARY KEY,
    sub_kategori_id CHAR(36) NOT NULL,
    ruang_lingkup_id CHAR(36) NOT NULL,
    pertanyaan_deteksi TEXT NOT NULL,
    index0 TEXT,
    index1 TEXT,
    index2 TEXT,
    index3 TEXT,
    index4 TEXT,
    index5 TEXT,

    CONSTRAINT fk_pertanyaan_deteksi_sub_kategori
        FOREIGN KEY (sub_kategori_id)
        REFERENCES sub_kategori(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_pertanyaan_deteksi_ruang_lingkup
        FOREIGN KEY (ruang_lingkup_id)
        REFERENCES ruang_lingkup(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
