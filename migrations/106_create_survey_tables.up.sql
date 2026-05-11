-- TABEL RESPONDEN
CREATE TABLE IF NOT EXISTS responden (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    id_perusahaan VARCHAR(255) NOT NULL,
    nama_lengkap VARCHAR(255) NOT NULL,
    jabatan VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    no_telepon VARCHAR(50) NOT NULL,
    sertifikat_training VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- TABEL RISIKO ELIGIBILITY (LANGKAH 1)
CREATE TABLE IF NOT EXISTS risiko_eligibility (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    responden_id BIGINT NOT NULL,
    risiko_id BIGINT NULL,
    pernah_terjadi BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_eligibility_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE,
    CONSTRAINT fk_eligibility_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE,
    UNIQUE KEY uk_eligibility (responden_id, risiko_id)
);

-- TABEL RISIKO ALASAN (LANGKAH 2A)
CREATE TABLE IF NOT EXISTS risiko_alasan (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    responden_id BIGINT NOT NULL,
    risiko_id BIGINT NULL,
    alasan TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_alasan_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE,
    CONSTRAINT fk_alasan_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE,
    UNIQUE KEY uk_alasan (responden_id, risiko_id)
);

-- TABEL RISIKO DAMPAK (LANGKAH 2B)
CREATE TABLE IF NOT EXISTS risiko_dampak (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    responden_id BIGINT NOT NULL,
    risiko_id BIGINT NULL,
    dampak_reputasi ENUM('Tidak Signifikan', 'Cukup Signifikan', 'Signifikan', 'Sangat Signifikan') NOT NULL,
    dampak_operasional ENUM('Tidak Signifikan', 'Cukup Signifikan', 'Signifikan', 'Sangat Signifikan') NOT NULL,
    dampak_finansial ENUM('Tidak Signifikan', 'Cukup Signifikan', 'Signifikan', 'Sangat Signifikan') NOT NULL,
    dampak_hukum ENUM('Tidak Signifikan', 'Cukup Signifikan', 'Signifikan', 'Sangat Signifikan') NOT NULL,
    frekuensi ENUM('Kecil', 'Sedang', 'Besar', 'Sangat Besar') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_dampak_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE,
    CONSTRAINT fk_dampak_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE,
    UNIQUE KEY uk_dampak (responden_id, risiko_id)
);

-- TABEL RISIKO PENGENDALIAN (LANGKAH 2C)
CREATE TABLE IF NOT EXISTS risiko_pengendalian (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    responden_id BIGINT NOT NULL,
    risiko_id BIGINT NULL,
    ada_pengendalian BOOLEAN NOT NULL,
    deskripsi_pengendalian TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_pengendalian_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE,
    CONSTRAINT fk_pengendalian_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE,
    UNIQUE KEY uk_pengendalian (responden_id, risiko_id)
);

-- TABEL SURVEY PROGRESS
CREATE TABLE IF NOT EXISTS survey_progress (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    responden_id BIGINT NOT NULL,
    risiko_id BIGINT NULL,
    langkah_saat_ini VARCHAR(100) NULL,
    selesai BOOLEAN NOT NULL DEFAULT FALSE,
    terakhir_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_progress_responden FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE,
    CONSTRAINT fk_progress_risiko FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE SET NULL,
    UNIQUE KEY uk_progress (responden_id)
);
