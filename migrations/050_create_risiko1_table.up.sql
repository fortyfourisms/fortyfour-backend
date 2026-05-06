DROP TABLE IF EXISTS survey_progress;
DROP TABLE IF EXISTS risiko_pengendalian;
DROP TABLE IF EXISTS risiko_dampak;
DROP TABLE IF EXISTS risiko_alasan;
DROP TABLE IF EXISTS risiko_eligibility;
DROP TABLE IF EXISTS risiko;

-- 1. MASTER RISIKO
CREATE TABLE risiko (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode VARCHAR(50) NOT NULL,
    nama VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    urutan INT DEFAULT 0,
    aktif BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_kode (kode)
) ENGINE=InnoDB;

-- 2. STEP 1: ELIGIBILITY
CREATE TABLE risiko_eligibility (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    risiko_id INT NULL,
    pernah_terjadi BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_eligibility (responden_id, risiko_id)
) ENGINE=InnoDB;

-- 3. STEP 2A: ALASAN
CREATE TABLE risiko_alasan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    risiko_id INT NULL,
    alasan TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_alasan (responden_id, risiko_id)
) ENGINE=InnoDB;

-- 4. STEP 2B: DAMPAK
CREATE TABLE risiko_dampak (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    risiko_id INT NULL,

    dampak_reputasi ENUM('Tidak Signifikan','Cukup Signifikan','Signifikan','Sangat Signifikan'),
    dampak_operasional ENUM('Tidak Signifikan','Cukup Signifikan','Signifikan','Sangat Signifikan'),
    dampak_finansial ENUM('Tidak Signifikan','Cukup Signifikan','Signifikan','Sangat Signifikan'),
    dampak_hukum ENUM('Tidak Signifikan','Cukup Signifikan','Signifikan','Sangat Signifikan'),

    frekuensi ENUM('Kecil','Sedang','Besar','Sangat Besar'),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_dampak (responden_id, risiko_id)
) ENGINE=InnoDB;

-- 5. STEP 2C: PENGENDALIAN
CREATE TABLE risiko_pengendalian (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    risiko_id INT NULL,

    ada_pengendalian BOOLEAN NOT NULL,
    deskripsi_pengendalian TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT chk_pengendalian CHECK (
        (ada_pengendalian = false) OR
        (ada_pengendalian = true AND deskripsi_pengendalian IS NOT NULL)
    ),

    UNIQUE KEY uk_pengendalian (responden_id, risiko_id)
) ENGINE=InnoDB;

-- 6. TRACKING PROGRESS 
CREATE TABLE survey_progress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    risiko_id INT,
    langkah_saat_ini VARCHAR(50),
    selesai BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    terakhir_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_progress (responden_id)
) ENGINE=InnoDB;