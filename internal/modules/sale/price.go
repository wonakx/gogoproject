package sale

import (
"encoding/json"
"fmt"
"strconv"
"strings"

"example.com/gogoproject/internal/platform/fault"
)

// Price stores DECIMAL(10,2) exactly as hundredths, never as a floating point value.
// JSON uses a decimal string, e.g. "12000.50".
type Price int64

const MaxPrice Price = 9999999999

func ParsePrice(s string) (Price, error) {
	parts := strings.Split(s, ".")
	if len(parts) > 2 || len(parts[0]) == 0 || len(parts[0]) > 8 {
		return 0, fault.ErrInvalid
	}
	for _, part := range parts {
		if part == "" {
			return 0, fault.ErrInvalid
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return 0, fault.ErrInvalid
			}
		}
	}
	fraction := "00"
	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			return 0, fault.ErrInvalid
		}
		fraction = (parts[1] + "0")[:2]
	}
	whole, _ := strconv.ParseInt(parts[0], 10, 64)
	frac, _ := strconv.ParseInt(fraction, 10, 64)
	return Price(whole*100 + frac), nil
}
func (p Price) String() string               { return fmt.Sprintf("%d.%02d", int64(p)/100, int64(p)%100) }
func (p Price) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }
func (p *Price) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParsePrice(s)
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}
