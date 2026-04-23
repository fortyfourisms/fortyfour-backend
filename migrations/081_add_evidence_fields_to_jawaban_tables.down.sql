ALTER TABLE jawaban_identifikasi
DROP COLUMN keterangan,
DROP COLUMN validasi,
DROP COLUMN evidence;

ALTER TABLE jawaban_proteksi
DROP COLUMN keterangan,
DROP COLUMN validasi,
DROP COLUMN evidence;

ALTER TABLE jawaban_deteksi
DROP COLUMN keterangan,
DROP COLUMN validasi,
DROP COLUMN evidence;

ALTER TABLE jawaban_gulih
DROP COLUMN keterangan,
DROP COLUMN validasi,
DROP COLUMN evidence;
