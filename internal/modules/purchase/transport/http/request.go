package httptransport

import (
	"encoding/json"
	"example.com/gogoproject/internal/platform/fault"
	"io"
	"net/http"
)

type createRequest struct {
	BuyerID  string `json:"buyer_id"`
	SaleID   string `json:"sale_id"`
	Quantity int64  `json:"quantity"`
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
