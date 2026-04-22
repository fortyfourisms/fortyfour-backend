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
	identifikasiRepo           repository.IdentifikasiRepositoryInterface
	pertanyaanIdentifikasiRepo repository.PertanyaanIdentifikasiRepositoryInterface
	jawabanIdentifikasiRepo    repository.JawabanIdentifikasiRepositoryInterface
	proteksiRepo               repository.ProteksiRepositoryInterface
	pertanyaanProteksiRepo     repository.PertanyaanProteksiRepositoryInterface
	jawabanProteksiRepo        repository.JawabanProteksiRepositoryInterface
	deteksiRepo                repository.DeteksiRepositoryInterface
	pertanyaanDeteksiRepo      repository.PertanyaanDeteksiRepositoryInterface
	jawabanDeteksiRepo         repository.JawabanDeteksiRepositoryInterface
	gulihRepo                  repository.GulihRepositoryInterface
	pertanyaanGulihRepo        repository.PertanyaanGulihRepositoryInterface
	jawabanGulihRepo           repository.JawabanGulihRepositoryInterface
	domainRepo                 repository.DomainRepositoryInterface
	ruangLingkupRepo           repository.RuangLingkupRepositoryInterface
	kategoriRepo               repository.KategoriRepositoryInterface
	subKategoriRepo            repository.SubKategoriRepositoryInterface
	auditLogRepo               repository.AuditLogRepositoryInterface
}

