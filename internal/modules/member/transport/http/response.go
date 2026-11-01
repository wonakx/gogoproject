package httptransport

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"example.com/gogoproject/internal/platform/fault"
)

func fail(w http.ResponseWriter, err error) {
	status, message := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, fault.ErrInvalid):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, fault.ErrNotFound):
		status, message = http.StatusNotFound, "not found"
	case errors.Is(err, fault.ErrConflict):
		status, message = http.StatusConflict, "duplicate member id"
	default:
		slog.Error("member request failed", "error", err)
	}
	respond(w, status, map[string]string{"error": message})
}

func respond(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
