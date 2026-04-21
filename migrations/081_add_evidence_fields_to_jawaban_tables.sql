ALTER TABLE jawaban_identifikasi
ADD COLUMN evidence TEXT NULL AFTER jawaban_identifikasi,
ADD COLUMN validasi ENUM('yes', 'no') NULL AFTER evidence,
ADD COLUMN keterangan TEXT NULL AFTER validasi;

ALTER TABLE jawaban_proteksi
ADD COLUMN evidence TEXT NULL AFTER jawaban_proteksi,
ADD COLUMN validasi ENUM('yes', 'no') NULL AFTER evidence,
ADD COLUMN keterangan TEXT NULL AFTER validasi;

ALTER TABLE jawaban_deteksi
ADD COLUMN evidence TEXT NULL AFTER jawaban_deteksi,
ADD COLUMN validasi ENUM('yes', 'no') NULL AFTER evidence,
ADD COLUMN keterangan TEXT NULL AFTER validasi;

ALTER TABLE jawaban_pulih
ADD COLUMN evidence TEXT NULL AFTER jawaban_pulih,
ADD COLUMN validasi ENUM('yes', 'no') NULL AFTER evidence,
ADD COLUMN keterangan TEXT NULL AFTER validasi;