package repository

import "database/sql"

type Risiko struct {
	ID                    int
	NamaRisiko            string
	Deskripsi             string

	PotensiKejadian       string
	DampakReputasi        string
	DampakOperasional     string
	DampakFinansial       string
	DampakHukum           string

	Frekuensi             string

	AdaPengendalian       string
	DeskripsiPengendalian string
}

type RisikoRepository struct {
	db *sql.DB
}

func NewRisikoRepository(db *sql.DB) *RisikoRepository {
	return &RisikoRepository{db: db}
}

// CREATE
func (r *RisikoRepository) Create(data Risiko) (Risiko, error) {

	query := `
	INSERT INTO risiko 
	(nama_risiko, deskripsi, potensi_kejadian, dampak_reputasi,
	dampak_operasional, dampak_finansial, dampak_hukum,
	frekuensi, ada_pengendalian, deskripsi_pengendalian)
	VALUES (?,?,?,?,?,?,?,?,?,?)
	`

	res, err := r.db.Exec(query,
		data.NamaRisiko,
		data.Deskripsi,
		data.PotensiKejadian,
		data.DampakReputasi,
		data.DampakOperasional,
		data.DampakFinansial,
		data.DampakHukum,
		data.Frekuensi,
		data.AdaPengendalian,
		data.DeskripsiPengendalian,
	)

	if err != nil {
		return data, err
	}

	id, _ := res.LastInsertId()
	data.ID = int(id)

	return data, nil
}

// GET ALL
func (r *RisikoRepository) GetAll() ([]Risiko, error) {

	rows, err := r.db.Query(`SELECT * FROM risiko`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Risiko

	for rows.Next() {
		var d Risiko

		err := rows.Scan(
			&d.ID,
			&d.NamaRisiko,
			&d.Deskripsi,
			&d.PotensiKejadian,
			&d.DampakReputasi,
			&d.DampakOperasional,
			&d.DampakFinansial,
			&d.DampakHukum,
			&d.Frekuensi,
			&d.AdaPengendalian,
			&d.DeskripsiPengendalian,
			new(interface{}),
		)

		if err != nil {
			return nil, err
		}

		list = append(list, d)
	}

	return list, nil
}