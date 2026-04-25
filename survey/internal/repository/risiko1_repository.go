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
	(responden_id, risiko_id, pernah_terjadi)
	VALUES (?, ?, ?)
	ON DUPLICATE KEY UPDATE 
	pernah_terjadi = VALUES(pernah_terjadi)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		m.PernahTerjadi,
	)

	return err
}

// STEP 2A - ALASAN
func (r *RisikoRepository) UpsertAlasan(m models.RisikoAlasan) error {
	query := `
	INSERT INTO risiko_alasan 
	(responden_id, risiko_id, alasan)
	VALUES (?, ?, ?)
	ON DUPLICATE KEY UPDATE 
	alasan = VALUES(alasan)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		m.Alasan,
	)

	return err
}

// STEP 2B - DAMPAK (FIX MAPPING)
func (r *RisikoRepository) UpsertDampak(m models.RisikoDampak) error {
	query := `
	INSERT INTO risiko_dampak
	(responden_id, risiko_id,
	dampak_reputasi, dampak_operasional, dampak_finansial, dampak_hukum,
	frekuensi)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	dampak_reputasi = VALUES(dampak_reputasi),
	dampak_operasional = VALUES(dampak_operasional),
	dampak_finansial = VALUES(dampak_finansial),
	dampak_hukum = VALUES(dampak_hukum),
	frekuensi = VALUES(frekuensi)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		models.MapImpactIntToString(m.DampakReputasi),
		models.MapImpactIntToString(m.DampakOperasional),
		models.MapImpactIntToString(m.DampakFinansial),
		models.MapImpactIntToString(m.DampakHukum),
		models.MapFrequencyIntToString(m.Frekuensi), 
	)

	return err
}

// STEP 2C - PENGENDALIAN
func (r *RisikoRepository) UpsertPengendalian(m models.RisikoPengendalian) error {
	query := `
	INSERT INTO risiko_pengendalian
	(responden_id, risiko_id, ada_pengendalian, deskripsi_pengendalian)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	ada_pengendalian = VALUES(ada_pengendalian),
	deskripsi_pengendalian = VALUES(deskripsi_pengendalian)
	`

	_, err := r.db.Exec(query,
		m.RespondenID,
		m.RisikoID,
		m.AdaPengendalian,
		m.DeskripsiPengendalian,
	)

	return err
}

// GET FULL RESPONSE (FIX SCAN + MAPPING)
func (r *RisikoRepository) FindByRespondentID(respondenID int) (map[string]interface{}, error) {

	query := `SELECT 
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
	WHERE e.responden_id = ?`

	row := r.db.QueryRow(query, respondenID)

	var (
		pernahTerjadi bool
		alasan        sql.NullString

		dampakReputasi    sql.NullString
		dampakOperasional sql.NullString
		dampakFinansial   sql.NullString
		dampakHukum       sql.NullString
		frekuensi         sql.NullString

		adaPengendalian       sql.NullBool
		deskripsiPengendalian sql.NullString
	)

	err := row.Scan(
		&pernahTerjadi,
		&alasan,
		&dampakReputasi,
		&dampakOperasional,
		&dampakFinansial,
		&dampakHukum,
		&frekuensi,
		&adaPengendalian,
		&deskripsiPengendalian,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// SAFE mapping (tidak jadi 0 kalau NULL)
	result := map[string]interface{}{
		"pernah_terjadi": pernahTerjadi,
		"alasan":         alasan.String,
		"ada_pengendalian":       adaPengendalian.Bool,
		"deskripsi_pengendalian": deskripsiPengendalian.String,
	}

	if dampakReputasi.Valid {
		result["dampak_reputasi"] = models.MapImpactStringToInt(dampakReputasi.String)
	}
	if dampakOperasional.Valid {
		result["dampak_operasional"] = models.MapImpactStringToInt(dampakOperasional.String)
	}
	if dampakFinansial.Valid {
		result["dampak_finansial"] = models.MapImpactStringToInt(dampakFinansial.String)
	}
	if dampakHukum.Valid {
		result["dampak_hukum"] = models.MapImpactStringToInt(dampakHukum.String)
	}
	if frekuensi.Valid {
		result["frekuensi"] = models.MapFrequencyStringToInt(frekuensi.String)
	}

	return result, nil
}

// PROGRESS
func (r *RisikoRepository) UpsertProgress(p models.SurveyProgress) error {
	query := `
	INSERT INTO survey_progress
	(responden_id, risiko_id, langkah_saat_ini, selesai)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	risiko_id = VALUES(risiko_id),
	langkah_saat_ini = VALUES(langkah_saat_ini),
	selesai = VALUES(selesai)
	`

	var risikoID interface{}
	if p.RisikoID.Valid {
		risikoID = p.RisikoID.Int64
	} else {
		risikoID = nil
	}

	var langkah interface{}
	if p.LangkahSaatIni.Valid {
		langkah = p.LangkahSaatIni.String
	} else {
		langkah = nil
	}

	_, err := r.db.Exec(query,
		p.RespondenID,
		risikoID,
		langkah,
		p.Selesai,
	)

	return err
}

func (r *RisikoRepository) GetProgress(respondenID int) (*models.SurveyProgress, error) {
	query := `
	SELECT id, responden_id, risiko_id, langkah_saat_ini, selesai, terakhir_update
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
		&p.Selesai,
		&p.TerakhirUpdate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			// insert default progress
			_, err := r.db.Exec(`
				INSERT INTO survey_progress (responden_id, langkah_saat_ini, selesai)
				VALUES (?, 'eligibility', false)
			`, respondenID)

			if err != nil {
				return nil, err
			}

			// ambil ulang
			return r.GetProgress(respondenID)
		}
		return nil, err
	}

	return &p, nil
}

func (r *RisikoRepository) ExistsRisiko(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM risiko WHERE id = ?)", id).Scan(&exists)
	return exists, err
}

// CUSTOM RISIKO
func (r *RisikoRepository) InsertCustomRisiko(respondenID int, nama string) (int, error) {

	result, err := r.db.Exec(`
		INSERT INTO risiko_custom (responden_id, nama_risiko)
		VALUES (?, ?)
	`, respondenID, nama)

	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *RisikoRepository) ExistsCustomRisiko(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM risiko_custom WHERE id = ?)
	`, id).Scan(&exists)
	return exists, err
}

// Exists Responden
func (r *RisikoRepository) ExistsResponden(id int) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM responden WHERE id = ?)",
		id,
	).Scan(&exists)

	return exists, err
}