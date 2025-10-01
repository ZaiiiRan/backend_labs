package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	connStr := cfg.DbSettings.MigrationConnectionString
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("ping: %v", err)
	}

	goose.SetTableName("goose_db_version")
	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatalf("goose up: %v", err)
	}
	fmt.Println("Migrations applied successfully")
}
