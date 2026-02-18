-- ILMATE
INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Logam', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Permesinan & alat pertanian', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Transportasi, maritim & pertahanan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Elektronika & telematika', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE';


-- Agro
INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Hasil hutan & perkebunan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Pangan & perikanan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Minuman, tembakau & bahan penyegar', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kemurgi, oleokimia & pakan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro';


-- IKFT
INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kimia hulu', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kimia hilir & farmasi', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Semen, keramik & nonlogam', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Tekstil, kulit & alas kaki', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT';
