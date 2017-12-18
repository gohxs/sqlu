package sqlu

import (
	"context"
	"database/sql"
	"fmt"
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
		tags := parseTag(f.Tag.Get("sqlu"))

		var fieldOptions string
		var typeName = f.Type.Name()
		if tags.fieldName == "" {
			continue
		}
		if tags.PrimaryKey {
			fieldOptions += " PRIMARY KEY"
		}
		if tags.Unique {
			fieldOptions += " UNIQUE"
		}
		if tags.NotNull {
			fieldOptions += " NOT NULL"
		}
		if tn, ok := typeMap[f.Type.Name()]; ok {
			typeName = tn
		}

		fieldEntry := fmt.Sprintf("\"%s\" %s %s", tags.fieldName, typeName, fieldOptions)
		fields = append(fields, fieldEntry)
	}
	qry := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS \"%s\" (%s)",
		table,
		strings.Join(fields, ","),
	)
	return db.ExecContext(ctx, qry)
}
