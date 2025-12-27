package app

import (
	"context"

	client "github.com/ZaiiiRan/backend_labs/order-service/internal/client/grpc"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/consumer"
	dalconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	kafkaconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer/kafka"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/logger"
	"go.uber.org/zap"
)

type ConsumerApp struct {
	cfg *config.ConsumerConfig
	log *zap.SugaredLogger

	orderCreatedConsumer       *kafkaconsumer.Consumer
	orderStatusChangedConsumer *kafkaconsumer.Consumer

	orderCreatedMessageProcessor       dalconsumer.MessageProcessor
	orderStatusChangedMessageProcessor dalconsumer.MessageProcessor

	omsClient *client.OmsGrpcClient
}

func NewConsumerApp() (*ConsumerApp, error) {
	cfg, err := config.LoadConsumerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	return &ConsumerApp{cfg: cfg, log: log}, nil
}

func (a *ConsumerApp) Run(ctx context.Context) error {
	if err := a.initOmsGrpcClient(); err != nil {
		return err
	}
	a.initOrderCreatedMessageProcessor()
	if err := a.initOrderCreatedConsumer(); err != nil {
		return err
	}
	a.initOrderStatusChangedMessageProcessor()
	if err := a.initOrderStatusChangedConsumer(); err != nil {
		return err
	}
	if err := a.startOrderCreatedConsumer(ctx); err != nil {
		return err
	}
	if err := a.startOrderStatusChangedConsumer(ctx); err != nil {
		return err
	}

	a.log.Infow("app.started")
	return nil
}

func (a *ConsumerApp) Stop() {
	a.log.Infow("app.stopping")

	a.orderCreatedConsumer.Close()
	a.orderStatusChangedConsumer.Close()

	a.omsClient.Close()

	a.log.Infow("app.stopped")
}

func (a *ConsumerApp) initOmsGrpcClient() error {
	client, err := client.NewOmsGrpcClient(a.cfg.OmsClientGrpcSettings)
	if err != nil {
		a.log.Errorw("app.create_oms_grpc_client_failed", "err", err)
		return err
	}
	a.omsClient = client
	return nil
}

func (a *ConsumerApp) initOrderCreatedMessageProcessor() {
	a.orderCreatedMessageProcessor = consumer.NewOrderCreatedMessageProcessor(a.omsClient, a.log)
}

func (a *ConsumerApp) initOrderStatusChangedMessageProcessor() {
	a.orderStatusChangedMessageProcessor = consumer.NewOrderStatusChangedMessageProcessor(a.omsClient, a.log)
}

func (a *ConsumerApp) initOrderCreatedConsumer() error {
	orderCreatedConsumer, err := kafkaconsumer.NewConsumer(
		&a.cfg.KafkaConsumerSettings, a.cfg.KafkaConsumerSettings.OrderCreatedTopic, a.orderCreatedMessageProcessor, a.log,
	)
	if err != nil {
		a.log.Errorw("app.init_order_created_consumer_failed", "err", err)
		return err
	}
	a.orderCreatedConsumer = orderCreatedConsumer
	return nil
}

func (a *ConsumerApp) initOrderStatusChangedConsumer() error {
	orderStatusChanged, err := kafkaconsumer.NewConsumer(
		&a.cfg.KafkaConsumerSettings, a.cfg.KafkaConsumerSettings.OrderStatusChangedTopic, a.orderStatusChangedMessageProcessor, a.log,
	)
	if err != nil {
		a.log.Errorw("app.init_order_status_changed_consumer_failed", "err", err)
		return err
	}
	a.orderStatusChangedConsumer = orderStatusChanged
	return nil
}

func (a *ConsumerApp) startOrderCreatedConsumer(ctx context.Context) error {
	go func() {
		if err := a.orderCreatedConsumer.Run(ctx); err != nil {
			a.log.Fatalw("app.start_order_created_consumer_failed", "err", err)
		}
	}()
	return nil
}

func (a *ConsumerApp) startOrderStatusChangedConsumer(ctx context.Context) error {
	go func() {
		if err := a.orderStatusChangedConsumer.Run(ctx); err != nil {
			a.log.Fatalw("app.start_order_status_changed_consumer_failed", "err", err)
		}
	}()
	return nil
}
