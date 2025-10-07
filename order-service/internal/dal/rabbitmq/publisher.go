package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	cfg  *config.RabbitMqSettings
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewPublisher(cfg *config.RabbitMqSettings) (*Publisher, error) {
	p := &Publisher{cfg: cfg}
	if err := p.connect(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Publisher) QueueOrderCreated() string {
	return p.cfg.OrderCreatedQueue
}

func (p *Publisher) PublishBatch(ctx context.Context, queue string, payloads []any) error {
	if len(payloads) == 0 {
		return nil
	}

	if p.ch == nil || p.conn == nil || p.conn.IsClosed() {
		if err := p.connect(); err != nil {
			return fmt.Errorf("rabbitmq reconnect: %w", err)
		}
	}

	_, err := p.ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	for _, msg := range payloads {
		body, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}

		err = p.ch.PublishWithContext(
			ctx,
			"",
			queue,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			return fmt.Errorf("publish: %w", err)
		}
	}

	return nil
}

func (p *Publisher) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Publisher) connect() error {
	vhost := url.PathEscape(p.cfg.VHost)
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", p.cfg.User, p.cfg.Password, p.cfg.Host, p.cfg.Port, vhost)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("amqp dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("channel open: %w", err)
	}

	p.conn = conn
	p.ch = ch

	return nil
}
