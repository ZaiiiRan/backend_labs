package app

import (
	client "github.com/ZaiiiRan/backend_labs/order-service/internal/client/grpc"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/consumer"
	dalconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	rabbitmqconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/logger"
	"go.uber.org/zap"
)

type ConsumerApp struct {
	cfg *config.ConsumerConfig
	log *zap.SugaredLogger

	orderCreatedRabbitmqClient       *rabbitmq.RabbitMqClient
	orderStatusChangedRabbitmqClient *rabbitmq.RabbitMqClient

	orderCreatedConsumer               *rabbitmqconsumer.Consumer
	orderStatusChangedConsumer         *rabbitmqconsumer.Consumer
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

func (a *ConsumerApp) Run() error {
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
	if err := a.startOrderCreatedConsumer(); err != nil {
		return err
	}
	if err := a.startOrderStatusChangedConsumer(); err != nil {
		return err
	}

	a.log.Infow("app.started")
	return nil
}

func (a *ConsumerApp) Stop() {
	a.log.Infow("app.stopping")

	a.orderCreatedConsumer.Stop()
	a.orderCreatedRabbitmqClient.Close()
	a.orderStatusChangedConsumer.Stop()
	a.orderStatusChangedRabbitmqClient.Close()
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
	rabbitMqClient, err := rabbitmq.NewRabbitMqClient(&a.cfg.OrderCreatedRabbitMqConsumerSettings.RabbitMqSettings)
	if err != nil {
		a.log.Errorw("app.rabbitmq_connect_failed", "err", err)
	}
	a.orderCreatedRabbitmqClient = rabbitMqClient

	orderCreatedConsumer, err := rabbitmqconsumer.NewConsumer(&a.cfg.OrderCreatedRabbitMqConsumerSettings, a.orderCreatedRabbitmqClient,
		a.orderCreatedMessageProcessor, a.log)
	if err != nil {
		a.log.Errorw("app.create_order_created_consumer_failed", "err", err)
	}
	a.orderCreatedConsumer = orderCreatedConsumer
	return nil
}

func (a *ConsumerApp) initOrderStatusChangedConsumer() error {
	rabbitMqClient, err := rabbitmq.NewRabbitMqClient(&a.cfg.OrderStatusChangedRabbitMqConsumerSettings.RabbitMqSettings)
	if err != nil {
		a.log.Errorw("app.rabbitmq_connect_failed", "err", err)
	}
	a.orderStatusChangedRabbitmqClient = rabbitMqClient

	orderStatusChangedConsumer, err := rabbitmqconsumer.NewConsumer(&a.cfg.OrderStatusChangedRabbitMqConsumerSettings, a.orderStatusChangedRabbitmqClient,
		a.orderStatusChangedMessageProcessor, a.log)
	if err != nil {
		a.log.Errorw("app.create_order_status_changed_consumer_failed", "err", err)
	}
	a.orderStatusChangedConsumer = orderStatusChangedConsumer
	return nil
}

func (a *ConsumerApp) startOrderCreatedConsumer() error {
	go func() {
		if err := a.orderCreatedConsumer.Start(); err != nil {
			a.log.Fatalw("app.start_order_created_consumer_failed", "err", err)
		}
	}()
	return nil
}

func (a *ConsumerApp) startOrderStatusChangedConsumer() error {
	go func() {
		if err := a.orderStatusChangedConsumer.Start(); err != nil {
			a.log.Fatalw("app.start_order_status_changed_consumer_failed", "err", err)
		}
	}()
	return nil
}
