package repository

import (
	"database/sql"
	"errors"
	"survey/internal/models"
)

var ErrNotFound = sql.ErrNoRows

type RisikoRepository struct {
	db *sql.DB
}

func NewRisikoRepository(db *sql.DB) *RisikoRepository {
	return &RisikoRepository{db: db}
}

// MASTER RISIKO
func (r *RisikoRepository) GetAllRisiko() ([]models.RisikoResponse, error) {
	rows, err := r.db.Query(`SELECT id, nama, deskripsi FROM risiko`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.RisikoResponse
	for rows.Next() {
		var m models.RisikoResponse
		if err := rows.Scan(&m.ID, &m.NamaRisiko, &m.Deskripsi); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}

// STEP 1 - ELIGIBILITY
func (r *RisikoRepository) UpsertEligibility(m models.RisikoEligibility) error {
	query := `
	INSERT INTO risiko_eligibility 
	(responden_id, risiko_id, pernah_terjadi, langkah_selanjutnya)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
	pernah_terjadi = VALUES(pernah_terjadi),
	langkah_selanjutnya = VALUES(langkah_selanjutnya)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		m.PernahTerjadi,
		m.LangkahSelanjutnya,
	)

	return err
}

// STEP 2a - ALASAN (JIKA TIDAK)
func (r *RisikoRepository) UpsertAlasan(m models.RisikoAlasan) error {
	query := `
	INSERT INTO risiko_alasan 
	(responden_id, risiko_id, alasan, selesai)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
	alasan = VALUES(alasan),
	selesai = VALUES(selesai)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		m.Alasan,
		m.Selesai,
	)

	return err
}

// STEP 2b - DAMPAK (JIKA YA)
func (r *RisikoRepository) UpsertDampak(m models.RisikoDampak) error {
	query := `
	INSERT INTO risiko_dampak
	(responden_id, risiko_id,
	dampak_reputasi, dampak_operasional, dampak_finansial, dampak_hukum,
	frekuensi, langkah_selanjutnya)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	dampak_reputasi = VALUES(dampak_reputasi),
	dampak_operasional = VALUES(dampak_operasional),
	dampak_finansial = VALUES(dampak_finansial),
	dampak_hukum = VALUES(dampak_hukum),
	frekuensi = VALUES(frekuensi),
	langkah_selanjutnya = VALUES(langkah_selanjutnya)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		m.DampakReputasi,
		m.DampakOperasional,
		m.DampakFinansial,
		m.DampakHukum,
		m.Frekuensi,
		m.LangkahSelanjutnya,
	)

	return err
}

// STEP 2c - PENGENDALIAN
func (r *RisikoRepository) UpsertPengendalian(m models.RisikoPengendalian) error {
	query := `
	INSERT INTO risiko_pengendalian
	(responden_id, risiko_id, ada_pengendalian, deskripsi_pengendalian, selesai)
	VALUES (?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	ada_pengendalian = VALUES(ada_pengendalian),
	deskripsi_pengendalian = VALUES(deskripsi_pengendalian),
	selesai = VALUES(selesai)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		m.AdaPengendalian,
		m.DeskripsiPengendalian,
		m.Selesai,
	)

	return err
}

// GET FULL RESPONSE (JOIN)
func (r *RisikoRepository) FindByRespondentID(respondenID int) (map[string]interface{}, error) {

	query := `
	SELECT 
		e.pernah_terjadi,
		a.alasan,
		d.dampak_reputasi, d.dampak_operasional, d.dampak_finansial, d.dampak_hukum,
		d.frekuensi,
		p.ada_pengendalian, p.deskripsi_pengendalian
	FROM risiko_eligibility e
	LEFT JOIN risiko_alasan a 
		ON e.responden_id = a.responden_id AND e.risiko_id = a.risiko_id
	LEFT JOIN risiko_dampak d 
		ON e.responden_id = d.responden_id AND e.risiko_id = d.risiko_id
	LEFT JOIN risiko_pengendalian p 
		ON e.responden_id = p.responden_id AND e.risiko_id = p.risiko_id
	WHERE e.responden_id = ?
	`

	row := r.db.QueryRow(query, respondenID)

	var result = make(map[string]interface{})

	err := row.Scan(
		&result["pernah_terjadi"],
		&result["alasan"],
		&result["dampak_reputasi"],
		&result["dampak_operasional"],
		&result["dampak_finansial"],
		&result["dampak_hukum"],
		&result["frekuensi"],
		&result["ada_pengendalian"],
		&result["deskripsi_pengendalian"],
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return result, nil
}

// PROGRESS
func (r *RisikoRepository) UpsertProgress(p models.SurveyProgress) error {
	query := `
	INSERT INTO survey_progress
	(responden_id, risiko_id, langkah_saat_ini, nomor_risiko, selesai)
	VALUES (?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	langkah_saat_ini = VALUES(langkah_saat_ini),
	nomor_risiko = VALUES(nomor_risiko),
	selesai = VALUES(selesai)
	`

	_, err := r.db.Exec(query,
		p.RespondenID,
		p.RisikoID,
		p.LangkahSaatIni,
		p.NomorRisiko,
		p.Selesai,
	)

	return err
}

func (r *RisikoRepository) GetProgress(respondenID int) (*models.SurveyProgress, error) {
	query := `
	SELECT id, responden_id, risiko_id, langkah_saat_ini, nomor_risiko, selesai, terakhir_update
	FROM survey_progress
	WHERE responden_id = ?
	`

	row := r.db.QueryRow(query, respondenID)

	var p models.SurveyProgress
	err := row.Scan(
		&p.ID,
		&p.RespondenID,
		&p.RisikoID,
		&p.LangkahSaatIni,
		&p.NomorRisiko,
		&p.Selesai,
		&p.TerakhirUpdate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}