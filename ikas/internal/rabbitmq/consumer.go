package rabbitmq

import (
	"context"
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/dto/dto_event"
	"ikas/internal/repository"
	"log"
	"strings"

	"fortyfour-backend/pkg/rabbitmq"
)

type Consumer struct {
	*rabbitmq.Consumer
	ikasRepo                   repository.IkasRepositoryInterface
	jawabanIdentifikasiRepo    repository.JawabanIdentifikasiRepositoryInterface
	pertanyaanIdentifikasiRepo repository.PertanyaanIdentifikasiRepositoryInterface
	jawabanProteksiRepo        repository.JawabanProteksiRepositoryInterface
	pertanyaanProteksiRepo     repository.PertanyaanProteksiRepositoryInterface
	jawabanDeteksiRepo         repository.JawabanDeteksiRepositoryInterface
	pertanyaanDeteksiRepo      repository.PertanyaanDeteksiRepositoryInterface
	jawabanGulihRepo           repository.JawabanGulihRepositoryInterface
	pertanyaanGulihRepo        repository.PertanyaanGulihRepositoryInterface
	domainRepo                 repository.DomainRepositoryInterface
	ruangLingkupRepo           repository.RuangLingkupRepositoryInterface
	kategoriRepo               repository.KategoriRepositoryInterface
	auditLogRepo               repository.AuditLogRepositoryInterface
}

func NewConsumer(
	c *rabbitmq.Consumer,
	ikasRepo repository.IkasRepositoryInterface,
	jawabanIdentifikasiRepo repository.JawabanIdentifikasiRepositoryInterface,
	pertanyaanIdentifikasiRepo repository.PertanyaanIdentifikasiRepositoryInterface,
	jawabanProteksiRepo repository.JawabanProteksiRepositoryInterface,
	pertanyaanProteksiRepo repository.PertanyaanProteksiRepositoryInterface,
	jawabanDeteksiRepo repository.JawabanDeteksiRepositoryInterface,
	pertanyaanDeteksiRepo repository.PertanyaanDeteksiRepositoryInterface,
	jawabanGulihRepo repository.JawabanGulihRepositoryInterface,
	pertanyaanGulihRepo repository.PertanyaanGulihRepositoryInterface,
	domainRepo repository.DomainRepositoryInterface,
	ruangLingkupRepo repository.RuangLingkupRepositoryInterface,
	kategoriRepo repository.KategoriRepositoryInterface,
	auditLogRepo repository.AuditLogRepositoryInterface,
) *Consumer {
	return &Consumer{
		Consumer:                   c,
		ikasRepo:                   ikasRepo,
		jawabanIdentifikasiRepo:    jawabanIdentifikasiRepo,
		pertanyaanIdentifikasiRepo: pertanyaanIdentifikasiRepo,
		jawabanProteksiRepo:        jawabanProteksiRepo,
		pertanyaanProteksiRepo:     pertanyaanProteksiRepo,
		jawabanDeteksiRepo:         jawabanDeteksiRepo,
		pertanyaanDeteksiRepo:      pertanyaanDeteksiRepo,
		jawabanGulihRepo:           jawabanGulihRepo,
		pertanyaanGulihRepo:        pertanyaanGulihRepo,
		domainRepo:                 domainRepo,
		ruangLingkupRepo:           ruangLingkupRepo,
		kategoriRepo:               kategoriRepo,
		auditLogRepo:               auditLogRepo,
	}
}

