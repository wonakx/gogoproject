package sale

import (
	"encoding/json"
	"testing"
)

func TestPriceExactDecimal(t *testing.T) {
	for _, tc := range []struct {
		input, output string
		units         Price
	}{{"0.01", "0.01", 1}, {"12000.5", "12000.50", 1200050}, {"99999999.99", "99999999.99", MaxPrice}} {
		p, err := ParsePrice(tc.input)
		if err != nil || p != tc.units {
			t.Fatalf("%s: %d %v", tc.input, p, err)
		}
		data, err := json.Marshal(p)
		if err != nil || string(data) != `"`+tc.output+`"` {
			t.Fatalf("wrong JSON: %s %v", data, err)
		}
		var decoded Price
		if err := json.Unmarshal(data, &decoded); err != nil || decoded != p {
			t.Fatalf("roundtrip: %d %v", decoded, err)
		}
	}
	for _, input := range []string{"100000000.00", "1.001", "-1.00", "1e2", "NaN", "", ".1", "1.", "+1"} {
		if _, err := ParsePrice(input); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
	var p Price
	for _, input := range []string{`12.50`, `null`, `"1.234"`} {
		if err := json.Unmarshal([]byte(input), &p); err == nil {
			t.Errorf("accepted JSON %s", input)
		}
	}
}
