package sqlu

type Schemer interface {
	scheme() Schema
	SchemaInit(s *Schema)
}

type SchemaInfo struct {
	Table string
	Field []Field
}
type Schema struct {
	*SchemaInfo
	ptrs []interface{}
}
type Field struct {
	Name string
	Type string
}

func (s *Schema) scheme() *SchemaInfo {
	if s.SchemaInfo == nil {

	}

}

func (s *Schema) SchemaInit(s *Scheme) {
	panic("please override SchemaInit")
}

type Model struct {
	Schema
	ID   int
	Name string
}
