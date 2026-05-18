-- Up migration
ALTER TABLE jawaban_identifikasi ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;
ALTER TABLE jawaban_identifikasi_buffer ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;

ALTER TABLE jawaban_proteksi ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;
ALTER TABLE jawaban_proteksi_buffer ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;

ALTER TABLE jawaban_deteksi ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;
ALTER TABLE jawaban_deteksi_buffer ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;

ALTER TABLE jawaban_gulih ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;
ALTER TABLE jawaban_gulih_buffer ADD COLUMN uuid CHAR(36) NOT NULL AFTER id;

-- Populate existing rows
-- Note: Using UUID() directly in the UPDATE statement ensures each row gets a unique value for every row in most modern MySQL versions.
UPDATE jawaban_identifikasi SET uuid = UUID() WHERE uuid = '';
UPDATE jawaban_identifikasi_buffer SET uuid = UUID() WHERE uuid = '';
UPDATE jawaban_proteksi SET uuid = UUID() WHERE uuid = '';
UPDATE jawaban_proteksi_buffer SET uuid = UUID() WHERE uuid = '';
UPDATE jawaban_deteksi SET uuid = UUID() WHERE uuid = '';
UPDATE jawaban_deteksi_buffer SET uuid = UUID() WHERE uuid = '';
UPDATE jawaban_gulih SET uuid = UUID() WHERE uuid = '';
UPDATE jawaban_gulih_buffer SET uuid = UUID() WHERE uuid = '';

-- Add Unique constraints
ALTER TABLE jawaban_identifikasi ADD UNIQUE (uuid);
ALTER TABLE jawaban_identifikasi_buffer ADD UNIQUE (uuid);
ALTER TABLE jawaban_proteksi ADD UNIQUE (uuid);
ALTER TABLE jawaban_proteksi_buffer ADD UNIQUE (uuid);
ALTER TABLE jawaban_deteksi ADD UNIQUE (uuid);
ALTER TABLE jawaban_deteksi_buffer ADD UNIQUE (uuid);
ALTER TABLE jawaban_gulih ADD UNIQUE (uuid);
ALTER TABLE jawaban_gulih_buffer ADD UNIQUE (uuid);
