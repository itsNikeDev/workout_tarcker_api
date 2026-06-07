package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
