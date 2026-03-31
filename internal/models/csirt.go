package models

type Csirt struct {
	ID                     string  `json:"id"`
	IdPerusahaan           string  `json:"id_perusahaan"`
	NamaCsirt              string  `json:"nama_csirt"`
	WebCsirt               string  `json:"web_csirt"`
	TeleponCsirt           *string `json:"telepon_csirt"`
	PhotoCsirt             *string `json:"photo_csirt"`
	FileRFC2350            *string `json:"file_rfc2350"`
	FilePublicKeyPGP       *string `json:"file_public_key_pgp"`
	FileStr                *string `json:"file_str"`
	TanggalRegistrasi      *string `json:"tanggal_registrasi"`
	TanggalKadaluarsa      *string `json:"tanggal_kadaluarsa"`
	TanggalRegistrasiUlang *string `json:"tanggal_registrasi_ulang"`
}
