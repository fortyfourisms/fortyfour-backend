package models

type Csirt struct {
	ID               string `json:"id"`
	IdPerusahaan     string `json:"id_perusahaan"`
	NamaCsirt        string `json:"nama_csirt"`
	WebCsirt         string `json:"web_csirt"`
	PhotoCsirt       string `json:"photo_csirt"`
	FileRFC2350      string `json:"file_rfc2350"`
	FilePublicKeyPGP string `json:"file_public_key_pgp"`
}
