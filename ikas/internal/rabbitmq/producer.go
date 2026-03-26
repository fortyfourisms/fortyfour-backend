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

// PublishJawabanProteksiCreated
func (p *Producer) PublishJawabanProteksiCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.proteksi.created", event)
}

// PublishJawabanProteksiUpdated
func (p *Producer) PublishJawabanProteksiUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.proteksi.updated", event)
}

// PublishJawabanProteksiDeleted
func (p *Producer) PublishJawabanProteksiDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.proteksi.deleted", event)
}

// PublishJawabanDeteksiCreated
func (p *Producer) PublishJawabanDeteksiCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.deteksi.created", event)
}

// PublishJawabanDeteksiUpdated
func (p *Producer) PublishJawabanDeteksiUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.deteksi.updated", event)
}

// PublishJawabanDeteksiDeleted
func (p *Producer) PublishJawabanDeteksiDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.deteksi.deleted", event)
}

// PublishJawabanGulihCreated
func (p *Producer) PublishJawabanGulihCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.gulih.created", event)
}

// PublishJawabanGulihUpdated
func (p *Producer) PublishJawabanGulihUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.gulih.updated", event)
}

// PublishJawabanGulihDeleted
func (p *Producer) PublishJawabanGulihDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "jawaban.gulih.deleted", event)
}

// PublishDomainCreated
func (p *Producer) PublishDomainCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "domain.created", event)
}

// PublishDomainUpdated
func (p *Producer) PublishDomainUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "domain.updated", event)
}

// PublishDomainDeleted
func (p *Producer) PublishDomainDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "domain.deleted", event)
}

func (p *Producer) PublishRuangLingkupCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "ruang_lingkup.created", event)
}

func (p *Producer) PublishRuangLingkupUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "ruang_lingkup.updated", event)
}

func (p *Producer) PublishRuangLingkupDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "ruang_lingkup.deleted", event)
}

func (p *Producer) PublishKategoriCreated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "kategori.created", event)
}

func (p *Producer) PublishKategoriUpdated(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "kategori.updated", event)
}

func (p *Producer) PublishKategoriDeleted(ctx context.Context, event interface{}) error {
	return p.Publish(ctx, "ikas.events", "kategori.deleted", event)
}

// PublishIkasAuditLog
func (p *Producer) PublishIkasAuditLog(ctx context.Context, log interface{}) error {
	return p.Publish(ctx, "ikas.events", "ikas.audit.log", log)
}
