package sqlu

import (
	"database/sql"
	"fmt"
	"strings"
)

func CreateQRY(s Schemer) string {
	schema := S{}
	s.Schema(&schema)
	// Each field
	createFields := make([]string, len(schema.Schema.Fields))
	for i, f := range schema.Schema.Fields {
		createFields[i] = f.Name + " " + f.Type
	}
	qry := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS \"%s\" (%s);",
		schema.Schema.Table,
		strings.Join(createFields, ","),
	)
	return qry
}

func Create(db Queryer, s Schemer) (sql.Result, error) {
	return db.Exec(CreateQRY(s))
}
