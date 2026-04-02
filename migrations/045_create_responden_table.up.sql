CREATE TABLE responden (
    id INT AUTO_INCREMENT PRIMARY KEY,
    
    nama_lengkap VARCHAR(150) NOT NULL,
    jabatan VARCHAR(100) NOT NULL,
    perusahaan VARCHAR(150) NOT NULL,
    
    email VARCHAR(150) NOT NULL UNIQUE,
    no_telepon VARCHAR(20) NOT NULL,
    
    sektor ENUM(
        'Industri Makanan dan Minuman',
        'Industri Tekstil dan Pakaian',
        'Industri Kimia',
        'Industri Otomotif',
        'Industri Elektronik',
        'Industri Farmasi',
        'Industri Alat Kesehatan',
        'Jasa Konstruksi',
        'Industri Keamanan Siber',
        'Industri Pertahanan',
        'Lainnya'
    ) NOT NULL,
    sektor_lainnya VARCHAR(150) DEFAULT NULL,
    sertifikat_training VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);