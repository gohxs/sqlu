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
		if ct.ScanType() == nil {
			scanParams = append(scanParams, &sql.NullString{})
			continue
		}
		p := reflect.New(ct.ScanType())
		scanParams = append(scanParams, p.Interface())
	}
	err = res.Scan(scanParams...)
	if err != nil {
		return nil, err
	}
	for i, ct := range colTypes {
		data[strings.ToLower(ct.Name())] = reflect.ValueOf(scanParams[i]).Elem().Interface()
	}
	// Copy?

	return data, nil

}

// Scan scan sqlu fields of a struct
// This will reflect the structure and match column names
/*func Scan(res *sql.Rows, data interface{}) error {
	Log.Println("Scan")

	rowMap, err := RowToMap(res)
	if err != nil {
		return err
	}
	// Go and copy
	valPtr := reflect.ValueOf(data)
	val := valPtr.Elem()
	typ := valPtr.Elem().Type()
	// Retrieve the fields into a map
	for i := 0; i < val.NumField(); i++ {
		fName := strings.ToLower(typ.Field(i).Name)
		v, ok := rowMap[fName]
		if !ok {
			continue // Ignore field
		}
		Log.Println("Assign:", typ.Field(i).Name, fName, reflect.ValueOf(v).Interface())
		val.Field(i).Set(reflect.ValueOf(v))
	}
	return nil
}*/

func Scan(res *sql.Rows, data interface{}) error {
	fields := map[string]interface{}{} // Pointer interface

	// Read values into a map

	valPtr := reflect.ValueOf(data)
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
		if tags.fieldName == "" {
			continue
		}
		fieldName = strings.ToLower(tags.fieldName) // Can conflict?
		fields[fieldName] = val.Field(i).Addr().Interface()
	}

	colTyp, err := res.ColumnTypes()
	if err != nil {
		return err
	}
	params := []interface{}{}
	for _, ct := range colTyp {
		colName := strings.ToLower(ct.Name())
		if ct.ScanType() == nil { // Assign zero to field, whatever zero is?
			/*if v, ok := fields[colName]; ok {
				reflect.ValueOf(v).Set(reflect.Zero(reflect.TypeOf(v)))
			}*/

			params = append(params, &sql.NullString{})
			continue
		}
		if fv, ok := fields[colName]; ok {
			params = append(params, fv)
		}
	}
	return res.Scan(params...)

}

// ScanRaw a full struct
// usefull when objects perfectly matches a table
func ScanRaw(res *sql.Rows, data interface{}) error {

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

// Data should be pointer
/*fields := map[string]interface{}{} // Pointer interface

	// Read values into a map

	valPtr := reflect.ValueOf(data)
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
		if tags.fieldName == "" {
			continue
		}
		fieldName = strings.ToLower(tags.fieldName) // Can conflict?
		Log.Println("Adding field:", fieldName)
		fields[fieldName] = val.Field(i).Addr().Interface()
	}

	colTyp, err := res.ColumnTypes()
	if err != nil {
		return err
	}
	params := []interface{}{}
	for _, c := range colTyp {
		colName := strings.ToLower(c.Name())
		v, ok := fields[strings.ToLower(colName)]
		Log.Println("Matchign column", colName, reflect.TypeOf(v).Elem().Name())
		if !ok {
			if c.ScanType() != nil {
				Log.Println("ScanType:", c.ScanType().Name())
				v = reflect.New(c.ScanType()) // Dumb pointer
			} else {
				v = ""
			}
		}
		params = append(params, v)
	}
	return res.Scan(params...)
}
*/
