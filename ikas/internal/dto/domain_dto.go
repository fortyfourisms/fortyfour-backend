package dto

import "time"

type CreateDomainRequest struct {
	NamaDomain string `json:"nama_domain"`
}

type UpdateDomainRequest struct {
	NamaDomain *string `json:"nama_domain,omitempty"`
}

type DomainResponse struct {
	ID         string    `json:"id"`
	NamaDomain string    `json:"nama_domain"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
