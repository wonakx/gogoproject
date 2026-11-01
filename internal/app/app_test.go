package app

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"example.com/gogoproject/internal/modules/member"
	"example.com/gogoproject/internal/modules/product"
	"example.com/gogoproject/internal/modules/purchase"
	"example.com/gogoproject/internal/modules/sale"
	"example.com/gogoproject/internal/platform/fault"
)

func fixture(t *testing.T, quantity int32) (*App, int64, int64, int64) {
	t.Helper()
	a := New()
	seller, err := a.Members.Register(member.CreateInput{ID: "seller"})
	if err != nil {
		t.Fatal(err)
	}
	buyer, err := a.Members.Register(member.CreateInput{ID: "buyer"})
	if err != nil {
		t.Fatal(err)
	}
	p, err := a.Products.Create(product.CreateInput{ProductCode: "KEYBOARD-001"})
	if err != nil {
		t.Fatal(err)
	}
	price, err := sale.ParsePrice("12000.50")
	if err != nil {
		t.Fatal(err)
	}
	listing, err := a.Sales.Create(sale.CreateInput{SellMemberNo: seller.MemberNo, ProductNo: p.ProductNo, Quantity: quantity, Price: price, Status: sale.StatusActive, StartAt: time.Now().UTC().Add(-time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	return a, seller.MemberNo, buyer.MemberNo, listing.SaleNo
}
func TestPurchaseAndHistory(t *testing.T) {
	a, seller, buyer, saleNo := fixture(t, 3)
	p, err := a.Purchases.Create(purchase.CreateInput{SaleNo: saleNo, BuyMemberNo: buyer, Quantity: 2})
	if err != nil {
		t.Fatal(err)
	}
	if p.PurchaseNo <= 0 || p.SellMemberNo != seller || p.BuyMemberNo != buyer || p.Quantity != 2 || p.PurchaseStatus != purchase.StatusCompleted || p.CreateAt.IsZero() || p.PurchaseAt.IsZero() {
		t.Fatalf("wrong purchase: %+v", p)
	}
	saved, err := a.Purchases.Get(p.PurchaseNo)
	if err != nil || saved != p {
		t.Fatalf("not persisted: %+v %v", saved, err)
	}
	histories, err := a.Purchases.History(p.PurchaseNo)
	if err != nil || len(histories) != 1 {
		t.Fatalf("history: %+v %v", histories, err)
	}
	h := histories[0]
	if h.PurchaseNo == nil || *h.PurchaseNo != p.PurchaseNo || h.Quantity != 2 || h.SellMemberNo != seller || h.PurchaseStatus != p.PurchaseStatus {
		t.Fatalf("wrong history: %+v", h)
	}
	*h.PurchaseNo = 999
	again, _ := a.Purchases.History(p.PurchaseNo)
	if *again[0].PurchaseNo != p.PurchaseNo {
		t.Fatal("history pointer escaped repository")
	}
	s, _ := a.Sales.Get(saleNo)
	if s.Quantity != 1 || s.UpdateAt == nil {
		t.Fatalf("wrong stock: %+v", s)
	}
}
func TestConcurrentPurchasesCannotOversell(t *testing.T) {
	a, _, buyer, saleNo := fixture(t, 10)
	var successes atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := a.Purchases.Create(purchase.CreateInput{SaleNo: saleNo, BuyMemberNo: buyer, Quantity: 1})
			if err == nil {
				successes.Add(1)
			} else if !errors.Is(err, fault.ErrConflict) {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()
	s, _ := a.Sales.Get(saleNo)
	if successes.Load() != 10 || s.Quantity != 0 {
		t.Fatalf("successes=%d quantity=%d", successes.Load(), s.Quantity)
	}
}
func TestInvalidPurchasesPreserveStock(t *testing.T) {
	a, seller, buyer, saleNo := fixture(t, 3)
	for _, tc := range []struct {
		buyer    int64
		quantity int32
		want     error
	}{
		{seller, 1, fault.ErrInvalid}, {buyer, 0, fault.ErrInvalid}, {buyer, -1, fault.ErrInvalid}, {buyer, 4, fault.ErrConflict}, {999, 1, fault.ErrNotFound},
	} {
		_, err := a.Purchases.Create(purchase.CreateInput{SaleNo: saleNo, BuyMemberNo: tc.buyer, Quantity: tc.quantity})
		if !errors.Is(err, tc.want) {
			t.Errorf("got %v want %v", err, tc.want)
		}
	}
	s, _ := a.Sales.Get(saleNo)
	if s.Quantity != 3 {
		t.Fatalf("quantity=%d", s.Quantity)
	}
}

type failingRepository struct{}

func (failingRepository) CreateWithHistory(purchase.Purchase) (purchase.Purchase, error) {
	return purchase.Purchase{}, errors.New("storage unavailable")
}
func (failingRepository) Get(int64) (purchase.Purchase, error) {
	return purchase.Purchase{}, fault.ErrNotFound
}
func (failingRepository) History(int64) ([]purchase.History, error) { return nil, fault.ErrNotFound }
func TestFailedPurchaseRestoresStock(t *testing.T) {
	a, _, buyer, saleNo := fixture(t, 3)
	service := purchase.NewService(failingRepository{}, directoryAdapter{a.Members, a.Products}, salesAdapter{a.Sales})
	if _, err := service.Create(purchase.CreateInput{SaleNo: saleNo, BuyMemberNo: buyer, Quantity: 2}); err == nil {
		t.Fatal("expected error")
	}
	s, _ := a.Sales.Get(saleNo)
	if s.Quantity != 3 {
		t.Fatalf("quantity=%d", s.Quantity)
	}
}
func TestSaleReferencesAndAvailability(t *testing.T) {
	a, seller, buyer, saleNo := fixture(t, 3)
	existing, _ := a.Sales.Get(saleNo)
	base := sale.CreateInput{SellMemberNo: seller, ProductNo: existing.ProductNo, Quantity: 1, Price: 100, Status: sale.StatusActive, StartAt: time.Now().Add(-time.Hour)}
	for _, tc := range []struct {
		modify func(*sale.CreateInput)
		want   error
	}{
		{func(in *sale.CreateInput) { in.ProductNo = 999 }, fault.ErrNotFound},
		{func(in *sale.CreateInput) { in.SellMemberNo = 999 }, fault.ErrNotFound},
		{func(in *sale.CreateInput) { in.Price = sale.MaxPrice + 1 }, fault.ErrInvalid},
		{func(in *sale.CreateInput) { in.Status = "unknown" }, fault.ErrInvalid},
		{func(in *sale.CreateInput) { end := in.StartAt.Add(-time.Hour); in.EndAt = &end }, fault.ErrInvalid},
	} {
		in := base
		tc.modify(&in)
		if _, err := a.Sales.Create(in); !errors.Is(err, tc.want) {
			t.Errorf("got %v want %v", err, tc.want)
		}
	}
	for _, modify := range []func(*sale.CreateInput){
		func(in *sale.CreateInput) { in.Status = sale.StatusClosed },
		func(in *sale.CreateInput) { in.StartAt = time.Now().Add(time.Hour) },
		func(in *sale.CreateInput) { end := time.Now().Add(-time.Minute); in.EndAt = &end },
	} {
		in := base
		modify(&in)
		s, err := a.Sales.Create(in)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := a.Purchases.Create(purchase.CreateInput{SaleNo: s.SaleNo, BuyMemberNo: buyer, Quantity: 1}); !errors.Is(err, fault.ErrConflict) {
			t.Errorf("unavailable sale: %v", err)
		}
	}
}
func TestProductCodeUniqueAndNullableFields(t *testing.T) {
	a := New()
	p, err := a.Products.Create(product.CreateInput{ProductCode: "SKU"})
	if err != nil {
		t.Fatal(err)
	}
	if p.ProductNo <= 0 || p.Name != nil || p.Description != nil || p.UpdateAt != nil || p.CreateAt.IsZero() {
		t.Fatalf("wrong product: %+v", p)
	}
	if _, err := a.Products.Create(product.CreateInput{ProductCode: "SKU"}); !errors.Is(err, fault.ErrConflict) {
		t.Fatal(err)
	}
}
