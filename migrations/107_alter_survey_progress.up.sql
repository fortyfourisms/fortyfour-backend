ALTER TABLE survey_progress
ADD COLUMN status ENUM('draft', 'submitted', 'edit_requested', 'edit_approved', 'edit_rejected') NOT NULL DEFAULT 'draft' AFTER selesai,
ADD COLUMN edit_request_reason TEXT NULL AFTER status,
ADD COLUMN edit_request_response TEXT NULL AFTER edit_request_reason,
ADD COLUMN submitted_at TIMESTAMP NULL AFTER edit_request_response,
ADD COLUMN edit_requested_at TIMESTAMP NULL AFTER submitted_at,
ADD COLUMN edit_reviewed_at TIMESTAMP NULL AFTER edit_requested_at,
ADD COLUMN edit_reviewed_by CHAR(36) NULL AFTER edit_reviewed_at;

UPDATE survey_progress
SET status = CASE
	WHEN selesai = TRUE THEN 'submitted'
	ELSE 'draft'
END;

ALTER TABLE survey_progress
ADD CONSTRAINT fk_survey_progress_edit_reviewed_by
FOREIGN KEY (edit_reviewed_by) REFERENCES users(id)
ON DELETE SET NULL ON UPDATE CASCADE;