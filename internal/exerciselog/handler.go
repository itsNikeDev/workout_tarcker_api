package exerciselog

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"workout_tracker_api/internal/httputil"
)

type Handler struct {
	svc *Service
}

type createRequest struct {
	UserID          string     `json:"user_id"`
	ExerciseID      string     `json:"exercise_id"`
	PerformedAt     *time.Time `json:"performed_at"`
	Sets            *int       `json:"sets"`
	Reps            *int       `json:"reps"`
	WeightKg        *float64   `json:"weight_kg"`
	DurationSeconds *int       `json:"duration_seconds"`
	Notes           *string    `json:"notes"`
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	el, err := h.svc.Create(r.Context(), CreateInput{
		UserID:          req.UserID,
		ExerciseID:      req.ExerciseID,
		PerformedAt:     req.PerformedAt,
		Sets:            req.Sets,
		Reps:            req.Reps,
		WeightKg:        req.WeightKg,
		DurationSeconds: req.DurationSeconds,
		Notes:           req.Notes,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrUserIDRequired), errors.Is(err, ErrExerciseIDRequired),
			errors.Is(err, ErrUserIDInvalid), errors.Is(err, ErrExerciseIDInvalid):
			httputil.Error(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidReference):
			httputil.Error(w, http.StatusUnprocessableEntity, err.Error())
		default:
			httputil.Error(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	httputil.JSON(w, http.StatusCreated, el)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if !httputil.IsValidUUID(userID) {
		httputil.Error(w, http.StatusBadRequest, "user_id must be a valid UUID")
		return
	}
	s, err := h.svc.GetStats(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	httputil.JSON(w, http.StatusOK, s)
}

func (h *Handler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if !httputil.IsValidUUID(userID) {
		httputil.Error(w, http.StatusBadRequest, "user_id must be a valid UUID")
		return
	}
	logs, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	if logs == nil {
		logs = []ExerciseLog{}
	}
	httputil.JSON(w, http.StatusOK, logs)
}
