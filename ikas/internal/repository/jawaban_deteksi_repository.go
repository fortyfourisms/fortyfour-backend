package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"ikas/internal/utils"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type JawabanDeteksiRepositoryInterface interface {
	Create(req dto.CreateJawabanDeteksiRequest) (int64, error)
	GetAll() ([]dto.JawabanDeteksiResponse, error)
	GetByID(id int) (*dto.JawabanDeteksiResponse, error)
	GetByIkasID(ikasID string) ([]dto.JawabanDeteksiResponse, error)
	GetByPertanyaan(pertanyaanID int) ([]dto.JawabanDeteksiResponse, error)
	Update(id int, req dto.UpdateJawabanDeteksiRequest) error
	Delete(id int) error
	CheckPertanyaanExists(pertanyaanID int) (bool, error)
	CheckIkasExists(ikasID string) (bool, error)
	CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error)
	RecalculateDeteksi(ikasID string) error
	UpsertToBuffer(req dto.CreateJawabanDeteksiRequest) error
	GetBufferCount(ikasID string) (int, error)
	FlushBuffer(ikasID string) error
}

type JawabanDeteksiRepository struct {
	db *sql.DB
}

func NewJawabanDeteksiRepository(db *sql.DB) *JawabanDeteksiRepository {
	return &JawabanDeteksiRepository{db: db}
}

const jawabanDeteksiSelectQuery = `
	SELECT 
		jd.id, 
		jd.ikas_id, 
		jd.jawaban_deteksi, 
		jd.evidence, 
		jd.validasi, 
		jd.keterangan, 
		jd.created_at, 
		jd.updated_at,
		pd.id, 
		pd.pertanyaan_deteksi, 
		sk.id, sk.nama_sub_kategori,
		k.id, k.nama_kategori,
		d.id, d.nama_domain
	FROM jawaban_deteksi jd
	JOIN pertanyaan_deteksi pd ON jd.pertanyaan_deteksi_id = pd.id
	JOIN sub_kategori sk ON pd.sub_kategori_id = sk.id
	JOIN kategori k ON sk.kategori_id = k.id
	JOIN domain d ON k.domain_id = d.id`

func scanJawabanDeteksi(row interface {
	Scan(dest ...any) error
}) (dto.JawabanDeteksiResponse, error) {
	var item dto.JawabanDeteksiResponse
	err := row.Scan(
		&item.ID,
		&item.IkasID,
		&item.JawabanDeteksi,
		&item.Evidence,
		&item.Validasi,
		&item.Keterangan,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.PertanyaanDeteksi.ID,
		&item.PertanyaanDeteksi.PertanyaanDeteksi,
		&item.PertanyaanDeteksi.SubKategori.ID,
		&item.PertanyaanDeteksi.SubKategori.NamaSubKategori,
		&item.PertanyaanDeteksi.SubKategori.Kategori.ID,
		&item.PertanyaanDeteksi.SubKategori.Kategori.NamaKategori,
		&item.PertanyaanDeteksi.SubKategori.Kategori.Domain.ID,
		&item.PertanyaanDeteksi.SubKategori.Kategori.Domain.NamaDomain,
	)
	return item, err
}

