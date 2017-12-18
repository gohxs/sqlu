package sqlu

// Insert a tabler
import (
	"context"
	"database/sql"
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
		var value interface{}
		tags := strings.Split(f.Tag.Get("sqlu"), ",")
		if len(tags) > tagField {
			fields = append(fields, "\""+tags[tagField]+"\"")
			if len(tags) > tagType && tags[tagType] == "createTimeStamp" {
				value = time.Now().UTC()
			} else {
				value = val.Field(i).Interface()
			}
		}
		values = append(values, value)
	}
	qry := fmt.Sprintf(
		"INSERT INTO \"%s\" (%s) values(%s)",
		table,
		strings.Join(fields, ","),
		strings.Repeat("?, ", len(fields)-1)+"?",
	)

	return db.ExecContext(ctx, qry, values...)
}
