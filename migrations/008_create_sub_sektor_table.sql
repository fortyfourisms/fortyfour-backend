CREATE TABLE sub_sektor (
    id CHAR(36) NOT NULL,
    id_sektor CHAR(36) NOT NULL,
    nama_sub_sektor VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id),
    INDEX idx_sub_sektor_sektor (id_sektor),
    
    CONSTRAINT fk_sub_sektor_sektor
        FOREIGN KEY (id_sektor)
        REFERENCES sektor(id)
        ON DELETE CASCADE
);
