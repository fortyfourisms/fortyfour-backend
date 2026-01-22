package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"

	"github.com/rollbar/rollbar-go"
)

type SdmCsirtRepository struct {
	db *sql.DB
}

func NewSdmCsirtRepository(db *sql.DB) *SdmCsirtRepository {
	return &SdmCsirtRepository{db: db}
}

func (r *SdmCsirtRepository) Create(req dto.CreateSdmCsirtRequest, id string) error {
	_, err := r.db.Exec(`
		INSERT INTO sdm_csirt
		(id, id_csirt, nama_personel, jabatan_csirt, jabatan_perusahaan, skill, sertifikasi)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		id,
		req.IdCsirt,
		utils.ValueOrEmpty(req.NamaPersonel),
		utils.ValueOrEmpty(req.JabatanCsirt),
		utils.ValueOrEmpty(req.JabatanPerusahaan),
		utils.ValueOrEmpty(req.Skill),
		utils.ValueOrEmpty(req.Sertifikasi),
	)
	return err
}

func (r *SdmCsirtRepository) GetAll() ([]dto.SdmCsirtResponse, error) {
	rows, err := r.db.Query(`
		SELECT
			s.id, s.nama_personel, s.jabatan_csirt, s.jabatan_perusahaan,
			s.skill, s.sertifikasi, s.created_at, s.updated_at,
			c.id, c.nama_csirt, c.web_csirt, c.telepon_csirt
		FROM sdm_csirt s
		LEFT JOIN csirt c ON s.id_csirt = c.id
	`)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var result []dto.SdmCsirtResponse

	for rows.Next() {
		var s dto.SdmCsirtResponse
		var c dto.CsirtMiniResponse
		var csirtID sql.NullString

		err := rows.Scan(
			&s.ID,
			&s.NamaPersonel,
			&s.JabatanCsirt,
			&s.JabatanPerusahaan,
			&s.Skill,
			&s.Sertifikasi,
			&s.CreatedAt,
			&s.UpdatedAt,
			&csirtID,
			&c.NamaCsirt,
			&c.WebCsirt,
			&c.TeleponCsirt,
		)
		if err != nil {
			return nil, err
		}

		if csirtID.Valid {
			c.ID = csirtID.String
			s.Csirt = &c
		}

		result = append(result, s)
	}

	return result, nil
}

func (r *SdmCsirtRepository) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	var s dto.SdmCsirtResponse
	var c dto.CsirtMiniResponse
	var csirtID sql.NullString

	err := r.db.QueryRow(`
		SELECT
			s.id, s.nama_personel, s.jabatan_csirt, s.jabatan_perusahaan,
			s.skill, s.sertifikasi, s.created_at, s.updated_at,
			c.id, c.nama_csirt, c.web_csirt, c.telepon_csirt
		FROM sdm_csirt s
		LEFT JOIN csirt c ON s.id_csirt = c.id
		WHERE s.id = ?
	`, id).Scan(
		&s.ID,
		&s.NamaPersonel,
		&s.JabatanCsirt,
		&s.JabatanPerusahaan,
		&s.Skill,
		&s.Sertifikasi,
		&s.CreatedAt,
		&s.UpdatedAt,
		&csirtID,
		&c.NamaCsirt,
		&c.WebCsirt,
		&c.TeleponCsirt,
	)

	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	if csirtID.Valid {
		c.ID = csirtID.String
		s.Csirt = &c
	}

	return &s, nil
}

func (r *SdmCsirtRepository) Update(id string, s dto.SdmCsirtResponse) error {
	_, err := r.db.Exec(`
		UPDATE sdm_csirt SET
			nama_personel=?,
			jabatan_csirt=?,
			jabatan_perusahaan=?,
			skill=?,
			sertifikasi=?
		WHERE id=?
	`,
		s.NamaPersonel,
		s.JabatanCsirt,
		s.JabatanPerusahaan,
		s.Skill,
		s.Sertifikasi,
		id,
	)
	return err
}

func (r *SdmCsirtRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM sdm_csirt WHERE id=?`, id)
	return err
}
