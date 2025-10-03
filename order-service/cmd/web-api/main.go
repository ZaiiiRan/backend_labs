// @title Order Service API
// @version 1.0
// @description API for managing orders
// @host localhost:5000
// @BasePath /
package main

import (
	"log"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
