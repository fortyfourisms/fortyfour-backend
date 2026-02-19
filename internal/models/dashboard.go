package models

// Models for dashboard aggregations (not mapping to single DB table).
type DashboardSectorCount struct {
	ID         string
	NamaSektor string
	Total      int64
	ThisMonth  int64
}
