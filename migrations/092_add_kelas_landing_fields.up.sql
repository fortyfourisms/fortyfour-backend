-- Tambah kolom baru pada tabel kelas untuk menampilkan data pada landing page publik
ALTER TABLE kelas
    ADD COLUMN kategori          VARCHAR(100) NULL AFTER thumbnail,
    ADD COLUMN durasi_jp         INT          NULL AFTER kategori,
    ADD COLUMN penyelenggara     VARCHAR(255) NULL AFTER durasi_jp,
    ADD COLUMN target_peserta    TEXT         NULL AFTER penyelenggara,
    ADD COLUMN syarat_pendaftaran TEXT        NULL AFTER target_peserta,
    ADD COLUMN informasi_umum    TEXT         NULL AFTER syarat_pendaftaran;
