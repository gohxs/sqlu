package sqlu

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

//Update update a table with a tabler
func Update(db SQLer, data Tabler) (sql.Result, error) {
	return TableUpdateContext(context.Background(), db, data.Table(), data)
}

//TableUpdateContext update a field based on name
func TableUpdateContext(ctx context.Context, db SQLer, table string, data interface{}) (sql.Result, error) {
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	fields := []string{}
	values := []interface{}{}
	keys := []string{} //?
	keyvals := []interface{}{}
	// Go through sqlu tags
	for i := 0; i < val.NumField(); i++ {
		f := typ.Field(i)
		if !val.Field(i).CanInterface() {
			continue
		}
		var value interface{}
		tags := parseTag(f.Tag.Get("sqlu"))
		if tags.fieldName == "" {
			continue
		}

		if tags.PrimaryKey {
			keys = append(keys, "\""+tags.fieldName+"\" = ?")
			keyvals = append(keyvals, val.Field(i).Interface())
			continue
		}

		fields = append(fields, "\""+tags.fieldName+"\" = ?")
		if tags.UpdateTimeStamp {
			value = time.Now().UTC()
		} else {
			value = val.Field(i).Interface()
		}
		values = append(values, value)
	}

	params := append([]interface{}{}, values...)
	params = append(params, keyvals...)
	qry := fmt.Sprintf(
		"UPDATE \"%s\" SET %s WHERE %s",
		table,
		strings.Join(fields, ","),
		strings.Join(keys, " AND "),
	)
	log.Println("Qry:", qry, params)

	return db.ExecContext(ctx, qry, params...)

}
