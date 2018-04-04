package sqlu

import (
	"database/sql"
	"fmt"
	"strings"
)

// Sample building facility
type sample struct {
	Schemer Schemer
	Fields  []string
}

// Sample returns a sample
func Sample(s Schemer, fields ...string) sample {
	return sample{s, fields}
}

// FindQRY creates a query
func FindQRY(samples ...sample) (string, []interface{}) {
	// Build a where
	params := []interface{}{}
	orClause := []string{}
	table := ""
	for i, sample := range samples {
		if i == 0 {
			table = sample.Schemer.Schema().Table
		}
		andClause := []string{}
		for _, sf := range sample.Fields {
			f, fieldI := sample.Schemer.Schema().fieldByName(sf)
			fields := sample.Schemer.Fields()
			params = append(params, fields[fieldI])
			andClause = append(andClause, fmt.Sprintf("%s=$%d", f.Name, len(params)))
		}
		orClause = append(orClause, "("+strings.Join(andClause, " AND ")+")")
	}
	// Build query:
	qry := fmt.Sprintf(
		`SELECT * FROM "%s" WHERE %s`,
		table,
		strings.Join(orClause, " OR "),
	)

	return qry, params

}

// Find returns a cursor
func Find(db SQLer, samples ...sample) (*sql.Rows, error) {
	qry, params := FindQRY(samples...)
	return db.Query(qry, params...)
}
