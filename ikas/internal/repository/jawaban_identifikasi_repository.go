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
	GetByUUID(uuid string) (*dto.JawabanIdentifikasiResponse, error)
	GetByIkasID(ikasID string) ([]dto.JawabanIdentifikasiResponse, error)
	GetByIkasIDFromBuffer(ikasID string) ([]dto.JawabanIdentifikasiResponse, error)
	GetByPerusahaanID(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error)
	GetByPertanyaan(pertanyaanID int) ([]dto.JawabanIdentifikasiResponse, error)
	GetByPertanyaanAndPerusahaan(pertanyaanID int, perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error)
	Update(uuid string, req dto.UpdateJawabanIdentifikasiRequest) error
	Delete(uuid string) error
	CheckPertanyaanExists(pertanyaanID int) (bool, error)
	CheckIkasExists(ikasID string) (bool, error)
	CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error)
	RecalculateIdentifikasi(ikasID string) error
	UpsertToBuffer(req dto.CreateJawabanIdentifikasiRequest) error
	GetBufferCount(ikasID string) (int, error)
	FlushBuffer(ikasID string) error
	CloneByIkasID(sourceID, targetID string) error
	GetIDByIkasAndPertanyaan(ikasID string, pertanyaanID int) (int, error)
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
		ji.uuid,
		ji.ikas_id,
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

