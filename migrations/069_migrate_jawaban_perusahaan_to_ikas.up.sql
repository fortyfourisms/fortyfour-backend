-- 1. Update Tabel Jawaban Domain (Identifikasi, Proteksi, Deteksi, Gulih)
-- Pola: Tambah ikas_id, map data, hapus perusahaan_id, tambah constraint baru

-- Jawaban Identifikasi
ALTER TABLE jawaban_identifikasi ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_identifikasi_id;
-- Map existing answers to ikas_id based on perusahaan_id (assuming current 1:1 relation)
UPDATE jawaban_identifikasi ji JOIN ikas i ON ji.perusahaan_id = i.id_perusahaan SET ji.ikas_id = i.id;
ALTER TABLE jawaban_identifikasi MODIFY COLUMN ikas_id CHAR(36) NOT NULL;
ALTER TABLE jawaban_identifikasi DROP FOREIGN KEY fk_jawaban_identifikasi_perusahaan;
ALTER TABLE jawaban_identifikasi DROP INDEX uq_jawaban_identifikasi_perusahaan;
ALTER TABLE jawaban_identifikasi DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_identifikasi ADD CONSTRAINT fk_jawaban_identifikasi_ikas FOREIGN KEY (ikas_id) REFERENCES ikas(id) ON DELETE CASCADE;
ALTER TABLE jawaban_identifikasi ADD CONSTRAINT uq_jawaban_identifikasi_ikas UNIQUE (ikas_id, pertanyaan_identifikasi_id);

-- Jawaban Proteksi
ALTER TABLE jawaban_proteksi ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_proteksi_id;
UPDATE jawaban_proteksi jp JOIN ikas i ON jp.perusahaan_id = i.id_perusahaan SET jp.ikas_id = i.id;
ALTER TABLE jawaban_proteksi MODIFY COLUMN ikas_id CHAR(36) NOT NULL;
ALTER TABLE jawaban_proteksi DROP FOREIGN KEY fk_jawaban_proteksi_perusahaan;
ALTER TABLE jawaban_proteksi DROP INDEX uq_jawaban_proteksi_perusahaan;
ALTER TABLE jawaban_proteksi DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_proteksi ADD CONSTRAINT fk_jawaban_proteksi_ikas FOREIGN KEY (ikas_id) REFERENCES ikas(id) ON DELETE CASCADE;
ALTER TABLE jawaban_proteksi ADD CONSTRAINT uq_jawaban_proteksi_ikas UNIQUE (ikas_id, pertanyaan_proteksi_id);

-- Jawaban Deteksi
ALTER TABLE jawaban_deteksi ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_deteksi_id;
UPDATE jawaban_deteksi jd JOIN ikas i ON jd.perusahaan_id = i.id_perusahaan SET jd.ikas_id = i.id;
ALTER TABLE jawaban_deteksi MODIFY COLUMN ikas_id CHAR(36) NOT NULL;
ALTER TABLE jawaban_deteksi DROP FOREIGN KEY fk_jawaban_deteksi_perusahaan;
ALTER TABLE jawaban_deteksi DROP INDEX uq_jawaban_deteksi_perusahaan;
ALTER TABLE jawaban_deteksi DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_deteksi ADD CONSTRAINT fk_jawaban_deteksi_ikas FOREIGN KEY (ikas_id) REFERENCES ikas(id) ON DELETE CASCADE;
ALTER TABLE jawaban_deteksi ADD CONSTRAINT uq_jawaban_deteksi_ikas UNIQUE (ikas_id, pertanyaan_deteksi_id);

-- Jawaban Gulih
ALTER TABLE jawaban_gulih ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_gulih_id;
UPDATE jawaban_gulih jg JOIN ikas i ON jg.perusahaan_id = i.id_perusahaan SET jg.ikas_id = i.id;
ALTER TABLE jawaban_gulih MODIFY COLUMN ikas_id CHAR(36) NOT NULL;
ALTER TABLE jawaban_gulih DROP FOREIGN KEY fk_jawaban_gulih_perusahaan;
ALTER TABLE jawaban_gulih DROP INDEX uq_jawaban_gulih_perusahaan;
ALTER TABLE jawaban_gulih DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_gulih ADD CONSTRAINT fk_jawaban_gulih_ikas FOREIGN KEY (ikas_id) REFERENCES ikas(id) ON DELETE CASCADE;
ALTER TABLE jawaban_gulih ADD CONSTRAINT uq_jawaban_gulih_ikas UNIQUE (ikas_id, pertanyaan_gulih_id);

-- 2. Update Tabel Buffer (Flush temporary data)
-- Karena data buffer bersifat sementara, kita bisa tambahkan ikas_id dan hapus perusahaan_id

