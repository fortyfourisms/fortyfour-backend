package dto

type CreateCsirtRequest struct {
	IdPerusahaan     string `json:"id_perusahaan"`
	NamaCsirt        string `json:"nama_csirt"`
	WebCsirt         string `json:"web_csirt"`
	FileRFC2350      string `json:"file_rfc2350"`
	FilePublicKeyPGP string `json:"file_public_key_pgp"`
}

type UpdateCsirtRequest struct {
	NamaCsirt        *string `json:"nama_csirt,omitempty"`
	WebCsirt         *string `json:"web_csirt,omitempty"`
	FileRFC2350      *string `json:"file_rfc2350,omitempty"`
	FilePublicKeyPGP *string `json:"file_public_key_pgp,omitempty"`
}
