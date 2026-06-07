package exerciselog

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidReference = errors.New("user or exercise not found")

type ExerciseLog struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	ExerciseID      string    `json:"exercise_id"`
	PerformedAt     time.Time `json:"performed_at"`
	Sets            *int      `json:"sets,omitempty"`
	Reps            *int      `json:"reps,omitempty"`
	WeightKg        *float64  `json:"weight_kg,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type Stats struct {
	Total     int64      `json:"total"`
	Today     int64      `json:"today"`
	Last7Days []DayCount `json:"last_7_days"`
}

type DayCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type Repository interface {
	Create(ctx context.Context, el *ExerciseLog) (*ExerciseLog, error)
	ListByUser(ctx context.Context, userID string) ([]ExerciseLog, error)
	GetStats(ctx context.Context, userID string) (*Stats, error)
}

type pgRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Create(ctx context.Context, el *ExerciseLog) (*ExerciseLog, error) {
	result := &ExerciseLog{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO exercise_logs
			(user_id, exercise_id, performed_at, sets, reps, weight_kg, duration_seconds, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, exercise_id, performed_at, sets, reps, weight_kg, duration_seconds, notes, created_at`,
		el.UserID, el.ExerciseID, el.PerformedAt,
		el.Sets, el.Reps, el.WeightKg, el.DurationSeconds, el.Notes,
	).Scan(
		&result.ID, &result.UserID, &result.ExerciseID, &result.PerformedAt,
		&result.Sets, &result.Reps, &result.WeightKg, &result.DurationSeconds,
		&result.Notes, &result.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrInvalidReference
		}
		return nil, err
	}
	return result, nil
}

func (r *pgRepository) GetStats(ctx context.Context, userID string) (*Stats, error) {
	var s Stats
	var raw json.RawMessage

	err := r.db.QueryRow(ctx, `
		WITH daily AS (
			SELECT
				date_trunc('day', performed_at)::date AS day,
				COUNT(*)                              AS count
			FROM exercise_logs
			WHERE user_id = $1
			  AND performed_at >= CURRENT_DATE - INTERVAL '6 days'
			GROUP BY day
		)
		SELECT
			(SELECT COUNT(*) FROM exercise_logs WHERE user_id = $1),
			(SELECT COUNT(*) FROM exercise_logs WHERE user_id = $1 AND performed_at >= CURRENT_DATE),
			COALESCE(
				(SELECT json_agg(json_build_object('date', day, 'count', count) ORDER BY day) FROM daily),
				'[]'::json
			)`,
		userID,
	).Scan(&s.Total, &s.Today, &raw)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, &s.Last7Days); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *pgRepository) ListByUser(ctx context.Context, userID string) ([]ExerciseLog, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, exercise_id, performed_at, sets, reps, weight_kg, duration_seconds, notes, created_at
		FROM exercise_logs
		WHERE user_id = $1
		ORDER BY performed_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ExerciseLog
	for rows.Next() {
		var el ExerciseLog
		if err := rows.Scan(
			&el.ID, &el.UserID, &el.ExerciseID, &el.PerformedAt,
			&el.Sets, &el.Reps, &el.WeightKg, &el.DurationSeconds,
			&el.Notes, &el.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, el)
	}
	return result, rows.Err()
}
