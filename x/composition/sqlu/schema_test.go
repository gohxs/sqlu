package sqlu_test

import "github.com/gohxs/sqlu/x/composition/sqlu"

type Model struct {
	sqlu.Schema
	ID   int
	Name string
}

func (m *Model) SchemaInit(s *Schema) {
	s.
		Fields(&m.ID, &m.Name).
		Field("id", "int").
		Field("name", "string")

}
