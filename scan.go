package sqlu

import (
	"database/sql"
	"reflect"
	"strings"
)

func RowToMap(res *sql.Rows) (map[string]interface{}, error) {

	data := map[string]interface{}{}
	scanParams := []interface{}{}
	colTypes, err := res.ColumnTypes()
	if err != nil {
		return nil, err
	}
	for _, ct := range colTypes {
		p := reflect.New(ct.ScanType())
		scanParams = append(scanParams, p.Interface())
	}
	err = res.Scan(scanParams...)
	if err != nil {
		return nil, err
	}
	for i, ct := range colTypes {
		data[ct.Name()] = reflect.ValueOf(scanParams[i]).Elem().Interface()
	}
	// Copy?

	return data, nil

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
		if !val.Field(i).CanInterface() {
			continue
		}
		var fieldName string
		tags := parseTag(f.Tag.Get("sqlu"))
		if tags.fieldName != "" {
			fieldName = strings.ToLower(tags.fieldName)
			val.Field(i).Addr().Interface()
		} else {
			fieldName = strings.ToLower(f.Type.Name())
		}
		fields[fieldName] = val.Field(i).Addr().Interface()
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
