package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/rabbitmq"
	"time"
)

// SSEBroadcaster defines the interface for SSE notifications to avoid import cycles
type SSEBroadcaster interface {
	NotifyCreate(resource string, data interface{}, userID string)
	NotifyUpdate(resource string, data interface{}, userID string)
	NotifyDelete(resource string, id interface{}, userID string)
}

type NotificationPusher interface {
	Push(userID string, notifType models.NotificationType, message string) error
}

// Consumer wrapper
type Consumer struct {
	*rabbitmq.Consumer
	sseService   SSEBroadcaster
	userRepo     repository.UserRepositoryInterface
	notifService NotificationPusher
	eventRepo    repository.EventRepositoryInterface
	beritaRepo   repository.BeritaRepositoryInterface
}

// NewConsumer
func NewConsumer(c *rabbitmq.Consumer, sseService SSEBroadcaster, userRepo repository.UserRepositoryInterface, notifService NotificationPusher, eventRepo repository.EventRepositoryInterface, beritaRepo repository.BeritaRepositoryInterface) *Consumer {
	return &Consumer{
		Consumer:     c,
		sseService:   sseService,
		userRepo:     userRepo,
		notifService: notifService,
		eventRepo:    eventRepo,
		beritaRepo:   beritaRepo,
	}
}

// Users
func (c *Consumer) ConsumeUserCreated(ctx context.Context) error {
	return c.Consume(ctx, "users.created", func(ctx context.Context, body []byte) error {
		var event dto_event.UserCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Created Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyCreate("users", event, "")
		}

		return nil
	})
}

func (c *Consumer) ConsumeUserUpdated(ctx context.Context) error {
	return c.Consume(ctx, "users.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.UserUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Updated Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyUpdate("users", event, "")
		}

		return nil
	})
}

func (c *Consumer) ConsumeUserDeleted(ctx context.Context) error {
	return c.Consume(ctx, "users.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.UserDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Deleted Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyDelete("users", event.ID, "")
		}

		return nil
	})
}

func (c *Consumer) ConsumeUserPasswordUpdated(ctx context.Context) error {
	return c.Consume(ctx, "users.password_updated", func(ctx context.Context, body []byte) error {
		var event dto_event.UserPasswordUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Password Updated Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyUpdate("password", event, "")
		}

		return nil
	})
}

// IKAS
func (c *Consumer) ConsumeIkasCreated(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.created", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Created Event (from RabbitMQ): %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyCreate("ikas", event, event.UserID)
		}

		return nil
	})
}

func (c *Consumer) ConsumeIkasUpdated(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Updated Event (from RabbitMQ): %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyUpdate("ikas", event, event.UserID)
		}

		return nil
	})
}

func (c *Consumer) ConsumeIkasDeleted(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Deleted Event (from RabbitMQ): %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyDelete("ikas", event.IkasID, event.UserID)
		}

		return nil
	})
}

func (c *Consumer) ConsumeIkasEditRequested(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.edit_requested", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasEditRequestedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Edit Requested Event: %+v", event)

		// 1. Persist to DB for all admins
		if c.userRepo != nil && c.notifService != nil {
			admins, _ := c.userRepo.FindAllAdmins()
			msg := "User " + event.Responden + " (Perusahaan: " + event.NamaPerusahaan + ") mengajukan permintaan edit data IKAS. Alasan: " + event.Reason
			for _, admin := range admins {
				_ = c.notifService.Push(admin.ID, models.NotifIkasEditRequested, msg)
			}
		}

		// 2. Real-time SSE (send to admin topic/all admins)
		if c.sseService != nil {
			c.sseService.NotifyCreate("ikas_request", event, "admin")
		}

		return nil
	})
}

func (c *Consumer) ConsumeIkasEditActioned(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.edit_actioned", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasEditActionedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Edit Actioned Event: %+v", event)

		// 1. Find user owner of this IKAS ID (this might require another event enrichment or repo call)
		// For now, we'll notify the 'ikas' resource which might be observed by the user
		msg := "Permintaan edit data IKAS Anda telah " + event.Status
		if event.Status == "rejected" && event.AdminReason != "" {
			msg += ". Alasan: " + event.AdminReason
		}

		// Ideally we find the UserID of the owner here to persist notification
		// If we don't have it in the event, we might need to fetch it or include it in event

		if c.sseService != nil {
			c.sseService.NotifyUpdate("ikas_action", event, "") // Broadcast or targeted
		}

		return nil
	})
}

// Csirt
func (c *Consumer) ConsumeCsirtCreated(ctx context.Context) error {
	return c.Consume(ctx, "csirt.created", func(ctx context.Context, body []byte) error {
		var event dto_event.CsirtCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("CSIRT Created Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyCreate("csirt", event, "")
		}

		return nil
	})
}

func (c *Consumer) ConsumeCsirtUpdated(ctx context.Context) error {
	return c.Consume(ctx, "csirt.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.CsirtUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("CSIRT Updated Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyUpdate("csirt", event, "")
		}

		return nil
	})
}

