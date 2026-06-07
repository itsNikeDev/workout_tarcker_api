package exercise

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Exercise struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, name string, description *string) (*Exercise, error)
	List(ctx context.Context) ([]Exercise, error)
}

type pgRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Create(ctx context.Context, name string, description *string) (*Exercise, error) {
	e := &Exercise{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO exercises (name, description) VALUES ($1, $2)
		 RETURNING id, name, description, created_at`,
		name, description,
	).Scan(&e.ID, &e.Name, &e.Description, &e.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}
	return e, nil
}

func (r *pgRepository) List(ctx context.Context) ([]Exercise, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, description, created_at FROM exercises ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Exercise
	for rows.Next() {
		var e Exercise
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}
