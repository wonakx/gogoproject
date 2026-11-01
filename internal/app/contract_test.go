package app

import (
"context"
"errors"
"net/http/httptest"
"strings"
"testing"
"time"

"example.com/gogoproject/internal/api"
"github.com/DATA-DOG/go-sqlmock"
"github.com/getkin/kin-openapi/openapi3filter"
"github.com/getkin/kin-openapi/routers/legacy"
"github.com/go-sql-driver/mysql"
)

func TestHTTPResponsesMatchOpenAPI(t *testing.T) {
	now := time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)
	memberRows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"member_no", "id", "email", "address", "phone_number", "create_at", "update_at"}).AddRow(1, "buyer", nil, nil, nil, now, nil)
	}
	purchaseRows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"purchase_no", "sale_no", "sell_member_no", "buy_member_no", "quantity", "purchase_status", "purchase_at", "create_at"}).AddRow(1, 2, 3, 4, 2, "COMPLETED", now, now)
	}
	for _, tc := range []struct {
		name, method, path, body string
		status                   int
		setup                    func(sqlmock.Sqlmock)
	}{
		{"health", "GET", "/health", "", 200, func(sqlmock.Sqlmock) {}},
		{"create member ignores trailing data and case variants", "POST", "/api/v1/members", `{"id":"buyer","ID":"must-not-override","email":null} garbage`, 201, func(m sqlmock.Sqlmock) {
			m.ExpectExec("INSERT INTO dealer.dt_member").WithArgs("buyer", nil, nil, nil, sqlmock.AnyArg(), nil).WillReturnResult(sqlmock.NewResult(1, 1))
		}},
		{"member nullable response", "GET", "/api/v1/members/1", "", 200, func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM dealer.dt_member").WithArgs(int64(1)).WillReturnRows(memberRows())
		}},
		{"member full update", "PUT", "/api/v1/members/1", `{"id":"buyer"} {}`, 200, func(m sqlmock.Sqlmock) {
			m.ExpectBegin()
			m.ExpectQuery("SELECT .* FOR UPDATE").WillReturnRows(memberRows())
			m.ExpectExec("UPDATE dealer.dt_member").WillReturnResult(sqlmock.NewResult(0, 1))
			m.ExpectCommit()
		}},
		{"member deletion has no body", "DELETE", "/api/v1/members/1", "", 204, func(m sqlmock.Sqlmock) {
			m.ExpectExec("DELETE FROM dealer.dt_member").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))
		}},
		{"product nullable response", "GET", "/api/v1/products/1", "", 200, func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM dealer.dt_product").WillReturnRows(sqlmock.NewRows([]string{"product_no", "product_code", "name", "description", "create_at", "update_at"}).AddRow(1, "P001", nil, nil, now, nil))
		}},
		{"sale price and zero remaining stock", "GET", "/api/v1/sales/1", "", 200, func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM dealer.dt_sale").WillReturnRows(sqlmock.NewRows([]string{"sale_no", "sell_member_no", "name", "description", "product_no", "quantity", "price", "status", "start_at", "end_at", "create_at", "update_at"}).AddRow(1, 2, nil, nil, 3, 0, []byte("12000.50"), "ACTIVE", now, nil, now, nil))
		}},
		{"purchase response", "GET", "/api/v1/purchases/1", "", 200, func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM dealer.dt_purchase WHERE").WillReturnRows(purchaseRows())
		}},
		{"history list", "GET", "/api/v1/purchases/1/histories", "", 200, func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM dealer.dt_purchase WHERE").WillReturnRows(purchaseRows())
			m.ExpectQuery("SELECT .* FROM dealer.dt_purchase_history").WillReturnRows(sqlmock.NewRows([]string{"purchase_history_no", "purchase_no", "sale_no", "sell_member_no", "buy_member_no", "quantity", "purchase_status", "purchase_at", "create_at"}).AddRow(1, 1, 2, 3, 4, 2, "COMPLETED", now, now))
		}},
		{"not found", "GET", "/api/v1/members/1", "", 404, func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM dealer.dt_member").WillReturnRows(sqlmock.NewRows([]string{"member_no", "id", "email", "address", "phone_number", "create_at", "update_at"}))
		}},
		{"duplicate", "POST", "/api/v1/members", `{"id":"buyer"}`, 409, func(m sqlmock.Sqlmock) {
			m.ExpectExec("INSERT INTO dealer.dt_member").WillReturnError(&mysql.MySQLError{Number: 1062})
		}},
		{"internal error is redacted", "GET", "/api/v1/members/1", "", 500, func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM dealer.dt_member").WillReturnError(errors.New("private database detail"))
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			tc.setup(mock)
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			New(db).Handler().ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Fatalf("got %d want %d: %s", w.Code, tc.status, w.Body.String())
			}
			if tc.status == 204 && w.Body.Len() != 0 {
				t.Fatal("204 response has a body")
			}
			if tc.status == 500 && strings.Contains(w.Body.String(), "private") {
				t.Fatal("internal error leaked")
			}
			spec, err := api.GetSpec()
			if err != nil {
				t.Fatal(err)
			}
			router, err := legacy.NewRouter(spec)
			if err != nil {
				t.Fatal(err)
			}
			route, params, err := router.FindRoute(req)
			if err != nil {
				t.Fatal(err)
			}
			input := (&openapi3filter.ResponseValidationInput{RequestValidationInput: &openapi3filter.RequestValidationInput{Request: req, PathParams: params, Route: route}, Status: w.Code, Header: w.Header(), Options: &openapi3filter.Options{IncludeResponseStatus: true}}).SetBodyBytes(w.Body.Bytes())
			if err := openapi3filter.ValidateResponse(context.Background(), input); err != nil {
				t.Fatal(err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestOpenAPIValidationRejectsBeforeDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	h := New(db).Handler()
	for _, tc := range []struct{ method, path, body string }{
		{"POST", "/api/v1/members", `{}`},
		{"POST", "/api/v1/members", `{"ID":"wrong-case"}`},
		{"POST", "/api/v1/members", `{"id":null}`},
		{"POST", "/api/v1/members", `{"id":"x","email":"` + strings.Repeat("a", 201) + `"}`},
		{"POST", "/api/v1/members", `{"id":"x","address":"` + strings.Repeat("a", 1<<20) + `"}`},
		{"POST", "/api/v1/members", `null`},
		{"POST", "/api/v1/products", `{"product_code":"P","unexpected":true}`},
		{"POST", "/api/v1/sales", `{"sell_member_no":1,"product_no":1,"quantity":1,"price":"1.001","status":"ACTIVE","start_at":"2026-01-01T00:00:00Z"}`},
		{"POST", "/api/v1/sales", `{"sell_member_no":1,"product_no":1,"quantity":1,"price":"1.00","status":"INVALID","start_at":"2026-01-01T00:00:00Z"}`},
		{"POST", "/api/v1/sales", `{"sell_member_no":1,"product_no":1,"quantity":1,"price":"1.00","status":"ACTIVE","start_at":"not-a-date"}`},
		{"POST", "/api/v1/purchases", `{"sale_no":1,"buy_member_no":2,"quantity":2147483648}`},
		{"POST", "/api/v1/purchases", `{"sale_no":1,"buy_member_no":2,"quantity":0}`},
		{"GET", "/api/v1/members/0", ""},
		{"GET", "/api/v1/members/9223372036854775808", ""},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
		if w.Code != 400 {
			t.Errorf("%s: got %d %s", tc.path, w.Code, w.Body.String())
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
