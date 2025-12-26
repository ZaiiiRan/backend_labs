package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/messages"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	client *rabbitmq.RabbitMqClient
	cfg    *config.RabbitMqPublisherSettings

	ch        *amqp.Channel
	onceSetup sync.Once
}

func NewPublisher(cfg *config.RabbitMqPublisherSettings, client *rabbitmq.RabbitMqClient) (*Publisher, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{
		cfg:    cfg,
		client: client,
		ch:     ch,
	}, nil
}

func (p *Publisher) PublishBatch(ctx context.Context, messages []messages.Message) error {
	if len(messages) == 0 {
		return nil
	}

	if err := p.configure(ctx); err != nil {
		return err
	}

	for _, msg := range messages {
		body, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("msg json marshal: %w", err)
		}

		err = p.ch.PublishWithContext(
			ctx,
			p.cfg.Exchange,
			msg.RoutingKey(),
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

func (p *Publisher) configure(ctx context.Context) error {
	var err error

	p.onceSetup.Do(func() {
		err = p.ensureChannel()
		if err != nil {
			return
		}

		err = p.ch.ExchangeDeclare(
			p.cfg.Exchange,
			"topic",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			err = fmt.Errorf("exchange declare: %w", err)
			return
		}

		for _, m := range p.cfg.ExchangeMappings {
			_, e := p.ch.QueueDeclare(
				m.Queue,
				false,
				false,
				false,
				false,
				nil,
			)
			if e != nil {
				err = fmt.Errorf("queue declare %s: %w", m.Queue, e)
				return
			}

			e = p.ch.QueueBind(
				m.Queue,
				m.RoutingKeyPattern,
				p.cfg.Exchange,
				false,
				nil,
			)
			if e != nil {
				err = fmt.Errorf("queue bind %s: %w", m.Queue, e)
				return
			}
		}
	})

	return err
}

func (p *Publisher) ensureChannel() error {
	if p.ch != nil && !p.ch.IsClosed() {
		return nil
	}
	ch, err := p.client.Channel()
	if err != nil {
		return err
	}
	p.ch = ch
	return nil
}
