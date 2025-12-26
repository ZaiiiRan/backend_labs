package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcclient "github.com/ZaiiiRan/backend_labs/order-generator/internal/client/grpc"
	"github.com/ZaiiiRan/backend_labs/order-generator/internal/config"
	"github.com/ZaiiiRan/backend_labs/order-generator/internal/generator"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := grpcclient.NewOmsGrpcClient(cfg.OmsClientGrpcSettings)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	customerIDs := []int64{1, 2, 3, 4, 5}

	gen := generator.NewGenerator(client, customerIDs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	gen.Start(ctx)
}
