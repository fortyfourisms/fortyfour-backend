-- Down migration
ALTER TABLE jawaban_identifikasi DROP COLUMN uuid;
ALTER TABLE jawaban_identifikasi_buffer DROP COLUMN uuid;

ALTER TABLE jawaban_proteksi DROP COLUMN uuid;
ALTER TABLE jawaban_proteksi_buffer DROP COLUMN uuid;

ALTER TABLE jawaban_deteksi DROP COLUMN uuid;
ALTER TABLE jawaban_deteksi_buffer DROP COLUMN uuid;

ALTER TABLE jawaban_gulih DROP COLUMN uuid;
ALTER TABLE jawaban_gulih_buffer DROP COLUMN uuid;
