package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type EventRegistrationQRPayload struct {
	EventID        string `json:"event_id"`
	EventTitle     string `json:"event_title"`
	RegistrationID string `json:"registration_id"`
	Nama           string `json:"nama"`
	Email          string `json:"email"`
	Perusahaan     string `json:"perusahaan"`
	Jabatan        string `json:"jabatan"`
	NoHP           string `json:"no_hp"`
	Sektor         string `json:"sektor"`
	QRToken        string `json:"qr_token"`
}

func (p EventRegistrationQRPayload) JSON() (string, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func GenerateQRCodePNG(payload string, size int) ([]byte, error) {
	content := strings.TrimSpace(payload)
	if content == "" {
		return nil, fmt.Errorf("payload QR tidak boleh kosong")
	}
	if size <= 0 {
		size = 256
	}

	code, err := qr.Encode(content, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	scaled, err := barcode.Scale(code, size, size)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
