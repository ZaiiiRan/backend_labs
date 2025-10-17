package app

import (
	client "github.com/ZaiiiRan/backend_labs/order-service/internal/client/http"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	consumer "github.com/ZaiiiRan/backend_labs/order-service/internal/consumer/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/logger"
	"go.uber.org/zap"
)

type ConsumerApp struct {
	cfg *config.ConsumerConfig
	log *zap.SugaredLogger

	rabbitmqClient *rabbitmq.RabbitMqClient

	orderCreatedConsumer *consumer.OrderCreatedConsumer

	omsClient *client.OmsHttpClient
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
	if err := a.initRabbitMqClient(); err != nil {
		return err
	}
	a.initOmsHttpClient()
	if err := a.initOrderCreatedConsumer(); err != nil {
		return err
	}
	if err := a.startOrderCreatedConsumer(); err != nil {
		return err
	}

	a.log.Infow("app.started")
	return nil
}

func (a *ConsumerApp) Stop() {
	a.log.Infow("app.stopping")

	a.rabbitmqClient.Close()

	a.log.Infow("app.stopped")
}

func (a *ConsumerApp) initRabbitMqClient() error {
	rabbitMqClient, err := rabbitmq.NewRabbitMqClient(&a.cfg.RabbitMqSettings)
	if err != nil {
		a.log.Errorw("app.rabbitmq_connect_failed", "err", err)
	}
	a.rabbitmqClient = rabbitMqClient
	return nil
}

func (a *ConsumerApp) initOmsHttpClient() {
	a.omsClient = client.NewOmsHttpClient(a.cfg.OmsClientHttpSettings)
}

func (a *ConsumerApp) initOrderCreatedConsumer() error {
	orderCreatedConsumer, err := consumer.NewOrderCreatedConsumer(a.rabbitmqClient, a.omsClient, a.cfg.RabbitMqSettings.OrderCreatedQueue, a.log)
	if err != nil {
		a.log.Errorw("app.create_order_created_consumer_failed", "err", err)
	}
	a.orderCreatedConsumer = orderCreatedConsumer
	return nil
}

func (a *ConsumerApp) startOrderCreatedConsumer() error {
	if err := a.orderCreatedConsumer.Start(); err != nil {
		a.log.Errorw("app.start_order_created_consumer_failed", "err", err)
		return err
	}
	return nil
}
