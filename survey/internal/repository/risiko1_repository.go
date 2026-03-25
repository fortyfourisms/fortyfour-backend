package repository

import (
	"database/sql"
	"survey/internal/models"
)

type RisikoRepository struct {
	db *sql.DB
}

func NewRisikoRepository(db *sql.DB) *RisikoRepository {
	return &RisikoRepository{db: db}
}

func (r *RisikoRepository) Create(data models.RisikoSurvey) error {

	query := `
	INSERT INTO risiko_survey (
		responden_id,
		risiko_ip,
		dampak_reputasi,
		dampak_operasional,
		dampak_finansial,
		dampak_hukum,
		frekuensi,
		ada_pengendalian,
		tindakan_pengendalian
	)
	VALUES (?,?,?,?,?,?,?,?,?)
	`

	_, err := r.db.Exec(
		query,
		data.RespondenID,
		data.RisikoIP,
		data.DampakReputasi,
		data.DampakOperasional,
		data.DampakFinansial,
		data.DampakHukum,
		data.Frekuensi,
		data.AdaPengendalian,
		data.TindakanPengendalian,
	)

	return err
}