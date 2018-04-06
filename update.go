package sqlu

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func Update(db Queryer, s Schemer, fnfields []string, where string, fnparams ...interface{}) (sql.Result, error) {
	schema := S{}
	s.Schema(&schema)
	// For Schemer fields
	fields := []string{}
	params := []interface{}{}
	fieldptr := schema.fields()
	for i, f := range schema.Schema.Fields {
		for _, sf := range fnfields {
			if f.Name == sf {
				fields = append(fields, f.Name+"= ?")
				params = append(params, fieldptr[i])
			}
		}
	}

	params = append(params, fnparams...)
	qry := fmt.Sprintf(
		"UPDATE \"%s\" SET %s WHERE %s",
		schema.Schema.Table,
		strings.Join(fields, ", "),
		where,
	)
	log.Println("Qry:", qry, params)
	return db.Exec(qry, params...)
}
