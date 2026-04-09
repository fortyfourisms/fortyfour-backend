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

// PublishUserCreated
func (p *Producer) PublishUserCreated(ctx context.Context, event dto_event.UserCreatedEvent) error {
	return p.Publish(ctx, "users.events", "users.created", event)
}

// PublishUserUpdated
func (p *Producer) PublishUserUpdated(ctx context.Context, event dto_event.UserUpdatedEvent) error {
	return p.Publish(ctx, "users.events", "users.updated", event)
}

// PublishUserDeleted
func (p *Producer) PublishUserDeleted(ctx context.Context, event dto_event.UserDeletedEvent) error {
	return p.Publish(ctx, "users.events", "users.deleted", event)
}

// PublishUserPasswordUpdated
func (p *Producer) PublishUserPasswordUpdated(ctx context.Context, event dto_event.UserPasswordUpdatedEvent) error {
	return p.Publish(ctx, "users.events", "users.password_updated", event)
}

// PublishCsirtCreated
func (p *Producer) PublishCsirtCreated(ctx context.Context, event dto_event.CsirtCreatedEvent) error {
	return p.Publish(ctx, "csirt.events", "csirt.created", event)
}

// PublishCsirtUpdated
func (p *Producer) PublishCsirtUpdated(ctx context.Context, event dto_event.CsirtUpdatedEvent) error {
	return p.Publish(ctx, "csirt.events", "csirt.updated", event)
}

// PublishCsirtDeleted
func (p *Producer) PublishCsirtDeleted(ctx context.Context, event dto_event.CsirtDeletedEvent) error {
	return p.Publish(ctx, "csirt.events", "csirt.deleted", event)
}
