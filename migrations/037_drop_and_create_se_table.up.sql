DROP TABLE IF EXISTS se;

CREATE TABLE se (
    id CHAR(36) PRIMARY KEY,
    id_perusahaan CHAR(36) NOT NULL,
    id_sub_sektor CHAR(36),
    id_csirt CHAR(36),

    nilai_investasi VARCHAR(1) NOT NULL CHECK (nilai_investasi IN ('A', 'B', 'C')),
    anggaran_operasional VARCHAR(1) NOT NULL CHECK (anggaran_operasional IN ('A', 'B', 'C')),
    kepatuhan_peraturan VARCHAR(1) NOT NULL CHECK (kepatuhan_peraturan IN ('A', 'B', 'C')),
    teknik_kriptografi VARCHAR(1) NOT NULL CHECK (teknik_kriptografi IN ('A', 'B', 'C')),
    jumlah_pengguna VARCHAR(1) NOT NULL CHECK (jumlah_pengguna IN ('A', 'B', 'C')),
    data_pribadi VARCHAR(1) NOT NULL CHECK (data_pribadi IN ('A', 'B', 'C')),
    klasifikasi_data VARCHAR(1) NOT NULL CHECK (klasifikasi_data IN ('A', 'B', 'C')),
    kekritisan_proses VARCHAR(1) NOT NULL CHECK (kekritisan_proses IN ('A', 'B', 'C')),
    dampak_kegagalan VARCHAR(1) NOT NULL CHECK (dampak_kegagalan IN ('A', 'B', 'C')),
    potensi_kerugian_dan_dampak_negatif VARCHAR(1) NOT NULL CHECK (potensi_kerugian_dan_dampak_negatif IN ('A', 'B', 'C')),

    nama_se VARCHAR(255) NOT NULL,
    ip_se VARCHAR(100) NOT NULL,
    as_number_se VARCHAR(100) NOT NULL,
    pengelola_se VARCHAR(255) NOT NULL,
    fitur_se TEXT,

    total_bobot INT NOT NULL,
    kategori_se ENUM('Strategis', 'Tinggi', 'Rendah') NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_se_perusahaan
        FOREIGN KEY (id_perusahaan) REFERENCES perusahaan(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_se_sub_sektor
        FOREIGN KEY (id_sub_sektor) REFERENCES sub_sektor(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_se_csirt
        FOREIGN KEY (id_csirt) REFERENCES csirt(id)
        ON DELETE RESTRICT
);

DROP TABLE IF EXISTS se_csirt;
