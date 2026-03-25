CREATE TABLE risiko_survey (
	id INT AUTO_INCREMENT PRIMARY KEY,
	responden_id INT,
	risiko_ip BOOLEAN,
	dampak_reputasi VARCHAR(50),
	dampak_operasional VARCHAR(50),
	dampak_finansial VARCHAR(50),
	dampak_hukum VARCHAR(50),
	frekuensi VARCHAR(50),
	ada_pengendalian BOOLEAN,
	tindakan_pengendalian TEXT
);