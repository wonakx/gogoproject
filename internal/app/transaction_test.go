package app

import (
"context"
"errors"
"net/http/httptest"
"strings"
"testing"
"time"

"example.com/gogoproject/internal/modules/purchase"
"example.com/gogoproject/internal/platform/fault"
"github.com/DATA-DOG/go-sqlmock"
)

func TestPurchaseUsesOneTransaction(t *testing.T) {
	for _, failure := range []string{"", "stock", "purchase", "history", "commit"} {
		t.Run(failure, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			a := New(db)
			now := time.Now().UTC().Truncate(time.Second)
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .* FROM dealer.dt_member WHERE member_no").WithArgs(int64(2)).WillReturnRows(sqlmock.NewRows([]string{"member_no", "id", "email", "address", "phone_number", "create_at", "update_at"}).AddRow(2, "buyer", nil, nil, nil, now, nil))
			quantity := 10
			if failure == "stock" {
				quantity = 0
			}
			mock.ExpectQuery("SELECT .* FROM dealer.dt_sale WHERE sale_no = .* FOR UPDATE").WithArgs(int64(3)).WillReturnRows(sqlmock.NewRows([]string{"sale_no", "sell_member_no", "name", "description", "product_no", "quantity", "price", "status", "start_at", "end_at", "create_at", "update_at"}).AddRow(3, 1, nil, nil, 4, quantity, []byte("123.45"), "ACTIVE", now.Add(-time.Hour), nil, now, nil))
			injected := errors.New("injected database failure")
			if failure != "stock" {
				mock.ExpectExec("UPDATE dealer.dt_sale SET quantity").WithArgs(int32(2), sqlmock.AnyArg(), int64(3)).WillReturnResult(sqlmock.NewResult(0, 1))
				insert := mock.ExpectExec("INSERT INTO dealer.dt_purchase ").WithArgs(int64(3), int64(1), int64(2), int32(2), "COMPLETED", sqlmock.AnyArg(), sqlmock.AnyArg())
				if failure == "purchase" {
					insert.WillReturnError(injected)
				} else {
					insert.WillReturnResult(sqlmock.NewResult(5, 1))
					history := mock.ExpectExec("INSERT INTO dealer.dt_purchase_history").WithArgs(int64(5), int64(3), int64(1), int64(2), int32(2), "COMPLETED", sqlmock.AnyArg(), sqlmock.AnyArg())
					if failure == "history" {
						history.WillReturnError(injected)
					} else {
						history.WillReturnResult(sqlmock.NewResult(6, 1))
					}
				}
			}
			switch failure {
			case "":
				mock.ExpectCommit()
			case "commit":
				mock.ExpectCommit().WillReturnError(injected)
			default:
				mock.ExpectRollback()
			}
			got, err := a.Purchases.Create(context.Background(), purchase.CreateInput{SaleNo: 3, BuyMemberNo: 2, Quantity: 2})
			if failure == "" {
				if err != nil || got.PurchaseNo != 5 {
					t.Fatalf("%+v %v", got, err)
				}
			} else {
				if err == nil || got.PurchaseNo != 0 {
					t.Fatalf("unexpected success: %+v %v", got, err)
				}
			}
			if failure == "stock" && !errors.Is(err, fault.ErrConflict) {
				t.Fatal(err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestHTTPInvalidInputDoesNotQueryDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	h := New(db).Handler()
	for _, tc := range []struct{ method, path, body string }{
		{"GET", "/api/v1/members/nope", ""},
		{"GET", "/api/v1/products/0", ""},
		{"POST", "/api/v1/members", `{"id":""}`},
		{"POST", "/api/v1/purchases", `{"quantity":1.5}`},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
		if w.Code != 400 {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
