package consumer

import (
	"encoding/json"
	"fmt"

	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	client "github.com/ZaiiiRan/backend_labs/order-service/internal/client/http"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	dto "github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/messages"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type OrderCreatedConsumer struct {
	client     *rabbitmq.RabbitMqClient
	httpClient *client.OmsHttpClient
	queue      string
	ch         *amqp.Channel
	stopCh     chan struct{}
	log        *zap.SugaredLogger
}

func NewOrderCreatedConsumer(
	client *rabbitmq.RabbitMqClient, 
	omsHttpClient *client.OmsHttpClient, 
	queue string, 
	log *zap.SugaredLogger,
) (*OrderCreatedConsumer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	return &OrderCreatedConsumer{
		client:     client,
		httpClient: omsHttpClient,
		queue:      queue,
		ch:         ch,
		stopCh:     make(chan struct{}),
		log: log,
	}, nil
}

func (c *OrderCreatedConsumer) Start() error {
	if err := c.ensureChannel(); err != nil {
		return err
	}

	_, err := c.ch.QueueDeclare(
		c.queue,
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
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		c.log.Infow("order_created_consumer.started", "queue", c.queue)
		for {
			select {
			case <-c.stopCh:
				c.log.Infow("order_created_consumer.stopped")
				return
			case msg, ok := <-msgs:
				if !ok {
					c.log.Warnw("order_created_consumer.channel_closed")
					return
				}
				c.handleMessage(msg)
			}
		}
	}()

	return nil
}

func (c *OrderCreatedConsumer) handleMessage(msg amqp.Delivery) {
	var order messages.OrderCreatedMessage
	if err := json.Unmarshal(msg.Body, &order); err != nil {
		c.log.Errorw("order_created_consumer.unmarshal_failed", "err", err)
		msg.Nack(false, false)
		return
	}

	c.log.Infow("order_created_consumer.message_received",
		"order_id", order.Id,
		"customer_id", order.CustomerID,
		"items_count", len(order.OrderItems),
	)

	req := &dto.V1CreateAuditLogOrderRequest{}
	for _, item := range order.OrderItems {
		req.Orders = append(req.Orders, dto.V1LogOrder{
			OrderId:     order.Id,
			OrderItemId: item.Id,
			CustomerId:  order.CustomerID,
			OrderStatus: bll.ORDER_STATUS_CREATED.String(),
		})
	}

	if _, err := c.httpClient.LogOrder(req); err != nil {
		c.log.Errorw("order_created_consumer.log_order_failed", "err", err)
		msg.Nack(false, true)
		return
	}

	if err := msg.Ack(false); err != nil {
		c.log.Errorw("order_created_consumer.ack_failed", "err", err)
		return
	}

	c.log.Infow("order_created_consumer.message_processed",
		"order_id", order.Id,
		"ack", true,
	)
}

func (c *OrderCreatedConsumer) Stop() {
	close(c.stopCh)
	if c.ch != nil {
		c.ch.Close()
	}
}

func (c *OrderCreatedConsumer) ensureChannel() error {
	if c.ch == nil || c.ch.IsClosed() {
		ch, err := c.client.Channel()
		if err != nil {
			return fmt.Errorf("reopen channel: %w", err)
		}
		c.ch = ch
	}
	return nil
}
