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
	// Map
	ptrs := s.SchemaFields()
	for i, fi := range r.fieldI {
		r.values[i] = ptrs[fi]
	}
	// map fields to the values
	return r.row.Scan(r.values...)
	// Load to temporary and copy/set to the one on param
	/*err := r.row.Scan(r.fields...)
	if err != nil {
		return err
	}
	reflect.ValueOf(s).Elem().Set(reflect.ValueOf(r.thing).Elem())
	return nil*/

	/*ptrs := s.SchemaFields()
	if len(r.cols) == len(ptrs) {
		return row.Scan(ptrs...)
	}

	schema := s.Schema()
	scanParams := make([]interface{}, len(r.cols))
	for i, cn := range r.cols {
		v := schema.fieldPtr(cn)
		if v == nil {
			return fmt.Errorf("Field '%s' not found on schemer", cn)
		}
		scanParams[i] = v

	}

	return row.Scan(scanParams...)*/
}
