package sqlu

import (
	"database/sql"
	"fmt"
	"strings"
)

func CreateQRY(s Schemer) string {
	schema := s.Schema()
	// Each field
	createFields := make([]string, len(schema.Fields))
	for i, f := range schema.Fields {
		createFields[i] = f.Name + " " + f.Type
	}
	qry := fmt.Sprintf(
		"CREATE TABLE \"%s\" (%s);",
		schema.Table,
		strings.Join(createFields, ","),
	)
	return qry
}

func Create(db SQLer, s Schemer) (sql.Result, error) {
	return db.Exec(CreateQRY(s))
}
