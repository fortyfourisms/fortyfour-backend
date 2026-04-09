package rabbitmq

import (
	"context"

	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/pkg/rabbitmq"
)

// Producer
type Producer struct {
	*rabbitmq.Producer
}

// NewProducer
func NewProducer(p *rabbitmq.Producer) *Producer {
	return &Producer{
		Producer: p,
	}
}

// Users
func (p *Producer) PublishUserCreated(ctx context.Context, event dto_event.UserCreatedEvent) error {
	return p.Publish(ctx, "users.events", "users.created", event)
}
func (p *Producer) PublishUserUpdated(ctx context.Context, event dto_event.UserUpdatedEvent) error {
	return p.Publish(ctx, "users.events", "users.updated", event)
}
func (p *Producer) PublishUserDeleted(ctx context.Context, event dto_event.UserDeletedEvent) error {
	return p.Publish(ctx, "users.events", "users.deleted", event)
}
func (p *Producer) PublishUserPasswordUpdated(ctx context.Context, event dto_event.UserPasswordUpdatedEvent) error {
	return p.Publish(ctx, "users.events", "users.password_updated", event)
}

// Csirt
func (p *Producer) PublishCsirtCreated(ctx context.Context, event dto_event.CsirtCreatedEvent) error {
	return p.Publish(ctx, "csirt.events", "csirt.created", event)
}
func (p *Producer) PublishCsirtUpdated(ctx context.Context, event dto_event.CsirtUpdatedEvent) error {
	return p.Publish(ctx, "csirt.events", "csirt.updated", event)
}
func (p *Producer) PublishCsirtDeleted(ctx context.Context, event dto_event.CsirtDeletedEvent) error {
	return p.Publish(ctx, "csirt.events", "csirt.deleted", event)
}

// Perusahaan
func (p *Producer) PublishPerusahaanCreated(ctx context.Context, event dto_event.PerusahaanCreatedEvent) error {
	return p.Publish(ctx, "perusahaan.events", "perusahaan.created", event)
}
func (p *Producer) PublishPerusahaanUpdated(ctx context.Context, event dto_event.PerusahaanUpdatedEvent) error {
	return p.Publish(ctx, "perusahaan.events", "perusahaan.updated", event)
}
func (p *Producer) PublishPerusahaanDeleted(ctx context.Context, event dto_event.PerusahaanDeletedEvent) error {
	return p.Publish(ctx, "perusahaan.events", "perusahaan.deleted", event)
}

// PIC
func (p *Producer) PublishPicCreated(ctx context.Context, event dto_event.PicCreatedEvent) error {
	return p.Publish(ctx, "pic.events", "pic.created", event)
}
func (p *Producer) PublishPicUpdated(ctx context.Context, event dto_event.PicUpdatedEvent) error {
	return p.Publish(ctx, "pic.events", "pic.updated", event)
}
func (p *Producer) PublishPicDeleted(ctx context.Context, event dto_event.PicDeletedEvent) error {
	return p.Publish(ctx, "pic.events", "pic.deleted", event)
}

// Jabatan
func (p *Producer) PublishJabatanCreated(ctx context.Context, event dto_event.JabatanCreatedEvent) error {
	return p.Publish(ctx, "jabatan.events", "jabatan.created", event)
}
func (p *Producer) PublishJabatanUpdated(ctx context.Context, event dto_event.JabatanUpdatedEvent) error {
	return p.Publish(ctx, "jabatan.events", "jabatan.updated", event)
}
func (p *Producer) PublishJabatanDeleted(ctx context.Context, event dto_event.JabatanDeletedEvent) error {
	return p.Publish(ctx, "jabatan.events", "jabatan.deleted", event)
}

// SDM CSIRT
func (p *Producer) PublishSdmCsirtCreated(ctx context.Context, event dto_event.SdmCsirtCreatedEvent) error {
	return p.Publish(ctx, "sdm_csirt.events", "sdm_csirt.created", event)
}
func (p *Producer) PublishSdmCsirtUpdated(ctx context.Context, event dto_event.SdmCsirtUpdatedEvent) error {
	return p.Publish(ctx, "sdm_csirt.events", "sdm_csirt.updated", event)
}
func (p *Producer) PublishSdmCsirtDeleted(ctx context.Context, event dto_event.SdmCsirtDeletedEvent) error {
	return p.Publish(ctx, "sdm_csirt.events", "sdm_csirt.deleted", event)
}

// Role
func (p *Producer) PublishRoleCreated(ctx context.Context, event dto_event.RoleCreatedEvent) error {
	return p.Publish(ctx, "role.events", "role.created", event)
}
func (p *Producer) PublishRoleUpdated(ctx context.Context, event dto_event.RoleUpdatedEvent) error {
	return p.Publish(ctx, "role.events", "role.updated", event)
}
func (p *Producer) PublishRoleDeleted(ctx context.Context, event dto_event.RoleDeletedEvent) error {
	return p.Publish(ctx, "role.events", "role.deleted", event)
}

// SE
func (p *Producer) PublishSeCreated(ctx context.Context, event dto_event.SeCreatedEvent) error {
	return p.Publish(ctx, "se.events", "se.created", event)
}
func (p *Producer) PublishSeUpdated(ctx context.Context, event dto_event.SeUpdatedEvent) error {
	return p.Publish(ctx, "se.events", "se.updated", event)
}
func (p *Producer) PublishSeDeleted(ctx context.Context, event dto_event.SeDeletedEvent) error {
	return p.Publish(ctx, "se.events", "se.deleted", event)
}
