-- DROP FK
ALTER TABLE risiko_eligibility DROP FOREIGN KEY fk_eligibility_responden, DROP FOREIGN KEY fk_eligibility_risiko;
ALTER TABLE risiko_alasan DROP FOREIGN KEY fk_alasan_responden, DROP FOREIGN KEY fk_alasan_risiko;
ALTER TABLE risiko_dampak DROP FOREIGN KEY fk_dampak_responden, DROP FOREIGN KEY fk_dampak_risiko;
ALTER TABLE risiko_pengendalian DROP FOREIGN KEY fk_pengendalian_responden, DROP FOREIGN KEY fk_pengendalian_risiko;
ALTER TABLE survey_progress DROP FOREIGN KEY fk_progress_responden, DROP FOREIGN KEY fk_progress_risiko;

-- MODIFY
ALTER TABLE risiko_eligibility MODIFY responden_id BIGINT NOT NULL, MODIFY risiko_id INT NOT NULL;
ALTER TABLE risiko_alasan MODIFY responden_id BIGINT NOT NULL, MODIFY risiko_id INT NOT NULL;
ALTER TABLE risiko_dampak MODIFY responden_id BIGINT NOT NULL, MODIFY risiko_id INT NOT NULL;
ALTER TABLE risiko_pengendalian MODIFY responden_id BIGINT NOT NULL, MODIFY risiko_id INT NOT NULL;
ALTER TABLE survey_progress MODIFY responden_id BIGINT NOT NULL;

-- ADD FK
ALTER TABLE risiko_eligibility ADD CONSTRAINT fk_eligibility_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE ON UPDATE CASCADE, ADD CONSTRAINT fk_eligibility_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE risiko_alasan ADD CONSTRAINT fk_alasan_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE ON UPDATE CASCADE, ADD CONSTRAINT fk_alasan_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE risiko_dampak ADD CONSTRAINT fk_dampak_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE ON UPDATE CASCADE, ADD CONSTRAINT fk_dampak_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE risiko_pengendalian ADD CONSTRAINT fk_pengendalian_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE ON UPDATE CASCADE, ADD CONSTRAINT fk_pengendalian_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE survey_progress ADD CONSTRAINT fk_progress_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE ON UPDATE CASCADE, ADD CONSTRAINT fk_progress_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE SET NULL ON UPDATE CASCADE;