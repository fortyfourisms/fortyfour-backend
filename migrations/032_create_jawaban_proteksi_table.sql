CREATE TABLE jawaban_proteksi (
    id CHAR(36) PRIMARY KEY,
    pertanyaan_proteksi_id CHAR(36) NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_proteksi TEXT NOT NULL,

    CONSTRAINT fk_jawaban_proteksi_pertanyaan
        FOREIGN KEY (pertanyaan_proteksi_id)
        REFERENCES pertanyaan_proteksi(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_proteksi_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_jawaban_proteksi_perusahaan
        UNIQUE (perusahaan_id, pertanyaan_proteksi_id)
);