func (c *Consumer) ConsumeIkasCreated(ctx context.Context) error {
	return c.Consume(ctx, "ikas.created", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from ikas.created: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
		}

		log.Printf("Processing IKAS Created for ID: %s", event.IkasID)

		// Validate mandatory fields (Poison Pill Prevention)
		if strings.TrimSpace(event.IDPerusahaan) == "" || strings.TrimSpace(event.Tanggal) == "" {
			log.Printf("❌ Skipping invalid message from ikas.created: id_perusahaan or tanggal is empty. ID: %s", event.IkasID)
			return nil // Acknowledge to remove from queue
		}

		req := dto.CreateIkasRequest{
			IDPerusahaan: event.IDPerusahaan,
			Tanggal:      event.Tanggal,
			Responden:    event.Responden,
			Telepon:      event.Telepon,
			Jabatan:      event.Jabatan,
			TargetNilai:  event.TargetNilai,
		}

		if err := c.ikasRepo.Create(req, event.IkasID, event.NilaiKematangan); err != nil {
			if strings.Contains(err.Error(), "Incorrect datetime value") {
				log.Printf("❌ Skipping message from ikas.created due to invalid date value: %v", err)
				return nil // Acknowledge to remove poison pill
			}
			return err
		}
		return nil
	})
}

func (c *Consumer) ConsumeIkasUpdated(ctx context.Context) error {
	return c.Consume(ctx, "ikas.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from ikas.updated: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
		}

		log.Printf("Processing IKAS Updated for ID: %s", event.IkasID)

		// Publish event is already validated at API level.
		// For updates, we allow partial fields.

		req := dto.UpdateIkasRequest{
			IDPerusahaan: event.IDPerusahaan,
			Tanggal:      event.Tanggal,
			Responden:    event.Responden,
			Telepon:      event.Telepon,
			Jabatan:      event.Jabatan,
			TargetNilai:  event.TargetNilai,
		}

		if err := c.ikasRepo.Update(event.IkasID, req); err != nil {
			if strings.Contains(err.Error(), "Incorrect datetime value") {
				log.Printf("❌ Skipping message from ikas.updated due to invalid date value: %v", err)
				return nil // Acknowledge to remove poison pill
			}
			return err
		}
		return nil
	})
}

func (c *Consumer) ConsumeIkasDeleted(ctx context.Context) error {
	return c.Consume(ctx, "ikas.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from ikas.deleted: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
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
		log.Printf("Raw Message from jawaban.identifikasi.created: %s", string(body))

		var req dto.CreateJawabanIdentifikasiRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from jawaban.identifikasi.created: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
		}

		// 1. Validate mandatory fields (Poison Pill Prevention)
		if req.JawabanIdentifikasi == nil {
			log.Printf("❌ Skipping invalid message from jawaban.identifikasi.created: jawaban_identifikasi is null. Perusahaan: %s, Question: %d", req.PerusahaanID, req.PertanyaanIdentifikasiID)
			return nil // Acknowledge to remove from queue
		}

		// 2. Save to buffer
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

// ConsumeJawabanProteksiCreated (Pola 2 Batch Write)
func (c *Consumer) ConsumeJawabanProteksiCreated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.proteksi.created", func(ctx context.Context, body []byte) error {
		log.Printf("Raw Message from jawaban.proteksi.created: %s", string(body))

		var req dto.CreateJawabanProteksiRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from jawaban.proteksi.created: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
		}

		// 1. Validate mandatory fields (Poison Pill Prevention)
		if req.JawabanProteksi == nil {
			log.Printf("❌ Skipping invalid message from jawaban.proteksi.created: jawaban_proteksi is null. Perusahaan: %s, Question: %d", req.PerusahaanID, req.PertanyaanProteksiID)
			return nil // Acknowledge to remove from queue
		}

		// 2. Save to buffer
		if err := c.jawabanProteksiRepo.UpsertToBuffer(req); err != nil {
			log.Printf("Error upserting to buffer: %v", err)
			return err
		}

		// 2. Check if all questions are answered
		totalQuestions, err := c.pertanyaanProteksiRepo.GetTotalCount()
		if err != nil {
			log.Printf("Error getting total questions: %v", err)
			return err
		}

		currentCount, err := c.jawabanProteksiRepo.GetBufferCount(req.PerusahaanID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for Perusahaan %s (%d/%d). Flushing buffer...", req.PerusahaanID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanProteksiRepo.FlushBuffer(req.PerusahaanID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for Perusahaan %s", req.PerusahaanID)
			return c.jawabanProteksiRepo.RecalculateProteksi(req.PerusahaanID)
		}

		log.Printf("Progress for Perusahaan %s: %d/%d", req.PerusahaanID, currentCount, totalQuestions)
		return nil
	})
}

