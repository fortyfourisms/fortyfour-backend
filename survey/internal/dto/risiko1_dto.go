package dto

import "survey/internal/models"

type EligibilityRequest struct {
	RespondenID   int  `json:"responden_id"`
	RisikoID      int  `json:"risiko_id"`
	PernahTerjadi bool `json:"pernah_terjadi"`
}

type AlasanRequest struct {
	RespondenID int    `json:"responden_id"`
	RisikoID    int    `json:"risiko_id"`
	Alasan      string `json:"alasan"`
}

type DampakRequest struct {
	RespondenID       int                   `json:"responden_id"`
	RisikoID          int                   `json:"risiko_id"`
	DampakReputasi    models.ImpactLevel    `json:"dampak_reputasi"`   
	DampakOperasional models.ImpactLevel    `json:"dampak_operasional"`
	DampakFinansial   models.ImpactLevel    `json:"dampak_finansial"`
	DampakHukum       models.ImpactLevel    `json:"dampak_hukum"`
	Frekuensi         models.FrequencyLevel `json:"frekuensi"`
}

type PengendalianRequest struct {
	RespondenID           int    `json:"responden_id"`
	RisikoID              int    `json:"risiko_id"`
	AdaPengendalian       bool   `json:"ada_pengendalian"`
	DeskripsiPengendalian string `json:"deskripsi_pengendalian,omitempty"` 
}

type NavigateRequest struct {             
	RespondenID int    `json:"responden_id"`
	Direction   string `json:"direction"`
	CurrentRisk int    `json:"current_risk"`
}