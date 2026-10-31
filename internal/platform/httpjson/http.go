// Package httpjson provides shared JSON transport helpers, without business routing.
package httpjson

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"example.com/gogoproject/internal/platform/fault"
)

func Create[T any](fn func(T) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var in T
		if err := dec.Decode(&in); err != nil {
			Fail(w, fault.ErrInvalid)
			return
		}
		if err := dec.Decode(new(any)); err != io.EOF {
			Fail(w, fault.ErrInvalid)
			return
		}
		result, err := fn(in)
		if err != nil {
			Fail(w, err)
			return
		}
		Respond(w, http.StatusCreated, result)
	}
}
func Get[T any](fn func(string) (T, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := fn(r.PathValue("id"))
		if err != nil {
			Fail(w, err)
			return
		}
		Respond(w, http.StatusOK, result)
	}
}
func Fail(w http.ResponseWriter, err error) {
	status, message := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, fault.ErrInvalid):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, fault.ErrNotFound):
		status, message = http.StatusNotFound, "not found"
	case errors.Is(err, fault.ErrConflict):
		status, message = http.StatusConflict, "duplicate email or insufficient stock"
	default:
		slog.Error("request failed", "error", err)
	}
	Respond(w, status, map[string]string{"error": message})
}
func Respond(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
