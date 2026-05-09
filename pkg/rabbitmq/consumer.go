package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Delivery is an alias for amqp.Delivery, exported so callers using the
// ConsumeBatch API don't need to import the amqp package directly.
type Delivery = amqp.Delivery

// MessageHandler handles incoming messages.
type MessageHandler func(ctx context.Context, body []byte) error

// BatchMessageHandler handles a batch of raw AMQP deliveries.
// The handler is responsible for processing each delivery and calling Ack/Nack.
type BatchMessageHandler func(ctx context.Context, msgs []amqp.Delivery)

// BatchConsumerConfig holds configuration for the batch consumer.
type BatchConsumerConfig struct {
	// BatchSize is the maximum number of messages to collect per batch.
	BatchSize int
	// BatchTimeout is the maximum time to wait while collecting a batch.
	// Prevents starvation when the queue has fewer messages than BatchSize.
	BatchTimeout time.Duration
	// BatchDelay is the pause between finishing one batch and starting the next.
	// This is the primary mechanism to avoid hitting downstream rate limits.
	BatchDelay time.Duration
	// UseAdaptive enables automatic throughput adjustment based on queue depth.
	UseAdaptive bool
}

// Consumer handles queue subscription and message processing.
type Consumer struct {
	rmq     *RabbitMQ     // Reference to RabbitMQ wrapper for robust channel management
	channel *amqp.Channel // Direct channel for legacy use
}

// NewConsumer creates a consumer that manages its own channels and handles reconnections.
func NewConsumer(rmq *RabbitMQ) *Consumer {
	return &Consumer{
		rmq: rmq,
	}
}

// NewConsumerWithChannel creates a basic consumer with a specific pre-opened channel.
func NewConsumerWithChannel(channel *amqp.Channel) *Consumer {
	return &Consumer{
		channel: channel,
	}
}

// Consume starts consuming messages from the specified queue.
// It uses dedicated channels per queue and automatically handles reconnections if created via NewConsumer.
func (c *Consumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	// If it's a channel-based consumer, use legacy logic
	if c.rmq == nil {
		return c.startChannelConsume(ctx, queueName, handler)
	}

	// Managed logic: Start a supervisor goroutine for this specific queue
	go c.supervisor(ctx, queueName, handler)
	return nil
}

// supervisor manages the lifecycle of a single queue consumer.
func (c *Consumer) supervisor(ctx context.Context, queueName string, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down supervisor for queue: %s", queueName)
			return
		default:
			// Attempt to consume
			err := c.runConsume(ctx, queueName, handler)
			if err != nil {
				log.Printf("Consumer for %s failed: %v. Retrying in 5 seconds...", queueName, err)

				// Wait before retry
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}
		}
	}
}

// runConsume performs the actual AMQP consumption. It blocks until the channel is closed or context is cancelled.
func (c *Consumer) runConsume(ctx context.Context, queueName string, handler MessageHandler) error {
	// 1. Get a fresh dedicated channel
	ch, err := c.rmq.NewChannel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// 2. Set QoS - prefetch 1 message at a time for better load balancing
	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// 3. Start consuming
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag (generated)
		false, // auto-ack: false (we handle manual ack for reliability)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Printf("Consumer registered and listening on queue: %s", queueName)

	// 4. Listen for channel closure
	notifyClose := ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-ctx.Done():
			return nil
		case closeErr := <-notifyClose:
			if closeErr != nil {
				return fmt.Errorf("channel closed unexpectedly: %v", closeErr)
			}
			return nil // Normal closure
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}

			// Process message
			if err := handler(ctx, msg.Body); err != nil {
				errStr := err.Error()
				if strings.Contains(errStr, "Error 1451") || strings.Contains(errStr, "foreign key constraint fails") {
					log.Printf("Dropping poison pill message from %s due to constraint error: %v", queueName, err)
					msg.Ack(false)
				} else {
					log.Printf("Error processing message from %s: %v", queueName, err)
					// Nack message and requeue so it's not lost
					msg.Nack(false, true)
				}
			} else {
				// Ack message on success
				msg.Ack(false)
			}
		}
	}
}

// startChannelConsume maintains support for specific channel-based consumers.
func (c *Consumer) startChannelConsume(ctx context.Context, queueName string, handler MessageHandler) error {
	if c.channel == nil {
		return fmt.Errorf("consumer has no channel")
	}

	// Set QoS
	_ = c.channel.Qos(1, 0, false)

	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Printf("Started channel-based consumer for queue: %s", queueName)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("Channel-based consumer closed for queue: %s", queueName)
					return
				}
				if err := handler(ctx, msg.Body); err != nil {
					log.Printf("Error in channel-based consumer %s: %v", queueName, err)
					msg.Nack(false, true)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// ConsumeBatch starts a batch consumer for the specified queue.
// It collects up to cfg.BatchSize messages (or waits up to cfg.BatchTimeout),
// calls the handler for the whole batch, then pauses for cfg.BatchDelay before
// reading the next batch. This distributes load and prevents rate-limit violations.
//
// The handler receives a slice of raw amqp.Delivery and is responsible for
// calling Ack or Nack on each message individually.
func (c *Consumer) ConsumeBatch(ctx context.Context, queueName string, cfg BatchConsumerConfig, handler BatchMessageHandler) error {
	if c.rmq == nil {
		return fmt.Errorf("batch consumer requires a managed RabbitMQ connection (use NewConsumer)")
	}

	// Apply sane defaults if caller left values at zero.
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.BatchTimeout <= 0 {
		cfg.BatchTimeout = 5 * time.Second
	}
	if cfg.BatchDelay < 0 {
		cfg.BatchDelay = 0
	}

	go c.supervisorBatch(ctx, queueName, cfg, handler)
	return nil
}

// supervisorBatch mirrors supervisor() but for batch consumers.
func (c *Consumer) supervisorBatch(ctx context.Context, queueName string, cfg BatchConsumerConfig, handler BatchMessageHandler) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down batch supervisor for queue: %s", queueName)
			return
		default:
			err := c.runBatchConsume(ctx, queueName, cfg, handler)
			if err != nil {
				log.Printf("Consumer for %s failed: %v. Retrying in 5 seconds...", queueName, err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}
		}
	}
}

