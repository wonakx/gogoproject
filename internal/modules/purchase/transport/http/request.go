package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/gogoproject/internal/platform/fault"
)

type createRequest struct {
	SaleNo      int64 `json:"sale_no"`
	BuyMemberNo int64 `json:"buy_member_no"`
	Quantity    int32 `json:"quantity"`
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
	no, err := strconv.ParseInt(r.PathValue("purchase_no"), 10, 64)
	if err != nil || no <= 0 {
		return 0, fault.ErrInvalid
	}
	return no, nil
}
