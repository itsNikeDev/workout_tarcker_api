package main

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"context"

	"workout_tracker_api/migrations"
)

// Usage: go run ./cmd/migrate [up|down|status|reset]
func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err := goose.RunContext(context.Background(), command, db, "."); err != nil {
		log.Fatalf("goose %s: %v", command, err)
	}
}
