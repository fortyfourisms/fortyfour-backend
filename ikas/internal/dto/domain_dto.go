package dto

import "time"

type CreateDomainRequest struct {
	NamaDomain string `json:"nama_domain"`
}

type UpdateDomainRequest struct {
	NamaDomain *string `json:"nama_domain,omitempty"`
}

type DomainResponse struct {
	ID         int       `json:"id"`
	NamaDomain string    `json:"nama_domain"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DomainMessageResponse struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message"`
}
