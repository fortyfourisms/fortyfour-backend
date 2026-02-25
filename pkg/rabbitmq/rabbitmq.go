package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// NewRabbitMQ
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{
		url: url,
	}

	if err := rmq.connect(); err != nil {
		return nil, err
	}

	return rmq, nil
}

// connect with retry and exponential backoff
func (r *RabbitMQ) connect() error {
	maxRetries := 5
	baseDelay := 2 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			delay := baseDelay * time.Duration(1<<(i-1)) // 2s, 4s, 8s, 16s, 32s
			log.Printf("⏳ RabbitMQ connection attempt %d/%d failed, retrying in %v...", i, maxRetries, delay)
			time.Sleep(delay)
		}

		var err error

		// create connection
		r.conn, err = amqp.Dial(r.url)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to RabbitMQ: %w", err)
			continue
		}

		// create channel
		r.channel, err = r.conn.Channel()
		if err != nil {
			r.conn.Close()
			lastErr = fmt.Errorf("failed to open channel: %w", err)
			continue
		}

		log.Println("Connected to RabbitMQ successfully")
		return nil
	}

	return fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, lastErr)
}

// GetChannel
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}

// Close
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("🔌 RabbitMQ connection closed")
}

// DeclareExchange
func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,  // name
		kind,  // type (direct, fanout, topic, headers)
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// BindQueue
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
}

// Reconnect ke RabbitMQ jika koneksi terputus
func (r *RabbitMQ) Reconnect() error {
	r.Close()
	time.Sleep(5 * time.Second)
	return r.connect()
}

// IsConnected
func (r *RabbitMQ) IsConnected() bool {
	return r.conn != nil && !r.conn.IsClosed()
}
