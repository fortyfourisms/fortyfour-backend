CREATE TABLE perusahaan (
    id CHAR(36) PRIMARY KEY,
    photo VARCHAR(255),
    nama_perusahaan VARCHAR(255) NOT NULL,
    sektor VARCHAR(255) NOT NULL,
    alamat TEXT,
    telepon VARCHAR(50),
    email VARCHAR(100) UNIQUE,
    website VARCHAR(255)
);