func (c *Consumer) ConsumeCsirtDeleted(ctx context.Context) error {
	return c.Consume(ctx, "csirt.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.CsirtDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("CSIRT Deleted Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyDelete("csirt", event.ID, "")
		}

		return nil
	})
}

// consumeGenericIkasEvent handles common CRUD events dynamically
func (c *Consumer) consumeGenericIkasEvent(ctx context.Context, queueName string, resource string, action string) error {
	return c.Consume(ctx, queueName, func(ctx context.Context, body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("[%s] %s Event (from RabbitMQ): %+v", resource, action, event)
		if c.sseService != nil {
			userID, _ := event["user_id"].(string)
			switch action {
			case "Created":
				c.sseService.NotifyCreate(resource, event, userID)
			case "Updated":
				c.sseService.NotifyUpdate(resource, event, userID)
			case "Deleted":
				c.sseService.NotifyDelete(resource, event, userID)
			}
		}

		return nil
	})
}

// StartAllConsumers
func (c *Consumer) StartAllConsumers(ctx context.Context) error {
	consumers := []func(context.Context) error{
		c.ConsumeUserCreated,
		c.ConsumeUserUpdated,
		c.ConsumeUserDeleted,
		c.ConsumeUserPasswordUpdated,
		c.ConsumeIkasCreated,
		c.ConsumeIkasUpdated,
		c.ConsumeIkasDeleted,
		c.ConsumeIkasEditRequested,
		c.ConsumeIkasEditActioned,
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.ruang_lingkup.created", "ruang_lingkup", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.ruang_lingkup.updated", "ruang_lingkup", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.ruang_lingkup.deleted", "ruang_lingkup", "Deleted")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.domain.created", "domain", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.domain.updated", "domain", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.domain.deleted", "domain", "Deleted")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.kategori.created", "kategori", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.kategori.updated", "kategori", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.kategori.deleted", "kategori", "Deleted")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.sub_kategori.created", "sub_kategori", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.sub_kategori.updated", "sub_kategori", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.sub_kategori.deleted", "sub_kategori", "Deleted")
		},
		c.ConsumeCsirtCreated,
		c.ConsumeCsirtUpdated,
		c.ConsumeCsirtDeleted,
		c.ConsumePerusahaanCreated,
		c.ConsumePerusahaanUpdated,
		c.ConsumePerusahaanDeleted,
		c.ConsumePicCreated,
		c.ConsumePicUpdated,
		c.ConsumePicDeleted,
		c.ConsumeSdmCsirtCreated,
		c.ConsumeSdmCsirtUpdated,
		c.ConsumeSdmCsirtDeleted,
		c.ConsumeRoleCreated,
		c.ConsumeRoleUpdated,
		c.ConsumeRoleDeleted,
		c.ConsumeSeCreated,
		c.ConsumeSeUpdated,
		c.ConsumeSeDeleted,
		c.ConsumeEventCreated,
		c.ConsumeEventUpdated,
		c.ConsumeEventDeleted,
		c.ConsumeBeritaCreated,
		c.ConsumeBeritaUpdated,
		c.ConsumeBeritaDeleted,
	}

	for _, consumer := range consumers {
		if err := consumer(ctx); err != nil {
			return err
		}
	}

	log.Println("All Main API consumers started successfully")
	return nil
}

