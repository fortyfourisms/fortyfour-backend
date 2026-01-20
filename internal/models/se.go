package models

import "time"

type SE struct {
	ID           string
	IDPerusahaan string
	IDSubSektor  string

	Q1  int
	Q2  int
	Q3  int
	Q4  int
	Q5  int
	Q6  int
	Q7  int
	Q8  int
	Q9  int
	Q10 int

	TotalBobot int
	KategoriSE string

	CreatedAt time.Time
	UpdatedAt time.Time
}