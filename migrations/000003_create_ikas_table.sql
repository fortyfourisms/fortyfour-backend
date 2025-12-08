CREATE TABLE IF NOT EXISTS ikas (
  `id` char(36) PRIMARY KEY,
  `id_stakeholder` char(36),
  `tanggal` datetime,
  `responden` varchar(255),
  `telepon` varchar(50),
  `jabatan` varchar(255),
  `nilai_kematangan` float,
  `target_nilai` float,
  `id_identifikasi` char(36),
  `id_proteksi` char(36),
  `id_deteksi` char(36),
  `id_gulih` char(36),
  FOREIGN KEY (`id_stakeholder`) REFERENCES `stakeholders` (`id`)
);