// ConsumeJawabanProteksiUpdated (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanProteksiUpdated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.proteksi.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanProteksiUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Proteksi Updated for ID: %d", event.ID)

		// 1. Update database
		if err := c.jawabanProteksiRepo.Update(event.ID, event.Request); err != nil {
			return err
		}

		// 2. Get PerusahaanID to recalculate
		resp, err := c.jawabanProteksiRepo.GetByID(event.ID)
		if err != nil {
			return err
		}

		// 3. Recalculate scores
		return c.jawabanProteksiRepo.RecalculateProteksi(resp.PerusahaanID)
	})
}

// ConsumeJawabanProteksiDeleted (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanProteksiDeleted(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.proteksi.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanProteksiDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Proteksi Deleted for ID: %d", event.ID)

		// 1. Delete from database
		if err := c.jawabanProteksiRepo.Delete(event.ID); err != nil {
			return err
		}

		// 2. Recalculate scores
		return c.jawabanProteksiRepo.RecalculateProteksi(event.PerusahaanID)
	})
}

// ConsumeJawabanDeteksiCreated (Pola 2 Batch Write)
func (c *Consumer) ConsumeJawabanDeteksiCreated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.deteksi.created", func(ctx context.Context, body []byte) error {
		log.Printf("Raw Message from jawaban.deteksi.created: %s", string(body))

		var req dto.CreateJawabanDeteksiRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from jawaban.deteksi.created: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
		}

		// 1. Validate mandatory fields (Poison Pill Prevention)
		if req.JawabanDeteksi == nil {
			log.Printf("❌ Skipping invalid message from jawaban.deteksi.created: jawaban_deteksi is null. Perusahaan: %s, Question: %d", req.PerusahaanID, req.PertanyaanDeteksiID)
			return nil // Acknowledge to remove from queue
		}

		// 2. Save to buffer
		if err := c.jawabanDeteksiRepo.UpsertToBuffer(req); err != nil {
			log.Printf("Error upserting to buffer: %v", err)
			return err
		}

		// 2. Check if all questions are answered
		totalQuestions, err := c.pertanyaanDeteksiRepo.GetTotalCount()
		if err != nil {
			log.Printf("Error getting total questions: %v", err)
			return err
		}

		currentCount, err := c.jawabanDeteksiRepo.GetBufferCount(req.PerusahaanID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for Perusahaan %s (%d/%d). Flushing buffer...", req.PerusahaanID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanDeteksiRepo.FlushBuffer(req.PerusahaanID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for Perusahaan %s", req.PerusahaanID)
			return c.jawabanDeteksiRepo.RecalculateDeteksi(req.PerusahaanID)
		}

		log.Printf("Progress for Perusahaan %s: %d/%d", req.PerusahaanID, currentCount, totalQuestions)
		return nil
	})
}

// ConsumeJawabanDeteksiUpdated (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanDeteksiUpdated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.deteksi.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanDeteksiUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Deteksi Updated for ID: %d", event.ID)

		// 1. Update database
		if err := c.jawabanDeteksiRepo.Update(event.ID, event.Request); err != nil {
			return err
		}

		// 2. Get PerusahaanID to recalculate
		resp, err := c.jawabanDeteksiRepo.GetByID(event.ID)
		if err != nil {
			return err
		}

		// 3. Recalculate scores
		return c.jawabanDeteksiRepo.RecalculateDeteksi(resp.PerusahaanID)
	})
}

