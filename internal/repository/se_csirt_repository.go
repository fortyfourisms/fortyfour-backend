package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/utils"
)

type SeCsirtRepository struct {
    db *sql.DB
}

func NewSeCsirtRepository(db *sql.DB) *SeCsirtRepository {
    return &SeCsirtRepository{db: db}
}

func (r *SeCsirtRepository) Create(req dto.CreateSeCsirtRequest, id string) error {
    _, err := r.db.Exec(`
        INSERT INTO se_csirt
        (id, id_csirt, nama_se, ip_se, as_number_se, pengelola_se, fitur_se, kategori_se, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `,
        id,
        utils.ValueOrEmpty(req.IdCsirt),
        utils.ValueOrEmpty(req.NamaSe),
        utils.ValueOrEmpty(req.IpSe),
        utils.ValueOrEmpty(req.AsNumberSe),
        utils.ValueOrEmpty(req.PengelolaSe),
        utils.ValueOrEmpty(req.FiturSe),
        utils.ValueOrEmpty(req.KategoriSe),
    )
    return err
}

func (r *SeCsirtRepository) GetAll() ([]dto.SeCsirtResponse, error) {
    rows, err := r.db.Query(`
        SELECT id, id_csirt, nama_se, ip_se, as_number_se, pengelola_se, fitur_se, kategori_se, created_at, updated_at
        FROM se_csirt
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    res := []dto.SeCsirtResponse{}
    for rows.Next() {
        var s dto.SeCsirtResponse
        rows.Scan(
            &s.ID,
            &s.IdCsirt,
            &s.NamaSe,
            &s.IpSe,
            &s.AsNumberSe,
            &s.PengelolaSe,
            &s.FiturSe,
            &s.KategoriSe,
            &s.CreatedAt,
            &s.UpdatedAt,
        )
        res = append(res, s)
    }
    return res, nil
}

func (r *SeCsirtRepository) GetByID(id string) (*dto.SeCsirtResponse, error) {
    var s dto.SeCsirtResponse
    row := r.db.QueryRow(`
        SELECT id, id_csirt, nama_se, ip_se, as_number_se, pengelola_se, fitur_se, kategori_se, created_at, updated_at
        FROM se_csirt WHERE id=?
    `, id)
    if err := row.Scan(
        &s.ID,
        &s.IdCsirt,
        &s.NamaSe,
        &s.IpSe,
        &s.AsNumberSe,
        &s.PengelolaSe,
        &s.FiturSe,
        &s.KategoriSe,
        &s.CreatedAt,
        &s.UpdatedAt,
    ); err != nil {
        return nil, err
    }
    return &s, nil
}

func (r *SeCsirtRepository) Update(id string, s dto.SeCsirtResponse) error {
    _, err := r.db.Exec(`
        UPDATE se_csirt SET
            nama_se=?, ip_se=?, as_number_se=?, pengelola_se=?, fitur_se=?, kategori_se=?, updated_at=CURRENT_TIMESTAMP
        WHERE id=?`,
        s.NamaSe,
        s.IpSe,
        s.AsNumberSe,
        s.PengelolaSe,
        s.FiturSe,
        s.KategoriSe,
        id,
    )
    return err
}

func (r *SeCsirtRepository) Delete(id string) error {
    _, err := r.db.Exec(`DELETE FROM se_csirt WHERE id=?`, id)
    return err
}
