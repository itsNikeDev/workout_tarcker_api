package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"workout_tracker_api/internal/config"
	"workout_tracker_api/internal/exercise"
	"workout_tracker_api/internal/exerciselog"
	"workout_tracker_api/internal/router"
	"workout_tracker_api/internal/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := newPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	r := router.New(
		exercise.NewHandler(exercise.NewService(exercise.NewRepository(pool))),
		user.NewHandler(user.NewService(user.NewRepository(pool))),
		exerciselog.NewHandler(exerciselog.NewService(exerciselog.NewRepository(pool))),
	)

	log.Printf("listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func newPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
