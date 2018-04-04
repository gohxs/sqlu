package sqlu

type RowScan struct {
	//rowScan Cache
	row RowScanner
	// Cache cols and see different
	cols []string

	fieldI  []int
	values  []interface{}
	started bool
}

func NewRowScanner(row RowScanner) *RowScan {
	return &RowScan{
		row:     row,
		started: false,

		cols:   []string{},
		fieldI: nil,

		values: nil,
		// to avoid alocations
	}
}

//Scan will use scanning thing
func (r *RowScan) Scan(s Schemer) error {

	if !r.started {
		schema := s.Schema()
		// Cache columns
		var err error
		r.cols, err = r.row.Columns()
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

	fields := s.Fields()
	for i, fi := range r.fieldI {
		r.values[i] = fields[fi]
	}
	return r.row.Scan(r.values...)
}
