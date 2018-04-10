package sqlu

import (
	"context"
	"database/sql"
	"errors"
)

// TxQueryer Query only (no transaction initiator)
type TxQueryer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Queryer *sql.DB and transaction initiator
type Queryer interface {
	Begin() (*sql.Tx, error)
	TxQueryer
}

// Get Queryer
func Q(db interface{}) Queryer {
	switch r := db.(type) {
	case *sql.DB:
		return r
	case *sql.Tx:
		return txWrap{r}
	default:
		return nil
	}
}

type txWrap struct {
	TxQueryer
}

func (txWrap) Begin() (*sql.Tx, error) {
	return nil, errors.New("not Tx")
}
