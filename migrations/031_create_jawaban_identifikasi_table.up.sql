CREATE TABLE IF NOT EXISTS jawaban_identifikasi (
    id CHAR(36) PRIMARY KEY,
    pertanyaan_identifikasi_id CHAR(36) NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_identifikasi TEXT NOT NULL,

    CONSTRAINT fk_jawaban_identifikasi_pertanyaan
        FOREIGN KEY (pertanyaan_identifikasi_id)
        REFERENCES pertanyaan_identifikasi(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_identifikasi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_jawaban_identifikasi_perusahaan
        UNIQUE (perusahaan_id, pertanyaan_identifikasi_id)
);
