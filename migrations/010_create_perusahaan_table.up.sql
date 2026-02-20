CREATE TABLE IF NOT EXISTS perusahaan (
    id CHAR(36) PRIMARY KEY,
    photo VARCHAR(255),
    nama_perusahaan VARCHAR(255) NOT NULL,
    sektor ENUM(
        'Teknologi',
        'Keuangan',
        'Pendidikan',
        'Kesehatan',
        'Manufaktur',
        'Layanan',
        'Transportasi',
        'Lainnya'
    ) NOT NULL,
    alamat TEXT,
    telepon VARCHAR(50),
    email VARCHAR(100) UNIQUE,
    website VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
