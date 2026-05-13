-- Drop FK dulu
ALTER TABLE event_registrations
DROP FOREIGN KEY fk_event_registrations_event;

-- Matikan auto_increment dan drop PK agar bisa diubah ke CHAR
ALTER TABLE events MODIFY id BIGINT;
ALTER TABLE events DROP PRIMARY KEY;

-- Ubah id events
ALTER TABLE events
MODIFY id CHAR(36) NOT NULL;

-- Generate UUID baru
UPDATE events
SET id = UUID();

-- Tambah slug
ALTER TABLE events
ADD COLUMN slug VARCHAR(255) NOT NULL AFTER id;

UPDATE events
SET slug = LOWER(REPLACE(judul, ' ', '-'));

-- Ubah event_id registrations
ALTER TABLE event_registrations
MODIFY event_id CHAR(36) NOT NULL;

-- Set Primary Key dan Unique
ALTER TABLE events ADD PRIMARY KEY (id);
ALTER TABLE events ADD UNIQUE KEY uq_events_slug (slug);

-- Restore Foreign Key
ALTER TABLE event_registrations ADD CONSTRAINT fk_event_registrations_event 
FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
