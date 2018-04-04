package sqlu

type Schemer interface {
	Schema() Schema
	SchemaFields() []interface{}
}

// Resulting internal schema
type Schema struct {
	Table     string
	Fields    []Field
	FieldOpts map[string]FieldOpt
}

type Field struct {
	Ptr  interface{}
	Name string
	Type string // or suffix
	// Maybe options
}
type FieldOpt struct {
	OmitEmpty bool
	IsKey     bool
}

func (s *Schema) Schema() *Schema {
	return s
}

// This could be cached somehow
func (s Schema) fieldNames() []string {
	names := make([]string, len(s.Fields))
	for i, f := range s.Fields {
		names[i] = f.Name
	}
	return names
}
func (s Schema) fieldTypes() []string {

	types := make([]string, len(s.Fields))
	for i, f := range s.Fields {
		types[i] = f.Type
	}
	return types
}

func (s Schema) fieldPtr(name string) interface{} {
	for _, f := range s.Fields {
		if f.Name == name {
			return f.Ptr
		}
	}
	return nil
}
