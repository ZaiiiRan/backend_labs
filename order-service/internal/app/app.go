package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/postgres"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	httpserver "github.com/ZaiiiRan/backend_labs/order-service/internal/server/http"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/server/http/controllers"
)

type App struct {
	cfg        *config.Config
	httpServer *httpserver.Server
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg}, nil
}

func (a *App) Run() error {
	ctx := context.Background()
	pool, err := postgres.NewPgxPool(ctx, a.cfg.DbSettings.ConnectionString)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	publisher, err := rabbitmq.NewPublisher(&a.cfg.RabbitMqSettings)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer publisher.Close()

	orderController := controllers.NewOrderController(pool, publisher)

	a.httpServer = httpserver.NewServer(a.cfg.Http.Port, orderController)

	go func() {
		log.Printf("HTTP server listening on %s", a.httpServer.Addr())
		if err := a.httpServer.Start(); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.httpServer.Stop(ctx); err != nil {
		return err
	}

	log.Println("Server exited properly")
	return nil
}
