package sqlu

type RowScan struct {
	err  error
	rows RowsScanner
	// Cache cols and see different

	schema  S
	fieldI  []int
	values  []interface{}
	started bool
}

func NewRowScanner(rows RowsScanner, err ...error) (*RowScan, error) {
	if len(err) != 0 && err[0] != nil {
		return nil, err[0]
	}
	return &RowScan{
		err:     nil,
		rows:    rows,
		started: false,

		// to avoid alocations
		fieldI: nil,
		values: nil,
	}, nil
}
func (r *RowScan) Next() bool {
	return r.rows.Next()
}
func (r *RowScan) Close() error {
	return r.rows.Close()
}

//Scan will use scanning thing
func (r *RowScan) Scan(s Schemer) error {
	if r.err != nil {
		return r.err
	}
	if !r.started {
		r.schema = S{}

		s.Schema(&r.schema)
		// Cache columns
		cols, err := r.rows.Columns()
		if err != nil {
			return err
		}
		r.fieldI = make([]int, len(cols))
		for ci, cn := range cols {
			// Find by name
			for fi, f := range r.schema.Schema.Fields {
				if f.Name == cn {
					r.fieldI[ci] = fi
				}
			}
		}
		r.values = make([]interface{}, len(cols))
		r.started = true
	}
	// Load Schema
	//s.Schema(&r.schema)
	//r.schema.LoadFields(r.values, r.fieldI...)
	reflectFields(s, r.values, r.fieldI...)
	return r.rows.Scan(r.values...)
}

func Scan(row RowScanner, s Schemer) error {
	schema := S{}
	s.Schema(&schema)
	fields := make([]interface{}, len(schema.Schema.Fields))
	schema.LoadFields(fields) // Load
	return row.Scan(fields...)
}
