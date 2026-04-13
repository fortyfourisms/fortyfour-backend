package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"ikas/internal/utils"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type JawabanProteksiRepositoryInterface interface {
	Create(req dto.CreateJawabanProteksiRequest) (int64, error)
	GetAll() ([]dto.JawabanProteksiResponse, error)
	GetByID(id int) (*dto.JawabanProteksiResponse, error)
	GetByIkasID(ikasID string) ([]dto.JawabanProteksiResponse, error)
	GetByPerusahaanID(perusahaanID string) ([]dto.JawabanProteksiResponse, error)
	GetByPertanyaan(pertanyaanID int) ([]dto.JawabanProteksiResponse, error)
	Update(id int, req dto.UpdateJawabanProteksiRequest) error
	Delete(id int) error
	CheckPertanyaanExists(pertanyaanID int) (bool, error)
	CheckIkasExists(ikasID string) (bool, error)
	CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error)
	RecalculateProteksi(ikasID string) error
	UpsertToBuffer(req dto.CreateJawabanProteksiRequest) error
	GetBufferCount(ikasID string) (int, error)
	FlushBuffer(ikasID string) error
}

type JawabanProteksiRepository struct {
	db *sql.DB
}

func NewJawabanProteksiRepository(db *sql.DB) *JawabanProteksiRepository {
	return &JawabanProteksiRepository{db: db}
}

const jawabanProteksiSelectQuery = `
	SELECT
		jp.id,
		jp.ikas_id,
		jp.jawaban_proteksi,
		jp.evidence,
		jp.validasi,
		jp.keterangan,
		jp.created_at,
		jp.updated_at,
		pp.id,
		pp.pertanyaan_proteksi,
		sk.id, sk.nama_sub_kategori,
		k.id, k.nama_kategori,
		d.id, d.nama_domain
	FROM jawaban_proteksi jp
	JOIN pertanyaan_proteksi pp ON jp.pertanyaan_proteksi_id = pp.id
	JOIN sub_kategori sk ON pp.sub_kategori_id = sk.id
	JOIN kategori k ON sk.kategori_id = k.id
	JOIN domain d ON k.domain_id = d.id`

func scanJawabanProteksi(row interface {
	Scan(dest ...any) error
}) (dto.JawabanProteksiResponse, error) {
	var item dto.JawabanProteksiResponse
	err := row.Scan(
		&item.ID,
		&item.IkasID,
		&item.JawabanProteksi,
		&item.Evidence,
		&item.Validasi,
		&item.Keterangan,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.PertanyaanProteksi.ID,
		&item.PertanyaanProteksi.PertanyaanProteksi,
		&item.PertanyaanProteksi.SubKategori.ID,
		&item.PertanyaanProteksi.SubKategori.NamaSubKategori,
		&item.PertanyaanProteksi.SubKategori.Kategori.ID,
		&item.PertanyaanProteksi.SubKategori.Kategori.NamaKategori,
		&item.PertanyaanProteksi.SubKategori.Kategori.Domain.ID,
		&item.PertanyaanProteksi.SubKategori.Kategori.Domain.NamaDomain,
	)
	return item, err
}

