// Package mysqlstore owns connections and shared transaction scope, not domain SQL.
package mysqlstore

import (
"context"
"database/sql"
"errors"
"fmt"
"time"

"example.com/gogoproject/internal/platform/fault"
"github.com/go-sql-driver/mysql"
)

type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Store struct{ db *sql.DB }

func New(db *sql.DB) *Store { return &Store{db: db} }

type transactionKey struct{ store *Store }

func (s *Store) Executor(ctx context.Context) Executor {
	if tx, ok := ctx.Value(transactionKey{s}).(*sql.Tx); ok {
		return tx
	}
	return s.db
}

// Within joins this store's transaction or starts one. Nested callbacks must propagate errors.
func (s *Store) Within(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(transactionKey{s}).(*sql.Tx); ok {
		return fn(ctx)
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = fn(context.WithValue(ctx, transactionKey{s}, tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func Error(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fault.ErrNotFound
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1062 || mysqlErr.Number == 1451 || mysqlErr.Number == 1452) {
		return fault.ErrConflict
	}
	return err
}

// Open requires a dealer DSN. Every connection uses UTC and DATETIME decoding.
func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("DB_DSN is required")
	}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_DSN: %w", err)
	}
	if cfg.DBName != "dealer" {
		return nil, errors.New("DB_DSN database must be dealer")
	}
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.MultiStatements = false
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 5 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 5 * time.Second
	}
	if cfg.Params == nil {
		cfg.Params = make(map[string]string)
	}
	cfg.Params["time_zone"] = "'+00:00'"
	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(3 * time.Minute)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect MySQL: %w", err)
	}
	if err := verifySchema(pingCtx, db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func verifySchema(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'dealer' AND TABLE_NAME IN ('dt_member', 'dt_product', 'dt_sale', 'dt_purchase', 'dt_purchase_history')")
	if err != nil {
		return fmt.Errorf("check dealer tables: %w", err)
	}
	expected := map[string]bool{"dt_member": false, "dt_product": false, "dt_sale": false, "dt_purchase": false, "dt_purchase_history": false}
	for rows.Next() {
		var name, engine string
		if err := rows.Scan(&name, &engine); err != nil {
			rows.Close()
			return err
		}
		if engine != "InnoDB" {
			rows.Close()
			return fmt.Errorf("dealer.%s must use InnoDB", name)
		}
		expected[name] = true
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for name, found := range expected {
		if !found {
			return fmt.Errorf("dealer.%s is missing; apply database schema first", name)
		}
	}
	// LIMIT 0 validates the approved added columns without reading application data.
	rows, err = db.QueryContext(ctx, "SELECT s.sell_member_no, p.quantity, h.quantity FROM dealer.dt_sale s, dealer.dt_purchase p, dealer.dt_purchase_history h LIMIT 0")
	if err != nil {
		return fmt.Errorf("dealer schema is incomplete; apply database migrations: %w", err)
	}
	return rows.Close()
}
