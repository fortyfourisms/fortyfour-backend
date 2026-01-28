package repository

import (
	"context"
	"database/sql"
	"fmt"

	"fortyfour-backend/internal/dto"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// CountPerSektor returns total and this_month counts grouped by sector.
// from/to are optional date strings "YYYY-MM-DD" for filtering this_month count.
func (r *DashboardRepository) CountPerSektor(ctx context.Context, from, to *string) ([]dto.SectorCount, error) {
	var args []interface{}
	var thisMonthExpr string
	if from != nil && to != nil {
		// Use BETWEEN with given dates (caller should ensure format YYYY-MM-DD)
		thisMonthExpr = "SUM(CASE WHEN p.created_at BETWEEN ? AND ? THEN 1 ELSE 0 END) AS this_month"
		args = append(args, *from, *to)
	} else {
		// Default: from first day of current month until now (MySQL expression)
		thisMonthExpr = "SUM(CASE WHEN p.created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS this_month"
	}

	query := fmt.Sprintf(`
		SELECT s.id, s.nama_sektor, COUNT(p.id) AS total, %s
		FROM perusahaan p
		JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
		JOIN sektor s ON ss.id_sektor = s.id
		GROUP BY s.id, s.nama_sektor
	`, thisMonthExpr)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []dto.SectorCount
	for rows.Next() {
		var s dto.SectorCount
		if err := rows.Scan(&s.ID, &s.Nama, &s.Total, &s.ThisMonth); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// IkasGlobalAgg returns simple ikas aggregate across all data (for summary).
// If ikas table or columns missing, this returns zeros or an error.
func (r *DashboardRepository) IkasGlobalAgg(ctx context.Context) (dto.IkasAgg, error) {
	var out dto.IkasAgg
	query := `SELECT COUNT(id) as total_ikas, AVG(nilai_kematangan) as avg_nilai_kematangan, AVG(target_nilai) as avg_target_nilai FROM ikas`
	row := r.db.QueryRowContext(ctx, query)
	if err := row.Scan(&out.Total, &out.AvgNilaiKematangan, &out.AvgTargetNilai); err != nil {
		// if no rows, return zero struct
		if err == sql.ErrNoRows {
			return out, nil
		}
		return out, err
	}
	return out, nil
}

// SeGlobalAgg returns simple se aggregate for summary
func (r *DashboardRepository) SeGlobalAgg(ctx context.Context) (dto.SeAgg, error) {
	var out dto.SeAgg
	query := `SELECT COUNT(id) as total_se FROM se`
	row := r.db.QueryRowContext(ctx, query)
	if err := row.Scan(&out.TotalSE); err != nil {
		if err == sql.ErrNoRows {
			return out, nil
		}
		return out, err
	}
	return out, nil
}
