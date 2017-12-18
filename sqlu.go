package sqlu

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//Errors
var (
	ErrNotPointer = errors.New("param is not a pointer")
)

const (
	tagField = 0
)

// SQLer interface for "sql" package
// Can be sql.DB or sql.Tx
type SQLer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Interface to return the table name of a struct
type Tabler interface {
	Table() string
}

// Insert a tabler
// just a struct with a metho
func Insert(db SQLer, data Tabler) (sql.Result, error) {
	return InsertContext(context.Background(), db, data)
}

// InsertContext insert a struct with Table() with context
func InsertContext(ctx context.Context, db SQLer, data Tabler) (sql.Result, error) {
	return TableInsertContext(ctx, db, data.Table(), data)
}

// Insert struct into specified table
// Warning table is not escaped
//
func TableInsert(db SQLer, table string, data interface{}) (sql.Result, error) {
	return TableInsertContext(context.Background(), db, table, data)
}

// InsertContext using a context
func TableInsertContext(ctx context.Context, db SQLer, table string, data interface{}) (sql.Result, error) {
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	fields := []string{}
	values := []interface{}{}
	// Go through sqlu tags
	for i := 0; i < val.NumField(); i++ {
		f := typ.Field(i)
		if !val.Field(i).CanInterface() {
			continue
		}

		tags := strings.Split(f.Tag.Get("sqlu"), ",")

		if len(tags) >= tagField {

			fields = append(fields, "\""+tags[tagField]+"\"")
			values = append(values, val.Field(i).Interface())
		}
	}
	qry := fmt.Sprintf(
		"INSERT INTO \"%s\" (%s) values(%s)",
		table,
		strings.Join(fields, ","),
		strings.Repeat("?, ", len(fields)-1)+"?",
	)

	return db.ExecContext(ctx, qry, values...)
}

// Scan a full object
// usefull when objects perfectly matches a table
func Scan(res *sql.Rows, data interface{}) error {

	valPtr := reflect.ValueOf(data)
	if valPtr.Type().Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	val := valPtr.Elem()
	params := []interface{}{}
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanInterface() {
			continue
		}
		params = append(params, val.Field(i).Addr().Interface())
	}
	return res.Scan(params...)
}

// ScanNamed fields into a struct
// This will reflect the structure and match column names
//
func ScanNamed(res *sql.Rows, data interface{}) error {
	// Data should be pointer
	fields := map[string]interface{}{} // Pointer interface

	valPtr := reflect.ValueOf(data)
	if valPtr.Type().Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	val := valPtr.Elem()
	typ := valPtr.Elem().Type()
	// Retrieve the fields into a map
	for i := 0; i < val.NumField(); i++ {
		f := typ.Field(i)
		tags := strings.Split(f.Tag.Get("sqlu"), ",")

		if !val.Field(i).CanInterface() {
			continue
		}
		// Parse fields
		if len(tags) >= tagField {
			fields[tags[tagField]] = val.Field(i).Addr().Interface()
		} else {
			fields[strings.ToLower(typ.Field(i).Name)] = val.Field(i).Addr().Interface()
		}

	}

	params := []interface{}{}
	colTyp, err := res.ColumnTypes()
	if err != nil {
		return err
	}
	for _, c := range colTyp {
		colName := strings.ToLower(c.Name())
		v, ok := fields[strings.ToLower(colName)]
		if !ok {
			if c.ScanType() != nil {
				v = reflect.New(c.ScanType()) // Dumb pointer
			} else {
				v = ""
			}
		}
		params = append(params, v)
	}
	return res.Scan(params...)
}
