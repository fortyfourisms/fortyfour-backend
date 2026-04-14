ALTER TABLE pertanyaan_identifikasi 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE pertanyaan_proteksi
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE pertanyaan_deteksi
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE pertanyaan_gulih 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE jawaban_identifikasi 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE jawaban_proteksi 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE jawaban_deteksi 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE jawaban_gulih 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;