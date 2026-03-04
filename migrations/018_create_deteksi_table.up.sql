CREATE TABLE IF NOT EXISTS deteksi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nilai_deteksi FLOAT,
    nilai_subdomain1 FLOAT NOT NULL,
    nilai_subdomain2 FLOAT NOT NULL,
    nilai_subdomain3 FLOAT NOT NULL
);
