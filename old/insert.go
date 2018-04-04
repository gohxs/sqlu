package sqlu

// Insert a tabler
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Insert insert a struct with Table() with context
func Insert(db SQLer, data Tabler) (sql.Result, error) {
	return InsertContext(context.Background(), db, data)
}

// InsertContext insert a struct with Table() with context
func InsertContext(ctx context.Context, db SQLer, data Tabler) (sql.Result, error) {
	return TableInsertContext(ctx, db, data.Table(), data)
}

// TableInsert struct into specified table
// Warning table is not escaped
func TableInsert(db SQLer, table string, data interface{}) (sql.Result, error) {
	return TableInsertContext(context.Background(), db, table, data)
}

// TableInsertContext using a context
func TableInsertContext(ctx context.Context, db SQLer, table string, data interface{}) (sql.Result, error) {
	val := reflect.ValueOf(data)
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	fields := []string{}
	values := []interface{}{}
	// Go through sqlu tags
	for i := 0; i < val.NumField(); i++ {
		f := typ.Field(i)
		if !val.Field(i).CanInterface() {
			continue
		}
		tags := parseTag(f.Tag.Get("sqlu"))
		if tags.fieldName == "" { // Not a sqlu, or malformed
			continue
		}

		var value = val.Field(i).Interface()
		if tags.OmitEmpty && value == reflect.Zero(val.Field(i).Type()).Interface() {
			continue
		}
		if tags.CreateTimeStamp || tags.UpdateTimeStamp {
			Log.Printf("[%s] - Create timestamp\n", tags.fieldName)
			value = time.Now().UTC()
		}

		fields = append(fields, "\""+tags.fieldName+"\"")
		values = append(values, value)
	}
	if len(fields) == 0 {
		return nil, errors.New("No fields")
	}
	qry := fmt.Sprintf(
		"INSERT INTO \"%s\" (%s) values(%s)",
		table,
		strings.Join(fields, ","),
		strings.Repeat("?, ", len(fields)-1)+"?",
	)
	Log.Println("QRY:", qry, values)

	return db.ExecContext(ctx, qry, values...)
}
