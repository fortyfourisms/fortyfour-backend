package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"
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
        (id, id_csirt, nama_personel, jabatan_csirt, jabatan_perusahaan, skill, sertifikasi, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `,
        id,
        utils.ValueOrEmpty(req.IdCsirt),
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
        SELECT id, id_csirt, nama_personel, jabatan_csirt, jabatan_perusahaan, skill, sertifikasi, created_at, updated_at
        FROM sdm_csirt
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    res := []dto.SdmCsirtResponse{}
    for rows.Next() {
        var s dto.SdmCsirtResponse
        rows.Scan(
            &s.ID,
            &s.IdCsirt,
            &s.NamaPersonel,
            &s.JabatanCsirt,
            &s.JabatanPerusahaan,
            &s.Skill,
            &s.Sertifikasi,
            &s.CreatedAt,
            &s.UpdatedAt,
        )
        res = append(res, s)
    }
    return res, nil
}

func (r *SdmCsirtRepository) GetByID(id string) (*dto.SdmCsirtResponse, error) {
    var s dto.SdmCsirtResponse
    row := r.db.QueryRow(`
        SELECT id, id_csirt, nama_personel, jabatan_csirt, jabatan_perusahaan, skill, sertifikasi, created_at, updated_at
        FROM sdm_csirt WHERE id=?
    `, id)
    if err := row.Scan(
        &s.ID,
        &s.IdCsirt,
        &s.NamaPersonel,
        &s.JabatanCsirt,
        &s.JabatanPerusahaan,
        &s.Skill,
        &s.Sertifikasi,
        &s.CreatedAt,
        &s.UpdatedAt,
    ); err != nil {
        return nil, err
    }
    return &s, nil
}

func (r *SdmCsirtRepository) Update(id string, s dto.SdmCsirtResponse) error {
    _, err := r.db.Exec(`
        UPDATE sdm_csirt SET
            nama_personel=?, jabatan_csirt=?, jabatan_perusahaan=?, skill=?, sertifikasi=?, updated_at=CURRENT_TIMESTAMP
        WHERE id=?`,
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
