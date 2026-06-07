package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"workout_tracker_api/internal/httputil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.svc.Create(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, ErrNameRequired) {
			httputil.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	httputil.JSON(w, http.StatusCreated, u)
}
