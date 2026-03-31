package dto

type CreateCsirtRequest struct {
	IdPerusahaan           string `json:"id_perusahaan"`
	NamaCsirt              string `json:"nama_csirt"`
	WebCsirt               string `json:"web_csirt"`
	TeleponCsirt           string `json:"telepon_csirt"`
	PhotoCsirt             string `json:"photo_csirt"`
	FileRFC2350            string `json:"file_rfc2350"`
	FilePublicKeyPGP       string `json:"file_public_key_pgp"`
	FileStr                string `json:"file_str"`
	TanggalRegistrasi      string `json:"tanggal_registrasi"`
	TanggalKadaluarsa      string `json:"tanggal_kadaluarsa"`
	TanggalRegistrasiUlang string `json:"tanggal_registrasi_ulang"`
}

type UpdateCsirtRequest struct {
	NamaCsirt              *string `json:"nama_csirt,omitempty"`
	WebCsirt               *string `json:"web_csirt,omitempty"`
	TeleponCsirt           *string `json:"telepon_csirt,omitempty"`
	PhotoCsirt             *string `json:"photo_csirt,omitempty"`
	FileRFC2350            *string `json:"file_rfc2350,omitempty"`
	FilePublicKeyPGP       *string `json:"file_public_key_pgp,omitempty"`
	FileStr                *string `json:"file_str,omitempty"`
	TanggalRegistrasi      *string `json:"tanggal_registrasi,omitempty"`
	TanggalKadaluarsa      *string `json:"tanggal_kadaluarsa,omitempty"`
	TanggalRegistrasiUlang *string `json:"tanggal_registrasi_ulang,omitempty"`
}

type CsirtResponse struct {
	ID                     string             `json:"id"`
	NamaCsirt              string             `json:"nama_csirt"`
	WebCsirt               string             `json:"web_csirt"`
	TeleponCsirt           *string            `json:"telepon_csirt"`
	PhotoCsirt             string             `json:"photo_csirt"`
	FileRFC2350            string             `json:"file_rfc2350"`
	FilePublicKeyPGP       string             `json:"file_public_key_pgp"`
	FileStr                string             `json:"file_str"`
	TanggalRegistrasi      string             `json:"tanggal_registrasi"`
	TanggalKadaluarsa      string             `json:"tanggal_kadaluarsa"`
	TanggalRegistrasiUlang string             `json:"tanggal_registrasi_ulang"`
	Perusahaan             PerusahaanResponse `json:"perusahaan"`
}

/*
========================================

	MINI RESPONSE (DIPAKAI SDM & SE)

========================================
*/
type CsirtMiniResponse struct {
	ID           string  `json:"id"`
	NamaCsirt    string  `json:"nama_csirt"`
	WebCsirt     *string `json:"web_csirt,omitempty"`
	TeleponCsirt *string `json:"telepon_csirt,omitempty"`
}
