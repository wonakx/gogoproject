package app

import (
	"errors"
	"math"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"example.com/gogoproject/internal/modules/purchase"
	"example.com/gogoproject/internal/platform/fault"
)

func fixture(t *testing.T, stock int64) (*App, string, string, string) {
	t.Helper()
	a := New()
	seller, err := a.Members.Register("Seller", "seller@example.com")
	if err != nil {
		t.Fatal(err)
	}
	buyer, err := a.Members.Register("Buyer", "buyer@example.com")
	if err != nil {
		t.Fatal(err)
	}
	p, err := a.Products.Create("Keyboard", "Mechanical")
	if err != nil {
		t.Fatal(err)
	}
	s, err := a.Sales.Create(seller.ID, p.ID, 12000, stock)
	if err != nil {
		t.Fatal(err)
	}
	return a, seller.ID, buyer.ID, s.ID
}
func TestPurchaseFlow(t *testing.T) {
	a, seller, buyer, saleID := fixture(t, 3)
	p, err := a.Purchases.Create(buyer, saleID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if p.TotalPrice != 24000 || p.SellerID != seller || p.ProductID == "" {
		t.Fatalf("wrong snapshot: %+v", p)
	}
	saved, err := a.Purchases.Get(p.ID)
	if err != nil || saved != p {
		t.Fatalf("purchase not saved: %+v %v", saved, err)
	}
	s, _ := a.Sales.Get(saleID)
	if s.Stock != 1 {
		t.Fatalf("stock = %d", s.Stock)
	}
}
func TestConcurrentPurchasesCannotOversell(t *testing.T) {
	a, _, buyer, saleID := fixture(t, 10)
	var successes atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := a.Purchases.Create(buyer, saleID, 1)
			if err == nil {
				successes.Add(1)
			} else if !errors.Is(err, fault.ErrConflict) {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()
	s, _ := a.Sales.Get(saleID)
	if successes.Load() != 10 || s.Stock != 0 {
		t.Fatalf("successes=%d stock=%d", successes.Load(), s.Stock)
	}
}
func TestInvalidPurchasesDoNotChangeStock(t *testing.T) {
	a, seller, buyer, saleID := fixture(t, 3)
	for _, tc := range []struct {
		buyer    string
		quantity int64
		want     error
	}{
		{seller, 1, fault.ErrInvalid}, {buyer, 0, fault.ErrInvalid}, {buyer, -1, fault.ErrInvalid},
		{buyer, 4, fault.ErrConflict}, {buyer, math.MaxInt64, fault.ErrInvalid}, {"missing", 1, fault.ErrNotFound},
	} {
		_, err := a.Purchases.Create(tc.buyer, saleID, tc.quantity)
		if !errors.Is(err, tc.want) {
			t.Errorf("got %v, want %v", err, tc.want)
		}
	}
	s, _ := a.Sales.Get(saleID)
	if s.Stock != 3 {
		t.Fatalf("stock = %d", s.Stock)
	}
}

type failingRepository struct{}

func (failingRepository) Create(purchase.Purchase) error { return errors.New("storage unavailable") }
func (failingRepository) Get(string) (purchase.Purchase, error) {
	return purchase.Purchase{}, fault.ErrNotFound
}
func TestFailedOrderRestoresStock(t *testing.T) {
	a, _, buyer, saleID := fixture(t, 3)
	service := purchase.NewService(failingRepository{}, directoryAdapter{a.Members, a.Products}, salesAdapter{a.Sales})
	if _, err := service.Create(buyer, saleID, 2); err == nil {
		t.Fatal("expected storage error")
	}
	s, _ := a.Sales.Get(saleID)
	if s.Stock != 3 {
		t.Fatalf("stock = %d", s.Stock)
	}
}
func TestReferencesAndUniqueEmail(t *testing.T) {
	a, seller, _, _ := fixture(t, 3)
	if _, err := a.Members.Register("Duplicate", "SELLER@example.com"); !errors.Is(err, fault.ErrConflict) {
		t.Fatal(err)
	}
	if _, err := a.Sales.Create(seller, "missing", 100, 1); !errors.Is(err, fault.ErrNotFound) {
		t.Fatal(err)
	}
}
func TestHTTPValidation(t *testing.T) {
	h := New().Handler()
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{"GET", "/health", "", 200},
		{"POST", "/api/v1/members", `{"name":"Member","email":"member@example.com"}`, 201},
		{"POST", "/api/v1/members", `{"name":"Member","email":"MEMBER@example.com"}`, 409},
		{"POST", "/api/v1/products", `{"name":"Keyboard","unknown":true}`, 400},
		{"POST", "/api/v1/products", `{"name":"Keyboard"} {}`, 400},
		{"POST", "/api/v1/products", `null`, 400},
		{"POST", "/api/v1/purchases", `{"buyer_id":"x","sale_id":"y","quantity":1.5}`, 400},
		{"GET", "/api/v1/members/missing", "", 404},
	} {
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != tc.status {
			t.Errorf("%s %s: got %d, want %d: %s", tc.method, tc.path, w.Code, tc.status, w.Body.String())
		}
		if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatal("missing JSON content type")
		}
	}
}

func TestVersionedResourceRoutes(t *testing.T) {
	a, seller, buyer, saleID := fixture(t, 3)
	order, err := a.Purchases.Create(buyer, saleID, 1)
	if err != nil {
		t.Fatal(err)
	}
	h := a.Handler()
	for _, path := range []string{
		"/api/v1/members/" + seller,
		"/api/v1/products/" + order.ProductID,
		"/api/v1/sales/" + saleID,
		"/api/v1/purchases/" + order.ID,
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 200 {
			t.Errorf("GET %s: %d %s", path, w.Code, w.Body.String())
		}
	}
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/members/" + seller, 404},
		{"GET", "/api/v2/members/" + seller, 404},
		{"DELETE", "/api/v1/members/" + seller, 405},
		{"GET", "/api/v1/health", 404},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status {
			t.Errorf("%s %s: got %d, want %d", tc.method, tc.path, w.Code, tc.status)
		}
	}
}
