package rabbitmq

import (
	"context"
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"log"

	"fortyfour-backend/pkg/rabbitmq"
)

type Consumer struct {
	*rabbitmq.Consumer
	ikasRepo                   repository.IkasRepositoryInterface
	jawabanIdentifikasiRepo    repository.JawabanIdentifikasiRepositoryInterface
	pertanyaanIdentifikasiRepo repository.PertanyaanIdentifikasiRepositoryInterface
}

func NewConsumer(
	c *rabbitmq.Consumer,
	ikasRepo repository.IkasRepositoryInterface,
	jawabanIdentifikasiRepo repository.JawabanIdentifikasiRepositoryInterface,
	pertanyaanIdentifikasiRepo repository.PertanyaanIdentifikasiRepositoryInterface,
) *Consumer {
	return &Consumer{
		Consumer:                   c,
		ikasRepo:                   ikasRepo,
		jawabanIdentifikasiRepo:    jawabanIdentifikasiRepo,
		pertanyaanIdentifikasiRepo: pertanyaanIdentifikasiRepo,
	}
}

func (c *Consumer) ConsumeIkasCreated(ctx context.Context) error {
	return c.Consume(ctx, "ikas.created", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing IKAS Created for ID: %s", event.IkasID)

		req := dto.CreateIkasRequest{
			IDPerusahaan: event.IDPerusahaan,
			Tanggal:      event.Tanggal,
			Responden:    event.Responden,
			Telepon:      event.Telepon,
			Jabatan:      event.Jabatan,
			TargetNilai:  event.TargetNilai,
		}

		return c.ikasRepo.Create(req, event.IkasID, event.NilaiKematangan)
	})
}

func (c *Consumer) ConsumeIkasUpdated(ctx context.Context) error {
	return c.Consume(ctx, "ikas.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing IKAS Updated for ID: %s", event.IkasID)

		req := dto.UpdateIkasRequest{
			IDPerusahaan: &event.IDPerusahaan,
			Tanggal:      &event.Tanggal,
			Responden:    &event.Responden,
			Telepon:      &event.Telepon,
			Jabatan:      &event.Jabatan,
			TargetNilai:  &event.TargetNilai,
		}

		return c.ikasRepo.Update(event.IkasID, req)
	})
}

func (c *Consumer) ConsumeIkasDeleted(ctx context.Context) error {
	return c.Consume(ctx, "ikas.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing IKAS Deleted for ID: %s", event.IkasID)

		return c.ikasRepo.Delete(event.IkasID)
	})
}

func (c *Consumer) ConsumeIkasImported(ctx context.Context) error {
	return c.Consume(ctx, "ikas.imported", func(ctx context.Context, body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Imported Event: %+v", event)
		return nil
	})
}

func (c *Consumer) ConsumeEmailNotifications(ctx context.Context) error {
	return c.Consume(ctx, "notifications.email", func(ctx context.Context, body []byte) error {
		var notification map[string]interface{}
		if err := json.Unmarshal(body, &notification); err != nil {
			return err
		}

		log.Printf("Email Notification Request: %+v", notification)
		return nil
	})
}

func (c *Consumer) ConsumeJawabanIdentifikasiCreated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.identifikasi.created", func(ctx context.Context, body []byte) error {
		var req dto.CreateJawabanIdentifikasiRequest
		if err := json.Unmarshal(body, &req); err != nil {
			return err
		}

		log.Printf("Buffering Jawaban Identifikasi for Perusahaan: %s, Question: %d", req.PerusahaanID, req.PertanyaanIdentifikasiID)

		// 1. Save to buffer
		if err := c.jawabanIdentifikasiRepo.UpsertToBuffer(req); err != nil {
			log.Printf("Error upserting to buffer: %v", err)
			return err
		}

		// 2. Check if all questions are answered
		totalQuestions, err := c.pertanyaanIdentifikasiRepo.GetTotalCount()
		if err != nil {
			log.Printf("Error getting total questions: %v", err)
			return err
		}

		currentCount, err := c.jawabanIdentifikasiRepo.GetBufferCount(req.PerusahaanID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for Perusahaan %s (%d/%d). Flushing buffer...", req.PerusahaanID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanIdentifikasiRepo.FlushBuffer(req.PerusahaanID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for Perusahaan %s", req.PerusahaanID)
			return c.jawabanIdentifikasiRepo.RecalculateIdentifikasi(req.PerusahaanID)
		}

		log.Printf("Progress for Perusahaan %s: %d/%d", req.PerusahaanID, currentCount, totalQuestions)
		return nil
	})
}

// ConsumeJawabanIdentifikasiUpdated (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanIdentifikasiUpdated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.identifikasi.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanIdentifikasiUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Identifikasi Updated for ID: %d", event.ID)

		// 1. Update database
		if err := c.jawabanIdentifikasiRepo.Update(event.ID, event.Request); err != nil {
			return err
		}

		// 2. Get PerusahaanID to recalculate
		resp, err := c.jawabanIdentifikasiRepo.GetByID(event.ID)
		if err != nil {
			return err
		}

		// 3. Recalculate scores
		return c.jawabanIdentifikasiRepo.RecalculateIdentifikasi(resp.PerusahaanID)
	})
}

// ConsumeJawabanIdentifikasiDeleted (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanIdentifikasiDeleted(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.identifikasi.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanIdentifikasiDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Identifikasi Deleted for ID: %d", event.ID)

		// 1. Delete from database
		if err := c.jawabanIdentifikasiRepo.Delete(event.ID); err != nil {
			return err
		}

		// 2. Recalculate scores
		return c.jawabanIdentifikasiRepo.RecalculateIdentifikasi(event.PerusahaanID)
	})
}

func (c *Consumer) StartAllConsumers(ctx context.Context) error {
	consumers := []func(context.Context) error{
		c.ConsumeIkasCreated,
		c.ConsumeIkasUpdated,
		c.ConsumeIkasDeleted,
		c.ConsumeIkasImported,
		c.ConsumeEmailNotifications,
		c.ConsumeJawabanIdentifikasiCreated,
		c.ConsumeJawabanIdentifikasiUpdated,
		c.ConsumeJawabanIdentifikasiDeleted,
	}

	for _, consumer := range consumers {
		if err := consumer(ctx); err != nil {
			return err
		}
	}

	log.Println("All consumers started successfully")
	return nil
}
