-- Hapus foreign key lama 
ALTER TABLE survey_progress
DROP FOREIGN KEY fk_survey_progress_edit_reviewed_by;

-- Hapus kolom lama yang sudah tidak dipakai
ALTER TABLE survey_progress
DROP COLUMN edit_reviewed_at,
DROP COLUMN edit_reviewed_by;

-- Tambahkan kolom baru
ALTER TABLE survey_progress
ADD COLUMN edit_approved_at TIMESTAMP NULL AFTER edit_requested_at,
ADD COLUMN edit_approved_by CHAR(36) NULL AFTER edit_approved_at,
ADD COLUMN edit_rejected_at TIMESTAMP NULL AFTER edit_approved_by,
ADD COLUMN edit_rejected_by CHAR(36) NULL AFTER edit_rejected_at;

-- Tambahkan foreign key baru
ALTER TABLE survey_progress
ADD CONSTRAINT fk_survey_progress_edit_approved_by
FOREIGN KEY (edit_approved_by) REFERENCES users(id)
ON DELETE SET NULL ON UPDATE CASCADE,
ADD CONSTRAINT fk_survey_progress_edit_rejected_by
FOREIGN KEY (edit_rejected_by) REFERENCES users(id)
ON DELETE SET NULL ON UPDATE CASCADE;