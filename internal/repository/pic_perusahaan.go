package repository

import (
	"database/sql"
	"fortyfour-backend/internal/models"
)

type PICPerusahaanRepository struct {
	DB *sql.DB
}

func NewPICPerusahaanRepository(db *sql.DB) *PICPerusahaanRepository {
	return &PICPerusahaanRepository{DB: db}
}

func (r *PICPerusahaanRepository) Create(pic models.PICPerusahaan) error {
	_, err := r.DB.Exec(`INSERT INTO pic_perusahaan (id, nama, telepon, id_perusahaan) VALUES (?, ?, ?, ?)`,
		pic.ID, pic.Nama, pic.Telepon, pic.IDPerusahaan)
	return err
}

func (r *PICPerusahaanRepository) GetAll() ([]models.PICPerusahaan, error) {
	rows, err := r.DB.Query(`SELECT * FROM pic_perusahaan`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pics []models.PICPerusahaan
	for rows.Next() {
		var pic models.PICPerusahaan
		rows.Scan(&pic.ID, &pic.Nama, &pic.Telepon, &pic.IDPerusahaan, new(string), new(string))
		pics = append(pics, pic)
	}
	return pics, nil
}

func (r *PICPerusahaanRepository) GetByID(id string) (models.PICPerusahaan, error) {
	var pic models.PICPerusahaan
	err := r.DB.QueryRow(`SELECT * FROM pic_perusahaan WHERE id=?`, id).
		Scan(&pic.ID, &pic.Nama, &pic.Telepon, &pic.IDPerusahaan, new(string), new(string))
	return pic, err
}

func (r *PICPerusahaanRepository) Update(id string, pic models.PICPerusahaan) error {
	_, err := r.DB.Exec(`UPDATE pic_perusahaan SET nama=?, telepon=?, id_perusahaan=? WHERE id=?`,
		pic.Nama, pic.Telepon, pic.IDPerusahaan, id)
	return err
}

func (r *PICPerusahaanRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM pic_perusahaan WHERE id=?`, id)
	return err
}

func (r *PICPerusahaanRepository) GetJoin() ([]models.PICPerusahaan, error) {
	rows, err := r.DB.Query(`
		SELECT pic_perusahaan.id, pic_perusahaan.nama, pic_perusahaan.telepon, pic_perusahaan.id_perusahaan, perusahaan.nama_perusahaan 
		FROM pic_perusahaan 
		JOIN perusahaan ON perusahaan.id = pic_perusahaan.id_perusahaan
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pics []models.PICPerusahaan
	for rows.Next() {
		var pic models.PICPerusahaan
		var namaCompany string
		rows.Scan(&pic.ID, &pic.Nama, &pic.Telepon, &pic.IDPerusahaan, &namaCompany)
		pic.Telepon = namaCompany // opsional kalau mau simpan di field lain silakan
		pics = append(pics, pic)
	}
	return pics, nil
}
