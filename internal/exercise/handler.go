package exercise

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

type createRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	e, err := h.svc.Create(r.Context(), req.Name, req.Description)
	if err != nil {
		switch {
		case errors.Is(err, ErrNameRequired):
			httputil.Error(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrAlreadyExists):
			httputil.Error(w, http.StatusConflict, err.Error())
		default:
			httputil.Error(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	httputil.JSON(w, http.StatusCreated, e)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.svc.List(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	if exercises == nil {
		exercises = []Exercise{}
	}
	httputil.JSON(w, http.StatusOK, exercises)
}
