-- MASTER RISIKO
CREATE TABLE risiko (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama_risiko VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- JAWABAN RISIKO
CREATE TABLE risiko_responden (
    id INT AUTO_INCREMENT PRIMARY KEY,
    responden_id INT NOT NULL,

    risiko_id INT NOT NULL,

    pernah_terjadi ENUM('ya','tidak') NOT NULL,

    dampak_reputasi ENUM('tidak','kecil','sedang','besar','sangat_besar'),
    dampak_operasional ENUM('tidak','kecil','sedang','besar','sangat_besar'),
    dampak_finansial ENUM('tidak','kecil','sedang','besar','sangat_besar'),
    dampak_hukum ENUM('tidak','kecil','sedang','besar','sangat_besar'),

    frekuensi ENUM('kecil','sedang','besar','sangat_besar'),

    ada_pengendalian ENUM('ya','tidak'),
    deskripsi_pengendalian TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);