// ConsumeJawabanDeteksiDeleted (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanDeteksiDeleted(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.deteksi.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanDeteksiDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Deteksi Deleted for ID: %d", event.ID)

		// 1. Delete from database
		if err := c.jawabanDeteksiRepo.Delete(event.ID); err != nil {
			return err
		}

		// 2. Recalculate scores
		return c.jawabanDeteksiRepo.RecalculateDeteksi(event.PerusahaanID)
	})
}

// ConsumeJawabanGulihCreated (Pola 2 Batch Write)
func (c *Consumer) ConsumeJawabanGulihCreated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.gulih.created", func(ctx context.Context, body []byte) error {
		log.Printf("Raw Message from jawaban.gulih.created: %s", string(body))

		var req dto.CreateJawabanGulihRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from jawaban.gulih.created: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
		}

		// 1. Validate mandatory fields (Poison Pill Prevention)
		if req.JawabanGulih == nil {
			log.Printf("❌ Skipping invalid message from jawaban.gulih.created: jawaban_gulih is null. Perusahaan: %s, Question: %d", req.PerusahaanID, req.PertanyaanGulihID)
			return nil // Acknowledge to remove from queue
		}

		// 2. Save to buffer
		if err := c.jawabanGulihRepo.UpsertToBuffer(req); err != nil {
			log.Printf("Error upserting to buffer: %v", err)
			return err
		}

		// 2. Check if all questions are answered
		totalQuestions, err := c.pertanyaanGulihRepo.GetTotalCount()
		if err != nil {
			log.Printf("Error getting total questions: %v", err)
			return err
		}

		currentCount, err := c.jawabanGulihRepo.GetBufferCount(req.PerusahaanID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for Perusahaan %s (%d/%d). Flushing buffer...", req.PerusahaanID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanGulihRepo.FlushBuffer(req.PerusahaanID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for Perusahaan %s", req.PerusahaanID)
			return c.jawabanGulihRepo.RecalculateGulih(req.PerusahaanID)
		}

		log.Printf("Progress for Perusahaan %s: %d/%d", req.PerusahaanID, currentCount, totalQuestions)
		return nil
	})
}

// ConsumeJawabanGulihUpdated (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanGulihUpdated(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.gulih.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanGulihUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Gulih Updated for ID: %d", event.ID)

		// 1. Update database
		if err := c.jawabanGulihRepo.Update(event.ID, event.Request); err != nil {
			return err
		}

		// 2. Get PerusahaanID to recalculate
		resp, err := c.jawabanGulihRepo.GetByID(event.ID)
		if err != nil {
			return err
		}

		// 3. Recalculate scores
		return c.jawabanGulihRepo.RecalculateGulih(resp.PerusahaanID)
	})
}

// ConsumeJawabanGulihDeleted (Pola 2 Asynchronous Write)
func (c *Consumer) ConsumeJawabanGulihDeleted(ctx context.Context) error {
	return c.Consume(ctx, "jawaban.gulih.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.JawabanGulihDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Processing Jawaban Gulih Deleted for ID: %d", event.ID)

		// 1. Delete from database
		if err := c.jawabanGulihRepo.Delete(event.ID); err != nil {
			return err
		}

		// 2. Recalculate scores
		return c.jawabanGulihRepo.RecalculateGulih(event.PerusahaanID)
	})
}

func (c *Consumer) ConsumeDomainCreated(ctx context.Context) error {
	return c.Consume(ctx, "domain.created", func(ctx context.Context, body []byte) error {
		var event dto_event.DomainCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from domain.created: %v", err)
			return nil
		}

		log.Printf("Processing Domain Created: %s", event.Request.NamaDomain)

		_, err := c.domainRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumeDomainUpdated(ctx context.Context) error {
	return c.Consume(ctx, "domain.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.DomainUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from domain.updated: %v", err)
			return nil
		}

		log.Printf("Processing Domain Updated for ID: %d", event.ID)

		return c.domainRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumeDomainDeleted(ctx context.Context) error {
	return c.Consume(ctx, "domain.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.DomainDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from domain.deleted: %v", err)
			return nil
		}

		log.Printf("Processing Domain Deleted for ID: %d", event.ID)

		return c.domainRepo.Delete(event.ID)
	})
}

