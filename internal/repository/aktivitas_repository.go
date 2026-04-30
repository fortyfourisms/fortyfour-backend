package repository

import (
	"database/sql"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"strings"

	"fortyfour-backend/pkg/logger"
)

type AktivitasRepository struct {
	db *sql.DB
}

func NewAktivitasRepository(db *sql.DB) *AktivitasRepository {
	return &AktivitasRepository{db: db}
}

func formatDate(dateStr string) string {
	if strings.Contains(dateStr, "T") {
		return strings.Split(dateStr, "T")[0]
	}
	return dateStr
}

func (r *AktivitasRepository) Create(req dto.CreateAktivitasRequest) (int64, error) {
	jenisJSON, err := json.Marshal(req.JenisAktivitas)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO aktivitas (perusahaan_id, judul, deskripsi, tanggal_mulai, tanggal_selesai, jenis_aktivitas) VALUES (?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query, req.PerusahaanID, req.Judul, req.Deskripsi, formatDate(req.TanggalMulai), formatDate(req.TanggalSelesai), string(jenisJSON))
	if err != nil {
		logger.Error(err, "operation failed")
		return 0, err
	}

	return res.LastInsertId()
}

func (r *AktivitasRepository) GetAll() ([]dto.AktivitasResponse, error) {
	query := `SELECT id, perusahaan_id, judul, deskripsi, tanggal_mulai, tanggal_selesai, jenis_aktivitas, created_at, updated_at FROM aktivitas ORDER BY tanggal_mulai DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	defer rows.Close()

	var result []dto.AktivitasResponse

	for rows.Next() {
		var item dto.AktivitasResponse
		var jenisJSON string
		if err := rows.Scan(&item.ID, &item.PerusahaanID, &item.Judul, &item.Deskripsi, &item.TanggalMulai, &item.TanggalSelesai, &jenisJSON, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logger.Error(err, "operation failed")
			continue
		}
		_ = json.Unmarshal([]byte(jenisJSON), &item.JenisAktivitas)
		result = append(result, item)
	}

	return result, nil
}

func (r *AktivitasRepository) GetByID(id int) (*dto.AktivitasResponse, error) {
	query := `SELECT id, perusahaan_id, judul, deskripsi, tanggal_mulai, tanggal_selesai, jenis_aktivitas, created_at, updated_at FROM aktivitas WHERE id = ?`

	var item dto.AktivitasResponse
	var jenisJSON string
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.PerusahaanID, &item.Judul, &item.Deskripsi, &item.TanggalMulai, &item.TanggalSelesai, &jenisJSON, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}

	_ = json.Unmarshal([]byte(jenisJSON), &item.JenisAktivitas)

	return &item, nil
}

func (r *AktivitasRepository) GetByPerusahaanID(perusahaanID string) ([]dto.AktivitasResponse, error) {
	query := `SELECT id, perusahaan_id, judul, deskripsi, tanggal_mulai, tanggal_selesai, jenis_aktivitas, created_at, updated_at FROM aktivitas WHERE perusahaan_id = ? ORDER BY tanggal_mulai DESC`

	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		logger.Error(err, "operation failed")
		return nil, err
	}
	defer rows.Close()

	var result []dto.AktivitasResponse

	for rows.Next() {
		var item dto.AktivitasResponse
		var jenisJSON string
		if err := rows.Scan(&item.ID, &item.PerusahaanID, &item.Judul, &item.Deskripsi, &item.TanggalMulai, &item.TanggalSelesai, &jenisJSON, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logger.Error(err, "operation failed")
			continue
		}
		_ = json.Unmarshal([]byte(jenisJSON), &item.JenisAktivitas)
		result = append(result, item)
	}

	return result, nil
}

func (r *AktivitasRepository) Update(id int, req dto.UpdateAktivitasRequest) error {
	query := "UPDATE aktivitas SET "
	args := []interface{}{}
	updates := []string{}

	if req.PerusahaanID != nil {
		updates = append(updates, "perusahaan_id=?")
		args = append(args, *req.PerusahaanID)
	}
	if req.Judul != nil {
		updates = append(updates, "judul=?")
		args = append(args, *req.Judul)
	}
	if req.Deskripsi != nil {
		updates = append(updates, "deskripsi=?")
		args = append(args, *req.Deskripsi)
	}
	if req.TanggalMulai != nil {
		updates = append(updates, "tanggal_mulai=?")
		args = append(args, formatDate(*req.TanggalMulai))
	}
	if req.TanggalSelesai != nil {
		updates = append(updates, "tanggal_selesai=?")
		args = append(args, formatDate(*req.TanggalSelesai))
	}
	if req.JenisAktivitas != nil {
		jenisJSON, err := json.Marshal(*req.JenisAktivitas)
		if err != nil {
			return err
		}
		updates = append(updates, "jenis_aktivitas=?")
		args = append(args, string(jenisJSON))
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

func (r *AktivitasRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM aktivitas WHERE id=?`, id)
	if err != nil {
		logger.Error(err, "operation failed")
		return err
	}

	return nil
}
