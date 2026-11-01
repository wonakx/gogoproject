// Package httpapi connects the generated router, contract validation and API documentation.
// It contains no domain input DTOs or business policies.
package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"example.com/gogoproject/internal/api"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/legacy"
)

// validatedBody lets generated strict handlers bind the exact validated JSON keys.
type validatedBody struct {
	io.ReadCloser
	value any
}

func (b validatedBody) ValidatedValue() any { return b.value }

func NewHandler(server api.StrictServerInterface) http.Handler {
	spec, err := api.GetSpec()
	if err != nil {
		panic(err)
	}
	if err := spec.Validate(context.Background()); err != nil {
		panic(err)
	}
	router, err := legacy.NewRouter(spec)
	if err != nil {
		panic(err)
	}
	invalid := func(w http.ResponseWriter, r *http.Request, err error) {
		writeError(w, http.StatusBadRequest, "invalid input")
	}
	strict := api.NewStrictHandlerWithOptions(server, nil, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: invalid,
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			slog.Error("API response failed", "error", err)
			writeError(w, 500, "internal server error")
		},
	})
	generated := api.HandlerWithOptions(strict, api.StdHTTPServerOptions{ErrorHandlerFunc: invalid})
	validated := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, params, err := router.FindRoute(r)
		if err != nil {
			generated.ServeHTTP(w, r)
			return
		} // Preserve generated 404/405 routing.
		var value any
		if route.Operation.RequestBody != nil {
			original := r.Body
			defer original.Close()
			dec := json.NewDecoder(http.MaxBytesReader(w, original, 1<<20))
			dec.UseNumber()
			// Intentionally read only the first JSON value; ignore all trailing data.
			if err := dec.Decode(&value); err != nil {
				invalid(w, r, err)
				return
			}
			body, err := json.Marshal(value)
			if err != nil {
				invalid(w, r, err)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
			r.ContentLength = int64(len(body))
			if r.Header.Get("Content-Type") == "" {
				r.Header.Set("Content-Type", "application/json")
			}
		}
		input := &openapi3filter.RequestValidationInput{Request: r, PathParams: params, Route: route, Options: &openapi3filter.Options{SkipSettingDefaults: true}}
		if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
			invalid(w, r, err)
			return
		}
		if route.Operation.RequestBody != nil {
			r.Body = validatedBody{ReadCloser: r.Body, value: value}
		}
		generated.ServeHTTP(w, r)
	})
	mux := documentationMux(true)
	mux.Handle("/", validated)
	return mux
}
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: message})
}
