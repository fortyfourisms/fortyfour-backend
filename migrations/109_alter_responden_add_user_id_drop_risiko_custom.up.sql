-- Tambahkan user_id ke tabel responden
ALTER TABLE responden
ADD COLUMN user_id CHAR(36) NOT NULL AFTER id;

-- Tambahkan unique constraint
ALTER TABLE responden
ADD CONSTRAINT uq_survey_responden_user UNIQUE (user_id);

-- Tambahkan foreign key ke users
ALTER TABLE responden
ADD CONSTRAINT fk_survey_responden_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- Drop tabel risiko_custom karena sudah tidak terpakai
DROP TABLE IF EXISTS risiko_custom;