func (r *JawabanProteksiRepository) Create(req dto.CreateJawabanProteksiRequest) (int64, error) {
	query := `INSERT INTO jawaban_proteksi
		(pertanyaan_proteksi_id, ikas_id, jawaban_proteksi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		req.PertanyaanProteksiID,
		req.IkasID,
		req.JawabanProteksi,
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

func (r *JawabanProteksiRepository) GetAll() ([]dto.JawabanProteksiResponse, error) {
	query := jawabanProteksiSelectQuery + ` ORDER BY jp.created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanProteksiResponse
	for rows.Next() {
		item, err := scanJawabanProteksi(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanProteksiRepository) GetByID(id int) (*dto.JawabanProteksiResponse, error) {
	query := jawabanProteksiSelectQuery + ` WHERE jp.id = ?`

	item, err := scanJawabanProteksi(r.db.QueryRow(query, id))
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return &item, nil
}

func (r *JawabanProteksiRepository) GetByIkasID(ikasID string) ([]dto.JawabanProteksiResponse, error) {
	query := jawabanProteksiSelectQuery + ` WHERE jp.ikas_id = ? ORDER BY jp.created_at ASC`

	rows, err := r.db.Query(query, ikasID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanProteksiResponse
	for rows.Next() {
		item, err := scanJawabanProteksi(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanProteksiRepository) GetByPerusahaanID(perusahaanID string) ([]dto.JawabanProteksiResponse, error) {
	query := jawabanProteksiSelectQuery + ` 
		JOIN ikas i ON jp.ikas_id = i.id 
		WHERE i.id_perusahaan = ? 
		ORDER BY jp.created_at ASC`

	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanProteksiResponse
	for rows.Next() {
		item, err := scanJawabanProteksi(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanProteksiRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanProteksiResponse, error) {
	query := jawabanProteksiSelectQuery + ` WHERE jp.pertanyaan_proteksi_id = ? ORDER BY jp.created_at ASC`

	rows, err := r.db.Query(query, pertanyaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.JawabanProteksiResponse
	for rows.Next() {
		item, err := scanJawabanProteksi(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *JawabanProteksiRepository) Update(id int, req dto.UpdateJawabanProteksiRequest) error {
	query := "UPDATE jawaban_proteksi SET "
	args := []interface{}{}
	updates := []string{}

	if req.JawabanProteksi != nil {
		updates = append(updates, "jawaban_proteksi=?")
		args = append(args, *req.JawabanProteksi)
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

func (r *JawabanProteksiRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM jawaban_proteksi WHERE id=?`, id)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *JawabanProteksiRepository) CheckPertanyaanExists(pertanyaanID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM pertanyaan_proteksi WHERE id = ?`, pertanyaanID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanProteksiRepository) CheckIkasExists(ikasID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ikas WHERE id = ?`, ikasID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanProteksiRepository) CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM jawaban_proteksi
		WHERE ikas_id = ? AND pertanyaan_proteksi_id = ?`
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

func (r *JawabanProteksiRepository) RecalculateProteksi(ikasID string) error {
	// Query rata-rata jawaban per kategori_id untuk assessment tertentu
	query := `
		SELECT k.id AS kategori_id, ROUND(AVG(jp.jawaban_proteksi), 2) AS avg_nilai
		FROM jawaban_proteksi jp
		JOIN pertanyaan_proteksi pp ON jp.pertanyaan_proteksi_id = pp.id
		JOIN sub_kategori sk ON pp.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		WHERE jp.ikas_id = ? AND jp.jawaban_proteksi IS NOT NULL
		GROUP BY k.id
		ORDER BY k.id`

	rows, err := r.db.Query(query, ikasID)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	defer rows.Close()

	// Map kategori_id ke nilai subdomain (default 0)
	// kategori 6→subdomain1, 7→subdomain2, ..., 11→subdomain6
	subdomain := map[int]float64{6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0}

	for rows.Next() {
		var kategoriID int
		var avgNilai float64
		if err := rows.Scan(&kategoriID, &avgNilai); err != nil {
			rollbar.Error(err)
			continue
		}
		if kategoriID >= 6 && kategoriID <= 11 {
			subdomain[kategoriID] = avgNilai
		}
	}

	// Hitung rata-rata keseluruhan (nilai_proteksi)
	nilaiProteksi := utils.RoundToTwo((subdomain[6] + subdomain[7] + subdomain[8] + subdomain[9] + subdomain[10] + subdomain[11]) / 6.0)

	upsertQuery := `
		INSERT INTO proteksi 
			(ikas_id, nilai_proteksi, nilai_subdomain1, nilai_subdomain2, 
			 nilai_subdomain3, nilai_subdomain4, nilai_subdomain5, nilai_subdomain6)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			nilai_proteksi = VALUES(nilai_proteksi),
			nilai_subdomain1 = VALUES(nilai_subdomain1),
			nilai_subdomain2 = VALUES(nilai_subdomain2),
			nilai_subdomain3 = VALUES(nilai_subdomain3),
			nilai_subdomain4 = VALUES(nilai_subdomain4),
			nilai_subdomain5 = VALUES(nilai_subdomain5),
			nilai_subdomain6 = VALUES(nilai_subdomain6)`

	_, err = r.db.Exec(upsertQuery,
		ikasID,
		nilaiProteksi,
		subdomain[6], subdomain[7], subdomain[8], subdomain[9], subdomain[10], subdomain[11],
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	var proteksiID int
	err = r.db.QueryRow(`SELECT id FROM proteksi WHERE ikas_id = ?`, ikasID).Scan(&proteksiID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	_, err = r.db.Exec(`UPDATE ikas SET id_proteksi = ? WHERE id = ?`, proteksiID, ikasID)
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

func (r *JawabanProteksiRepository) UpsertToBuffer(req dto.CreateJawabanProteksiRequest) error {
	query := `INSERT INTO jawaban_proteksi_buffer 
		(pertanyaan_proteksi_id, ikas_id, jawaban_proteksi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		jawaban_proteksi = VALUES(jawaban_proteksi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	_, err := r.db.Exec(query,
		req.PertanyaanProteksiID,
		req.IkasID,
		req.JawabanProteksi,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	return err
}

func (r *JawabanProteksiRepository) GetBufferCount(ikasID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jawaban_proteksi_buffer WHERE ikas_id = ?`
	err := r.db.QueryRow(query, ikasID).Scan(&count)
	return count, err
}

func (r *JawabanProteksiRepository) FlushBuffer(ikasID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	moveQuery := `INSERT INTO jawaban_proteksi 
		(pertanyaan_proteksi_id, ikas_id, jawaban_proteksi, evidence, validasi, keterangan)
		SELECT pertanyaan_proteksi_id, ikas_id, jawaban_proteksi, evidence, validasi, keterangan
		FROM jawaban_proteksi_buffer WHERE ikas_id = ?
		ON DUPLICATE KEY UPDATE 
		jawaban_proteksi = VALUES(jawaban_proteksi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	if _, err := tx.Exec(moveQuery, ikasID); err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM jawaban_proteksi_buffer WHERE ikas_id = ?`, ikasID); err != nil {
		return err
	}

	return tx.Commit()
}
