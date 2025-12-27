package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	dalconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Consumer struct {
	cfg       *config.KafkaConsumerSettings
	consumer  *kafka.Consumer
	topic     string
	processor dalconsumer.MessageProcessor
	log       *zap.SugaredLogger

	mu     sync.Mutex
	buffer []batchMessage
	timer  *time.Timer
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

		"enable.auto.commit": false,

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
		buffer:    make([]batchMessage, 0, cfg.BatchSize),
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	if err := c.consumer.Subscribe(c.topic, nil); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	c.log.Infow("consumer.kafka_consumer.started", "topic", c.topic)

	batchTimeout := time.Duration(c.cfg.BatchTimeoutMs) * time.Millisecond
	c.resetTimer(batchTimeout)

	for {
		select {
		case <-ctx.Done():
			c.flush("shutdown")
			return nil
		default:
			msg, err := c.consumer.ReadMessage(500 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok {
					if kafkaErr.Code() == kafka.ErrTimedOut {
						continue
					}
					c.log.Errorw("consumer.kafka_consume_error", "err", kafkaErr, "topic", c.topic)
				}
				continue
			}
			c.handleMessage(msg)
		}
	}
}

func (c *Consumer) Close() {
	if c.timer != nil {
		c.timer.Stop()
	}
	c.consumer.Close()
}

func (c *Consumer) handleMessage(msg *kafka.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buffer = append(c.buffer, batchMessage{
		msg: msg.TopicPartition,
		val: dalconsumer.Message{
			Key:  string(msg.Key),
			Body: msg.Value,
		},
	})

	if len(c.buffer) >= c.cfg.BatchSize {
		go c.processBatch("size_limit")
	}
}

func (c *Consumer) processBatch(trigger string) {
	c.mu.Lock()
	if len(c.buffer) == 0 {
		c.mu.Unlock()
		return
	}

	batch := make([]batchMessage, len(c.buffer))
	copy(batch, c.buffer)
	c.buffer = c.buffer[:0]
	c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.log.Infow(
		"consumer.process_started",
		"topic", c.topic,
		"count", len(batch),
		"trigger", trigger,
	)

	requeue, err := c.processor.ProcessMessage(ctx, extract(batch))
	if err != nil {
		c.log.Errorw(
			"consumer.process_failed",
			"topic", c.topic,
			"err", err,
			"need_requeue", requeue,
		)
		return
	}

	last := batch[len(batch)-1].msg
	last.Offset++

	if _, err := c.consumer.CommitOffsets([]kafka.TopicPartition{last}); err != nil {
		c.log.Errorw("consumer.commit_failed", "err", err)
	}
}

func (c *Consumer) resetTimer(d time.Duration) {
	if c.timer != nil {
		c.timer.Stop()
	}

	c.timer = time.AfterFunc(d, func() {
		c.processBatch("timeout")
		c.resetTimer(d)
	})
}

func (c *Consumer) flush(reason string) {
	c.log.Infow("consumer.flush", "topic", c.topic, "reason", reason)
	c.processBatch(reason)
}

func extract(batch []batchMessage) []dalconsumer.Message {
	res := make([]dalconsumer.Message, len(batch))
	for i, b := range batch {
		res[i] = b.val
	}
	return res
}
