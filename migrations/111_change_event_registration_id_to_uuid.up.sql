-- Matikan auto_increment dan drop PK agar bisa diubah ke CHAR
ALTER TABLE event_registrations MODIFY id BIGINT;
ALTER TABLE event_registrations DROP PRIMARY KEY;

-- Ubah id event_registrations
ALTER TABLE event_registrations
MODIFY id CHAR(36) NOT NULL;

-- Generate UUID baru untuk data yang sudah ada
UPDATE event_registrations
SET id = UUID();

-- Set kembali sebagai Primary Key
ALTER TABLE event_registrations ADD PRIMARY KEY (id);
