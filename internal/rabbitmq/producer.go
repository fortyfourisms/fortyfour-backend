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
