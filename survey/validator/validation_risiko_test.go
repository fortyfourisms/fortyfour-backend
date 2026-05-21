package validation

import (
	"survey/internal/dto"
	"testing"
)

func stringPtr(v string) *string { return &v }

// STEP 1 — ELIGIBILITY
func TestValidateEligibilityRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.EligibilityRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: dto.EligibilityRequest{
				RisikoID: stringPtr("uuid-1"),
			},
			wantErr: nil,
		},
		{
			name: "missing responden_id is allowed for owned endpoint",
			req: dto.EligibilityRequest{
				RespondenID: "uuid-resp",
				RisikoID:    stringPtr("uuid-1"),
			},
			wantErr: nil,
		},
		{
			name: "missing risiko_id",
			req: dto.EligibilityRequest{
				RespondenID: "uuid-resp",
				RisikoID:    nil,
			},
			wantErr: ErrMissingRisikoID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEligibilityRequest(tt.req)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// STEP 2A — ALASAN
func TestValidateAlasanRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.AlasanRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: dto.AlasanRequest{
				RisikoID: stringPtr("uuid-1"),
				Alasan:   "Valid alasan",
			},
			wantErr: nil,
		},
		{
			name: "missing alasan",
			req: dto.AlasanRequest{
				RisikoID: stringPtr("uuid-1"),
				Alasan:   "",
			},
			wantErr: ErrMissingReason,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAlasanRequest(tt.req)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// STEP 2B — DAMPAK
func TestValidateDampakRequest(t *testing.T) {
	validReq := dto.DampakRequest{
		RisikoID:          stringPtr("uuid-1"),
		DampakReputasi:    2,
		DampakOperasional: 2,
		DampakFinansial:   2,
		DampakHukum:       2,
		Frekuensi:         2,
	}

	tests := []struct {
		name    string
		req     dto.DampakRequest
		wantErr error
	}{
		{
			name:    "valid request",
			req:     validReq,
			wantErr: nil,
		},
		{
			name: "invalid dampak reputasi",
			req: func() dto.DampakRequest {
				r := validReq
				r.DampakReputasi = 5
				return r
			}(),
			wantErr: ErrInvalidImpact,
		},
		{
			name: "invalid frekuensi",
			req: func() dto.DampakRequest {
				r := validReq
				r.Frekuensi = 0
				return r
			}(),
			wantErr: ErrInvalidFreq,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDampakRequest(tt.req)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// STEP 2C — PENGENDALIAN
func TestValidatePengendalianRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.PengendalianRequest
		wantErr error
	}{
		{
			name: "valid tanpa pengendalian",
			req: dto.PengendalianRequest{
				RisikoID:        stringPtr("uuid-1"),
				AdaPengendalian: false,
			},
			wantErr: nil,
		},
		{
			name: "valid dengan pengendalian",
			req: dto.PengendalianRequest{
				RisikoID:              stringPtr("uuid-1"),
				AdaPengendalian:       true,
				DeskripsiPengendalian: "Kontrol tersedia",
			},
			wantErr: nil,
		},
		{
			name: "missing deskripsi saat ada pengendalian",
			req: dto.PengendalianRequest{
				RisikoID:              stringPtr("uuid-1"),
				AdaPengendalian:       true,
				DeskripsiPengendalian: "",
			},
			wantErr: ErrMissingControl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePengendalianRequest(tt.req)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
