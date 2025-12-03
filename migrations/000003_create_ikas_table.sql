CREATE TABLE `ikas` (
  `id` integer PRIMARY KEY,
  `id_stakeholder` integer,
  `tanggal` datetime,
  `responden` varchar(255),
  `telepon` integer,
  `jabatan` varchar(255),
  `nilai_kematangan` float,
  `target_nilai` float,
  `id_identifikasi` integer,
  `id_proteksi` integer,
  `id_deteksi` integer,
  `id_gulih` integer
);