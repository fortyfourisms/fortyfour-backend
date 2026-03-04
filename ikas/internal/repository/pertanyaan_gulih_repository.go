package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanGulihRepositoryInterface interface {
	Create(req dto.CreatePertanyaanGulihRequest) (int64, error)
	GetAll() ([]dto.PertanyaanGulihResponse, error)
	GetByID(id int) (*dto.PertanyaanGulihResponse, error)
	Update(id int, req dto.UpdatePertanyaanGulihRequest) error
	Delete(id int) error
	CheckSubKategoriExists(subKategoriID int) (bool, error)
	CheckRuangLingkupExists(ruangLingkupID int) (bool, error)
}

type PertanyaanGulihRepository struct {
	db *sql.DB
}

func NewPertanyaanGulihRepository(db *sql.DB) *PertanyaanGulihRepository {
	return &PertanyaanGulihRepository{db: db}
}

func (r *PertanyaanGulihRepository) Create(req dto.CreatePertanyaanGulihRequest) (int64, error) {
	query := `INSERT INTO pertanyaan_gulih
		(sub_kategori_id, ruang_lingkup_id, pertanyaan_gulih, index0, index1, index2, index3, index4, index5)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		req.SubKategoriID,
		req.RuangLingkupID,
		req.PertanyaanGulih,
		req.Index0,
		req.Index1,
		req.Index2,
		req.Index3,
		req.Index4,
		req.Index5,
	)
	if err != nil {
		rollbar.Error(err)
		return 0, err
	}

	return res.LastInsertId()
}

func (r *PertanyaanGulihRepository) GetAll() ([]dto.PertanyaanGulihResponse, error) {
	query := `
		SELECT
			pg.id, pg.pertanyaan_gulih,
			pg.index0, pg.index1, pg.index2, pg.index3, pg.index4, pg.index5,
			pg.created_at, pg.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_gulih pg
		JOIN sub_kategori sk ON pg.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pg.ruang_lingkup_id = rl.id
		ORDER BY pg.created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.PertanyaanGulihResponse
	for rows.Next() {
		var item dto.PertanyaanGulihResponse
		if err := rows.Scan(
			&item.ID,
			&item.PertanyaanGulih,
			&item.Index0,
			&item.Index1,
			&item.Index2,
			&item.Index3,
			&item.Index4,
			&item.Index5,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.SubKategori.ID,
			&item.SubKategori.NamaSubKategori,
			&item.SubKategori.Kategori.ID,
			&item.SubKategori.Kategori.NamaKategori,
			&item.SubKategori.Kategori.Domain.ID,
			&item.SubKategori.Kategori.Domain.NamaDomain,
			&item.RuangLingkup.ID,
			&item.RuangLingkup.NamaRuangLingkup,
		); err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *PertanyaanGulihRepository) GetByID(id int) (*dto.PertanyaanGulihResponse, error) {
	query := `
		SELECT
			pg.id, pg.pertanyaan_gulih,
			pg.index0, pg.index1, pg.index2, pg.index3, pg.index4, pg.index5,
			pg.created_at, pg.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_gulih pg
		JOIN sub_kategori sk ON pg.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pg.ruang_lingkup_id = rl.id
		WHERE pg.id = ?`

	var item dto.PertanyaanGulihResponse
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.PertanyaanGulih,
		&item.Index0,
		&item.Index1,
		&item.Index2,
		&item.Index3,
		&item.Index4,
		&item.Index5,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.SubKategori.ID,
		&item.SubKategori.NamaSubKategori,
		&item.SubKategori.Kategori.ID,
		&item.SubKategori.Kategori.NamaKategori,
		&item.SubKategori.Kategori.Domain.ID,
		&item.SubKategori.Kategori.Domain.NamaDomain,
		&item.RuangLingkup.ID,
		&item.RuangLingkup.NamaRuangLingkup,
	)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return &item, nil
}

func (r *PertanyaanGulihRepository) Update(id int, req dto.UpdatePertanyaanGulihRequest) error {
	query := "UPDATE pertanyaan_gulih SET "
	args := []interface{}{}
	updates := []string{}

	if req.SubKategoriID != nil {
		updates = append(updates, "sub_kategori_id=?")
		args = append(args, *req.SubKategoriID)
	}
	if req.RuangLingkupID != nil {
		updates = append(updates, "ruang_lingkup_id=?")
		args = append(args, *req.RuangLingkupID)
	}
	if req.PertanyaanGulih != nil {
		updates = append(updates, "pertanyaan_gulih=?")
		args = append(args, *req.PertanyaanGulih)
	}
	if req.Index0 != nil {
		updates = append(updates, "index0=?")
		args = append(args, *req.Index0)
	}
	if req.Index1 != nil {
		updates = append(updates, "index1=?")
		args = append(args, *req.Index1)
	}
	if req.Index2 != nil {
		updates = append(updates, "index2=?")
		args = append(args, *req.Index2)
	}
	if req.Index3 != nil {
		updates = append(updates, "index3=?")
		args = append(args, *req.Index3)
	}
	if req.Index4 != nil {
		updates = append(updates, "index4=?")
		args = append(args, *req.Index4)
	}
	if req.Index5 != nil {
		updates = append(updates, "index5=?")
		args = append(args, *req.Index5)
	}

	if len(updates) == 0 {
		return nil
	}

	query += strings.Join(updates, ", ")
	query += " WHERE id=?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *PertanyaanGulihRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM pertanyaan_gulih WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *PertanyaanGulihRepository) CheckSubKategoriExists(subKategoriID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM sub_kategori WHERE id = ?`, subKategoriID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *PertanyaanGulihRepository) CheckRuangLingkupExists(ruangLingkupID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ruang_lingkup WHERE id = ?`, ruangLingkupID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}
