-- REVERTING UUID TO INT IS DESTRUCTIVE IF UUID DATA WAS GENERATED
-- THIS IS A SKELETON REVERT SCRIPT

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

ALTER TABLE risiko MODIFY id INT AUTO_INCREMENT;
ALTER TABLE responden MODIFY id BIGINT AUTO_INCREMENT;

-- ... other alters ...

-- RE-ADD FKs
-- ...