func (c *Consumer) ConsumeRuangLingkupCreated(ctx context.Context) error {
	return c.Consume(ctx, "ruang_lingkup.created", func(ctx context.Context, body []byte) error {
		var event dto_event.RuangLingkupCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from ruang_lingkup.created: %v", err)
			return nil
		}
		log.Printf("Processing Ruang Lingkup Created: %s", event.Request.NamaRuangLingkup)
		_, err := c.ruangLingkupRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumeRuangLingkupUpdated(ctx context.Context) error {
	return c.Consume(ctx, "ruang_lingkup.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.RuangLingkupUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from ruang_lingkup.updated: %v", err)
			return nil
		}
		log.Printf("Processing Ruang Lingkup Updated for ID: %d", event.ID)
		return c.ruangLingkupRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumeRuangLingkupDeleted(ctx context.Context) error {
	return c.Consume(ctx, "ruang_lingkup.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.RuangLingkupDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from ruang_lingkup.deleted: %v", err)
			return nil
		}
		log.Printf("Processing Ruang Lingkup Deleted for ID: %d", event.ID)
		return c.ruangLingkupRepo.Delete(event.ID)
	})
}

func (c *Consumer) ConsumeKategoriCreated(ctx context.Context) error {
	return c.Consume(ctx, "kategori.created", func(ctx context.Context, body []byte) error {
		var event dto_event.KategoriCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from kategori.created: %v", err)
			return nil
		}
		log.Printf("Processing Kategori Created: %s", event.Request.NamaKategori)
		_, err := c.kategoriRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumeKategoriUpdated(ctx context.Context) error {
	return c.Consume(ctx, "kategori.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.KategoriUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from kategori.updated: %v", err)
			return nil
		}
		log.Printf("Processing Kategori Updated for ID: %d", event.ID)
		return c.kategoriRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumeKategoriDeleted(ctx context.Context) error {
	return c.Consume(ctx, "kategori.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.KategoriDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from kategori.deleted: %v", err)
			return nil
		}
		log.Printf("Processing Kategori Deleted for ID: %d", event.ID)
		return c.kategoriRepo.Delete(event.ID)
	})
}

func (c *Consumer) ConsumeIkasAuditLog(ctx context.Context) error {
	return c.Consume(ctx, "ikas.audit_logs", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasAuditLogEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from ikas.audit_logs: %v. Body: %s", err, string(body))
			return nil // Acknowledge to remove invalid JSON from queue
		}

		log.Printf("Processing IKAS Audit Log for ID: %s, User: %s", event.IkasID, event.UserID)

		if err := c.auditLogRepo.SaveAuditLog(event); err != nil {
			log.Printf("Error saving audit log: %v", err)
			return err
		}
		return nil
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
		c.ConsumeJawabanProteksiCreated,
		c.ConsumeJawabanProteksiUpdated,
		c.ConsumeJawabanProteksiDeleted,
		c.ConsumeJawabanDeteksiCreated,
		c.ConsumeJawabanDeteksiUpdated,
		c.ConsumeJawabanDeteksiDeleted,
		c.ConsumeJawabanGulihCreated,
		c.ConsumeJawabanGulihUpdated,
		c.ConsumeJawabanGulihDeleted,
		c.ConsumeDomainCreated,
		c.ConsumeDomainUpdated,
		c.ConsumeDomainDeleted,
		c.ConsumeRuangLingkupCreated,
		c.ConsumeRuangLingkupUpdated,
		c.ConsumeRuangLingkupDeleted,
		c.ConsumeKategoriCreated,
		c.ConsumeKategoriUpdated,
		c.ConsumeKategoriDeleted,
		c.ConsumeIkasAuditLog,
	}

	for _, consumer := range consumers {
		if err := consumer(ctx); err != nil {
			return err
		}
	}

	log.Println("All consumers started successfully")
	return nil
}
