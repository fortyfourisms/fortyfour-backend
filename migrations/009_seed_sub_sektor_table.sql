-- ILMATE
INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Elektronik', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Otomotif', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Keamanan Siber', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE';


-- Industri Agro
INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Agro Bisnis', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Industri Agro';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Konstruksi', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Industri Agro';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Jasa', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Industri Agro';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Surveyor', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Industri Agro';


-- IKFT
INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Tekstil', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kimia', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kawasan Industri', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Farmasi', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';
