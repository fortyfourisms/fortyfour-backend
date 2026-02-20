CREATE TABLE IF NOT EXISTS proteksi (
    id CHAR(36) PRIMARY KEY,
    nilai_proteksi FLOAT,
    nilai_subdomain1 FLOAT NOT NULL,
    nilai_subdomain2 FLOAT NOT NULL,
    nilai_subdomain3 FLOAT NOT NULL,
    nilai_subdomain4 FLOAT NOT NULL,
    nilai_subdomain5 FLOAT NOT NULL,
    nilai_subdomain6 FLOAT NOT NULL
);