package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	client *rabbitmq.RabbitMqClient
	queue  string
	ch     *amqp.Channel
}

func NewPublisher(client *rabbitmq.RabbitMqClient, queue string) (*Publisher, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{
		client: client,
		queue:  queue,
		ch:     ch,
	}, nil
}

func (p *Publisher) PublishBatch(ctx context.Context, payloads []any) error {
	if len(payloads) == 0 {
		return nil
	}

	if err := p.ensureChannel(); err != nil {
		return err
	}

	_, err := p.ch.QueueDeclare(
		p.queue,
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
			return fmt.Errorf("msg json marshal: %w", err)
		}

		err = p.ch.PublishWithContext(
			ctx,
			"",
			p.queue,
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
}

func (p *Publisher) ensureChannel() error {
	if p.ch == nil || p.ch.IsClosed() {
		ch, err := p.client.Channel()
		if err != nil {
			return fmt.Errorf("reopen channel: %w", err)
		}
		p.ch = ch
	}
	return nil
}
