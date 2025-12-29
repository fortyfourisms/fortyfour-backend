package utils

// GetKategoriTingkatKematangan menentukan kategori berdasarkan nilai
func GetKategoriTingkatKematangan(nilai float64) string {
	if nilai >= 0 && nilai <= 1.5 {
		return "Level 1 - Awal"
	} else if nilai >= 1.51 && nilai <= 2.5 {
		return "Level 2 - Berulang"
	} else if nilai >= 2.51 && nilai <= 3.5 {
		return "Level 3 - Terdefinisi"
	} else if nilai >= 3.51 && nilai <= 4.5 {
		return "Level 4 - Terkelola"
	} else if nilai >= 4.51 && nilai <= 5 {
		return "Level 5 - Inovatif"
	}
	return "Tidak Terdefinisi"
}
