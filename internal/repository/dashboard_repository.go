package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"fortyfour-backend/internal/dto"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// buildDateRange mengkonversi filter Year/Quarter menjadi from/to string "YYYY-MM-DD".
// Prioritas: from/to eksplisit > year+quarter > year saja.
// Jika tidak ada filter, from/to dikembalikan nil (default bulan berjalan di query).
func buildDateRange(f dto.DashboardFilter) (from, to *string) {
	// Jika from/to eksplisit sudah ada, gunakan langsung
	if f.From != nil && f.To != nil {
		return f.From, f.To
	}

	// Bangun dari year + quarter
	if f.Year != nil {
		y := *f.Year
		var start, end string
		if f.Quarter != nil {
			switch *f.Quarter {
			case "1":
				start, end = y+"-01-01", y+"-03-31"
			case "2":
				start, end = y+"-04-01", y+"-06-30"
			case "3":
				start, end = y+"-07-01", y+"-09-30"
			case "4":
				start, end = y+"-10-01", y+"-12-31"
			default:
				// Quarter tidak valid, fallback ke full year
				start, end = y+"-01-01", y+"-12-31"
			}
		} else {
			start, end = y+"-01-01", y+"-12-31"
		}
		return &start, &end
	}

	return nil, nil
}

// CountPerSektor returns total dan this_month counts grouped by sektor.
// Filter opsional: from/to, year, quarter, sub_sektor_id.
func (r *DashboardRepository) CountPerSektor(ctx context.Context, f dto.DashboardFilter) ([]dto.SectorCount, error) {
	from, to := buildDateRange(f)

	var args []interface{}
	var thisMonthExpr string
	if from != nil && to != nil {
		thisMonthExpr = "SUM(CASE WHEN p.created_at BETWEEN ? AND ? THEN 1 ELSE 0 END) AS this_month"
		args = append(args, *from, *to)
	} else {
		thisMonthExpr = "SUM(CASE WHEN p.created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS this_month"
	}

	// Filter sub_sektor opsional
	var whereClause string
	if f.SubSektorID != nil {
		whereClause = "WHERE p.id_sub_sektor = ?"
		args = append(args, *f.SubSektorID)
	}

	query := fmt.Sprintf(`
		SELECT s.id, s.nama_sektor, COUNT(p.id) AS total, %s
		FROM perusahaan p
		JOIN sub_sektor ss ON p.id_sub_sektor = ss.id
		JOIN sektor s ON ss.id_sektor = s.id
		%s
		GROUP BY s.id, s.nama_sektor
	`, thisMonthExpr, whereClause)

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

// TODO: re-enable ikas summary when ikas table is ready
// IkasGlobalAgg returns simple ikas aggregate across all data (for summary).
// If ikas table or columns missing, this returns zeros or an error.
// func (r *DashboardRepository) IkasGlobalAgg(ctx context.Context) (dto.IkasAgg, error) {
// 	var out dto.IkasAgg
// 	query := `SELECT COUNT(id) as total_ikas, AVG(nilai_kematangan) as avg_nilai_kematangan, AVG(target_nilai) as avg_target_nilai FROM ikas`
// 	row := r.db.QueryRowContext(ctx, query)
// 	if err := row.Scan(&out.Total, &out.AvgNilaiKematangan, &out.AvgTargetNilai); err != nil {
// 		if err == sql.ErrNoRows {
// 			return out, nil
// 		}
// 		return out, err
// 	}
// 	return out, nil
// }

// SeGlobalAgg returns se aggregate dengan breakdown per kategori dan this_month.
// Filter opsional: from/to, year, quarter, sub_sektor_id, kategori_se.
func (r *DashboardRepository) SeGlobalAgg(ctx context.Context, f dto.DashboardFilter) (dto.SeAgg, error) {
	var out dto.SeAgg

	from, to := buildDateRange(f)

	var args []interface{}
	var conditions []string

	// Filter sub_sektor
	if f.SubSektorID != nil {
		conditions = append(conditions, "id_sub_sektor = ?")
		args = append(args, *f.SubSektorID)
	}

	// Filter kategori SE
	if f.KategoriSE != nil {
		conditions = append(conditions, "kategori_se = ?")
		args = append(args, *f.KategoriSE)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Ekspresi this_month berdasarkan filter tanggal
	var thisMonthExpr string
	if from != nil && to != nil {
		thisMonthExpr = "SUM(CASE WHEN created_at BETWEEN ? AND ? THEN 1 ELSE 0 END)"
		// args untuk this_month ditambahkan setelah where args
		args = append(args, *from, *to)
	} else {
		thisMonthExpr = "SUM(CASE WHEN created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END)"
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(id)                                                       AS total_se,
			%s                                                              AS this_month,
			SUM(CASE WHEN kategori_se = 'Strategis' THEN 1 ELSE 0 END)    AS strategis,
			SUM(CASE WHEN kategori_se = 'Tinggi'    THEN 1 ELSE 0 END)    AS tinggi,
			SUM(CASE WHEN kategori_se = 'Rendah'    THEN 1 ELSE 0 END)    AS rendah
		FROM se
		%s
	`, thisMonthExpr, whereClause)

	row := r.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&out.TotalSE, &out.ThisMonth, &out.Strategis, &out.Tinggi, &out.Rendah); err != nil {
		if err == sql.ErrNoRows {
			return out, nil
		}
		return out, err
	}
	return out, nil
}

// SeStatusCount menghitung perusahaan yang sudah/belum mengisi KSE.
// Filter opsional: sub_sektor_id.
func (r *DashboardRepository) SeStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.SeStatusCount, error) {
	var out dto.SeStatusCount

	var args []interface{}
	var whereClause string
	if f.SubSektorID != nil {
		whereClause = "WHERE p.id_sub_sektor = ?"
		args = append(args, *f.SubSektorID)
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(p.id)                                         AS total_perusahaan,
			COUNT(se.id_perusahaan)                             AS sudah_mengisi_kse,
			COUNT(p.id) - COUNT(se.id_perusahaan)              AS belum_mengisi_kse
		FROM perusahaan p
		LEFT JOIN (
			SELECT DISTINCT id_perusahaan FROM se
		) se ON p.id = se.id_perusahaan
		%s
	`, whereClause)

	row := r.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&out.TotalPerusahaan, &out.SudahMengisiKSE, &out.BelumMengisiKSE); err != nil {
		if err == sql.ErrNoRows {
			return out, nil
		}
		return out, err
	}
	return out, nil
}

// TODO: re-enable ikas status when ikas table is ready
// IkasStatusCount menghitung perusahaan yang sudah/belum mengisi IKAS.
// func (r *DashboardRepository) IkasStatusCount(ctx context.Context, f dto.DashboardFilter) (dto.IkasStatusCount, error) {
// 	var out dto.IkasStatusCount
// 	var args []interface{}
// 	var whereClause string
// 	if f.SubSektorID != nil {
// 		whereClause = "WHERE p.id_sub_sektor = ?"
// 		args = append(args, *f.SubSektorID)
// 	}
// 	query := fmt.Sprintf(`
// 		SELECT
// 			COUNT(p.id)                                         AS total_perusahaan,
// 			COUNT(ik.id_perusahaan)                             AS sudah_mengisi_ikas,
// 			COUNT(p.id) - COUNT(ik.id_perusahaan)              AS belum_mengisi_ikas
// 		FROM perusahaan p
// 		LEFT JOIN (
// 			SELECT DISTINCT id_perusahaan FROM ikas
// 		) ik ON p.id = ik.id_perusahaan
// 		%s
// 	`, whereClause)
// 	row := r.db.QueryRowContext(ctx, query, args...)
// 	if err := row.Scan(&out.TotalPerusahaan, &out.SudahMengisiIKAS, &out.BelumMengisiIKAS); err != nil {
// 		if err == sql.ErrNoRows {
// 			return out, nil
// 		}
// 		return out, err
// 	}
// 	return out, nil
// }
