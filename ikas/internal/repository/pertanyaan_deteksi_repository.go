package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanDeteksiRepositoryInterface interface {
	Create(req dto.CreatePertanyaanDeteksiRequest, id string) error
	GetAll() ([]dto.PertanyaanDeteksiResponse, error)
	GetByID(id string) (*dto.PertanyaanDeteksiResponse, error)
	Update(id string, req dto.UpdatePertanyaanDeteksiRequest) error
	Delete(id string) error
	CheckSubKategoriExists(subKategoriID string) (bool, error)
	CheckRuangLingkupExists(ruangLingkupID string) (bool, error)
}

type PertanyaanDeteksiRepository struct {
	db *sql.DB
}

func NewPertanyaanDeteksiRepository(db *sql.DB) *PertanyaanDeteksiRepository {
	return &PertanyaanDeteksiRepository{db: db}
}

func (r *PertanyaanDeteksiRepository) Create(req dto.CreatePertanyaanDeteksiRequest, id string) error {
	query := `INSERT INTO pertanyaan_deteksi
		(id, sub_kategori_id, ruang_lingkup_id, pertanyaan_deteksi, index0, index1, index2, index3, index4, index5)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		req.SubKategoriID,
		req.RuangLingkupID,
		req.PertanyaanDeteksi,
		req.Index0,
		req.Index1,
		req.Index2,
		req.Index3,
		req.Index4,
		req.Index5,
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *PertanyaanDeteksiRepository) GetAll() ([]dto.PertanyaanDeteksiResponse, error) {
	query := `
		SELECT
			pd.id, pd.pertanyaan_deteksi,
			pd.index0, pd.index1, pd.index2, pd.index3, pd.index4, pd.index5,
			pd.created_at, pd.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_deteksi pd
		JOIN sub_kategori sk ON pd.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pd.ruang_lingkup_id = rl.id
		ORDER BY pd.created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.PertanyaanDeteksiResponse
	for rows.Next() {
		var item dto.PertanyaanDeteksiResponse
		if err := rows.Scan(
			&item.ID,
			&item.PertanyaanDeteksi,
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

func (r *PertanyaanDeteksiRepository) GetByID(id string) (*dto.PertanyaanDeteksiResponse, error) {
	query := `
		SELECT
			pd.id, pd.pertanyaan_deteksi,
			pd.index0, pd.index1, pd.index2, pd.index3, pd.index4, pd.index5,
			pd.created_at, pd.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_deteksi pd
		JOIN sub_kategori sk ON pd.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pd.ruang_lingkup_id = rl.id
		WHERE pd.id = ?`

	var item dto.PertanyaanDeteksiResponse
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.PertanyaanDeteksi,
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

func (r *PertanyaanDeteksiRepository) Update(id string, req dto.UpdatePertanyaanDeteksiRequest) error {
	query := "UPDATE pertanyaan_deteksi SET "
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
	if req.PertanyaanDeteksi != nil {
		updates = append(updates, "pertanyaan_deteksi=?")
		args = append(args, *req.PertanyaanDeteksi)
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

func (r *PertanyaanDeteksiRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM pertanyaan_deteksi WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *PertanyaanDeteksiRepository) CheckSubKategoriExists(subKategoriID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM sub_kategori WHERE id = ?`, subKategoriID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *PertanyaanDeteksiRepository) CheckRuangLingkupExists(ruangLingkupID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ruang_lingkup WHERE id = ?`, ruangLingkupID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}
