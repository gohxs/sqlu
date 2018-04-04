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

func NewRowScan(row RowScanner) *RowScan {
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
func (r RowScan) Scan(s Schemer) error {

	if !r.started {
		// Cache columns
		var err error
		r.cols, err = r.row.Columns()
		if err != nil {
			return err
		}
		schema := s.Schema()
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
	schema := s.Schema() // slow?
	// Map
	for i, fi := range r.fieldI {
		r.values[i] = schema.Fields[fi].Ptr
	}
	// map fields to the values
	return r.row.Scan(r.values...)
}
