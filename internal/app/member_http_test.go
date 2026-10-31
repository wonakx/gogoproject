package app

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/gogoproject/internal/modules/member"
	"example.com/gogoproject/internal/platform/fault"
)

func TestMemberHTTPUpdateAndDelete(t *testing.T) {
	a, seller, buyer, _ := fixture(t, 3)
	h := a.Handler()
	request := func(method, id, body string, status int) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, "/api/v1/members/"+id, strings.NewReader(body)))
		if w.Code != status {
			t.Fatalf("%s: got %d, want %d: %s", method, w.Code, status, w.Body.String())
		}
		return w
	}
	request("PUT", buyer, `{"name":"Changed","email":"seller@example.com"}`, 409)
	for _, body := range []string{
		`{"id":"` + seller + `","name":"Changed","email":"new@example.com"}`,
		`{"name":"Changed","email":"bad"}`,
		`{"name":"Changed"}`,
		`null`,
	} {
		request("PUT", buyer, body, 400)
	}
	unchanged, err := a.Members.Get(buyer)
	if err != nil || unchanged.Name != "Buyer" || unchanged.Email != "buyer@example.com" {
		t.Fatalf("invalid request changed member: %+v %v", unchanged, err)
	}
	w := request("PUT", buyer, `{"name":" Changed ","email":" NEW@example.com "} garbage`, 200)
	var updated member.Member
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatal(err)
	}
	if updated.ID != buyer || updated.Name != "Changed" || updated.Email != "new@example.com" {
		t.Fatalf("unexpected updated member: %+v", updated)
	}
	if _, err := a.Members.Register("Reuse old email", "buyer@example.com"); err != nil {
		t.Fatal(err)
	}
	request("PUT", "missing", `{"name":"Changed","email":"another@example.com"}`, 404)
	w = request("DELETE", buyer, "", 204)
	if w.Body.Len() != 0 {
		t.Fatalf("204 must have no body: %s", w.Body.String())
	}
	if _, err := a.Members.Get(buyer); !errors.Is(err, fault.ErrNotFound) {
		t.Fatalf("deleted member still exists: %v", err)
	}
	request("DELETE", buyer, "", 404)
	if _, err := a.Members.Register("Reuse deleted email", "new@example.com"); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Members.Get(seller); err != nil {
		t.Fatalf("unrelated member changed: %v", err)
	}
}

func TestModuleHTTPInputHandling(t *testing.T) {
	a, seller, buyer, saleID := fixture(t, 3)
	h := a.Handler()
	sale, err := a.Sales.Get(saleID)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		path, body string
		status     int
	}{
		{"products", `{"name":"Mouse","description":"Wireless"}`, 201},
		{"products", `{"name":" "}`, 400},
		{"sales", `{"seller_id":"` + seller + `","product_id":"` + sale.ProductID + `","unit_price":100,"stock":2}`, 201},
		{"sales", `{"seller_id":"` + seller + `","product_id":"` + sale.ProductID + `","unit_price":0,"stock":2}`, 400},
		{"sales", `{"unknown":true}`, 400},
		{"purchases", `{"buyer_id":"` + buyer + `","sale_id":"` + saleID + `","quantity":2}`, 201},
		{"purchases", `{"buyer_id":"` + buyer + `","sale_id":"` + saleID + `","quantity":2}`, 409},
		{"purchases", `{"buyer_id":"` + buyer + `","sale_id":"` + saleID + `","quantity":0}`, 400},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("POST", "/api/v1/"+tc.path, strings.NewReader(tc.body)))
		if w.Code != tc.status {
			t.Errorf("%s: got %d, want %d: %s", tc.path, w.Code, tc.status, w.Body.String())
		}
	}
}
