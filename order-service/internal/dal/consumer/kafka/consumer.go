package consumer

import (
	"context"
	"fmt"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Consumer struct {
	cfg       *config.KafkaConsumerSettings
	consumer  *kafka.Consumer
	topic     string
	processor consumer.MessageProcessor
	log       *zap.SugaredLogger
}

func NewConsumer(
	cfg *config.KafkaConsumerSettings,
	topic string,
	processor consumer.MessageProcessor,
	log *zap.SugaredLogger,
) (*Consumer, error) {
	kc, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaSettings.BootstrapServers,
		"group.id":          cfg.GroupId,

		"auto.offset.reset": "latest",

		"enable.auto.commit":      true,
		"auto.commit.interval.ms": 5000,

		"session.timeout.ms":    60000,
		"heartbeat.interval.ms": 3000,
		"max.poll.interval.ms":  300000,
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{
		cfg:       cfg,
		consumer:  kc,
		topic:     topic,
		processor: processor,
		log:       log,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	if err := c.consumer.Subscribe(c.topic, nil); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	c.log.Infow("consumer.kafka_consumer.started", "topic", c.topic)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok {
					c.log.Errorw("consumer.kafka_consume_error", "err", kafkaErr, "topic", c.topic)
				}
				continue
			}

			m := consumer.Message{
				Key:  string(msg.Key),
				Body: msg.Value,
			}

			_, err = c.processor.ProcessMessage(ctx, []consumer.Message{m})
			if err != nil {
				c.log.Errorw(
					"consumer.kafka_message_processing_failed",
					"topic", c.topic,
					"key", m.Key,
					"err", err,
				)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.consumer.Close()
}
