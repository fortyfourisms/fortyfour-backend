INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Logam', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Logam' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Permesinan & alat pertanian', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Permesinan & alat pertanian' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Transportasi, maritim & pertahanan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Transportasi, maritim & pertahanan' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Elektronika & telematika', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'ILMATE'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Elektronika & telematika' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Hasil hutan & perkebunan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Hasil hutan & perkebunan' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Pangan & perikanan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Pangan & perikanan' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Minuman, tembakau & bahan penyegar', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Minuman, tembakau & bahan penyegar' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kemurgi, oleokimia & pakan', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'Agro'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Kemurgi, oleokimia & pakan' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kimia hulu', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Kimia hulu' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Kimia hilir & farmasi', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Kimia hilir & farmasi' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Semen, keramik & nonlogam', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Semen, keramik & nonlogam' AND ss.id_sektor = s.id);

INSERT INTO sub_sektor (id, id_sektor, nama_sub_sektor, created_at, updated_at)
SELECT UUID(), s.id, 'Tekstil, kulit & alas kaki', NOW(), NOW()
FROM sektor s WHERE s.nama_sektor = 'IKFT'
AND NOT EXISTS (SELECT 1 FROM sub_sektor ss WHERE ss.nama_sub_sektor = 'Tekstil, kulit & alas kaki' AND ss.id_sektor = s.id);
