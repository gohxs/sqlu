package sqlu

import (
	"database/sql"
	"fmt"
	"strings"
)

// TableInsertContext using a context
func Insert(db SQLer, s Schemer) (sql.Result, error) {
	schema := s.Schema()

	fieldNames := make([]string, len(schema.Fields))
	fieldPtrs := make([]interface{}, len(schema.Fields))
	for i, f := range schema.Fields {
		fieldNames[i] = f.Name
		fieldPtrs[i] = f.Ptr
	}
	//fieldNames := schema.FieldNames
	//fieldPtrs := schema.FieldPtrs
	qry := fmt.Sprintf(
		"INSERT INTO \"%s\" (%s) values(%s)",
		schema.Table,
		strings.Join(fieldNames, ","),
		strings.Repeat("?, ", len(fieldNames)-1)+"?",
	)
	return db.Exec(qry, fieldPtrs...)
}
