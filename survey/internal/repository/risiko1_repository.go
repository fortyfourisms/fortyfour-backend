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
	rows, err := r.db.Query(`
		SELECT id, nama, COALESCE(deskripsi, '')
		FROM risiko
		WHERE aktif = TRUE
		ORDER BY urutan ASC
	`)
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

func (r *RisikoRepository) GetByID(id string) (*models.Risiko, error) {

	row := r.db.QueryRow(`
		SELECT id, kode, nama, deskripsi, urutan, aktif, created_at, updated_at
		FROM risiko
		WHERE id = ?
	`, id)

	var m models.Risiko

	err := row.Scan(
		&m.ID,
		&m.Kode,
		&m.Nama,
		&m.Deskripsi,
		&m.Urutan,
		&m.Aktif,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *RisikoRepository) GetRisikoIDByUrutan(urutan int) (string, error) {
	row := r.db.QueryRow(`
		SELECT id
		FROM risiko
		WHERE urutan = ? AND aktif = TRUE
		LIMIT 1
	`, urutan)

	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *RisikoRepository) GetUrutanByRisikoID(id string) (int, error) {
	row := r.db.QueryRow(`
		SELECT urutan
		FROM risiko
		WHERE id = ? AND aktif = TRUE
		LIMIT 1
	`, id)

	var urutan int
	if err := row.Scan(&urutan); err != nil {
		return 0, err
	}

	return urutan, nil
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

	var risikoID interface{}
	if m.RisikoID != nil {
		risikoID = *m.RisikoID
	}

	_, err := r.db.Exec(query,
		m.RespondenID,
		risikoID,
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

	var risikoID interface{}
	if m.RisikoID != nil {
		risikoID = *m.RisikoID
	}

	_, err := r.db.Exec(query,
		m.RespondenID,
		risikoID,
		m.Alasan,
	)

	return err
}

// STEP 2B - DAMPAK
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

	var risikoID interface{}
	if m.RisikoID != nil {
		risikoID = *m.RisikoID
	}

	_, err := r.db.Exec(query,
		m.RespondenID,
		risikoID,
		m.DampakReputasi,
		m.DampakOperasional,
		m.DampakFinansial,
		m.DampakHukum,
		m.Frekuensi,
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

	var risikoID interface{}
	if m.RisikoID != nil {
		risikoID = *m.RisikoID
	}

	_, err := r.db.Exec(query,
		m.RespondenID,
		risikoID,
		m.AdaPengendalian,
		m.DeskripsiPengendalian,
	)

	return err
}

// GET FULL RISIKO
func (r *RisikoRepository) FindByRespondentID(respondenID string) (map[string]interface{}, error) {

	query := `
	SELECT 
		e.responden_id,
		e.risiko_id,
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
	ORDER BY e.risiko_id ASC
	`

	rows, err := r.db.Query(query, respondenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]map[string]interface{}, 0)

	for rows.Next() {
		var (
			respondentID  string
			risikoID      sql.NullString
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

		if err := rows.Scan(
			&respondentID,
			&risikoID,
			&pernahTerjadi,
			&alasan,
			&dampakReputasi,
			&dampakOperasional,
			&dampakFinansial,
			&dampakHukum,
			&frekuensi,
			&adaPengendalian,
			&deskripsiPengendalian,
		); err != nil {
			return nil, err
		}

		item := map[string]interface{}{
			"responden_id":           respondentID,
			"pernah_terjadi":         pernahTerjadi,
			"alasan":                 alasan.String,
			"ada_pengendalian":       adaPengendalian.Bool,
			"deskripsi_pengendalian": deskripsiPengendalian.String,
		}

		if risikoID.Valid {
			item["risiko_id"] = risikoID.String
		}
		if dampakReputasi.Valid {
			item["dampak_reputasi"] = dampakReputasi.String
		}
		if dampakOperasional.Valid {
			item["dampak_operasional"] = dampakOperasional.String
		}
		if dampakFinansial.Valid {
			item["dampak_finansial"] = dampakFinansial.String
		}
		if dampakHukum.Valid {
			item["dampak_hukum"] = dampakHukum.String
		}
		if frekuensi.Valid {
			item["frekuensi"] = frekuensi.String
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrNotFound
	}

	result := map[string]interface{}{
		"responden_id": respondenID,
		"items":        items,
	}

	// Keep the first item at the top level for backward compatibility.
	for key, value := range items[0] {
		if key == "responden_id" {
			continue
		}
		result[key] = value
	}

	return result, nil
}

// PROGRESS
func (r *RisikoRepository) UpsertProgress(p models.SurveyProgress) error {

	query := `
	INSERT INTO survey_progress
	(responden_id, risiko_id, langkah_saat_ini, selesai, status, edit_request_reason, edit_request_response, submitted_at, edit_requested_at, edit_approved_at, edit_approved_by, edit_rejected_at, edit_rejected_by)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	risiko_id = VALUES(risiko_id),
	langkah_saat_ini = VALUES(langkah_saat_ini),
	selesai = VALUES(selesai),
	status = VALUES(status),
	edit_request_reason = VALUES(edit_request_reason),
	edit_request_response = VALUES(edit_request_response),
	submitted_at = VALUES(submitted_at),
	edit_requested_at = VALUES(edit_requested_at),
	edit_approved_at = VALUES(edit_approved_at),
	edit_approved_by = VALUES(edit_approved_by),
	edit_rejected_at = VALUES(edit_rejected_at),
	edit_rejected_by = VALUES(edit_rejected_by)
	`

	var risikoID interface{}
	if p.RisikoID.Valid {
		risikoID = p.RisikoID.String
	}

	var langkah interface{}
	if p.LangkahSaatIni.Valid {
		langkah = p.LangkahSaatIni.String
	}

	status := p.Status
	if status == "" {
		status = "draft"
	}

	var editReason interface{}
	if p.EditReason.Valid {
		editReason = p.EditReason.String
	}

	var editResponse interface{}
	if p.EditResponse.Valid {
		editResponse = p.EditResponse.String
	}

	var submittedAt interface{}
	if p.SubmittedAt.Valid {
		submittedAt = p.SubmittedAt.Time
	}

	var editRequestedAt interface{}
	if p.EditRequestedAt.Valid {
		editRequestedAt = p.EditRequestedAt.Time
	}

	var editApprovedAt interface{}
	if p.EditApprovedAt.Valid {
		editApprovedAt = p.EditApprovedAt.Time
	}

	var editApprovedBy interface{}
	if p.EditApprovedBy.Valid {
		editApprovedBy = p.EditApprovedBy.String
	}

	var editRejectedAt interface{}
	if p.EditRejectedAt.Valid {
		editRejectedAt = p.EditRejectedAt.Time
	}

	var editRejectedBy interface{}
	if p.EditRejectedBy.Valid {
		editRejectedBy = p.EditRejectedBy.String
	}

	_, err := r.db.Exec(query,
		p.RespondenID,
		risikoID,
		langkah,
		p.Selesai,
		status,
		editReason,
		editResponse,
		submittedAt,
		editRequestedAt,
		editApprovedAt,
		editApprovedBy,
		editRejectedAt,
		editRejectedBy,
	)

	return err
}

func (r *RisikoRepository) GetProgress(respondenID string) (*models.SurveyProgress, error) {

	row := r.db.QueryRow(`
		SELECT id, responden_id, risiko_id, langkah_saat_ini, selesai,
			COALESCE(status, 'draft'),
			edit_request_reason,
			edit_request_response,
			submitted_at,
			edit_requested_at,
			edit_approved_at,
			edit_approved_by,
			edit_rejected_at,
			edit_rejected_by,
			terakhir_update
		FROM survey_progress
		WHERE responden_id = ?
	`, respondenID)

	var p models.SurveyProgress

	err := row.Scan(
		&p.ID,
		&p.RespondenID,
		&p.RisikoID,
		&p.LangkahSaatIni,
		&p.Selesai,
		&p.Status,
		&p.EditReason,
		&p.EditResponse,
		&p.SubmittedAt,
		&p.EditRequestedAt,
		&p.EditApprovedAt,
		&p.EditApprovedBy,
		&p.EditRejectedAt,
		&p.EditRejectedBy,
		&p.TerakhirUpdate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := r.db.Exec(`
				INSERT INTO survey_progress (responden_id, langkah_saat_ini, selesai, status)
				VALUES (?, 'eligibility', false, 'draft')
			`, respondenID)
			if err != nil {
				return nil, err
			}
			return r.GetProgress(respondenID)
		}
		return nil, err
	}

	return &p, nil
}

func (r *RisikoRepository) GetRespondentIDByUserID(userID string) (string, error) {
	row := r.db.QueryRow(`
		SELECT id
		FROM responden
		WHERE user_id = ?
	`, userID)

	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *RisikoRepository) InsertCustomRisiko(respondenID string, nama string) (int, error) {
	result, err := r.db.Exec(`
		INSERT INTO risiko_custom (responden_id, nama)
		VALUES (?, ?)
	`, respondenID, nama)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// EXISTS
func (r *RisikoRepository) ExistsRisiko(id string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM risiko WHERE id = ?)`, id).Scan(&exists)
	return exists, err
}

func (r *RisikoRepository) ExistsResponden(id string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM responden WHERE id = ?)`, id).Scan(&exists)
	return exists, err
}