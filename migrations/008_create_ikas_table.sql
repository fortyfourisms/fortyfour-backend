CREATE TABLE IF NOT EXISTS ikas (
  `id` char(36) PRIMARY KEY,
  `id_perusahaan` char(36) NOT NULL,
  `tanggal` datetime,
  `responden` varchar(255) NOT NULL,
  `telepon` varchar(50),
  `jabatan` varchar(255),
  `nilai_kematangan` float NOT NULL,
  `target_nilai` float NOT NULL,
  `id_identifikasi` char(36) NOT NULL,
  `id_proteksi` char(36) NOT NULL,
  `id_deteksi` char(36) NOT NULL,
  `id_gulih` char(36) NOT NULL,
  FOREIGN KEY (`id_perusahaan`) REFERENCES `perusahaan` (`id`)
);
