CREATE TABLE IF NOT EXISTS se (
    id CHAR(36) PRIMARY KEY,
    id_perusahaan CHAR(36) NOT NULL,
    id_sub_sektor CHAR(36) NOT NULL,

    q1  VARCHAR(1) NOT NULL,
    q2  VARCHAR(1) NOT NULL,
    q3  VARCHAR(1) NOT NULL,
    q4  VARCHAR(1) NOT NULL,
    q5  VARCHAR(1) NOT NULL,
    q6  VARCHAR(1) NOT NULL,
    q7  VARCHAR(1) NOT NULL,
    q8  VARCHAR(1) NOT NULL,
    q9  VARCHAR(1) NOT NULL,
    q10 VARCHAR(1) NOT NULL,

    total_bobot INT NOT NULL,
    kategori_se ENUM('Strategis','Tinggi','Rendah') NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_se_perusahaan
        FOREIGN KEY (id_perusahaan) REFERENCES perusahaan(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_se_sub_sektor
        FOREIGN KEY (id_sub_sektor) REFERENCES sub_sektor(id)
        ON DELETE RESTRICT
);
