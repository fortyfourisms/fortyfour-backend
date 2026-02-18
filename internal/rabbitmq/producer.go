package rabbitmq

import (
	"context"

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
func (p *Producer) PublishUserCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "users.events", "users.created", event)
}

// PublishUserUpdated
func (p *Producer) PublishUserUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "users.events", "users.updated", event)
}

// PublishUserDeleted
func (p *Producer) PublishUserDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "users.events", "users.deleted", event)
}

// PublishUserPasswordUpdated
func (p *Producer) PublishUserPasswordUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "users.events", "users.password_updated", event)
}
