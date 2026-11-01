package purchase

import (
"context"

"example.com/gogoproject/internal/platform/mysqlstore"
)

type MySQLRepository struct{ store *mysqlstore.Store }

func NewMySQLRepository(store *mysqlstore.Store) *MySQLRepository {
	return &MySQLRepository{store: store}
}

var _ Repository = (*MySQLRepository)(nil)

func (r *MySQLRepository) CreateWithHistory(ctx context.Context, p Purchase) (Purchase, error) {
	err := r.store.Within(ctx, func(ctx context.Context) error {
		result, err := r.store.Executor(ctx).ExecContext(ctx, "INSERT INTO dealer.dt_purchase (sale_no, sell_member_no, buy_member_no, quantity, purchase_status, purchase_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?)", p.SaleNo, p.SellMemberNo, p.BuyMemberNo, p.Quantity, p.PurchaseStatus, p.PurchaseAt, p.CreateAt)
		if err != nil {
			return mysqlstore.Error(err)
		}
		p.PurchaseNo, err = result.LastInsertId()
		if err != nil {
			return err
		}
		_, err = r.store.Executor(ctx).ExecContext(ctx, "INSERT INTO dealer.dt_purchase_history (purchase_no, sale_no, sell_member_no, buy_member_no, quantity, purchase_status, purchase_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", p.PurchaseNo, p.SaleNo, p.SellMemberNo, p.BuyMemberNo, p.Quantity, p.PurchaseStatus, p.PurchaseAt, p.CreateAt)
		return mysqlstore.Error(err)
	})
	if err != nil {
		return Purchase{}, err
	}
	return p, nil
}
func (r *MySQLRepository) Get(ctx context.Context, no int64) (Purchase, error) {
	var p Purchase
	err := r.store.Executor(ctx).QueryRowContext(ctx, "SELECT purchase_no, sale_no, sell_member_no, buy_member_no, quantity, purchase_status, purchase_at, create_at FROM dealer.dt_purchase WHERE purchase_no = ?", no).Scan(&p.PurchaseNo, &p.SaleNo, &p.SellMemberNo, &p.BuyMemberNo, &p.Quantity, &p.PurchaseStatus, &p.PurchaseAt, &p.CreateAt)
	return p, mysqlstore.Error(err)
}
func (r *MySQLRepository) History(ctx context.Context, no int64) ([]History, error) {
	if _, err := r.Get(ctx, no); err != nil {
		return nil, err
	}
	rows, err := r.store.Executor(ctx).QueryContext(ctx, "SELECT purchase_history_no, purchase_no, sale_no, sell_member_no, buy_member_no, quantity, purchase_status, purchase_at, create_at FROM dealer.dt_purchase_history WHERE purchase_no = ? ORDER BY purchase_history_no", no)
	if err != nil {
		return nil, mysqlstore.Error(err)
	}
	defer rows.Close()
	result := make([]History, 0)
	for rows.Next() {
		var h History
		if err := rows.Scan(&h.PurchaseHistoryNo, &h.PurchaseNo, &h.SaleNo, &h.SellMemberNo, &h.BuyMemberNo, &h.Quantity, &h.PurchaseStatus, &h.PurchaseAt, &h.CreateAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, rows.Err()
}
