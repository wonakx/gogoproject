package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/gogoproject/internal/platform/fault"
)

type createRequest struct {
	ID          string  `json:"id"`
	Email       *string `json:"email"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phone_number"`
}
type updateRequest struct {
	ID          string  `json:"id"`
	Email       *string `json:"email"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phone_number"`
}

func readCreateRequest(w http.ResponseWriter, r *http.Request) (createRequest, error) {
	var in createRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&in); err != nil {
		return in, fault.ErrInvalid
	}
	return in, nil
}
func readUpdateRequest(w http.ResponseWriter, r *http.Request) (updateRequest, error) {
	var in updateRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&in); err != nil {
		return in, fault.ErrInvalid
	}
	return in, nil
}
func readNumber(r *http.Request) (int64, error) {
	no, err := strconv.ParseInt(r.PathValue("member_no"), 10, 64)
	if err != nil || no <= 0 {
		return 0, fault.ErrInvalid
	}
	return no, nil
}
