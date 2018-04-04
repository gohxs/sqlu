package sqlu

import (
	"database/sql"
	"fmt"
)

func Get(db SQLer, s Schemer, qry string, params ...interface{}) error {
	row, err := db.Query(qry, params...)
	if err != nil {
		return err
	}
	defer row.Close()
	if !row.Next() {
		return sql.ErrNoRows
	}

	colTypes, err := row.ColumnTypes()
	if err != nil {
		return err
	}

	ptr := s.SchemaFields()
	if len(colTypes) == len(ptr) {
		return row.Scan(ptr...)
	}

	schema := s.Schema()
	scanParams := make([]interface{}, len(colTypes))
	for i, c := range colTypes {

		v := schema.fieldPtr(c.Name())
		if v == nil {
			return fmt.Errorf("Field '%s' not found on schemer", c.Name())
		}
		scanParams[i] = v

	}

	return row.Scan(scanParams...)
}
