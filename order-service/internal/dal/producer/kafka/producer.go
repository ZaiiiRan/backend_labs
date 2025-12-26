package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
	once     sync.Once
	cfg      *config.KafkaProducerSettings
}

func NewKafkaProducer(cfg *config.KafkaProducerSettings) (*Producer, error) {
	return &Producer{cfg: cfg}, nil
}

func (p *Producer) Produce(ctx context.Context, topic string, messages []Message) error {
	if len(messages) == 0 {
		return nil
	}

	if err := p.init(); err != nil {
		return err
	}

	for _, msg := range messages {
		value, err := json.Marshal(msg.Value)
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}

		err = p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(msg.Key),
			Value: value,
		}, nil)

		if err != nil {
			return fmt.Errorf("produce: %w", err)
		}
	}

	return nil
}

func (p *Producer) Close() {
	if p.producer != nil {
		p.producer.Flush(10000)
		p.producer.Close()
	}
}

func (p *Producer) init() error {
	var err error

	p.once.Do(func() {
		p.producer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": p.cfg.KafkaSettings.BootstrapServers,
			"client.id":         p.cfg.KafkaSettings.ClientId,
			"linger.ms":         100,
			"compression.type":  "snappy",
			"partitioner":       "consistent",
		})
		if err != nil {
			return
		}
	})

	return err
}
