package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"workout_tracker_api/internal/exercise"
	"workout_tracker_api/internal/exerciselog"
	"workout_tracker_api/internal/user"
)

func New(
	exerciseHandler *exercise.Handler,
	userHandler *user.Handler,
	exerciseLogHandler *exerciselog.Handler,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/exercises", exerciseHandler.Create)
		r.Get("/exercises", exerciseHandler.List)
		r.Post("/users", userHandler.Create)
		r.Post("/exercise-logs", exerciseLogHandler.Create)
		r.Get("/users/{user_id}/exercise-logs", exerciseLogHandler.ListByUser)
		r.Get("/users/{user_id}/stats", exerciseLogHandler.GetStats)
	})

	return r
}
