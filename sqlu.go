package sqlu

import (
	"context"
	"database/sql"
)

// SQLer wrapper for db, tx
type SQLer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type RowScanner interface {
	Scan(...interface{}) error
	Columns() ([]string, error)
	Next() bool
}
