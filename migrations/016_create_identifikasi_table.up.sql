CREATE TABLE IF NOT EXISTS identifikasi (
    id CHAR(36) PRIMARY KEY,
    nilai_identifikasi FLOAT,
    nilai_subdomain1 FLOAT NOT NULL,
    nilai_subdomain2 FLOAT NOT NULL,
    nilai_subdomain3 FLOAT NOT NULL,
    nilai_subdomain4 FLOAT NOT NULL,
    nilai_subdomain5 FLOAT NOT NULL
);