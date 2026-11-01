package app

import (
"crypto/rand"
"encoding/json"
"fmt"
"net/http/httptest"
"strings"
"testing"
)

func TestSchemaHTTPFlow(t *testing.T) {
	h := integrationApp(t).Handler()
	token := rand.Text()
	request := func(method, path, body string, status int) map[string]any {
		t.Helper()
		body = strings.NewReplacer(`"seller"`, `"seller-`+token+`"`, `"buyer"`, `"buyer-`+token+`"`, `"renamed"`, `"renamed-`+token+`"`, `P001`, `P001-`+token).Replace(body)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
		if w.Code != status {
			t.Fatalf("%s %s: got %d want %d: %s", method, path, w.Code, status, w.Body.String())
		}
		if status == 204 {
			if w.Body.Len() != 0 {
				t.Fatal("204 has body")
			}
			return nil
		}
		var out map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		return out
	}
	seller := request("POST", "/api/v1/members", `{"id":"seller","email":null,"address":"Seoul","phone_number":null}`, 201)
	buyer := request("POST", "/api/v1/members", `{"id":"buyer"} ignored`, 201)
	if seller["email"] != nil || seller["update_at"] != nil {
		t.Fatal("nullable fields not null")
	}
	request("POST", "/api/v1/members", `{"id":"seller"}`, 409)
	product := request("POST", "/api/v1/products", `{"product_code":"P001","name":null,"description":null}`, 201)
	request("POST", "/api/v1/products", `{"product_code":"P001"}`, 409)
	sale := request("POST", "/api/v1/sales", fmt.Sprintf(`{"sell_member_no":%.0f,"product_no":%.0f,"quantity":5,"price":"123.45","status":"ACTIVE","start_at":"2020-01-01T00:00:00Z"}`, seller["member_no"], product["product_no"]), 201)
	if sale["price"] != "123.45" {
		t.Fatalf("wrong price: %v", sale["price"])
	}
	purchase := request("POST", "/api/v1/purchases", fmt.Sprintf(`{"sale_no":%.0f,"buy_member_no":%.0f,"quantity":2}`, sale["sale_no"], buyer["member_no"]), 201)
	if purchase["sell_member_no"] != seller["member_no"] || purchase["quantity"] != float64(2) {
		t.Fatalf("wrong purchase: %+v", purchase)
	}
	got := request("GET", fmt.Sprintf("/api/v1/sales/%.0f", sale["sale_no"]), "", 200)
	if got["quantity"] != float64(3) {
		t.Fatalf("wrong stock: %+v", got)
	}
	request("GET", fmt.Sprintf("/api/v1/purchases/%.0f", purchase["purchase_no"]), "", 200)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", fmt.Sprintf("/api/v1/purchases/%.0f/histories", purchase["purchase_no"]), nil))
	var history []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &history); err != nil || w.Code != 200 || len(history) != 1 {
		t.Fatalf("history: %s %v", w.Body.String(), err)
	}
	if history[0]["quantity"] != float64(2) || history[0]["purchase_no"] != purchase["purchase_no"] {
		t.Fatalf("wrong history: %+v", history)
	}
	updated := request("PUT", fmt.Sprintf("/api/v1/members/%.0f", buyer["member_no"]), `{"id":"renamed","email":null,"address":null,"phone_number":"010"} {}`, 200)
	if updated["id"] != "renamed-"+token || updated["update_at"] == nil {
		t.Fatalf("wrong update: %+v", updated)
	}
	request("DELETE", fmt.Sprintf("/api/v1/members/%.0f", buyer["member_no"]), "", 204)
	request("GET", fmt.Sprintf("/api/v1/members/%.0f", buyer["member_no"]), "", 404)
}
func TestHTTPInvalidInputs(t *testing.T) {
	h := integrationApp(t).Handler()
	token := rand.Text()
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{"GET", "/health", "", 200},
		{"GET", "/api/v1/members/not-a-number", "", 400},
		{"GET", "/api/v1/members/9223372036854775808", "", 400},
		{"GET", "/api/v1/products/0", "", 400},
		{"GET", "/api/v1/sales/-1", "", 400},
		{"GET", "/api/v1/purchases/9223372036854775807", "", 404},
		{"POST", "/api/v1/members", `{"id":""}`, 400},
		{"POST", "/api/v1/products", `{"product_code":"P-` + token + `"} garbage`, 201},
		{"POST", "/api/v1/sales", `{"quantity":2147483648}`, 400},
		{"POST", "/api/v1/sales", `{"price":12.34}`, 400},
		{"POST", "/api/v1/purchases", `{"quantity":1.5}`, 400},
		{"POST", "/api/v1/purchases", `null`, 400},
		{"GET", "/members/1", "", 404},
		{"GET", "/api/v2/members/1", "", 404},
		{"PATCH", "/api/v1/members/1", "", 405},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
		if w.Code != tc.status {
			t.Errorf("%s %s: got %d want %d: %s", tc.method, tc.path, w.Code, tc.status, w.Body.String())
		}
	}
}
