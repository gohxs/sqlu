package sqlu

import (
	"context"
	"database/sql"
	"errors"
)

//Errors
var (
	ErrNotPointer = errors.New("param is not a pointer")
)

const (
	tagField     = 0
	tagOmitEmpty = 1
	tagType      = 2
)

// SQLer interface for "sql" package
// Can be sql.DB or sql.Tx
type SQLer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Tabler Interface to return the table name of a struct
type Tabler interface {
	Table() string
}
