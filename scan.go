package sqlu

type RowScan struct {
	err  error
	rows RowsScanner
	// Cache cols and see different
	cols []string

	fieldI  []int
	values  []interface{}
	started bool
}

func NewRowScanner(rows RowsScanner, err ...error) *RowScan {
	if len(err) != 0 && err[0] != nil {
		return &RowScan{err: err[0]}
	}
	return &RowScan{
		err:     nil,
		rows:    rows,
		started: false,

		cols:   []string{},
		fieldI: nil,

		values: nil,
		// to avoid alocations
	}
}
func (r *RowScan) Next() bool {
	return r.rows.Next()
}

//Scan will use scanning thing
func (r *RowScan) Scan(s FieldMapper) error {
	if r.err != nil {
		return r.err
	}

	if !r.started {
		schema := s.Schema()
		// Cache columns
		var err error
		r.cols, err = r.rows.Columns()
		if err != nil {
			return err
		}
		r.values = make([]interface{}, len(r.cols))
		for _, cn := range r.cols {
			for i, f := range schema.Fields {
				if f.Name == cn {
					r.fieldI = append(r.fieldI, i)
				}
			}
		}
		r.started = true
	}

	// Map fields to values
	fields := s.Fields()
	for i, fi := range r.fieldI {
		r.values[i] = fields[fi]
	}
	return r.rows.Scan(r.values...)
}

func Scan(row RowScanner, s FieldMapper) error {
	return row.Scan(s.Fields()...)
}
