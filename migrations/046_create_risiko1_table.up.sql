CREATE TABLE risiko (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama_risiko VARCHAR(255) NOT NULL,
    deskripsi TEXT,

    potensi_kejadian ENUM('Ya','Tidak') NOT NULL,
    dampak_reputasi ENUM('Tidak','Kecil','Sedang','Besar','Sangat Besar'),
    dampak_operasional ENUM('Tidak','Kecil','Sedang','Besar','Sangat Besar'),
    dampak_finansial ENUM('Tidak','Kecil','Sedang','Besar','Sangat Besar'),
    dampak_hukum ENUM('Tidak','Kecil','Sedang','Besar','Sangat Besar'),

    frekuensi ENUM('Kecil','Sedang','Besar','Sangat Besar'),

    ada_pengendalian ENUM('Ya','Tidak'),
    deskripsi_pengendalian TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);