DELETE FROM `pertanyaan_identifikasi`;
DELETE FROM `pertanyaan_proteksi`;
DELETE FROM `pertanyaan_deteksi`;
DELETE FROM `pertanyaan_gulih`;

DELETE FROM `sub_kategori`;
DELETE FROM `kategori`;
DELETE FROM `domain`;
DELETE FROM `ruang_lingkup`;


ALTER TABLE `domain` AUTO_INCREMENT = 1;
ALTER TABLE `pertanyaan_identifikasi` AUTO_INCREMENT = 1;
ALTER TABLE `pertanyaan_proteksi` AUTO_INCREMENT = 1;
ALTER TABLE `pertanyaan_deteksi` AUTO_INCREMENT = 1;
ALTER TABLE `pertanyaan_gulih` AUTO_INCREMENT = 1;
ALTER TABLE `sub_kategori` AUTO_INCREMENT = 1;
ALTER TABLE `kategori` AUTO_INCREMENT = 1;
ALTER TABLE `domain` AUTO_INCREMENT = 1;
ALTER TABLE `ruang_lingkup` AUTO_INCREMENT = 1;