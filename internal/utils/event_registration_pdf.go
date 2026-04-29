package utils

import (
	"bytes"

	"fortyfour-backend/internal/models"

	"github.com/go-pdf/fpdf"
)

func GenerateEventRegistrationPDF(eventTitle string, reg *models.EventRegistration, qrPNG []byte) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.SetAutoPageBreak(true, pageMargin)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, toSafe("Registrasi Event"), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(0, 8, toSafe(eventTitle), "", 1, "C", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("Arial", "", 10)
	rows := [][]string{
		{"Nama", reg.Nama},
		{"Email", reg.Email},
		{"Perusahaan", reg.Perusahaan},
		{"Jabatan", reg.Jabatan},
		{"No HP", reg.NoHP},
		{"Sektor", reg.Sektor},
		{"QR Token", reg.QRToken},
	}
	printRows(pdf, rows)

	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, "QR Code", "", 1, "C", false, 0, "")

	opt := fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	pdf.RegisterImageOptionsReader("event-registration-qr", opt, bytes.NewReader(qrPNG))
	x := (210.0 - 70.0) / 2.0
	pdf.ImageOptions("event-registration-qr", x, pdf.GetY()+2, 70, 70, false, opt, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
