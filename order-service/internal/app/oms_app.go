package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/postgres"
	publisher "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/publisher/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/logger"
	grpcserver "github.com/ZaiiiRan/backend_labs/order-service/internal/server/grpc"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/server/grpc/services"
	grpcgateway "github.com/ZaiiiRan/backend_labs/order-service/internal/server/grpc_gateway"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type OmsApp struct {
	cfg *config.ServerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	rabbitmqClient *rabbitmq.RabbitMqClient

	omsPublisher *publisher.Publisher

	orderService *services.OrderService

	grpcServer  *grpcserver.Server
	grpcGateway *grpcgateway.Server
}

func NewOmsApp() (*OmsApp, error) {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	return &OmsApp{cfg: cfg, log: log}, nil
}

func (a *OmsApp) Run(ctx context.Context) error {
	if err := a.initPostgresClient(ctx); err != nil {
		return err
	}
	if err := a.initRabbitMqClient(); err != nil {
		return err
	}
	if err := a.initPublishers(); err != nil {
		return err
	}
	a.initOrderService()
	if err := a.initGrpcServer(); err != nil {
		return err
	}
	a.startGrpcServer()
	if err := a.initGrpcGateway(ctx); err != nil {
		return err
	}
	a.startGrpcGateway()
	a.log.Infow("app.started")
	return nil
}

func (a *OmsApp) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	a.postgresClient.Close()
	a.omsPublisher.Close()
	a.rabbitmqClient.Close()
	
	a.grpcServer.Stop(shCtx)
	a.grpcGateway.Stop(shCtx)

	a.log.Infow("app.stopped")
}

func (a *OmsApp) initPostgresClient(ctx context.Context) error {
	pgClient, err := postgres.NewPostgresClient(ctx, a.cfg.DbSettings.ConnectionString)
	if err != nil {
		a.log.Errorw("app.postgres_connect_failed", "err", err)
		return err
	}
	a.postgresClient = pgClient
	return nil
}

func (a *OmsApp) initRabbitMqClient() error {
	rabbitMqClient, err := rabbitmq.NewRabbitMqClient(&a.cfg.OmsRabbitMqPublisherSettings.RabbitMqSettings)
	if err != nil {
		a.log.Errorw("app.rabbitmq_connect_failed", "err", err)
	}
	a.rabbitmqClient = rabbitMqClient
	return nil
}

func (a *OmsApp) initPublishers() error {
	orderCreatedPublisher, err := publisher.NewPublisher(&a.cfg.OmsRabbitMqPublisherSettings, a.rabbitmqClient)
	if err != nil {
		a.log.Errorw("app.create_order_created_publisher_failed", "err", err)
		return err
	}
	a.omsPublisher = orderCreatedPublisher
	return nil
}

func (a *OmsApp) initOrderService() {
	a.orderService = services.NewOrderService(a.postgresClient, a.omsPublisher, a.log)
}

func (a *OmsApp) initGrpcServer() error {
	srv, err := grpcserver.NewServer(a.cfg.Grpc.Port, a.orderService)
	if err != nil {
		a.log.Errorw("app.grpc.server_init_failed", "err", err)
		return err
	}

	a.grpcServer = srv
	return nil
}

func (a *OmsApp) startGrpcServer() {
	go func() {
		a.log.Infow("app.grpc.serve_start", "port", a.cfg.Grpc.Port)
		if err := a.grpcServer.Start(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			a.log.Fatalw("app.grpc.serve_error", "err", err)
		}
	}()
}

func (a *OmsApp) initGrpcGateway(ctx context.Context) error {
	srv, err := grpcgateway.NewServer(ctx, a.cfg.Http.Port, a.cfg.Grpc.Port)
	if err != nil {
		a.log.Errorw("app.http.gateway_init_failed", "err", err)
		return err
	}
	a.grpcGateway = srv
	return nil
}

func (a *OmsApp) startGrpcGateway() {
	go func() {
		a.log.Infow("app.http.gateway_start", "port", a.cfg.Http.Port)
		if err := a.grpcGateway.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalw("app.http.gateway_error", "err", err)
		}
	}()
}
