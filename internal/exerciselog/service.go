package exerciselog

import (
	"context"
	"errors"
	"time"

	"workout_tracker_api/internal/httputil"
)

var (
	ErrUserIDRequired     = errors.New("user_id is required")
	ErrUserIDInvalid      = errors.New("user_id must be a valid UUID")
	ErrExerciseIDRequired = errors.New("exercise_id is required")
	ErrExerciseIDInvalid  = errors.New("exercise_id must be a valid UUID")
)

type CreateInput struct {
	UserID          string
	ExerciseID      string
	PerformedAt     *time.Time
	Sets            *int
	Reps            *int
	WeightKg        *float64
	DurationSeconds *int
	Notes           *string
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*ExerciseLog, error) {
	if in.UserID == "" {
		return nil, ErrUserIDRequired
	}
	if !httputil.IsValidUUID(in.UserID) {
		return nil, ErrUserIDInvalid
	}
	if in.ExerciseID == "" {
		return nil, ErrExerciseIDRequired
	}
	if !httputil.IsValidUUID(in.ExerciseID) {
		return nil, ErrExerciseIDInvalid
	}

	performedAt := time.Now().UTC()
	if in.PerformedAt != nil {
		performedAt = *in.PerformedAt
	}

	return s.repo.Create(ctx, &ExerciseLog{
		UserID:          in.UserID,
		ExerciseID:      in.ExerciseID,
		PerformedAt:     performedAt,
		Sets:            in.Sets,
		Reps:            in.Reps,
		WeightKg:        in.WeightKg,
		DurationSeconds: in.DurationSeconds,
		Notes:           in.Notes,
	})
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]ExerciseLog, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) GetStats(ctx context.Context, userID string) (*Stats, error) {
	return s.repo.GetStats(ctx, userID)
}
