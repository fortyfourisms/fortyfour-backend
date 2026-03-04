package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type PertanyaanIdentifikasiRepositoryInterface interface {
	Create(req dto.CreatePertanyaanIdentifikasiRequest) (int64, error)
	GetAll() ([]dto.PertanyaanIdentifikasiResponse, error)
	GetByID(id int) (*dto.PertanyaanIdentifikasiResponse, error)
	Update(id int, req dto.UpdatePertanyaanIdentifikasiRequest) error
	Delete(id int) error
	CheckSubKategoriExists(subKategoriID int) (bool, error)
	CheckRuangLingkupExists(ruangLingkupID int) (bool, error)
}

type PertanyaanIdentifikasiRepository struct {
	db *sql.DB
}

func NewPertanyaanIdentifikasiRepository(db *sql.DB) *PertanyaanIdentifikasiRepository {
	return &PertanyaanIdentifikasiRepository{db: db}
}

func (r *PertanyaanIdentifikasiRepository) Create(req dto.CreatePertanyaanIdentifikasiRequest) (int64, error) {
	query := `INSERT INTO pertanyaan_identifikasi (
		sub_kategori_id, ruang_lingkup_id, pertanyaan_identifikasi, index0, index1, index2, index3, index4, index5
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		req.SubKategoriID, req.RuangLingkupID, req.PertanyaanIdentifikasi,
		req.Index0, req.Index1, req.Index2, req.Index3, req.Index4, req.Index5)
	if err != nil {
		rollbar.Error(err)
		return 0, err
	}

	return res.LastInsertId()
}

func (r *PertanyaanIdentifikasiRepository) GetAll() ([]dto.PertanyaanIdentifikasiResponse, error) {
	query := `
		SELECT 
			pi.id, pi.pertanyaan_identifikasi,
			pi.index0, pi.index1, pi.index2, pi.index3, pi.index4, pi.index5,
			pi.created_at, pi.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_identifikasi pi
		JOIN sub_kategori sk ON pi.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pi.ruang_lingkup_id = rl.id
		ORDER BY pi.created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.PertanyaanIdentifikasiResponse

	for rows.Next() {
		var item dto.PertanyaanIdentifikasiResponse
		if err := rows.Scan(
			&item.ID,
			&item.PertanyaanIdentifikasi,
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

func (r *PertanyaanIdentifikasiRepository) GetByID(id int) (*dto.PertanyaanIdentifikasiResponse, error) {
	query := `
		SELECT 
			pi.id, pi.pertanyaan_identifikasi,
			pi.index0, pi.index1, pi.index2, pi.index3, pi.index4, pi.index5,
			pi.created_at, pi.updated_at,
			sk.id, sk.nama_sub_kategori,
			k.id, k.nama_kategori,
			d.id, d.nama_domain,
			rl.id, rl.nama_ruang_lingkup
		FROM pertanyaan_identifikasi pi
		JOIN sub_kategori sk ON pi.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		JOIN domain d ON k.domain_id = d.id
		JOIN ruang_lingkup rl ON pi.ruang_lingkup_id = rl.id
		WHERE pi.id = ?`

	var item dto.PertanyaanIdentifikasiResponse
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.PertanyaanIdentifikasi,
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

func (r *PertanyaanIdentifikasiRepository) Update(id int, req dto.UpdatePertanyaanIdentifikasiRequest) error {
	query := "UPDATE pertanyaan_identifikasi SET "
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

	if req.PertanyaanIdentifikasi != nil {
		updates = append(updates, "pertanyaan_identifikasi=?")
		args = append(args, *req.PertanyaanIdentifikasi)
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

func (r *PertanyaanIdentifikasiRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM pertanyaan_identifikasi WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *PertanyaanIdentifikasiRepository) CheckSubKategoriExists(subKategoriID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM sub_kategori WHERE id = ?`, subKategoriID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *PertanyaanIdentifikasiRepository) CheckRuangLingkupExists(ruangLingkupID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ruang_lingkup WHERE id = ?`, ruangLingkupID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}
