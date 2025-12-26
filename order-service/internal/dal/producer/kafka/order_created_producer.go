package producer

import (
	"context"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
)

type OrderCreatedProducer struct {
	producer *Producer
	topic    string
}

func NewOrderCreatedProducer(cfg *config.KafkaProducerSettings, producer *Producer) *OrderCreatedProducer {
	return &OrderCreatedProducer{
		producer: producer,
		topic:    cfg.OrderCreatedTopic,
	}
}

func (p *OrderCreatedProducer) Produce(ctx context.Context, messages []Message) error {
	return p.producer.Produce(ctx, p.topic, messages)
}
