package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/services"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/postgres"
	repositories "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/repositories/postgres"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
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

	uow := unitofwork.New(pool)
	orderRepo := repositories.NewOrderRepository(uow)
	orderItemRepo := repositories.NewOrderItemRepository(uow)

	orderService := services.NewOrderService(uow, orderRepo, orderItemRepo)
	orderController := controllers.NewOrderController(orderService)

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
