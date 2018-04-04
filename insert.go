package sqlu

import (
	"database/sql"
	"fmt"
	"strings"
)

func InsertQRY(s Schemer) (string, []interface{}) {
	schema := s.Schema()

	fieldNames := make([]string, len(schema.Fields))
	fieldParam := make([]string, len(schema.Fields))
	fieldPtrs := make([]interface{}, len(schema.Fields))
	fields := s.Fields()
	for i, f := range schema.Fields {
		fieldNames[i] = f.Name
		fieldParam[i] = fmt.Sprintf("$%d", i+1)
		fieldPtrs[i] = fields[i]
	}
	qry := fmt.Sprintf(
		"INSERT INTO \"%s\" (%s) values(%s)",
		schema.Table,
		strings.Join(fieldNames, ","),
		strings.Join(fieldParam, ","),
	)
	return qry, fieldPtrs
}

func Insert(db SQLer, s Schemer) (sql.Result, error) {
	q, f := InsertQRY(s)
	return db.Exec(q, f...)
}
