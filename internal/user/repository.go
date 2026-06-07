package user

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, name string) (*User, error)
}

type pgRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Create(ctx context.Context, name string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (name) VALUES ($1) RETURNING id, name, created_at`,
		name,
	).Scan(&u.ID, &u.Name, &u.CreatedAt)
	return u, err
}