ALTER TABLE jawaban_identifikasi_buffer ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_identifikasi_id;
-- Note: Indexing patterns vary, adjusting based on existing structure
ALTER TABLE jawaban_identifikasi_buffer DROP INDEX perusahaan_id;
ALTER TABLE jawaban_identifikasi_buffer DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_identifikasi_buffer ADD UNIQUE KEY uq_jawaban_id_buffer (ikas_id, pertanyaan_identifikasi_id);

ALTER TABLE jawaban_proteksi_buffer ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_proteksi_id;
ALTER TABLE jawaban_proteksi_buffer DROP INDEX uq_jawaban_proteksi_buffer;
ALTER TABLE jawaban_proteksi_buffer DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_proteksi_buffer ADD UNIQUE KEY uq_jawaban_proteksi_buffer (ikas_id, pertanyaan_proteksi_id);

ALTER TABLE jawaban_deteksi_buffer ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_deteksi_id;
ALTER TABLE jawaban_deteksi_buffer DROP INDEX uq_jawaban_deteksi_buffer;
ALTER TABLE jawaban_deteksi_buffer DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_deteksi_buffer ADD UNIQUE KEY uq_jawaban_deteksi_buffer (ikas_id, pertanyaan_deteksi_id);

ALTER TABLE jawaban_gulih_buffer ADD COLUMN IF NOT EXISTS ikas_id CHAR(36) AFTER pertanyaan_gulih_id;
ALTER TABLE jawaban_gulih_buffer DROP INDEX uq_jawaban_gulih_buffer;
ALTER TABLE jawaban_gulih_buffer DROP COLUMN perusahaan_id;
ALTER TABLE jawaban_gulih_buffer ADD UNIQUE KEY uq_jawaban_gulih_buffer (ikas_id, pertanyaan_gulih_id);

-- -- 3. Update Tabel Domain Summary (Identifikasi, Proteksi, Deteksi, Gulih)
-- -- Kita harus menghapus constraint unik perusahaan agar bisa punya banyak data summary (1 per asesmen)

-- -- Identifikasi Table
-- ALTER TABLE identifikasi ADD COLUMN ikas_id CHAR(36) AFTER perusahaan_id;
-- UPDATE identifikasi iden JOIN ikas i ON iden.perusahaan_id = i.id_perusahaan SET iden.ikas_id = i.id;
-- ALTER TABLE identifikasi DROP FOREIGN KEY fk_identifikasi_perusahaan;
-- ALTER TABLE identifikasi DROP INDEX uq_identifikasi_perusahaan;
-- ALTER TABLE identifikasi DROP COLUMN perusahaan_id;
-- ALTER TABLE identifikasi ADD UNIQUE KEY uq_identifikasi_ikas (ikas_id);

-- -- Proteksi Table
-- ALTER TABLE proteksi ADD COLUMN ikas_id CHAR(36) AFTER perusahaan_id;
-- UPDATE proteksi prot JOIN ikas i ON prot.perusahaan_id = i.id_perusahaan SET prot.ikas_id = i.id;
-- ALTER TABLE proteksi DROP FOREIGN KEY fk_proteksi_perusahaan;
-- ALTER TABLE proteksi DROP INDEX uq_proteksi_perusahaan;
-- ALTER TABLE proteksi DROP COLUMN perusahaan_id;
-- ALTER TABLE proteksi ADD UNIQUE KEY uq_proteksi_ikas (ikas_id);

-- -- Deteksi Table
-- ALTER TABLE deteksi ADD COLUMN ikas_id CHAR(36) AFTER perusahaan_id;
-- UPDATE deteksi det JOIN ikas i ON det.perusahaan_id = i.id_perusahaan SET det.ikas_id = i.id;
-- ALTER TABLE deteksi DROP INDEX uq_deteksi_perusahaan;
-- ALTER TABLE deteksi DROP COLUMN perusahaan_id;
-- ALTER TABLE deteksi ADD UNIQUE KEY uq_deteksi_ikas (ikas_id);

-- -- Gulih Table
-- ALTER TABLE gulih ADD COLUMN ikas_id CHAR(36) AFTER perusahaan_id;
-- UPDATE gulih g JOIN ikas i ON g.perusahaan_id = i.id_perusahaan SET g.ikas_id = i.id;
-- ALTER TABLE gulih DROP INDEX uq_gulih_perusahaan;
-- ALTER TABLE gulih DROP COLUMN perusahaan_id;
-- ALTER TABLE gulih ADD UNIQUE KEY uq_gulih_ikas (ikas_id);
