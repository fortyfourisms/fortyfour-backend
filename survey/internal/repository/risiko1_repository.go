package repository

import "database/sql"

type RisikoRepository struct {
	db *sql.DB
}

func NewRisikoRepository(db *sql.DB) *RisikoRepository {
	return &RisikoRepository{db: db}
}

// MASTER RISIKO
func (r *RisikoRepository) GetAllRisiko() (*sql.Rows, error) {
	return r.db.Query(`SELECT id, nama_risiko, deskripsi FROM risiko`)
}

// JAWABAN RISIKO
func (r *RisikoRepository) CreateJawaban(
	req map[string]interface{},
) error {

	query := `
	INSERT INTO risiko_responden
	(responden_id, risiko_id, pernah_terjadi,
	dampak_reputasi, dampak_operasional, dampak_finansial, dampak_hukum,
	frekuensi, ada_pengendalian, deskripsi_pengendalian)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		req["responden_id"],
		req["risiko_id"],
		req["pernah_terjadi"],
		req["dampak_reputasi"],
		req["dampak_operasional"],
		req["dampak_finansial"],
		req["dampak_hukum"],
		req["frekuensi"],
		req["ada_pengendalian"],
		req["deskripsi_pengendalian"],
	)

	return err
}