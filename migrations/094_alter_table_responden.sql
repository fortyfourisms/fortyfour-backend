ALTER TABLE responden
ADD COLUMN nama VARCHAR(100) AFTER user_id,
MODIFY no_telepon VARCHAR(30) NOT NULL,
MODIFY sektor ENUM(
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
) NOT NULL;

INSERT INTO users (id)
SELECT DISTINCT r.user_id
FROM responden r
LEFT JOIN users u ON r.user_id = u.id
WHERE u.id IS NULL;

ALTER TABLE responden
ADD CONSTRAINT fk_responden_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE
ON UPDATE CASCADE;