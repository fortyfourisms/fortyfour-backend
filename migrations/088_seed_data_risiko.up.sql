INSERT INTO risiko (kode, nama, deskripsi, urutan) VALUES

-- 1
('ip-theft',
'Pencurian Intellectual Property Perusahaan',
'Intellectual Property atau Hak Kekayaan Intelektual mencakup paten, hak cipta, merek dagang, desain industri, rahasia dagang, serta inovasi lainnya yang menjadi aset strategis perusahaan. Dalam era Industri 4.0, penggunaan teknologi digital meningkatkan risiko pencurian HAKI melalui serangan siber, insider threat, maupun kebocoran data yang tidak disengaja.',
1),

-- 2
('data-breach',
'Kebocoran Data',
'Kebocoran data terjadi ketika informasi sensitif perusahaan atau pelanggan tersebar secara tidak sah atau diperjualbelikan. Hal ini dapat disebabkan oleh peretasan, kelalaian karyawan, kegagalan sistem keamanan, atau kesalahan manusia.',
2),

-- 3
('kerusakan-perangkat',
'Kerusakan Perangkat Fisik pada Fasilitas Manufaktur',
'Kerusakan perangkat fisik mencakup kegagalan mesin atau infrastruktur produksi akibat faktor internal seperti pemeliharaan buruk atau kesalahan teknis, maupun faktor eksternal seperti bencana atau kecelakaan.',
3),

-- 4
('kehilangan-peralatan',
'Kehilangan Peralatan Fisik Fasilitas Manufaktur',
'Kehilangan peralatan mencakup hilangnya mesin atau alat produksi akibat pencurian, kesalahan manajemen inventaris, masalah logistik, atau bencana yang menyebabkan kerusakan permanen.',
4),

-- 5
('human-error',
'Human Error',
'Human error adalah kesalahan tenaga kerja dalam operasional, baik disengaja maupun tidak. Dapat terjadi akibat kurangnya keterampilan, kelelahan, sistem kerja yang buruk, atau tekanan kerja tinggi.',
5),

-- 6
('third-party-breach',
'Pelanggaran Keamanan oleh Pihak Ketiga',
'Pelanggaran keamanan oleh pihak ketiga terjadi ketika vendor atau mitra bisnis menyebabkan kebocoran data atau akses tidak sah. Hal ini sering disebabkan oleh lemahnya kontrol akses dan standar keamanan yang tidak memadai.',
6),

-- 7
('management-lack-tech',
'Kurangnya Inisiatif Manajemen dalam Implementasi Teknologi',
'Kurangnya dukungan manajemen dalam adopsi teknologi dapat menghambat efisiensi dan daya saing. Hal ini bisa disebabkan oleh minimnya pemahaman, keterbatasan sumber daya, atau resistensi terhadap perubahan.',
7),

-- 8
('malware-attack',
'Serangan Virus dan Malware',
'Serangan malware mencakup virus, trojan, ransomware, spyware, dan lainnya yang bertujuan merusak, mencuri, atau mengganggu sistem perusahaan.',
8),

-- 9
('ddos-attack',
'Serangan DDoS (Distributed Denial of Service)',
'Serangan DDoS bertujuan membuat sistem tidak dapat diakses dengan membanjiri trafik secara berlebihan menggunakan botnet.',
9),

-- 10
('phishing',
'Serangan Phishing',
'Phishing adalah upaya penipuan untuk mendapatkan informasi sensitif melalui email, pesan, atau situs palsu yang menyerupai sumber terpercaya.',
10),

-- 11
('zero-day',
'Serangan Zero-Day',
'Serangan zero-day mengeksploitasi celah keamanan yang belum diketahui atau belum diperbaiki oleh pengembang.',
11),

-- 12
('ransomware',
'Serangan Ransomware',
'Ransomware mengenkripsi data perusahaan dan meminta tebusan agar data dapat diakses kembali. Jika tidak dibayar, data dapat hilang atau diperjualbelikan.',
12),

-- 13
('brute-force',
'Serangan Brute Force',
'Serangan brute force mencoba berbagai kombinasi password secara otomatis hingga menemukan yang benar, biasanya karena lemahnya sistem autentikasi.',
13);