func (r *JawabanDeteksiRepository) Create(req dto.CreateJawabanDeteksiRequest) (int64, error) {
	query := `INSERT INTO jawaban_deteksi 
		(pertanyaan_deteksi_id, ikas_id, jawaban_deteksi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		req.PertanyaanDeteksiID,
		req.IkasID,
		req.JawabanDeteksi,
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

func (r *JawabanDeteksiRepository) GetAll() ([]dto.JawabanDeteksiResponse, error) {
	query := jawabanDeteksiSelectQuery + ` ORDER BY jd.created_at ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var results []dto.JawabanDeteksiResponse
	for rows.Next() {
		item, err := scanJawabanDeteksi(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *JawabanDeteksiRepository) GetByID(id int) (*dto.JawabanDeteksiResponse, error) {
	query := jawabanDeteksiSelectQuery + ` WHERE jd.id = ?`
	row := r.db.QueryRow(query, id)
	item, err := scanJawabanDeteksi(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		rollbar.Error(err)
		return nil, err
	}
	return &item, nil
}

func (r *JawabanDeteksiRepository) GetByIkasID(ikasID string) ([]dto.JawabanDeteksiResponse, error) {
	query := jawabanDeteksiSelectQuery + ` WHERE jd.ikas_id = ? ORDER BY jd.created_at ASC`

	rows, err := r.db.Query(query, ikasID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var results []dto.JawabanDeteksiResponse
	for rows.Next() {
		item, err := scanJawabanDeteksi(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *JawabanDeteksiRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanDeteksiResponse, error) {
	query := jawabanDeteksiSelectQuery + ` WHERE jd.pertanyaan_deteksi_id = ? ORDER BY jd.created_at ASC`
	rows, err := r.db.Query(query, pertanyaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var results []dto.JawabanDeteksiResponse
	for rows.Next() {
		item, err := scanJawabanDeteksi(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *JawabanDeteksiRepository) Update(id int, req dto.UpdateJawabanDeteksiRequest) error {
	var updates []string
	var args []interface{}

	if req.JawabanDeteksi != nil {
		updates = append(updates, "jawaban_deteksi = ?")
		args = append(args, req.JawabanDeteksi)
	}
	if req.Evidence != nil {
		updates = append(updates, "evidence = ?")
		args = append(args, req.Evidence)
	}
	if req.Validasi != nil {
		updates = append(updates, "validasi = ?")
		args = append(args, req.Validasi)
	}
	if req.Keterangan != nil {
		updates = append(updates, "keterangan = ?")
		args = append(args, req.Keterangan)
	}

	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE jawaban_deteksi SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		rollbar.Error(err)
	}
	return err
}

func (r *JawabanDeteksiRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM jawaban_deteksi WHERE id = ?`, id)
	if err != nil {
		rollbar.Error(err)
	}
	return err
}

func (r *JawabanDeteksiRepository) CheckPertanyaanExists(pertanyaanID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM pertanyaan_deteksi WHERE id = ?`, pertanyaanID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanDeteksiRepository) CheckIkasExists(ikasID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ikas WHERE id = ?`, ikasID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanDeteksiRepository) CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM jawaban_deteksi
		WHERE ikas_id = ? AND pertanyaan_deteksi_id = ?`
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

func (r *JawabanDeteksiRepository) RecalculateDeteksi(ikasID string) error {
	// Query rata-rata jawaban per kategori_id untuk assessment tertentu
	query := `
		SELECT k.id AS kategori_id, ROUND(AVG(jd.jawaban_deteksi), 2) AS avg_nilai
		FROM jawaban_deteksi jd
		JOIN pertanyaan_deteksi pd ON jd.pertanyaan_deteksi_id = pd.id
		JOIN sub_kategori sk ON pd.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		WHERE jd.ikas_id = ? AND jd.jawaban_deteksi IS NOT NULL
		GROUP BY k.id
		ORDER BY k.id`

	rows, err := r.db.Query(query, ikasID)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	defer rows.Close()

	// Map kategori_id ke nilai subdomain (default 0)
	// kategori 12→subdomain1, 13→subdomain2, 14→subdomain3
	subdomain := map[int]float64{12: 0, 13: 0, 14: 0}

	for rows.Next() {
		var kategoriID int
		var avgNilai float64
		if err := rows.Scan(&kategoriID, &avgNilai); err != nil {
			rollbar.Error(err)
			continue
		}
		if kategoriID >= 12 && kategoriID <= 14 {
			subdomain[kategoriID] = avgNilai
		}
	}

	// Hitung rata-rata keseluruhan (nilai_deteksi)
	nilaiDeteksi := utils.RoundToTwo((subdomain[12] + subdomain[13] + subdomain[14]) / 3.0)

	upsertQuery := `
		INSERT INTO deteksi 
			(ikas_id, nilai_deteksi, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			nilai_deteksi = VALUES(nilai_deteksi),
			nilai_subdomain1 = VALUES(nilai_subdomain1),
			nilai_subdomain2 = VALUES(nilai_subdomain2),
			nilai_subdomain3 = VALUES(nilai_subdomain3)`

	_, err = r.db.Exec(upsertQuery,
		ikasID,
		nilaiDeteksi,
		subdomain[12], subdomain[13], subdomain[14],
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	var deteksiID int
	err = r.db.QueryRow(`SELECT id FROM deteksi WHERE ikas_id = ?`, ikasID).Scan(&deteksiID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	_, err = r.db.Exec(`UPDATE ikas SET id_deteksi = ? WHERE id = ?`, deteksiID, ikasID)
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

func (r *JawabanDeteksiRepository) UpsertToBuffer(req dto.CreateJawabanDeteksiRequest) error {
	query := `INSERT INTO jawaban_deteksi_buffer 
		(pertanyaan_deteksi_id, ikas_id, jawaban_deteksi, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		jawaban_deteksi = VALUES(jawaban_deteksi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	_, err := r.db.Exec(query,
		req.PertanyaanDeteksiID,
		req.IkasID,
		req.JawabanDeteksi,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	return err
}

func (r *JawabanDeteksiRepository) GetBufferCount(ikasID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jawaban_deteksi_buffer WHERE ikas_id = ?`
	err := r.db.QueryRow(query, ikasID).Scan(&count)
	return count, err
}

func (r *JawabanDeteksiRepository) FlushBuffer(ikasID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	moveQuery := `INSERT INTO jawaban_deteksi 
		(pertanyaan_deteksi_id, ikas_id, jawaban_deteksi, evidence, validasi, keterangan)
		SELECT pertanyaan_deteksi_id, ikas_id, jawaban_deteksi, evidence, validasi, keterangan
		FROM jawaban_deteksi_buffer WHERE ikas_id = ?
		ON DUPLICATE KEY UPDATE 
		jawaban_deteksi = VALUES(jawaban_deteksi),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	if _, err := tx.Exec(moveQuery, ikasID); err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM jawaban_deteksi_buffer WHERE ikas_id = ?`, ikasID); err != nil {
		return err
	}

	return tx.Commit()
}
