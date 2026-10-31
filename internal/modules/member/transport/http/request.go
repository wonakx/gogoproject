package httptransport

import (
	"encoding/json"
	"io"
	"net/http"

	"example.com/gogoproject/internal/platform/fault"
)

type createRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type updateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// readCreateRequest checks HTTP syntax only; the service validates business values.
func readCreateRequest(w http.ResponseWriter, r *http.Request) (createRequest, error) {
	var in createRequest

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&in); err != nil {
		return in, fault.ErrInvalid
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return in, fault.ErrInvalid
	}

	return in, nil
}

// readUpdateRequest checks HTTP syntax only; the service validates business values.
func readUpdateRequest(w http.ResponseWriter, r *http.Request) (updateRequest, error) {
	var in updateRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		return in, fault.ErrInvalid
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return in, fault.ErrInvalid
	}
	return in, nil
}