func NewConsumer(
	c *rabbitmq.Consumer,
	ikasRepo repository.IkasRepositoryInterface,
	identifikasiRepo repository.IdentifikasiRepositoryInterface,
	pertanyaanIdentifikasiRepo repository.PertanyaanIdentifikasiRepositoryInterface,
	jawabanIdentifikasiRepo repository.JawabanIdentifikasiRepositoryInterface,
	proteksiRepo repository.ProteksiRepositoryInterface,
	pertanyaanProteksiRepo repository.PertanyaanProteksiRepositoryInterface,
	jawabanProteksiRepo repository.JawabanProteksiRepositoryInterface,
	deteksiRepo repository.DeteksiRepositoryInterface,
	pertanyaanDeteksiRepo repository.PertanyaanDeteksiRepositoryInterface,
	jawabanDeteksiRepo repository.JawabanDeteksiRepositoryInterface,
	gulihRepo repository.GulihRepositoryInterface,
	pertanyaanGulihRepo repository.PertanyaanGulihRepositoryInterface,
	jawabanGulihRepo repository.JawabanGulihRepositoryInterface,
	domainRepo repository.DomainRepositoryInterface,
	ruangLingkupRepo repository.RuangLingkupRepositoryInterface,
	kategoriRepo repository.KategoriRepositoryInterface,
	subKategoriRepo repository.SubKategoriRepositoryInterface,
	auditLogRepo repository.AuditLogRepositoryInterface,
) *Consumer {
	return &Consumer{
		Consumer:                   c,
		ikasRepo:                   ikasRepo,
		identifikasiRepo:           identifikasiRepo,
		pertanyaanIdentifikasiRepo: pertanyaanIdentifikasiRepo,
		jawabanIdentifikasiRepo:    jawabanIdentifikasiRepo,
		proteksiRepo:               proteksiRepo,
		pertanyaanProteksiRepo:     pertanyaanProteksiRepo,
		jawabanProteksiRepo:        jawabanProteksiRepo,
		deteksiRepo:                deteksiRepo,
		pertanyaanDeteksiRepo:      pertanyaanDeteksiRepo,
		jawabanDeteksiRepo:         jawabanDeteksiRepo,
		gulihRepo:                  gulihRepo,
		pertanyaanGulihRepo:        pertanyaanGulihRepo,
		jawabanGulihRepo:           jawabanGulihRepo,
		domainRepo:                 domainRepo,
		ruangLingkupRepo:           ruangLingkupRepo,
		kategoriRepo:               kategoriRepo,
		subKategoriRepo:            subKategoriRepo,
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
		// Hanya id_perusahaan yang wajib ada (foreign key). tanggal boleh kosong karena nullable di DB.
		if strings.TrimSpace(event.IDPerusahaan) == "" {
			log.Printf("❌ Skipping invalid message from ikas.created: id_perusahaan is empty. ID: %s", event.IkasID)
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
			log.Printf("❌ Skipping invalid message from jawaban.identifikasi.created: jawaban_identifikasi is null. Ikas: %s, Question: %d", req.IkasID, req.PertanyaanIdentifikasiID)
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

		currentCount, err := c.jawabanIdentifikasiRepo.GetBufferCount(req.IkasID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for IkasID %s (%d/%d). Flushing buffer...", req.IkasID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanIdentifikasiRepo.FlushBuffer(req.IkasID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for IkasID %s", req.IkasID)
			return c.jawabanIdentifikasiRepo.RecalculateIdentifikasi(req.IkasID)
		}

		log.Printf("Progress for IkasID %s: %d/%d", req.IkasID, currentCount, totalQuestions)
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
		return c.jawabanIdentifikasiRepo.RecalculateIdentifikasi(resp.IkasID)
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
		return c.jawabanIdentifikasiRepo.RecalculateIdentifikasi(event.IkasID)
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
			log.Printf("❌ Skipping invalid message from jawaban.proteksi.created: jawaban_proteksi is null. Ikas: %s, Question: %d", req.IkasID, req.PertanyaanProteksiID)
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

		currentCount, err := c.jawabanProteksiRepo.GetBufferCount(req.IkasID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for IkasID %s (%d/%d). Flushing buffer...", req.IkasID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanProteksiRepo.FlushBuffer(req.IkasID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for IkasID %s", req.IkasID)
			return c.jawabanProteksiRepo.RecalculateProteksi(req.IkasID)
		}

		log.Printf("Progress for IkasID %s: %d/%d", req.IkasID, currentCount, totalQuestions)
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
		return c.jawabanProteksiRepo.RecalculateProteksi(resp.IkasID)
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
		return c.jawabanProteksiRepo.RecalculateProteksi(event.IkasID)
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
			log.Printf("❌ Skipping invalid message from jawaban.deteksi.created: jawaban_deteksi is null. Ikas: %s, Question: %d", req.IkasID, req.PertanyaanDeteksiID)
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

		currentCount, err := c.jawabanDeteksiRepo.GetBufferCount(req.IkasID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for IkasID %s (%d/%d). Flushing buffer...", req.IkasID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanDeteksiRepo.FlushBuffer(req.IkasID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for IkasID %s", req.IkasID)
			return c.jawabanDeteksiRepo.RecalculateDeteksi(req.IkasID)
		}

		log.Printf("Progress for IkasID %s: %d/%d", req.IkasID, currentCount, totalQuestions)
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
		return c.jawabanDeteksiRepo.RecalculateDeteksi(resp.IkasID)
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
		return c.jawabanDeteksiRepo.RecalculateDeteksi(event.IkasID)
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
			log.Printf("❌ Skipping invalid message from jawaban.gulih.created: jawaban_gulih is null. Ikas: %s, Question: %d", req.IkasID, req.PertanyaanGulihID)
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

		currentCount, err := c.jawabanGulihRepo.GetBufferCount(req.IkasID)
		if err != nil {
			log.Printf("Error getting buffer count: %v", err)
			return err
		}

		if currentCount >= totalQuestions {
			log.Printf("All questions answered for IkasID %s (%d/%d). Flushing buffer...", req.IkasID, currentCount, totalQuestions)
			// 3. Flush buffer to main table
			if err := c.jawabanGulihRepo.FlushBuffer(req.IkasID); err != nil {
				log.Printf("Error flushing buffer: %v", err)
				return err
			}
			// 4. Recalculate scores
			log.Printf("Recalculating scores for IkasID %s", req.IkasID)
			return c.jawabanGulihRepo.RecalculateGulih(req.IkasID)
		}

		log.Printf("Progress for IkasID %s: %d/%d", req.IkasID, currentCount, totalQuestions)
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
		return c.jawabanGulihRepo.RecalculateGulih(resp.IkasID)
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
		return c.jawabanGulihRepo.RecalculateGulih(event.IkasID)
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

func (c *Consumer) ConsumeSubKategoriCreated(ctx context.Context) error {
	return c.Consume(ctx, "sub_kategori.created", func(ctx context.Context, body []byte) error {
		var event dto_event.SubKategoriCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from sub_kategori.created: %v", err)
			return nil
		}
		log.Printf("Processing Sub Kategori Created: %s", event.Request.NamaSubKategori)
		_, err := c.subKategoriRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumeSubKategoriUpdated(ctx context.Context) error {
	return c.Consume(ctx, "sub_kategori.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.SubKategoriUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from sub_kategori.updated: %v", err)
			return nil
		}
		log.Printf("Processing Sub Kategori Updated for ID: %d", event.ID)
		return c.subKategoriRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumeSubKategoriDeleted(ctx context.Context) error {
	return c.Consume(ctx, "sub_kategori.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.SubKategoriDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from sub_kategori.deleted: %v", err)
			return nil
		}
		log.Printf("Processing Sub Kategori Deleted for ID: %d", event.ID)
		return c.subKategoriRepo.Delete(event.ID)
	})
}

