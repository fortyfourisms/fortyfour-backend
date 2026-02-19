package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type KategoriRepository struct {
	db *sql.DB
}

func NewKategoriRepository(db *sql.DB) *KategoriRepository {
	return &KategoriRepository{db: db}
}

func (r *KategoriRepository) Create(req dto.CreateKategoriRequest, id string) error {
	query := `INSERT INTO kategori (id, domain_id, nama_kategori) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, id, req.DomainID, req.NamaKategori)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *KategoriRepository) GetAll() ([]dto.KategoriResponse, error) {
	query := `SELECT id, domain_id, nama_kategori, created_at, updated_at 
	          FROM kategori 
	          ORDER BY nama_kategori ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.KategoriResponse

	for rows.Next() {
		var item dto.KategoriResponse
		if err := rows.Scan(&item.ID, &item.DomainID, &item.NamaKategori, &item.CreatedAt, &item.UpdatedAt); err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *KategoriRepository) GetByID(id string) (*dto.KategoriResponse, error) {
	query := `SELECT id, domain_id, nama_kategori, created_at, updated_at 
	          FROM kategori 
	          WHERE id = ?`

	var item dto.KategoriResponse
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.DomainID, &item.NamaKategori, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return &item, nil
}

func (r *KategoriRepository) Update(id string, req dto.UpdateKategoriRequest) error {
	query := "UPDATE kategori SET "
	args := []interface{}{}
	updates := []string{}

	if req.DomainID != nil {
		updates = append(updates, "domain_id=?")
		args = append(args, *req.DomainID)
	}

	if req.NamaKategori != nil {
		updates = append(updates, "nama_kategori=?")
		args = append(args, *req.NamaKategori)
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

func (r *KategoriRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM kategori WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *KategoriRepository) CheckDuplicateName(domainID string, namaKategori string, excludeID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM kategori 
	          WHERE domain_id = ? 
	          AND LOWER(TRIM(nama_kategori)) = LOWER(TRIM(?))`
	args := []interface{}{domainID, namaKategori}

	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}

	return count > 0, nil
}

func (r *KategoriRepository) CheckDomainExists(domainID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM domain WHERE id = ?`
	err := r.db.QueryRow(query, domainID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}

	return count > 0, nil
}
