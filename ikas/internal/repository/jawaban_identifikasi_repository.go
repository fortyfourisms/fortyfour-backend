package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"ikas/internal/utils"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type JawabanIdentifikasiRepositoryInterface interface {
	Create(req dto.CreateJawabanIdentifikasiRequest) (int64, error)
	GetAll() ([]dto.JawabanIdentifikasiResponse, error)
	GetByID(id int) (*dto.JawabanIdentifikasiResponse, error)
	GetByPerusahaan(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error)
	GetByPertanyaan(pertanyaanID int) ([]dto.JawabanIdentifikasiResponse, error)
	Update(id int, req dto.UpdateJawabanIdentifikasiRequest) error
	Delete(id int) error
	CheckPertanyaanExists(pertanyaanID int) (bool, error)
	CheckPerusahaanExists(perusahaanID string) (bool, error)
	CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error)
	RecalculateIdentifikasi(perusahaanID string) error
	UpsertToBuffer(req dto.CreateJawabanIdentifikasiRequest) error
	GetBufferCount(perusahaanID string) (int, error)
	FlushBuffer(perusahaanID string) error
}

type JawabanIdentifikasiRepository struct {
	db *sql.DB
}

func NewJawabanIdentifikasiRepository(db *sql.DB) *JawabanIdentifikasiRepository {
	return &JawabanIdentifikasiRepository{db: db}
}

const jawabanIdentifikasiSelectQuery = `
	SELECT
		ji.id,
		ji.perusahaan_id,
		ji.jawaban_identifikasi,
		ji.evidence,
		ji.validasi,
		ji.keterangan,
		ji.created_at,
		ji.updated_at,
		pi.id,
		pi.pertanyaan_identifikasi,
		sk.id, sk.nama_sub_kategori,
		k.id, k.nama_kategori,
		d.id, d.nama_domain
	FROM jawaban_identifikasi ji
	JOIN pertanyaan_identifikasi pi ON ji.pertanyaan_identifikasi_id = pi.id
	JOIN sub_kategori sk ON pi.sub_kategori_id = sk.id
	JOIN kategori k ON sk.kategori_id = k.id
	JOIN domain d ON k.domain_id = d.id`

func scanJawaban(row interface {
	Scan(dest ...any) error
}) (dto.JawabanIdentifikasiResponse, error) {
	var item dto.JawabanIdentifikasiResponse
	err := row.Scan(
		&item.ID,
		&item.PerusahaanID,
		&item.JawabanIdentifikasi,
		&item.Evidence,
		&item.Validasi,
		&item.Keterangan,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.PertanyaanIdentifikasi.ID,
		&item.PertanyaanIdentifikasi.PertanyaanIdentifikasi,
		&item.PertanyaanIdentifikasi.SubKategori.ID,
		&item.PertanyaanIdentifikasi.SubKategori.NamaSubKategori,
		&item.PertanyaanIdentifikasi.SubKategori.Kategori.ID,
		&item.PertanyaanIdentifikasi.SubKategori.Kategori.NamaKategori,
		&item.PertanyaanIdentifikasi.SubKategori.Kategori.Domain.ID,
		&item.PertanyaanIdentifikasi.SubKategori.Kategori.Domain.NamaDomain,
	)
	return item, err
}

