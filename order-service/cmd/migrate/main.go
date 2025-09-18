package main

import (
	"log"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := goose.OpenDBWithDriver("pgx", cfg.Db.MigrationConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}
}
