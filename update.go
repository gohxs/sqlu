package sqlu

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func Update(db SQLer, s Schemer, fnfields []string, where string, fnparams ...interface{}) (sql.Result, error) {
	schema := s.Schema()
	// For Schemer fields
	fields := []string{}
	params := []interface{}{}
	fieldptr := s.Fields()
	for i, f := range schema.Fields {
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
		schema.Table,
		strings.Join(fields, ", "),
		where,
	)
	log.Println("Qry:", qry, params)
	return db.Exec(qry, params...)
}
