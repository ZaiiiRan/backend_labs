// @title Order Service API
// @version 1.0
// @description API for managing orders
// @host localhost:5000
// @BasePath /
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		os.Exit(1)
	}

	<-ctx.Done()
	a.Stop(context.Background())
}