func (r *JawabanIdentifikasiRepository) Create(req dto.CreateJawabanIdentifikasiRequest) (int64, error) {
	query := `INSERT INTO jawaban_identifikasi
		(pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		req.PertanyaanIdentifikasiID,
		req.PerusahaanID,
		req.JawabanIdentifikasi,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	if err != nil {
		rollbar.Error(err)
		return 0, err
	}
	return res.LastInsertId()
}

func (r *JawabanIdentifikasiRepository) GetAll() ([]dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + ` ORDER BY ji.created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanIdentifikasiResponse
	for rows.Next() {
		item, err := scanJawaban(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanIdentifikasiRepository) GetByID(id int) (*dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + ` WHERE ji.id = ?`

	item, err := scanJawaban(r.db.QueryRow(query, id))
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return &item, nil
}

func (r *JawabanIdentifikasiRepository) GetByPerusahaan(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + ` WHERE ji.perusahaan_id = ? ORDER BY ji.created_at ASC`

	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanIdentifikasiResponse
	for rows.Next() {
		item, err := scanJawaban(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanIdentifikasiRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + ` WHERE ji.pertanyaan_identifikasi_id = ? ORDER BY ji.created_at ASC`

	rows, err := r.db.Query(query, pertanyaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanIdentifikasiResponse
	for rows.Next() {
		item, err := scanJawaban(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanIdentifikasiRepository) Update(id int, req dto.UpdateJawabanIdentifikasiRequest) error {
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

func (r *JawabanIdentifikasiRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM jawaban_identifikasi WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *JawabanIdentifikasiRepository) CheckPertanyaanExists(pertanyaanID int) (bool, error) {
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

func (r *JawabanIdentifikasiRepository) CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM jawaban_identifikasi
		WHERE perusahaan_id = ? AND pertanyaan_identifikasi_id = ?`
	args := []interface{}{perusahaanID, pertanyaanID}

	if excludeID != 0 {
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

func (r *JawabanIdentifikasiRepository) RecalculateIdentifikasi(perusahaanID string) error {
	// Query rata-rata jawaban per kategori_id untuk perusahaan tertentu
	query := `
		SELECT k.id AS kategori_id, ROUND(AVG(ji.jawaban_identifikasi), 2) AS avg_nilai
		FROM jawaban_identifikasi ji
		JOIN pertanyaan_identifikasi pi ON ji.pertanyaan_identifikasi_id = pi.id
		JOIN sub_kategori sk ON pi.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		WHERE ji.perusahaan_id = ? AND ji.jawaban_identifikasi IS NOT NULL
		GROUP BY k.id
		ORDER BY k.id`

	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	defer rows.Close()

	// Map kategori_id ke nilai subdomain (default 0)
	subdomain := map[int]float64{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}

	for rows.Next() {
		var kategoriID int
		var avgNilai float64
		if err := rows.Scan(&kategoriID, &avgNilai); err != nil {
			rollbar.Error(err)
			continue
		}
		if kategoriID >= 1 && kategoriID <= 5 {
			subdomain[kategoriID] = avgNilai
		}
	}

	// Hitung rata-rata keseluruhan (nilai_identifikasi)
	nilaiIdentifikasi := utils.RoundToTwo((subdomain[1] + subdomain[2] + subdomain[3] + subdomain[4] + subdomain[5]) / 5.0)

	// Upsert ke tabel identifikasi
	upsertQuery := `
		INSERT INTO identifikasi 
			(perusahaan_id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, 
			 nilai_subdomain3, nilai_subdomain4, nilai_subdomain5)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			nilai_identifikasi = VALUES(nilai_identifikasi),
			nilai_subdomain1 = VALUES(nilai_subdomain1),
			nilai_subdomain2 = VALUES(nilai_subdomain2),
			nilai_subdomain3 = VALUES(nilai_subdomain3),
			nilai_subdomain4 = VALUES(nilai_subdomain4),
			nilai_subdomain5 = VALUES(nilai_subdomain5)`

	_, err = r.db.Exec(upsertQuery,
		perusahaanID,
		nilaiIdentifikasi,
		subdomain[1], subdomain[2], subdomain[3], subdomain[4], subdomain[5],
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	// Ambil id identifikasi yang baru di-upsert
	var identifikasiID int
	err = r.db.QueryRow(`SELECT id FROM identifikasi WHERE perusahaan_id = ?`, perusahaanID).Scan(&identifikasiID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	// Update tabel ikas agar id_identifikasi menunjuk ke identifikasi yang baru dihitung
	_, err = r.db.Exec(`UPDATE ikas SET id_identifikasi = ? WHERE id_perusahaan = ?`, identifikasiID, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	updateKematanganQuery := `
		UPDATE ikas i
		LEFT JOIN identifikasi iden ON i.id_identifikasi = iden.id
		LEFT JOIN proteksi prot ON i.id_proteksi = prot.id
		LEFT JOIN deteksi det ON i.id_deteksi = det.id
		LEFT JOIN gulih g ON i.id_gulih = g.id
		SET i.nilai_kematangan = ROUND(
			COALESCE(iden.nilai_identifikasi, 0) * 0.25 + 
			COALESCE(prot.nilai_proteksi, 0) * 0.30 + 
			COALESCE(det.nilai_deteksi, 0) * 0.25 + 
			COALESCE(g.nilai_gulih, 0) * 0.20
		, 2)
		WHERE i.id_perusahaan = ? AND (
			iden.id IS NOT NULL OR prot.id IS NOT NULL OR det.id IS NOT NULL OR g.id IS NOT NULL
		)`

	_, err = r.db.Exec(updateKematanganQuery, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *JawabanIdentifikasiRepository) UpsertToBuffer(req dto.CreateJawabanIdentifikasiRequest) error {
	query := `INSERT INTO jawaban_identifikasi_buffer 
		(pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		jawaban_identifikasi = VALUES(jawaban_identifikasi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	_, err := r.db.Exec(query,
		req.PertanyaanIdentifikasiID,
		req.PerusahaanID,
		req.JawabanIdentifikasi,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	return err
}

func (r *JawabanIdentifikasiRepository) GetBufferCount(perusahaanID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jawaban_identifikasi_buffer WHERE perusahaan_id = ?`
	err := r.db.QueryRow(query, perusahaanID).Scan(&count)
	return count, err
}

func (r *JawabanIdentifikasiRepository) FlushBuffer(perusahaanID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Move from buffer to main table
	moveQuery := `INSERT INTO jawaban_identifikasi 
		(pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi, evidence, validasi, keterangan)
		SELECT pertanyaan_identifikasi_id, perusahaan_id, jawaban_identifikasi, evidence, validasi, keterangan
		FROM jawaban_identifikasi_buffer WHERE perusahaan_id = ?
		ON DUPLICATE KEY UPDATE 
		jawaban_identifikasi = VALUES(jawaban_identifikasi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	if _, err := tx.Exec(moveQuery, perusahaanID); err != nil {
		return err
	}

	// 2. Delete from buffer
	if _, err := tx.Exec(`DELETE FROM jawaban_identifikasi_buffer WHERE perusahaan_id = ?`, perusahaanID); err != nil {
		return err
	}

	return tx.Commit()
}
