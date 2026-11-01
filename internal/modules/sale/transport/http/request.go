package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"example.com/gogoproject/internal/modules/sale"
	"example.com/gogoproject/internal/platform/fault"
)

type createRequest struct {
	SellMemberNo int64      `json:"sell_member_no"`
	ProductNo    int64      `json:"product_no"`
	Name         *string    `json:"name"`
	Description  *string    `json:"description"`
	Quantity     int32      `json:"quantity"`
	Price        sale.Price `json:"price"`
	Status       string     `json:"status"`
	StartAt      time.Time  `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
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
	no, err := strconv.ParseInt(r.PathValue("sale_no"), 10, 64)
	if err != nil || no <= 0 {
		return 0, fault.ErrInvalid
	}
	return no, nil
}
