DROP TABLE IF EXISTS survey_progress;
DROP TABLE IF EXISTS risiko_pengendalian;
DROP TABLE IF EXISTS risiko_dampak;
DROP TABLE IF EXISTS risiko_alasan;
DROP TABLE IF EXISTS risiko_eligibility;
DROP TABLE IF EXISTS assessment;
DROP TABLE IF EXISTS risiko;

-- 1. TABEL RISIKO (Master)
CREATE TABLE risiko (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode VARCHAR(50) NOT NULL,
    nama VARCHAR(255) NOT NULL,
    nama_en VARCHAR(255),
    deskripsi TEXT,
    urutan INT DEFAULT 0,
    aktif BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_kode (kode)
) ENGINE=InnoDB CHARSET=utf8mb4;

-- 2. TABEL ASSESSMENT (Header Sesi)
CREATE TABLE assessment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    kode_assessment VARCHAR(50) NOT NULL,
    tanggal_assessment DATE NOT NULL,
    status ENUM('berjalan','selesai','arsip') DEFAULT 'berjalan',
    total_risiko INT DEFAULT 10,
    risiko_selesai INT DEFAULT 0,
    persen_selesai DECIMAL(5,2) DEFAULT 0.00,
    catatan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_kode_assessment (kode_assessment),
    FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8mb4;

-- 3. TABEL RISIKO_ELIGIBILITY (Step 1)
CREATE TABLE risiko_eligibility (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    assessment_id INT NOT NULL,
    risiko_id INT NOT NULL,
    pernah_terjadi BOOLEAN NOT NULL,
    langkah_selanjutnya VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_eligibility (responden_id, assessment_id, risiko_id),
    FOREIGN KEY (responden_id) REFERENCES responden(id) ON DELETE CASCADE,
    FOREIGN KEY (assessment_id) REFERENCES assessment(id) ON DELETE CASCADE,
    FOREIGN KEY (risiko_id) REFERENCES risiko(id) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8mb4;

-- 4. TABEL RISIKO_ALASAN (Step 2a - Alur Tidak)
CREATE TABLE risiko_alasan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    assessment_id INT NOT NULL,
    risiko_id INT NOT NULL,
    alasan TEXT NOT NULL,
    selesai BOOLEAN DEFAULT TRUE,
    selesai_pada TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_alasan (responden_id, assessment_id, risiko_id),
    FOREIGN KEY (assessment_id) REFERENCES assessment(id) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8mb4;

-- 5. TABEL RISIKO_DAMPAK (Step 2b - Alur Ya)
CREATE TABLE risiko_dampak (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    assessment_id INT NOT NULL,
    risiko_id INT NOT NULL,
    
    dampak_reputasi TINYINT CHECK (dampak_reputasi BETWEEN 1 AND 4),
    dampak_operasional TINYINT CHECK (dampak_operasional BETWEEN 1 AND 4),
    dampak_finansial TINYINT CHECK (dampak_finansial BETWEEN 1 AND 4),
    dampak_hukum TINYINT CHECK (dampak_hukum BETWEEN 1 AND 4),
    frekuensi TINYINT CHECK (frekuensi BETWEEN 1 AND 4),
    label_frekuensi VARCHAR(50),
    
    -- Auto-calculate
    total_skor_dampak TINYINT AS (
        COALESCE(dampak_reputasi,0) + COALESCE(dampak_operasional,0) + 
        COALESCE(dampak_finansial,0) + COALESCE(dampak_hukum,0)
    ) STORED,
    
    level_risiko VARCHAR(20) AS (
        CASE 
            WHEN (COALESCE(dampak_reputasi,0) + COALESCE(dampak_operasional,0) + 
                  COALESCE(dampak_finansial,0) + COALESCE(dampak_hukum,0)) <= 8 THEN 'Rendah'
            WHEN (COALESCE(dampak_reputasi,0) + COALESCE(dampak_operasional,0) + 
                  COALESCE(dampak_finansial,0) + COALESCE(dampak_hukum,0)) <= 12 THEN 'Sedang'
            WHEN (COALESCE(dampak_reputasi,0) + COALESCE(dampak_operasional,0) + 
                  COALESCE(dampak_finansial,0) + COALESCE(dampak_hukum,0)) <= 16 THEN 'Tinggi'
            ELSE 'Sangat Tinggi'
        END
    ) STORED,
    
    langkah_selanjutnya VARCHAR(50) DEFAULT 'pengendalian',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_dampak (responden_id, assessment_id, risiko_id),
    FOREIGN KEY (assessment_id) REFERENCES assessment(id) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8mb4;

-- 6. TABEL RISIKO_PENGENDALIAN (Step 2c - Alur Ya)
CREATE TABLE risiko_pengendalian (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    assessment_id INT NOT NULL,
    risiko_id INT NOT NULL,
    ada_pengendalian BOOLEAN NOT NULL,
    deskripsi_pengendalian TEXT,
    selesai BOOLEAN DEFAULT TRUE,
    selesai_pada TIMESTAMP NULL,
    langkah_selanjutnya VARCHAR(50) DEFAULT 'selesai',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_deskripsi CHECK (
        (ada_pengendalian = false) OR 
        (ada_pengendalian = true AND deskripsi_pengendalian IS NOT NULL 
         AND LENGTH(TRIM(deskripsi_pengendalian)) > 0)
    ),
    
    UNIQUE KEY uk_pengendalian (responden_id, assessment_id, risiko_id),
    FOREIGN KEY (assessment_id) REFERENCES assessment(id) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8mb4;

-- 7. TABEL SURVEY_PROGRESS (Tracking)
CREATE TABLE survey_progress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,
    assessment_id INT NOT NULL,
    risiko_id INT,
    langkah_saat_ini VARCHAR(50) NOT NULL,
    nomor_risiko INT DEFAULT 1,
    history_langkah JSON,
    selesai BOOLEAN DEFAULT FALSE,
    risiko_selesai JSON,
    terakhir_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_progress (responden_id, assessment_id),
    FOREIGN KEY (assessment_id) REFERENCES assessment(id) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8mb4;