// runBatchConsume is the core batch-processing loop.
// It reads up to cfg.BatchSize messages within cfg.BatchTimeout, delivers the
// collected batch to the handler, then sleeps for cfg.BatchDelay.
func (c *Consumer) runBatchConsume(ctx context.Context, queueName string, cfg BatchConsumerConfig, handler BatchMessageHandler) error {
	// 1. Open a dedicated channel for this queue.
	ch, err := c.rmq.NewChannel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// 2. Set QoS to batchSize so RabbitMQ pre-fetches at most that many messages.
	//    This is the key lever: the broker won't push more than batchSize unacknowledged
	//    messages to us at once, naturally rate-limiting inbound pressure.
	if err := ch.Qos(cfg.BatchSize, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// 3. Start consuming.
	deliveries, err := ch.Consume(
		queueName,
		"",    // consumer tag (auto-generated)
		false, // auto-ack: false — we ack/nack manually per message
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Printf("Batch consumer registered for queue: %s (batchSize=%d, delay=%s, adaptive=%v)",
		queueName, cfg.BatchSize, cfg.BatchDelay, cfg.UseAdaptive)

	// 4. Listen for channel-level errors.
	notifyClose := ch.NotifyClose(make(chan *amqp.Error, 1))

	for {
		// Collect a batch of messages.
		currentBatchSize := cfg.BatchSize
		currentDelay := cfg.BatchDelay

		// Adaptive logic: Adjust parameters based on current queue depth.
		if cfg.UseAdaptive {
			q, err := ch.QueueDeclare(
				queueName,
				false, // durable (ignored when passive is true)
				false, // delete when unused (ignored when passive is true)
				false, // exclusive (ignored when passive is true)
				false, // no-wait (ignored when passive is true)
				amqp.Table{"passive": true},
			)
			if err == nil {
				if q.Messages > 5000 {
					// Emergency Mode
					currentBatchSize = 10
					currentDelay = 1 * time.Second
					log.Printf("Adaptive Mode [Emergency]: Queue %s has %d messages. Throttling applied.", queueName, q.Messages)
				} else if q.Messages > 1000 {
					// Busy Mode
					currentBatchSize = 20
					currentDelay = 500 * time.Millisecond
					log.Printf("Adaptive Mode [Busy]: Queue %s has %d messages. Adjusting throughput.", queueName, q.Messages)
				}
			}
		}

		batch := make([]amqp.Delivery, 0, currentBatchSize)
		batchTimer := time.NewTimer(cfg.BatchTimeout)

	collectLoop:
		for len(batch) < currentBatchSize {
			select {
			case <-ctx.Done():
				batchTimer.Stop()
				// Flush whatever we collected before shutdown.
				if len(batch) > 0 {
					log.Printf("Context cancelled: flushing partial batch of %d messages from %s", len(batch), queueName)
					handler(ctx, batch)
				}
				return nil

			case closeErr := <-notifyClose:
				batchTimer.Stop()
				// Flush partial batch before propagating the error.
				if len(batch) > 0 {
					log.Printf("Channel closing: flushing partial batch of %d messages from %s", len(batch), queueName)
					handler(ctx, batch)
				}
				if closeErr != nil {
					return fmt.Errorf("channel closed unexpectedly: %v", closeErr)
				}
				return nil

			case msg, ok := <-deliveries:
				if !ok {
					batchTimer.Stop()
					if len(batch) > 0 {
						handler(ctx, batch)
					}
					return fmt.Errorf("delivery channel closed")
				}
				batch = append(batch, msg)

			case <-batchTimer.C:
				// Timeout reached — process whatever we have so far.
				break collectLoop
			}
		}

		batchTimer.Stop()

		if len(batch) == 0 {
			// Nothing collected during the timeout window; loop again.
			continue
		}

		log.Printf("Processing batch of %d messages from queue: %s", len(batch), queueName)

		// 5. Deliver the batch to the caller's handler.
		//    The handler acks/nacks each delivery individually.
		handler(ctx, batch)

		// 6. Throttle: pause before reading the next batch.
		//    This is what prevents downstream rate-limit violations.
		if currentDelay > 0 {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(currentDelay):
			}
		}
	}
}

// DropOrRequeue decides whether a failed message should be dropped (poison pill)
// or requeued for retry, based on the error content.
// It is exported so downstream packages (e.g. ikas/internal/rabbitmq) can reuse
// the same consistent error-handling logic inside their batch handlers.
func DropOrRequeue(msg amqp.Delivery, queueName string, err error) {
	errStr := err.Error()
	if strings.Contains(errStr, "Error 1451") || strings.Contains(errStr, "foreign key constraint fails") {
		log.Printf("Dropping poison pill from %s (constraint error): %v", queueName, err)
		msg.Ack(false)
	} else {
		log.Printf("Error processing message from %s: %v — requeuing", queueName, err)
		msg.Nack(false, true)
	}
}
