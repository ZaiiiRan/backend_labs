package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	client "github.com/ZaiiiRan/backend_labs/order-service/internal/client/grpc"
	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/messages"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type OrderCreatedConsumer struct {
	client     *rabbitmq.RabbitMqClient
	grpcClient *client.OmsGrpcClient
	cfg        *config.RabbitMqConsumerSettings
	ch         *amqp.Channel
	stopCh     chan struct{}
	log        *zap.SugaredLogger
}

func NewOrderCreatedConsumer(
	cfg *config.RabbitMqConsumerSettings,
	client *rabbitmq.RabbitMqClient,
	omsGrpcClient *client.OmsGrpcClient,
	log *zap.SugaredLogger,
) (*OrderCreatedConsumer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	return &OrderCreatedConsumer{
		client:     client,
		grpcClient: omsGrpcClient,
		cfg:        cfg,
		ch:         ch,
		stopCh:     make(chan struct{}),
		log:        log,
	}, nil
}

func (c *OrderCreatedConsumer) Start() error {
	if err := c.ensureChannel(); err != nil {
		return err
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

	if err := c.ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("set qos: %w", err)
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

	go func() {
		c.log.Infow("order_created_consumer.started", "queue", c.cfg.Queue)
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

	req := &pb.AuditLogOrderBatchCreateRequest{}
	for _, item := range order.OrderItems {
		req.Orders = append(req.Orders, &pb.LogOrder{
			OrderId:     order.Id,
			OrderItemId: item.Id,
			CustomerId:  order.CustomerID,
			OrderStatus: bll.ORDER_STATUS_CREATED.String(),
		})
	}

	if _, err := c.grpcClient.LogOrder(context.Background(), req); err != nil {
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
