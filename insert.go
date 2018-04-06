package sqlu

import (
	"database/sql"
	"fmt"
	"strings"
)

func InsertQRY(s Schemer) (string, []interface{}) {
	schema := S{}
	s.Schema(&schema)

	fieldNames := make([]string, len(schema.Schema.Fields))
	fieldParam := make([]string, len(schema.Schema.Fields))
	fieldPtrs := make([]interface{}, len(schema.Schema.Fields))
	for i, f := range schema.Schema.Fields {
		fieldNames[i] = f.Name
		fieldParam[i] = fmt.Sprintf("$%d", i+1)
	}
	reflectFields(s, fieldPtrs)
	//schema.LoadFields(fieldPtrs)
	qry := fmt.Sprintf(
		"INSERT INTO \"%s\" (%s) values(%s)",
		schema.Schema.Table,
		strings.Join(fieldNames, ","),
		strings.Join(fieldParam, ","),
	)
	return qry, fieldPtrs
}

func Insert(db Queryer, s Schemer) (sql.Result, error) {
	q, f := InsertQRY(s)
	return db.Exec(q, f...)
}