// Perusahaan
func (c *Consumer) ConsumePerusahaanCreated(ctx context.Context) error {
	return c.Consume(ctx, "perusahaan.created", func(ctx context.Context, body []byte) error {
		var event dto_event.PerusahaanCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyCreate("perusahaan", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumePerusahaanUpdated(ctx context.Context) error {
	return c.Consume(ctx, "perusahaan.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.PerusahaanUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyUpdate("perusahaan", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumePerusahaanDeleted(ctx context.Context) error {
	return c.Consume(ctx, "perusahaan.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.PerusahaanDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyDelete("perusahaan", event.ID, "system")
		return nil
	})
}

// PIC
func (c *Consumer) ConsumePicCreated(ctx context.Context) error {
	return c.Consume(ctx, "pic.created", func(ctx context.Context, body []byte) error {
		var event dto_event.PicCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyCreate("pic", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumePicUpdated(ctx context.Context) error {
	return c.Consume(ctx, "pic.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.PicUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyUpdate("pic", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumePicDeleted(ctx context.Context) error {
	return c.Consume(ctx, "pic.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.PicDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyDelete("pic", event.ID, "system")
		return nil
	})
}

// SDM CSIRT
func (c *Consumer) ConsumeSdmCsirtCreated(ctx context.Context) error {
	return c.Consume(ctx, "sdm_csirt.created", func(ctx context.Context, body []byte) error {
		var event dto_event.SdmCsirtCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyCreate("sdm_csirt", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumeSdmCsirtUpdated(ctx context.Context) error {
	return c.Consume(ctx, "sdm_csirt.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.SdmCsirtUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyUpdate("sdm_csirt", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumeSdmCsirtDeleted(ctx context.Context) error {
	return c.Consume(ctx, "sdm_csirt.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.SdmCsirtDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyDelete("sdm_csirt", event.ID, "system")
		return nil
	})
}

// Role
func (c *Consumer) ConsumeRoleCreated(ctx context.Context) error {
	return c.Consume(ctx, "role.created", func(ctx context.Context, body []byte) error {
		var event dto_event.RoleCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyCreate("role", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumeRoleUpdated(ctx context.Context) error {
	return c.Consume(ctx, "role.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.RoleUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyUpdate("role", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumeRoleDeleted(ctx context.Context) error {
	return c.Consume(ctx, "role.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.RoleDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyDelete("role", event.ID, "system")
		return nil
	})
}

// SE
func (c *Consumer) ConsumeSeCreated(ctx context.Context) error {
	return c.Consume(ctx, "se.created", func(ctx context.Context, body []byte) error {
		var event dto_event.SeCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyCreate("se", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumeSeUpdated(ctx context.Context) error {
	return c.Consume(ctx, "se.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.SeUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyUpdate("se", event, "system")
		return nil
	})
}

func (c *Consumer) ConsumeSeDeleted(ctx context.Context) error {
	return c.Consume(ctx, "se.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.SeDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		c.sseService.NotifyDelete("se", event.ID, "system")
		return nil
	})
}

// Event
func (c *Consumer) ConsumeEventCreated(ctx context.Context) error {
	return c.Consume(ctx, "event.created", func(ctx context.Context, body []byte) error {
		var event dto_event.EventCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		t, _ := time.Parse(time.RFC3339, event.Request.Tanggal)
		model := &models.Event{
			Judul:     event.Request.Judul,
			Deskripsi: event.Request.Deskripsi,
			Lokasi:    event.Request.Lokasi,
			Tanggal:   t,
		}

		if err := c.eventRepo.Create(model); err != nil {
			log.Printf("Error creating event from RabbitMQ: %v", err)
			return err
		}

		if c.sseService != nil {
			c.sseService.NotifyCreate("event", model, "system")
		}

		return nil
	})
}

func (c *Consumer) ConsumeEventUpdated(ctx context.Context) error {
	return c.Consume(ctx, "event.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.EventUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		existing, err := c.eventRepo.FindByID(event.ID)
		if err != nil || existing == nil {
			return err
		}

		if event.Request.Judul != nil {
			existing.Judul = *event.Request.Judul
		}
		if event.Request.Deskripsi != nil {
			existing.Deskripsi = *event.Request.Deskripsi
		}
		if event.Request.Lokasi != nil {
			existing.Lokasi = *event.Request.Lokasi
		}
		if event.Request.Tanggal != nil {
			t, _ := time.Parse(time.RFC3339, *event.Request.Tanggal)
			existing.Tanggal = t
		}

		if err := c.eventRepo.Update(existing); err != nil {
			return err
		}

		if c.sseService != nil {
			c.sseService.NotifyUpdate("event", existing, "system")
		}

		return nil
	})
}

func (c *Consumer) ConsumeEventDeleted(ctx context.Context) error {
	return c.Consume(ctx, "event.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.EventDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		if err := c.eventRepo.Delete(event.ID); err != nil {
			return err
		}

		if c.sseService != nil {
			c.sseService.NotifyDelete("event", event.ID, "system")
		}

		return nil
	})
}

// Berita
func (c *Consumer) ConsumeBeritaCreated(ctx context.Context) error {
	return c.Consume(ctx, "berita.created", func(ctx context.Context, body []byte) error {
		var event dto_event.BeritaCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		model := &models.Berita{
			Judul:     event.Request.Judul,
			Deskripsi: event.Request.Deskripsi,
			AuthorID:  event.AuthorID,
		}

		if err := c.beritaRepo.Create(model); err != nil {
			log.Printf("Error creating berita from RabbitMQ: %v", err)
			return err
		}

		if c.sseService != nil {
			c.sseService.NotifyCreate("berita", model, "system")
		}

		return nil
	})
}

func (c *Consumer) ConsumeBeritaUpdated(ctx context.Context) error {
	return c.Consume(ctx, "berita.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.BeritaUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		existing, err := c.beritaRepo.FindByID(event.ID)
		if err != nil || existing == nil {
			return err
		}

		if event.Request.Judul != nil {
			existing.Judul = *event.Request.Judul
		}
		if event.Request.Deskripsi != nil {
			existing.Deskripsi = *event.Request.Deskripsi
		}

		if err := c.beritaRepo.Update(existing); err != nil {
			return err
		}

		if c.sseService != nil {
			c.sseService.NotifyUpdate("berita", existing, "system")
		}

		return nil
	})
}

func (c *Consumer) ConsumeBeritaDeleted(ctx context.Context) error {
	return c.Consume(ctx, "berita.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.BeritaDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		if err := c.beritaRepo.Delete(event.ID); err != nil {
			return err
		}

		if c.sseService != nil {
			c.sseService.NotifyDelete("berita", event.ID, "system")
		}

		return nil
	})
}
