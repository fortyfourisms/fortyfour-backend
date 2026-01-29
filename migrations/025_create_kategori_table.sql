CREATE TABLE kategori (
    id CHAR(36) PRIMARY KEY,
    domain_id CHAR(36) NOT NULL,
    nama_kategori VARCHAR(50) NOT NULL,

    CONSTRAINT fk_kategori_domain
        FOREIGN KEY (domain_id)
        REFERENCES domain(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
