package sqlu

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

//Errors
var (
	ErrNotPointer = errors.New("param is not a pointer")
)

const (
	tagField = 1
)

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
		// Parse fields
		if len(tags) > tagField {
			fields[tags[tagField]] = val.Field(i).Addr()
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
