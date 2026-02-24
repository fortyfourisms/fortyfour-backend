package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type JawabanIdentifikasiRepositoryInterface interface {
	Create(req dto.CreateJawabanIdentifikasiRequest, id string) error
	GetAll() ([]dto.JawabanIdentifikasiResponse, error)
	GetByID(id string) (*dto.JawabanIdentifikasiResponse, error)
	GetByPerusahaan(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error)
	GetByPertanyaan(pertanyaanID string) ([]dto.JawabanIdentifikasiResponse, error)
	Update(id string, req dto.UpdateJawabanIdentifikasiRequest) error
	Delete(id string) error
	CheckPertanyaanExists(pertanyaanID string) (bool, error)
	CheckPerusahaanExists(perusahaanID string) (bool, error)
	CheckDuplicate(perusahaanID, pertanyaanID, excludeID string) (bool, error)
}

type JawabanIdentifikasiRepository struct {
	db *sql.DB
}

func NewJawabanIdentifikasiRepository(db *sql.DB) *JawabanIdentifikasiRepository {
	return &JawabanIdentifikasiRepository{db: db}
}

func (r *JawabanIdentifikasiRepository) Create(req dto.CreateJawabanIdentifikasiRequest, id string) error {
	query := `INSERT INTO jawaban_identifikasi
		(id, pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		req.PertanyaanIdentifikasiID,
		req.PerusahaanID,
		req.JawabanIdentifikasi,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *JawabanIdentifikasiRepository) GetAll() ([]dto.JawabanIdentifikasiResponse, error) {
	query := `SELECT id, pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi,
		evidence, validasi, keterangan, created_at, updated_at
		FROM jawaban_identifikasi
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanIdentifikasiResponse
	for rows.Next() {
		var item dto.JawabanIdentifikasiResponse
		if err := rows.Scan(
			&item.ID,
			&item.PertanyaanIdentifikasiID,
			&item.PerusahaanID,
			&item.JawabanIdentifikasi,
			&item.Evidence,
			&item.Validasi,
			&item.Keterangan,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanIdentifikasiRepository) GetByID(id string) (*dto.JawabanIdentifikasiResponse, error) {
	query := `SELECT id, pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi,
		evidence, validasi, keterangan, created_at, updated_at
		FROM jawaban_identifikasi
		WHERE id = ?`

	var item dto.JawabanIdentifikasiResponse
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.PertanyaanIdentifikasiID,
		&item.PerusahaanID,
		&item.JawabanIdentifikasi,
		&item.Evidence,
		&item.Validasi,
		&item.Keterangan,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return &item, nil
}

func (r *JawabanIdentifikasiRepository) GetByPerusahaan(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	query := `SELECT id, pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi,
		evidence, validasi, keterangan, created_at, updated_at
		FROM jawaban_identifikasi
		WHERE perusahaan_id = ?
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanIdentifikasiResponse
	for rows.Next() {
		var item dto.JawabanIdentifikasiResponse
		if err := rows.Scan(
			&item.ID,
			&item.PertanyaanIdentifikasiID,
			&item.PerusahaanID,
			&item.JawabanIdentifikasi,
			&item.Evidence,
			&item.Validasi,
			&item.Keterangan,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanIdentifikasiRepository) GetByPertanyaan(pertanyaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	query := `SELECT id, pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi,
		evidence, validasi, keterangan, created_at, updated_at
		FROM jawaban_identifikasi
		WHERE pertanyaan_identifikasi_id = ?
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, pertanyaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanIdentifikasiResponse
	for rows.Next() {
		var item dto.JawabanIdentifikasiResponse
		if err := rows.Scan(
			&item.ID,
			&item.PertanyaanIdentifikasiID,
			&item.PerusahaanID,
			&item.JawabanIdentifikasi,
			&item.Evidence,
			&item.Validasi,
			&item.Keterangan,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanIdentifikasiRepository) Update(id string, req dto.UpdateJawabanIdentifikasiRequest) error {
	query := "UPDATE jawaban_identifikasi SET "
	args := []interface{}{}
	updates := []string{}

	if req.JawabanIdentifikasi != nil {
		updates = append(updates, "jawaban_identifikasi=?")
		args = append(args, *req.JawabanIdentifikasi)
	}

	if req.Evidence != nil {
		updates = append(updates, "evidence=?")
		args = append(args, *req.Evidence)
	}

	if req.Validasi != nil {
		updates = append(updates, "validasi=?")
		args = append(args, *req.Validasi)
	}

	if req.Keterangan != nil {
		updates = append(updates, "keterangan=?")
		args = append(args, *req.Keterangan)
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

func (r *JawabanIdentifikasiRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM jawaban_identifikasi WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *JawabanIdentifikasiRepository) CheckPertanyaanExists(pertanyaanID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM pertanyaan_identifikasi WHERE id = ?`, pertanyaanID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanIdentifikasiRepository) CheckPerusahaanExists(perusahaanID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM perusahaan WHERE id = ?`, perusahaanID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanIdentifikasiRepository) CheckDuplicate(perusahaanID, pertanyaanID, excludeID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM jawaban_identifikasi
		WHERE perusahaan_id = ? AND pertanyaan_identifikasi_id = ?`
	args := []interface{}{perusahaanID, pertanyaanID}

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
