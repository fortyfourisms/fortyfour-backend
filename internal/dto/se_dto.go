package dto

import "time"

// ===== REQUEST =====

type CreateSERequest struct {
	IDPerusahaan *string `json:"id_perusahaan" validate:"required"`
	IDSubSektor  *string `json:"id_sub_sektor" validate:"required"`

	Q1  string `json:"q1"` // A / B / C
	Q2  string `json:"q2"`
	Q3  string `json:"q3"`
	Q4  string `json:"q4"`
	Q5  string `json:"q5"`
	Q6  string `json:"q6"`
	Q7  string `json:"q7"`
	Q8  string `json:"q8"`
	Q9  string `json:"q9"`
	Q10 string `json:"q10"`
}

type UpdateSERequest struct {
	Q1  string `json:"q1"`
	Q2  string `json:"q2"`
	Q3  string `json:"q3"`
	Q4  string `json:"q4"`
	Q5  string `json:"q5"`
	Q6  string `json:"q6"`
	Q7  string `json:"q7"`
	Q8  string `json:"q8"`
	Q9  string `json:"q9"`
	Q10 string `json:"q10"`
}

// ===== RESPONSE =====

type SEResponse struct {
    ID         string    `json:"id"`
    Q1         string       `json:"q1"`
    Q2         string       `json:"q2"`
    Q3         string       `json:"q3"`
    Q4         string       `json:"q4"`
    Q5         string       `json:"q5"`
    Q6         string       `json:"q6"`
    Q7         string       `json:"q7"`
    Q8         string       `json:"q8"`
    Q9         string       `json:"q9"`
    Q10        string       `json:"q10"`
    TotalBobot string       `json:"total_bobot"`
    KategoriSE string    `json:"kategori_se"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`

    Perusahaan *PerusahaanMiniResponse `json:"perusahaan"`
    SubSektor  *SubSektorMiniResponse  `json:"sub_sektor"`
}

type PerusahaanMiniResponse struct {
    ID             string `json:"id"`
    NamaPerusahaan string `json:"nama_perusahaan"`
}

type SubSektorMiniResponse struct {
    ID             string `json:"id"`
    NamaSubSektor  string `json:"nama_sub_sektor"`
    IDSektor       string `json:"id_sektor"`
    NamaSektor     string `json:"nama_sektor"`
}
