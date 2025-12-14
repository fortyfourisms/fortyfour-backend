package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

type IkasRepository struct {
	db *sql.DB
}

func NewIkasRepository(db *sql.DB) *IkasRepository {
	return &IkasRepository{db: db}
}

func (r *IkasRepository) Create(req dto.CreateIkasRequest, id string) error {
	query := `INSERT INTO ikas
				(id, id_perusahaan, tanggal, responden, telepon, jabatan,
				nilai_kematangan, target_nilai, id_identifikasi, id_proteksi,
				id_deteksi, id_gulih)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		req.IDPerusahaan,
		req.Tanggal,
		req.Responden,
		req.Telepon,
		req.Jabatan,
		req.NilaiKematangan,
		req.TargetNilai,
		req.IDIdentifikasi,
		req.IDProteksi,
		req.IDDeteksi,
		req.IDGulih,
	)

	return err
}

func (r *IkasRepository) GetAll() ([]models.Ikas, error) {
	rows, err := r.db.Query(`
		SELECT id, id_perusahaan, tanggal, responden, telepon, jabatan,
       	nilai_kematangan, target_nilai, id_identifikasi, id_proteksi,
       	id_deteksi, id_gulih
		FROM ikas`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Ikas

	for rows.Next() {
		var i models.Ikas
		err := rows.Scan(
			&i.ID,
			&i.IDPerusahaan,
			&i.Tanggal,
			&i.Responden,
			&i.Telepon,
			&i.Jabatan,
			&i.NilaiKematangan,
			&i.TargetNilai,
			&i.IDIdentifikasi,
			&i.IDProteksi,
			&i.IDDeteksi,
			&i.IDGulih,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}

	return result, nil
}

func (r *IkasRepository) GetByID(id string) (*models.Ikas, error) {
	rows := r.db.QueryRow(`
		SELECT id, id_perusahaan, tanggal, responden, telepon, jabatan,
      	nilai_kematangan, target_nilai, id_identifikasi, id_proteksi,
       	id_deteksi, id_gulih
		FROM ikas
		WHERE id = ?`, id)

	var i models.Ikas
	err := rows.Scan(
		&i.ID,
		&i.IDPerusahaan,
		&i.Tanggal,
		&i.Responden,
		&i.Telepon,
		&i.Jabatan,
		&i.NilaiKematangan,
		&i.TargetNilai,
		&i.IDIdentifikasi,
		&i.IDProteksi,
		&i.IDDeteksi,
		&i.IDGulih,
	)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *IkasRepository) Update(id string, i models.Ikas) error {
	query := `
		UPDATE ikas SET
			id_perusahaan=?,
			tanggal=?,
			responden=?,
			telepon=?,
			jabatan=?,
			nilai_kematangan=?,
			target_nilai=?,
			id_identifikasi=?,
			id_proteksi=?,
			id_deteksi=?,
			id_gulih=?
		WHERE id=?`

	_, err := r.db.Exec(query,
		i.IDPerusahaan,
		i.Tanggal,
		i.Responden,
		i.Telepon,
		i.Jabatan,
		i.NilaiKematangan,
		i.TargetNilai,
		i.IDIdentifikasi,
		i.IDProteksi,
		i.IDDeteksi,
		i.IDGulih,
		i.ID,
	)

	return err
}

func (r *IkasRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM ikas WHERE id=?`, id)
	return err
}