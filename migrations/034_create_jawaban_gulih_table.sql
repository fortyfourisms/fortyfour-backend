CREATE TABLE jawaban_gulih (
    id CHAR(36) PRIMARY KEY,
    pertanyaan_gulih_id CHAR(36) NOT NULL,
    perusahaan_id CHAR(36) NOT NULL,
    jawaban_gulih TEXT NOT NULL,

    CONSTRAINT fk_jawaban_gulih_pertanyaan
        FOREIGN KEY (pertanyaan_gulih_id)
        REFERENCES pertanyaan_gulih(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_gulih_perusahaan
        FOREIGN KEY (perusahaan_id)
        REFERENCES perusahaan(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_jawaban_gulih_perusahaan
        UNIQUE (perusahaan_id, pertanyaan_gulih_id)
);
