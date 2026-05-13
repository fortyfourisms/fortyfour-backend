-- 1. DROP FOREIGN KEYS FIRST
ALTER TABLE risiko_eligibility DROP FOREIGN KEY fk_eligibility_responden;
ALTER TABLE risiko_eligibility DROP FOREIGN KEY fk_eligibility_risiko;
ALTER TABLE risiko_alasan DROP FOREIGN KEY fk_alasan_responden;
ALTER TABLE risiko_alasan DROP FOREIGN KEY fk_alasan_risiko;
ALTER TABLE risiko_dampak DROP FOREIGN KEY fk_dampak_responden;
ALTER TABLE risiko_dampak DROP FOREIGN KEY fk_dampak_risiko;
ALTER TABLE risiko_pengendalian DROP FOREIGN KEY fk_pengendalian_responden;
ALTER TABLE risiko_pengendalian DROP FOREIGN KEY fk_pengendalian_risiko;
ALTER TABLE survey_progress DROP FOREIGN KEY fk_progress_responden;
ALTER TABLE survey_progress DROP FOREIGN KEY fk_progress_risiko;

-- 2. ALTER PRIMARY TABLES
-- RISIKO
ALTER TABLE risiko MODIFY id VARCHAR(36) NOT NULL;
ALTER TABLE risiko MODIFY id VARCHAR(36) NOT NULL DEFAULT (UUID());

-- RESPONDEN
ALTER TABLE responden MODIFY id VARCHAR(36) NOT NULL;
ALTER TABLE responden MODIFY id VARCHAR(36) NOT NULL DEFAULT (UUID());

-- 3. ALTER STEP TABLES (IDs and FKs)
ALTER TABLE risiko_eligibility MODIFY id VARCHAR(36) NOT NULL DEFAULT (UUID());
ALTER TABLE risiko_eligibility MODIFY responden_id VARCHAR(36) NOT NULL;
ALTER TABLE risiko_eligibility MODIFY risiko_id VARCHAR(36) NULL;

ALTER TABLE risiko_alasan MODIFY id VARCHAR(36) NOT NULL DEFAULT (UUID());
ALTER TABLE risiko_alasan MODIFY responden_id VARCHAR(36) NOT NULL;
ALTER TABLE risiko_alasan MODIFY risiko_id VARCHAR(36) NULL;

ALTER TABLE risiko_dampak MODIFY id VARCHAR(36) NOT NULL DEFAULT (UUID());
ALTER TABLE risiko_dampak MODIFY responden_id VARCHAR(36) NOT NULL;
ALTER TABLE risiko_dampak MODIFY risiko_id VARCHAR(36) NULL;

ALTER TABLE risiko_pengendalian MODIFY id VARCHAR(36) NOT NULL DEFAULT (UUID());
ALTER TABLE risiko_pengendalian MODIFY responden_id VARCHAR(36) NOT NULL;
ALTER TABLE risiko_pengendalian MODIFY risiko_id VARCHAR(36) NULL;

ALTER TABLE survey_progress MODIFY id VARCHAR(36) NOT NULL DEFAULT (UUID());
ALTER TABLE survey_progress MODIFY responden_id VARCHAR(36) NOT NULL;
ALTER TABLE survey_progress MODIFY risiko_id VARCHAR(36) NULL;

-- 4. RE-ADD FOREIGN KEYS
ALTER TABLE risiko_eligibility ADD CONSTRAINT fk_eligibility_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE;
ALTER TABLE risiko_eligibility ADD CONSTRAINT fk_eligibility_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE;

ALTER TABLE risiko_alasan ADD CONSTRAINT fk_alasan_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE;
ALTER TABLE risiko_alasan ADD CONSTRAINT fk_alasan_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE;

ALTER TABLE risiko_dampak ADD CONSTRAINT fk_dampak_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE;
ALTER TABLE risiko_dampak ADD CONSTRAINT fk_dampak_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE;

ALTER TABLE risiko_pengendalian ADD CONSTRAINT fk_pengendalian_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE;
ALTER TABLE risiko_pengendalian ADD CONSTRAINT fk_pengendalian_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE;

ALTER TABLE survey_progress ADD CONSTRAINT fk_progress_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE;
ALTER TABLE survey_progress ADD CONSTRAINT fk_progress_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE SET NULL;
