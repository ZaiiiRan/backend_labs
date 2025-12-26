package producer

import (
	"context"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
)

type OrderStatusChangedProducer struct {
	producer *Producer
	topic    string
}

func NewOrderStatusChangedProducer(cfg *config.KafkaProducerSettings, producer *Producer) *OrderStatusChangedProducer {
	return &OrderStatusChangedProducer{
		producer: producer,
		topic:    cfg.OrderStatusChangedTopic,
	}
}

func (p *OrderStatusChangedProducer) Produce(ctx context.Context, messages []Message) error {
	return p.producer.Produce(ctx, p.topic, messages)
}