const jawabanIdentifikasiBufferSelectQuery = `
	SELECT
		ji.id,
		ji.uuid,
		ji.ikas_id,
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
	FROM jawaban_identifikasi_buffer ji
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
		&item.UUID,
		&item.IkasID,
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
		(uuid, pertanyaan_identifikasi_id, ikas_id, jawaban_identifikasi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		utils.GenerateUUID(),
		req.PertanyaanIdentifikasiID,
		req.IkasID,
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

func (r *JawabanIdentifikasiRepository) GetByUUID(uuid string) (*dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + ` WHERE ji.uuid = ?`

	item, err := scanJawaban(r.db.QueryRow(query, uuid))
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return &item, nil
}

func (r *JawabanIdentifikasiRepository) GetByIkasID(ikasID string) ([]dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + ` WHERE ji.ikas_id = ? ORDER BY ji.created_at ASC`

	rows, err := r.db.Query(query, ikasID)
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

func (r *JawabanIdentifikasiRepository) GetByIkasIDFromBuffer(ikasID string) ([]dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiBufferSelectQuery + ` WHERE ji.ikas_id = ? ORDER BY ji.created_at ASC`

	rows, err := r.db.Query(query, ikasID)
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

func (r *JawabanIdentifikasiRepository) GetByPerusahaanID(perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + ` 
		JOIN ikas i ON ji.ikas_id = i.id 
		WHERE i.id_perusahaan = ? 
		ORDER BY ji.created_at ASC`

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

// GetByPertanyaanAndPerusahaan returns answers filtered by both question ID and company ID.
// This enforces multi-tenancy at the database level for non-admin users.
func (r *JawabanIdentifikasiRepository) GetByPertanyaanAndPerusahaan(pertanyaanID int, perusahaanID string) ([]dto.JawabanIdentifikasiResponse, error) {
	query := jawabanIdentifikasiSelectQuery + `
		JOIN ikas i ON ji.ikas_id = i.id
		WHERE ji.pertanyaan_identifikasi_id = ? AND i.id_perusahaan = ?
		ORDER BY ji.created_at ASC`

	rows, err := r.db.Query(query, pertanyaanID, perusahaanID)
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

func (r *JawabanIdentifikasiRepository) Update(uuid string, req dto.UpdateJawabanIdentifikasiRequest) error {
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
	query += " WHERE uuid=?"
	args = append(args, uuid)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *JawabanIdentifikasiRepository) Delete(uuid string) error {
	_, err := r.db.Exec(`DELETE FROM jawaban_identifikasi WHERE uuid=?`, uuid)
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

func (r *JawabanIdentifikasiRepository) CheckIkasExists(ikasID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ikas WHERE id = ?`, ikasID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanIdentifikasiRepository) CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM jawaban_identifikasi
		WHERE ikas_id = ? AND pertanyaan_identifikasi_id = ?`
	args := []interface{}{ikasID, pertanyaanID}

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

func (r *JawabanIdentifikasiRepository) RecalculateIdentifikasi(ikasID string) error {
	// Query rata-rata jawaban per kategori_id untuk assessment tertentu
	query := `
		SELECT k.id AS kategori_id, ROUND(AVG(ji.jawaban_identifikasi), 2) AS avg_nilai
		FROM jawaban_identifikasi ji
		JOIN pertanyaan_identifikasi pi ON ji.pertanyaan_identifikasi_id = pi.id
		JOIN sub_kategori sk ON pi.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		WHERE ji.ikas_id = ? AND ji.jawaban_identifikasi IS NOT NULL
		GROUP BY k.id
		ORDER BY k.id`

	rows, err := r.db.Query(query, ikasID)
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
			(ikas_id, nilai_identifikasi, nilai_subdomain1, nilai_subdomain2, 
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
		ikasID,
		nilaiIdentifikasi,
		subdomain[1], subdomain[2], subdomain[3], subdomain[4], subdomain[5],
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	// Ambil id identifikasi yang baru di-upsert
	var identifikasiID int
	err = r.db.QueryRow(`SELECT id FROM identifikasi WHERE ikas_id = ?`, ikasID).Scan(&identifikasiID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	// Update tabel ikas agar id_identifikasi menunjuk ke identifikasi yang baru dihitung
	_, err = r.db.Exec(`UPDATE ikas SET id_identifikasi = ? WHERE id = ?`, identifikasiID, ikasID)
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
		WHERE i.id = ? AND (
			iden.id IS NOT NULL OR prot.id IS NOT NULL OR det.id IS NOT NULL OR g.id IS NOT NULL
		)`

	_, err = r.db.Exec(updateKematanganQuery, ikasID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *JawabanIdentifikasiRepository) UpsertToBuffer(req dto.CreateJawabanIdentifikasiRequest) error {
	// 1. Cek dulu apakah data ini SUDAH ADA di tabel utama
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM jawaban_identifikasi WHERE ikas_id = ? AND pertanyaan_identifikasi_id = ?)`
	err := r.db.QueryRow(checkQuery, req.IkasID, req.PertanyaanIdentifikasiID).Scan(&exists)
	if err != nil {
		return err
	}

	// 2. Jika sudah ada di utama, jangan masukkan ke buffer.
	// Sebaliknya, pastikan di buffer bersih (cleanup sampah)
	if exists {
		_, _ = r.db.Exec(`DELETE FROM jawaban_identifikasi_buffer WHERE ikas_id = ? AND pertanyaan_identifikasi_id = ?`, req.IkasID, req.PertanyaanIdentifikasiID)
		return nil
	}

	// 3. Jika belum ada di utama, baru masukkan ke buffer seperti biasa
	query := `INSERT INTO jawaban_identifikasi_buffer 
		(uuid, pertanyaan_identifikasi_id, ikas_id, jawaban_identifikasi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		jawaban_identifikasi = VALUES(jawaban_identifikasi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	_, err = r.db.Exec(query,
		utils.GenerateUUID(),
		req.PertanyaanIdentifikasiID,
		req.IkasID,
		req.JawabanIdentifikasi,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	return err
}

func (r *JawabanIdentifikasiRepository) GetBufferCount(ikasID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jawaban_identifikasi_buffer WHERE ikas_id = ?`
	err := r.db.QueryRow(query, ikasID).Scan(&count)
	return count, err
}

func (r *JawabanIdentifikasiRepository) FlushBuffer(ikasID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Move from buffer to main table
	moveQuery := `INSERT INTO jawaban_identifikasi 
		(uuid, pertanyaan_identifikasi_id, ikas_id, jawaban_identifikasi, evidence, validasi, keterangan)
		SELECT uuid, pertanyaan_identifikasi_id, ikas_id, jawaban_identifikasi, evidence, validasi, keterangan
		FROM jawaban_identifikasi_buffer WHERE ikas_id = ?
		ON DUPLICATE KEY UPDATE 
		jawaban_identifikasi = VALUES(jawaban_identifikasi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	if _, err := tx.Exec(moveQuery, ikasID); err != nil {
		return err
	}

	// 2. Delete from buffer
	if _, err := tx.Exec(`DELETE FROM jawaban_identifikasi_buffer WHERE ikas_id = ?`, ikasID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *JawabanIdentifikasiRepository) CloneByIkasID(sourceID, targetID string) error {
	query := `
		INSERT INTO jawaban_identifikasi 
			(uuid, pertanyaan_identifikasi_id, ikas_id, jawaban_identifikasi, evidence, validasi, keterangan)
		SELECT 
			uuid, pertanyaan_identifikasi_id, ?, jawaban_identifikasi, evidence, validasi, keterangan
		FROM jawaban_identifikasi 
		WHERE ikas_id = ?`

	_, err := r.db.Exec(query, targetID, sourceID)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *JawabanIdentifikasiRepository) GetIDByIkasAndPertanyaan(ikasID string, pertanyaanID int) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM jawaban_identifikasi WHERE ikas_id = ? AND pertanyaan_identifikasi_id = ?`, ikasID, pertanyaanID).Scan(&id)
	return id, err
}
