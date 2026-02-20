CREATE TABLE IF NOT EXISTS sub_kategori (
    id CHAR(36) PRIMARY KEY,
    kategori_id CHAR(36) NOT NULL,
    nama_sub_kategori VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_sub_kategori_kategori
        FOREIGN KEY (kategori_id)
        REFERENCES kategori(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
