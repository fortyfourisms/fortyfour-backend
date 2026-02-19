CREATE TABLE IF NOT EXISTS ikas (
    id CHAR(36) PRIMARY KEY,
    id_perusahaan CHAR(36),
    id_identifikasi CHAR(36),
    id_proteksi CHAR(36),
    id_deteksi CHAR(36),
    id_gulih CHAR(36),

    tanggal DATETIME,
    responden VARCHAR(255) NOT NULL,
    telepon VARCHAR(50) NOT NULL,
    jabatan VARCHAR(255) NOT NULL,
    nilai_kematangan FLOAT NOT NULL,
    target_nilai FLOAT NOT NULL,

    FOREIGN KEY (id_perusahaan) REFERENCES perusahaan(id),
    FOREIGN KEY (id_identifikasi) REFERENCES identifikasi(id),
    FOREIGN KEY (id_proteksi) REFERENCES proteksi(id),
    FOREIGN KEY (id_deteksi) REFERENCES deteksi(id),
    FOREIGN KEY (id_gulih) REFERENCES gulih(id)
);
