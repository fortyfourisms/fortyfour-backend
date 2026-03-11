CREATE TABLE IF NOT EXISTS pertanyaan_deteksi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sub_kategori_id INT NOT NULL,
    ruang_lingkup_id INT NOT NULL,
    pertanyaan_deteksi TEXT NOT NULL,
    index0 TEXT,
    index1 TEXT,
    index2 TEXT,
    index3 TEXT,
    index4 TEXT,
    index5 TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

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
