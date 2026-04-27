ALTER TABLE ikas ADD COLUMN edit_request_status VARCHAR(20) DEFAULT 'none';
ALTER TABLE ikas ADD COLUMN edit_request_reason TEXT NULL;
