package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanProteksiRepositoryInterface interface {
	Create(req dto.CreatePertanyaanProteksiRequest) (int64, error)
	GetAll() ([]dto.PertanyaanProteksiResponse, error)
	GetByID(id int) (*dto.PertanyaanProteksiResponse, error)
	Update(id int, req dto.UpdatePertanyaanProteksiRequest) error
	Delete(id int) error
	CheckSubKategoriExists(subKategoriID int) (bool, error)
	CheckRuangLingkupExists(ruangLingkupID int) (bool, error)
}

type PertanyaanProteksiRepository struct {
	db *sql.DB
}

func NewPertanyaanProteksiRepository(db *sql.DB) *PertanyaanProteksiRepository {
	return &PertanyaanProteksiRepository{db: db}
}

func (r *PertanyaanProteksiRepository) Create(req dto.CreatePertanyaanProteksiRequest) (int64, error) {
	query := `INSERT INTO pertanyaan_proteksi
		(sub_kategori_id, ruang_lingkup_id, pertanyaan_proteksi, index0, index1, index2, index3, index4, index5)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		req.SubKategoriID,
		req.RuangLingkupID,
		req.PertanyaanProteksi,
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

func (r *PertanyaanProteksiRepository) GetAll() ([]dto.PertanyaanProteksiResponse, error) {
	query := `
		SELECT
			pp.id, pp.pertanyaan_proteksi,
			pp.index0, pp.index1, pp.index2, pp.index3, pp.index4, pp.index5,
			pp.created_at, pp.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_proteksi pp
		JOIN sub_kategori sk ON pp.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pp.ruang_lingkup_id = rl.id
		ORDER BY pp.created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.PertanyaanProteksiResponse
	for rows.Next() {
		var item dto.PertanyaanProteksiResponse
		if err := rows.Scan(
			&item.ID,
			&item.PertanyaanProteksi,
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

func (r *PertanyaanProteksiRepository) GetByID(id int) (*dto.PertanyaanProteksiResponse, error) {
	query := `
		SELECT
			pp.id, pp.pertanyaan_proteksi,
			pp.index0, pp.index1, pp.index2, pp.index3, pp.index4, pp.index5,
			pp.created_at, pp.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_proteksi pp
		JOIN sub_kategori sk ON pp.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pp.ruang_lingkup_id = rl.id
		WHERE pp.id = ?`

	var item dto.PertanyaanProteksiResponse
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.PertanyaanProteksi,
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

func (r *PertanyaanProteksiRepository) Update(id int, req dto.UpdatePertanyaanProteksiRequest) error {
	query := "UPDATE pertanyaan_proteksi SET "
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
	if req.PertanyaanProteksi != nil {
		updates = append(updates, "pertanyaan_proteksi=?")
		args = append(args, *req.PertanyaanProteksi)
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

func (r *PertanyaanProteksiRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM pertanyaan_proteksi WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *PertanyaanProteksiRepository) CheckSubKategoriExists(subKategoriID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM sub_kategori WHERE id = ?`, subKategoriID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *PertanyaanProteksiRepository) CheckRuangLingkupExists(ruangLingkupID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ruang_lingkup WHERE id = ?`, ruangLingkupID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}
