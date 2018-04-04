package sqlu

var (
	OptOmitEmpty Opt
	OptKey       Opt
)

type Opt struct{}

type FieldMapper interface {
	Schemer
	Fields() []interface{}
}
type Schemer interface {
	Schema() *Schema
}

// Field table information
type Field struct {
	Name      string
	Type      string // or suffix
	OmitEmpty bool
	IsKey     bool
	// Maybe options
}

// Schema represents a database schema
type Schema struct { // Internal
	Table  string
	Fields []Field
}

/////////////////////////////////
// Builder one
///////////////////////

// AddField SchemaBuilder
func (s *Schema) Field(name string, typ string, opts ...Opt) *Schema {
	f := Field{
		Name: name,
		Type: typ,
	}
	for _, o := range opts {
		switch o {
		case OptOmitEmpty:
			f.OmitEmpty = true
		case OptKey:
			f.IsKey = true
		}
	}
	s.Fields = append(s.Fields, f)
	return s
}

//////////////////////////////////////////
// builder two
/////////////////////////////

var schemaCache = map[string]*Schema{}

func BuildSchema(name string, init func(s *Schema)) *Schema {
	// Cached schema
	if schema, ok := schemaCache[name]; ok {
		return schema
	}
	s := &Schema{Table: name}
	init(s)
	schemaCache[name] = s
	return s
}

// just an alias
func Fields(ptrs ...interface{}) []interface{} {
	return ptrs
}

// Helper functions
func (s *Schema) fieldByName(name string) (*Field, int) {
	for i, f := range s.Fields {
		if f.Name == name {
			return &f, i
		}
	}
	return nil, -1
}
