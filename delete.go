package sqlu

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Delete execute a delete statement on a tabler
func Delete(db SQLer, data Tabler) (sql.Result, error) {
	return TableDeleteContext(context.Background(), db, data.Table(), data)
}

// DeleteContext execute delete on a tabler interface with context
func DeleteContext(ctx context.Context, db SQLer, data Tabler) (sql.Result, error) {
	return TableDeleteContext(ctx, db, data.Table(), data)
}

//TableDeleteContext build and execute delete statement from a struct
func TableDeleteContext(ctx context.Context, db SQLer, table string, data interface{}) (sql.Result, error) {
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	keys := []string{} //?
	keyvals := []interface{}{}
	// Go through sqlu tags
	for i := 0; i < val.NumField(); i++ {
		f := typ.Field(i)
		if !val.Field(i).CanInterface() {
			continue
		}
		tags := parseTag(f.Tag.Get("sqlu"))
		if tags.fieldName == "" {
			continue
		}
		if tags.PrimaryKey {
			keys = append(keys, "\""+tags.fieldName+"\" = ?")
			keyvals = append(keyvals, val.Field(i).Interface())
			continue
		}
	}
	if len(keys) == 0 {
		return nil, errors.New("Key(s) not defiend")
	}

	qry := fmt.Sprintf(
		"DELETE FROM \"%s\" WHERE %s",
		table,
		strings.Join(keys, " AND "),
	)
	Log.Println("Qry:", qry, keyvals)

	return db.ExecContext(ctx, qry, keyvals...)

}
