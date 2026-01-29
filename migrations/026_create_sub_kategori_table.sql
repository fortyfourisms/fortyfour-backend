CREATE TABLE sub_kategori (
    id CHAR(36) PRIMARY KEY,
    kategori_id CHAR(36) NOT NULL,
    nama_sub_kategori VARCHAR(50) NOT NULL,

    CONSTRAINT fk_sub_kategori_kategori
        FOREIGN KEY (kategori_id)
        REFERENCES kategori(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
