package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type SubKategoriRepository struct {
	db *sql.DB
}

func NewSubKategoriRepository(db *sql.DB) *SubKategoriRepository {
	return &SubKategoriRepository{db: db}
}

func (r *SubKategoriRepository) Create(req dto.CreateSubKategoriRequest, id string) error {
	query := `INSERT INTO sub_kategori (id, kategori_id, nama_sub_kategori) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, id, req.KategoriID, req.NamaSubKategori)
	if err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}

func (r *SubKategoriRepository) GetAll() ([]dto.SubKategoriResponse, error) {
	query := `SELECT id, kategori_id, nama_sub_kategori, created_at, updated_at 
	          FROM sub_kategori 
	          ORDER BY nama_sub_kategori ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	defer rows.Close()

	var result []dto.SubKategoriResponse

	for rows.Next() {
		var item dto.SubKategoriResponse
		if err := rows.Scan(&item.ID, &item.KategoriID, &item.NamaSubKategori, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logger.Error(err, "operation failed")
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *SubKategoriRepository) GetByID(id string) (*dto.SubKategoriResponse, error) {
	query := `SELECT id, kategori_id, nama_sub_kategori, created_at, updated_at 
	          FROM sub_kategori 
	          WHERE id = ?`

	var item dto.SubKategoriResponse
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.KategoriID, &item.NamaSubKategori, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	return &item, nil
}

func (r *SubKategoriRepository) Update(id string, req dto.UpdateSubKategoriRequest) error {
	query := "UPDATE sub_kategori SET "
	args := []interface{}{}
	updates := []string{}

	if req.KategoriID != nil {
		updates = append(updates, "kategori_id=?")
		args = append(args, *req.KategoriID)
	}

	if req.NamaSubKategori != nil {
		updates = append(updates, "nama_sub_kategori=?")
		args = append(args, *req.NamaSubKategori)
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

func (r *SubKategoriRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM sub_kategori WHERE id=?`, id)
	if err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}

func (r *SubKategoriRepository) CheckDuplicateName(kategoriID string, namaSubKategori string, excludeID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM sub_kategori 
	          WHERE kategori_id = ? 
	          AND LOWER(TRIM(nama_sub_kategori)) = LOWER(TRIM(?))`
	args := []interface{}{kategoriID, namaSubKategori}

	if excludeID != "" {
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

func (r *SubKategoriRepository) CheckKategoriExists(kategoriID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM kategori WHERE id = ?`
	err := r.db.QueryRow(query, kategoriID).Scan(&count)
	if err != nil {
		logger.Error(err, "operation failed")
		return false, err
	}

	return count > 0, nil
}
