package sqlu

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

var (
	typeMap = map[string]string{
		"int":  "integer",
		"Time": "datetime",
	}
)

// Create create a table based on Tabler interface
func Create(db SQLer, data Tabler) (sql.Result, error) {
	return CreateContext(context.Background(), db, data)
}

// CreateContext same as Create but with context
func CreateContext(ctx context.Context, db SQLer, data Tabler) (sql.Result, error) {
	return TableCreateContext(ctx, db, data.Table(), data)
}

// TableCreateContext Generate SQL Create a named table from a struct
func TableCreateContext(ctx context.Context, db SQLer, table string, data interface{}) (sql.Result, error) {
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	fields := []string{}
	// Go through sqlu tags
	for i := 0; i < val.NumField(); i++ {
		f := typ.Field(i)
		if !val.Field(i).CanInterface() {
			continue
		}
		tags := strings.Split(f.Tag.Get("sqlu"), ",")
		// Translate types by string somehow
		if len(tags) >= tagField {
			typeName, ok := typeMap[f.Type.Name()]
			if !ok {
				typeName = f.Type.Name()
			}

			fieldEntry := fmt.Sprintf("\"%s\" %s", tags[tagField], typeName)
			fields = append(fields, fieldEntry)
		}
	}
	qry := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS \"%s\" (%s)",
		table,
		strings.Join(fields, ","),
	)
	log.Println("Query:", qry)
	return db.ExecContext(ctx, qry)
}
