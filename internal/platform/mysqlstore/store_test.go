package mysqlstore

import (
"context"
"database/sql"
"errors"
"testing"

"example.com/gogoproject/internal/platform/fault"
"github.com/DATA-DOG/go-sqlmock"
"github.com/go-sql-driver/mysql"
)

func TestNestedTransactionUsesSameExecutor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := New(db)
	mock.ExpectBegin()
	mock.ExpectCommit()
	err = store.Within(context.Background(), func(ctx context.Context) error {
		outer := store.Executor(ctx)
		if outer == db {
			t.Fatal("executor is not a transaction")
		}
		return store.Within(ctx, func(inner context.Context) error {
			if store.Executor(inner) != outer {
				t.Fatal("nested transaction differs")
			}
			return nil
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
func TestSchemaRequiresTransactionalTablesAndAddedColumns(t *testing.T) {
	for _, variant := range []string{"ok", "missing", "nontransactional", "columns"} {
		t.Run(variant, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			rows := sqlmock.NewRows([]string{"TABLE_NAME", "ENGINE"})
			for _, name := range []string{"dt_member", "dt_product", "dt_sale", "dt_purchase", "dt_purchase_history"} {
				if variant == "missing" && name == "dt_sale" {
					continue
				}
				engine := "InnoDB"
				if variant == "nontransactional" && name == "dt_sale" {
					engine = "MyISAM"
				}
				rows.AddRow(name, engine)
			}
			mock.ExpectQuery("SELECT TABLE_NAME, ENGINE").WillReturnRows(rows).RowsWillBeClosed()
			if variant == "ok" {
				mock.ExpectQuery("SELECT s.sell_member_no").WillReturnRows(sqlmock.NewRows([]string{"sell_member_no", "quantity", "quantity"})).RowsWillBeClosed()
			}
			if variant == "columns" {
				mock.ExpectQuery("SELECT s.sell_member_no").WillReturnError(&mysql.MySQLError{Number: 1054})
			}
			err = verifySchema(context.Background(), db)
			if (variant == "ok") != (err == nil) {
				t.Fatalf("got %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestConfigurationAndErrors(t *testing.T) {
	if _, err := Open(context.Background(), ""); err == nil {
		t.Fatal("accepted missing DSN")
	}
	if _, err := Open(context.Background(), "user:pass@tcp(localhost:3306)/wrong"); err == nil {
		t.Fatal("accepted wrong database")
	}
	if !errors.Is(Error(sql.ErrNoRows), fault.ErrNotFound) {
		t.Fatal("not-found mapping")
	}
	if !errors.Is(Error(&mysql.MySQLError{Number: 1062}), fault.ErrConflict) {
		t.Fatal("duplicate mapping")
	}
	if Error(context.Canceled) != context.Canceled {
		t.Fatal("cancellation lost")
	}
}
