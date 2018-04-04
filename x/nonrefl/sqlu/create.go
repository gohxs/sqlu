package sqlu

import (
	"database/sql"
	"fmt"
	"strings"
)

func Create(db SQLer, s Schemer) (sql.Result, error) {
	schema := s.Schema()
	// Each field

	createFields := make([]string, len(schema.Fields))
	for i, f := range schema.Fields {
		createFields[i] = f.Name + " " + f.Type
		/*if schema.FieldOpts == nil {
			continue
		}
		opt, ok := schema.FieldOpts[f.Name]
		if !ok {
			continue
		}
		if opt.IsKey {
			createFields[i] += " primary key"
		}*/
	}
	qry := fmt.Sprintf(
		"CREATE TABLE \"%s\" (%s);",
		schema.Table,
		strings.Join(createFields, ","),
	)

	return db.Exec(qry)
}