func (c *Consumer) ConsumePertanyaanIdentifikasiCreated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_identifikasi.created", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanIdentifikasiCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_identifikasi.created: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Identifikasi Created: %s", event.Request.PertanyaanIdentifikasi)
		_, err := c.pertanyaanIdentifikasiRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumePertanyaanIdentifikasiUpdated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_identifikasi.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanIdentifikasiUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_identifikasi.updated: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Identifikasi Updated for ID: %d", event.ID)
		return c.pertanyaanIdentifikasiRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumePertanyaanIdentifikasiDeleted(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_identifikasi.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanIdentifikasiDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_identifikasi.deleted: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Identifikasi Deleted for ID: %d", event.ID)
		return c.pertanyaanIdentifikasiRepo.Delete(event.ID)
	})
}

func (c *Consumer) ConsumePertanyaanProteksiCreated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_proteksi.created", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanProteksiCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_proteksi.created: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Proteksi Created: %s", event.Request.PertanyaanProteksi)
		_, err := c.pertanyaanProteksiRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumePertanyaanProteksiUpdated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_proteksi.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanProteksiUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_proteksi.updated: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Proteksi Updated for ID: %d", event.ID)
		return c.pertanyaanProteksiRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumePertanyaanProteksiDeleted(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_proteksi.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanProteksiDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_proteksi.deleted: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Proteksi Deleted for ID: %d", event.ID)
		return c.pertanyaanProteksiRepo.Delete(event.ID)
	})
}

func (c *Consumer) ConsumePertanyaanDeteksiCreated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_deteksi.created", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanDeteksiCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_deteksi.created: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Deteksi Created: %s", event.Request.PertanyaanDeteksi)
		_, err := c.pertanyaanDeteksiRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumePertanyaanDeteksiUpdated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_deteksi.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanDeteksiUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_deteksi.updated: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Deteksi Updated for ID: %d", event.ID)
		return c.pertanyaanDeteksiRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumePertanyaanDeteksiDeleted(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_deteksi.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanDeteksiDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_deteksi.deleted: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Deteksi Deleted for ID: %d", event.ID)
		return c.pertanyaanDeteksiRepo.Delete(event.ID)
	})
}

func (c *Consumer) ConsumePertanyaanGulihCreated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_gulih.created", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanGulihCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_gulih.created: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Gulih Created: %s", event.Request.PertanyaanGulih)
		_, err := c.pertanyaanGulihRepo.Create(event.Request)
		return err
	})
}

func (c *Consumer) ConsumePertanyaanGulihUpdated(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_gulih.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanGulihUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_gulih.updated: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Gulih Updated for ID: %d", event.ID)
		return c.pertanyaanGulihRepo.Update(event.ID, event.Request)
	})
}

func (c *Consumer) ConsumePertanyaanGulihDeleted(ctx context.Context) error {
	return c.Consume(ctx, "pertanyaan_gulih.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.PertanyaanGulihDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ Fatal: Unmarshal error from pertanyaan_gulih.deleted: %v", err)
			return nil
		}
		log.Printf("Processing Pertanyaan Gulih Deleted for ID: %d", event.ID)
		return c.pertanyaanGulihRepo.Delete(event.ID)
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
		c.ConsumePertanyaanIdentifikasiCreated,
		c.ConsumePertanyaanIdentifikasiUpdated,
		c.ConsumePertanyaanIdentifikasiDeleted,
		c.ConsumePertanyaanProteksiCreated,
		c.ConsumePertanyaanProteksiUpdated,
		c.ConsumePertanyaanProteksiDeleted,
		c.ConsumePertanyaanDeteksiCreated,
		c.ConsumePertanyaanDeteksiUpdated,
		c.ConsumePertanyaanDeteksiDeleted,
		c.ConsumePertanyaanGulihCreated,
		c.ConsumePertanyaanGulihUpdated,
		c.ConsumePertanyaanGulihDeleted,
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
