package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/gogoproject/internal/platform/fault"
)

type createRequest struct {
	ProductCode string  `json:"product_code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

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
func readNumber(r *http.Request) (int64, error) {
	no, err := strconv.ParseInt(r.PathValue("product_no"), 10, 64)
	if err != nil || no <= 0 {
		return 0, fault.ErrInvalid
	}
	return no, nil
}
