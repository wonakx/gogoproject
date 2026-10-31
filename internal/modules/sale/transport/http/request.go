package httptransport

import (
	"encoding/json"
	"example.com/gogoproject/internal/platform/fault"
	"net/http"
)

type createRequest struct {
	SellerID  string `json:"seller_id"`
	ProductID string `json:"product_id"`
	UnitPrice int64  `json:"unit_price"`
	Stock     int64  `json:"stock"`
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
	return in, nil
}
