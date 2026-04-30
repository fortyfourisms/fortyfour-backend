CREATE TABLE responden (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,

    no_telepon VARCHAR(20) NOT NULL,

    sektor ENUM(
        'Keamanan Siber',
        'Jasa Konsultan dan Sertifikasi Keamanan Informasi',
        'Industri Pulp dan Kertas',
        'Jasa Konstruksi',
        'Jasa Sertifikasi, Inspeksi, Pengujian, dan Survey',
        'Jasa Industri',
        'Komponen Kendaraan Listrik',
        'Alat Kesehatan',
        'Peralatan Listrik',
        'Industri Kecil dan Menengah Furnitur, dan Bahan Bangunan',
        'Logam',
        'Permesinan dan Alat Mesin Pertanian',
        'Maritim, Alat Transportasi, dan Alat Pertahanan',
        'Elektronika dan Telematika',
        'Hasil Hutan dan Perkebunan',
        'Makanan, Hasil Laut, dan Perikanan',
        'Minuman, Tembakau, dan Bahan Penyegar',
        'Kemurgi, Oleokimia, dan Pakan',
        'Kimia hulu',
        'Kawasan Industri',
        'Semen, Keramik, dan Pengolahan Bahan Galian Nonlogam',
        'Tekstil, Kulit, dan Alas Kaki',
        'Industri Aneka',
        'Industri Bahan dan Produk Farmasi',
        'Kimia Hilir',
        'Lainnya'
    ) NOT NULL,

    sektor_lainnya VARCHAR(150),
    sertifikat_training VARCHAR(500) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY unique_user (user_id),

    CONSTRAINT fk_responden_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);