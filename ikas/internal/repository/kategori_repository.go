package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type KategoriRepository struct {
	db *sql.DB
}

func NewKategoriRepository(db *sql.DB) *KategoriRepository {
	return &KategoriRepository{db: db}
}

func (r *KategoriRepository) Create(req dto.CreateKategoriRequest) (int64, error) {
	query := `INSERT INTO kategori (domain_id, nama_kategori) VALUES (?, ?)`

	res, err := r.db.Exec(query, req.DomainID, req.NamaKategori)
	if err != nil {
		logger.Error(err, "operation failed")
		return 0, err
	}

	return res.LastInsertId()
}

func (r *KategoriRepository) GetAll() ([]dto.KategoriResponse, error) {
	query := `SELECT id, domain_id, nama_kategori, created_at, updated_at 
	          FROM kategori 
	          ORDER BY nama_kategori ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	defer rows.Close()

	var result []dto.KategoriResponse

	for rows.Next() {
		var item dto.KategoriResponse
		if err := rows.Scan(&item.ID, &item.DomainID, &item.NamaKategori, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logger.Error(err, "operation failed")
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *KategoriRepository) GetByID(id int) (*dto.KategoriResponse, error) {
	query := `SELECT id, domain_id, nama_kategori, created_at, updated_at 
	          FROM kategori 
	          WHERE id = ?`

	var item dto.KategoriResponse
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.DomainID, &item.NamaKategori, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	return &item, nil
}

func (r *KategoriRepository) Update(id int, req dto.UpdateKategoriRequest) error {
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
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}

func (r *KategoriRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM kategori WHERE id=?`, id)
	if err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}

func (r *KategoriRepository) CheckDuplicateName(domainID int, namaKategori string, excludeID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM kategori 
	          WHERE domain_id = ? 
	          AND LOWER(TRIM(nama_kategori)) = LOWER(TRIM(?))`
	args := []interface{}{domainID, namaKategori}

	if excludeID != 0 {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		logger.Error(err, "operation failed")
		return false, err
	}

	return count > 0, nil
}

func (r *KategoriRepository) CheckDomainExists(domainID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM domain WHERE id = ?`
	err := r.db.QueryRow(query, domainID).Scan(&count)
	if err != nil {
		logger.Error(err, "operation failed")
		return false, err
	}

	return count > 0, nil
}
