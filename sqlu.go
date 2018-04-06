package sqlu

type RowScanner interface {
	Scan(...interface{}) error
}
type RowsScanner interface {
	RowScanner
	Columns() ([]string, error)
	Next() bool
	Close() error
}
