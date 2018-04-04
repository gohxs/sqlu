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
type DBer interface {
	Begin() (*sql.Tx, error)
	SQLer
}

type RowScanner interface {
	Scan(...interface{}) error
}
type RowsScanner interface {
	RowScanner
	Columns() ([]string, error)
	Next() bool
}
