package repository

import (
	"database/sql"
	"fortyfour-backend/internal/models"
)

type IkasRepository interface {
	Create(ikas *models.Ikas) error
	GetAll() ([]models.Ikas, error)
	GetByID(id int) (*models.Ikas, error)
	Update(id int, ikas *models.Ikas) error
	Delete(id int) error
}

type ikasRepository struct {
	db *sql.DB
}

func NewIkasRepository(db *sql.DB) IkasRepository {
	return &ikasRepository{db: db}
}

func (r *ikasRepository) Create(ikas *models.Ikas) error {
	query := `INSERT INTO ikas (id_stakeholder, tanggal, responden, telepon, jabatan, 
				nilai_kematangan, target_nilai, id_identifikasi, id_proteksi, id_deteksi, id_gulih) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, ikas.IDStakeholder, ikas.Tanggal, ikas.Responden,
		ikas.Telepon, ikas.Jabatan, ikas.NilaiKematangan, ikas.TargetNilai,
		ikas.IDIdentifikasi, ikas.IDProteksi, ikas.IDDeteksi, ikas.IDGulih)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	ikas.ID = int(id)
	return nil
}

func (r *ikasRepository) GetAll() ([]models.Ikas, error) {
	rows, err := r.db.Query("SELECT * FROM ikas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ikasList []models.Ikas
	for rows.Next() {
		var ikas models.Ikas
		err := rows.Scan(&ikas.ID, &ikas.IDStakeholder, &ikas.Tanggal, &ikas.Responden,
			&ikas.Telepon, &ikas.Jabatan, &ikas.NilaiKematangan, &ikas.TargetNilai,
			&ikas.IDIdentifikasi, &ikas.IDProteksi, &ikas.IDDeteksi, &ikas.IDGulih)
		if err != nil {
			return nil, err
		}
		ikasList = append(ikasList, ikas)
	}
	return ikasList, nil
}

func (r *ikasRepository) GetByID(id int) (*models.Ikas, error) {
	var ikas models.Ikas
	query := "SELECT * FROM ikas WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&ikas.ID, &ikas.IDStakeholder, &ikas.Tanggal,
		&ikas.Responden, &ikas.Telepon, &ikas.Jabatan, &ikas.NilaiKematangan,
		&ikas.TargetNilai, &ikas.IDIdentifikasi, &ikas.IDProteksi, &ikas.IDDeteksi, &ikas.IDGulih)

	if err != nil {
		return nil, err
	}
	return &ikas, nil
}

func (r *ikasRepository) Update(id int, ikas *models.Ikas) error {
	query := `UPDATE ikas SET id_stakeholder=?, tanggal=?, responden=?, telepon=?, 
				jabatan=?, nilai_kematangan=?, target_nilai=?, id_identifikasi=?, 
				id_proteksi=?, id_deteksi=?, id_gulih=? WHERE id=?`

	_, err := r.db.Exec(query, ikas.IDStakeholder, ikas.Tanggal, ikas.Responden,
		ikas.Telepon, ikas.Jabatan, ikas.NilaiKematangan, ikas.TargetNilai,
		ikas.IDIdentifikasi, ikas.IDProteksi, ikas.IDDeteksi, ikas.IDGulih, id)

	return err
}

func (r *ikasRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM ikas WHERE id=?", id)
	return err
}