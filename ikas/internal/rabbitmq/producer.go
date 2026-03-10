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

// PublishIkasCreated
func (p *Producer) PublishIkasCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "ikas.created", event)
}

// PublishIkasUpdated
func (p *Producer) PublishIkasUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "ikas.updated", event)
}

// PublishIkasDeleted
func (p *Producer) PublishIkasDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "ikas.deleted", event)
}

// PublishIkasImported
func (p *Producer) PublishIkasImported(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "ikas.imported", event)
}

// PublishEmailNotification
func (p *Producer) PublishEmailNotification(ctx context.Context, notification interface{}) error {
	return p.Publish(ctx, "ikas.events", "notification.email", notification)
}

// PublishJawabanIdentifikasiCreated
func (p *Producer) PublishJawabanIdentifikasiCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.identifikasi.created", event)
}

// PublishJawabanIdentifikasiUpdated
func (p *Producer) PublishJawabanIdentifikasiUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.identifikasi.updated", event)
}

// PublishJawabanIdentifikasiDeleted
func (p *Producer) PublishJawabanIdentifikasiDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.identifikasi.deleted", event)
}
