ALTER TABLE ikas
ADD CONSTRAINT ikas_ibfk_2 FOREIGN KEY (id_identifikasi) REFERENCES identifikasi(id),
ADD CONSTRAINT ikas_ibfk_3 FOREIGN KEY (id_proteksi) REFERENCES proteksi(id),
ADD CONSTRAINT ikas_ibfk_4 FOREIGN KEY (id_deteksi) REFERENCES deteksi(id),
ADD CONSTRAINT ikas_ibfk_5 FOREIGN KEY (id_gulih) REFERENCES gulih(id);
