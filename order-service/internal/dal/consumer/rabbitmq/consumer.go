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
	reconnectMu      sync.Mutex
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

	c := &Consumer{
		cfg:              cfg,
		client:           client,
		ch:               ch,
		log:              log,
		stopCh:           make(chan struct{}),
		messageProcessor: messageProcessor,
	}

	return c, nil
}

func (c *Consumer) Start() error {
	return c.runConsumeLoop()
}

func (c *Consumer) runConsumeLoop() error {
	start := time.Now()
	reconnectAttempts := 0
	reconnectTimeout := time.Duration(c.cfg.RabbitMqSettings.ReconnectTimeoutSeconds) * time.Second
	maxReconnects := c.cfg.RabbitMqSettings.MaxReconnectAttempts

	for {
		err := c.consume()
		if err == nil {
			return nil
		}

		select {
		case <-c.stopCh:
			return nil
		default:
		}

		reconnectAttempts++
		elapsed := time.Since(start)

		if elapsed > reconnectTimeout {
			c.log.Errorw("consumer.reconnect_timeout_exceeded",
				"queue", c.cfg.Queue,
				"elapsed", elapsed,
				"timeout", reconnectTimeout)
			return fmt.Errorf("reconnect timeout exceeded")
		}

		if reconnectAttempts > maxReconnects {
			c.log.Errorw("consumer.max_reconnects_exceeded",
				"queue", c.cfg.Queue,
				"attempts", reconnectAttempts)
			return fmt.Errorf("max reconnect attempts exceeded")
		}

		backoff := time.Duration(reconnectAttempts) * time.Second
		c.log.Warnw("consumer.reconnect_attempt",
			"queue", c.cfg.Queue,
			"attempt", reconnectAttempts,
			"backoff", backoff,
			"err", err)

		select {
		case <-time.After(backoff):
			continue
		case <-c.stopCh:
			return nil
		}
	}
}

func (c *Consumer) consume() error {
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
	c.resetTimer(batchTimeout)

	notifyClose := c.ch.NotifyClose(make(chan *amqp.Error))
	notifyCancel := c.ch.NotifyCancel(make(chan string))

	c.log.Infow("consumer.started", "queue", c.cfg.Queue)

	for {
		select {
		case <-c.stopCh:
			c.log.Infow("consumer.stopped", "queue", c.cfg.Queue)
			return nil
		case msg, ok := <-msgs:
			if !ok {
				c.log.Warnw("consumer.msg_channel_closed", "queue", c.cfg.Queue)
				return fmt.Errorf("channel closed")
			}
			c.handleMessage(msg)
		case err := <-notifyClose:
			if err != nil {
				c.log.Warnw("consumer.channel_closed_by_server", "queue", c.cfg.Queue, "err", err)
			} else {
				c.log.Warnw("consumer.channel_closed_by_server", "queue", c.cfg.Queue)
			}
			return fmt.Errorf("channel closed")
		case tag := <-notifyCancel:
			c.log.Warnw("consumer.canceled_by_server", "queue", c.cfg.Queue, "tag", tag)
			return fmt.Errorf("consumer canceled by server")
		}
	}
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

	requeue, err := c.messageProcessor.ProcessMessage(ctx, batch)
	if err != nil {
		c.log.Errorw("consumer.process_failed", "queue", c.cfg.Queue, "err", err, "need_requeue", requeue)
		last := batch[len(batch)-1]
		c.nack(last.DeliveryTag, true, requeue)
		return
	}

	last := batch[len(batch)-1]
	if err := c.ack(last.DeliveryTag, true); err != nil {
		c.log.Errorw("consumer.ack_failed", "queue", c.cfg.Queue, "err", err)
	}
}

func (c *Consumer) ack(tag uint64, multiple bool) error {
	if err := c.ensureChannel(); err != nil {
		return err
	}
	return c.ch.Ack(tag, multiple)
}

func (c *Consumer) nack(tag uint64, multiple, requeue bool) error {
	if err := c.ensureChannel(); err != nil {
		return err
	}
	return c.ch.Nack(tag, multiple, requeue)
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

func (c *Consumer) resetTimer(d time.Duration) {
	if c.timer != nil {
		c.timer.Stop()
	}
	c.timer = time.AfterFunc(d, func() {
		c.processBatch("timeout")
		c.resetTimer(d)
	})
}

func (c *Consumer) ensureChannel() error {
	c.reconnectMu.Lock()
	defer c.reconnectMu.Unlock()

	if c.ch != nil && !c.ch.IsClosed() {
		return nil
	}

	ch, err := c.client.Channel()
	if err != nil {
		return fmt.Errorf("reopen channel: %w", err)
	}
	c.ch = ch
	return nil
}
