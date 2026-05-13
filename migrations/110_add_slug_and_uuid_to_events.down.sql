-- Drop FK
ALTER TABLE event_registrations
DROP FOREIGN KEY fk_event_registrations_event;

-- Balikkan id events ke BIGINT
ALTER TABLE events MODIFY id BIGINT AUTO_INCREMENT;
ALTER TABLE events DROP PRIMARY KEY;
ALTER TABLE events MODIFY id BIGINT AUTO_INCREMENT PRIMARY KEY;

-- Hapus slug
ALTER TABLE events
DROP COLUMN slug;

-- Balikkan event_id registrations
ALTER TABLE event_registrations
MODIFY event_id BIGINT NOT NULL;

-- Restore FK
ALTER TABLE event_registrations ADD CONSTRAINT fk_event_registrations_event 
FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
