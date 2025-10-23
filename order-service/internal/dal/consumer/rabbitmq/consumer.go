package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	cfg              *config.RabbitMqConsumerSettings
	client           *rabbitmq.RabbitMqClient
	ch               *amqp.Channel
	log              *zap.SugaredLogger
	stopCh           chan struct{}
	buffer           []consumer.MessageInfo
	mu               sync.Mutex
	timer            *time.Timer
	messageProcessor consumer.MessageProcessor
}

func NewConsumer(
	cfg *config.RabbitMqConsumerSettings,
	client *rabbitmq.RabbitMqClient,
	messageProcessor consumer.MessageProcessor,
	log *zap.SugaredLogger,
) (*Consumer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	return &Consumer{
		cfg:              cfg,
		client:           client,
		ch:               ch,
		log:              log,
		stopCh:           make(chan struct{}),
		messageProcessor: messageProcessor,
	}, nil
}

func (c *Consumer) Start() error {
	if err := c.ensureChannel(); err != nil {
		return err
	}

	if err := c.ch.Qos(c.cfg.BatchSize*2, 0, false); err != nil {
		return fmt.Errorf("set qos: %w", err)
	}

	_, err := c.ch.QueueDeclare(
		c.cfg.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	msgs, err := c.ch.Consume(
		c.cfg.Queue,
		c.cfg.Consumer,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	batchTimeout := time.Duration(c.cfg.BatchTimeoutSeconds) * time.Second
	c.timer = time.AfterFunc(batchTimeout, func() {
		c.processBatch("timeout")
	})

	go func() {
		c.log.Infow("consumer.started", "queue", c.cfg.Queue)
		for {
			select {
			case <-c.stopCh:
				c.log.Infow("consumer.stopped", "queue", c.cfg.Queue)
				return
			case msg, ok := <-msgs:
				if !ok {
					c.log.Warnw("consumer.channel_closed", "queue", c.cfg.Queue)
					return
				}
				c.handleMessage(msg)
			}
		}
	}()

	return nil
}

func (c *Consumer) handleMessage(msg amqp.Delivery) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buffer = append(c.buffer, consumer.MessageInfo{
		DeliveryTag: msg.DeliveryTag,
		Body:        msg.Body,
		ReceivedAt:  time.Now(),
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

	batch := make([]consumer.MessageInfo, len(c.buffer))
	copy(batch, c.buffer)
	c.buffer = c.buffer[:0]
	c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.log.Infow("consumer.process_started", "queue", c.cfg.Queue, "count", len(batch), "trigger", trigger)

	if err := c.messageProcessor.ProcessMessage(ctx, batch); err != nil {
		c.log.Errorw("consumer.process_failed", "queue", c.cfg.Queue, "err", err)
		last := batch[len(batch)-1]
		c.ch.Nack(last.DeliveryTag, true, true)
		return
	}

	last := batch[len(batch)-1]
	if err := c.ch.Ack(last.DeliveryTag, true); err != nil {
		c.log.Errorw("consumer.ack_failed", "queue", c.cfg.Queue, "err", err)
	}
}

func (c *Consumer) Stop() {
	close(c.stopCh)
	if c.timer != nil {
		c.timer.Stop()
	}

	time.Sleep(500 * time.Millisecond)

	if c.ch != nil && !c.ch.IsClosed() {
		c.ch.Close()
	}
}

func (c *Consumer) ensureChannel() error {
	if c.ch == nil || c.ch.IsClosed() {
		ch, err := c.client.Channel()
		if err != nil {
			return fmt.Errorf("reopen channel: %w", err)
		}
		c.ch = ch
	}
	return nil
}
