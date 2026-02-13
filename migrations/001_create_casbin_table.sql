CREATE TABLE casbin_rule (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ptype VARCHAR(100),
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100),
    UNIQUE KEY idx_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
);


INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES
-- Admin: full access
('p', 'admin', '*', '*', NULL, NULL, NULL),

-- User: read-only posts
('p', 'user', '/api/posts', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/perusahaan', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/pic', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/jabatan', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/identifikasi', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/gulih', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/proteksi', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/deteksi', 'GET', NULL, NULL, NULL),
('p', 'user', '/api/ikas', 'GET', NULL, NULL, NULL);

