ALTER TABLE perusahaan
    DROP FOREIGN KEY fk_perusahaan_sub_sektor;

ALTER TABLE perusahaan
    DROP COLUMN id_sub_sektor;

ALTER TABLE perusahaan
    ADD COLUMN sektor ENUM(
        'Teknologi',
        'Keuangan',
        'Pendidikan',
        'Kesehatan',
        'Manufaktur',
        'Layanan',
        'Transportasi',
        'Lainnya'
    ) NOT NULL AFTER nama_perusahaan;
