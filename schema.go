package sqlu

import (
	"reflect"
	"strings"
)

var (
	OptOmitEmpty Opt
	OptKey       Opt
)

type Opt struct{}

type Schemer interface {
	Schema(s *S)
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

type S struct {
	Schema     *Schema // Meta internal?
	LoadFields func([]interface{}, ...int)
}
type InitFunc func(m Schemer, s *Schema)

func (s *S) BuildSchema(name string, m Schemer, init InitFunc, fields func([]interface{}, ...int)) {

	s.LoadFields = fields
	// Field loader
	// Cached schema
	var ok bool
	if s.Schema, ok = schemaCache[name]; ok {
		return
	}
	s.Schema = &Schema{Table: name}
	init(m, s.Schema)
	schemaCache[name] = s.Schema

}

func (s S) fields() []interface{} {
	fields := make([]interface{}, len(s.Schema.Fields))
	s.LoadFields(fields)
	return fields
}

// StructSchema retuns a schema builder for the specific struct

func ReflectSchema(m Schemer, s *Schema) {
	typ := reflect.TypeOf(m).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type.Kind() == reflect.Struct {
			continue
		}
		fname := strings.ToLower(typ.Field(i).Name)
		tname := typ.Field(i).Type.Name()
		s.Field(fname, tname)
	}
}

// Fields utility to pass pointers as fields
func Fields(ptrs ...interface{}) func([]interface{}, ...int) {
	return func(values []interface{}, indexes ...int) {
		if len(indexes) == 0 {
			for i := range values {
				values[i] = ptrs[i]
			}
			return
		}
		//fields := []interface{}{}
		for i, fi := range indexes {
			values[i] = ptrs[fi]
		}
	}
}

func ReflectFieldsx(s Schemer, values []interface{}, indexes ...int) {
	reflectFields(s, values, indexes...)
}

func ReflectFields(s interface{}) func([]interface{}, ...int) {
	return func(values []interface{}, indexes ...int) {
		reflectFields(s, values, indexes...)
		/*val := reflect.ValueOf(s).Elem()
		if len(indexes) == 0 {
			for i := 0; i < val.NumField(); i++ {
				if val.Field(i).Kind() == reflect.Struct {
					continue
				}
				values[i] = val.Field(i).Addr().Interface()
			}
			return
		}
		// Should be a pointer
		for i, fi := range indexes {
			if val.Field(i).Kind() == reflect.Struct { // or anything else
				continue
			}
			values[i] = val.Field(fi).Addr().Interface()
		}*/
	}
}
func reflectFields(s interface{}, values []interface{}, indexes ...int) {
	val := reflect.ValueOf(s).Elem()
	if len(indexes) == 0 {
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).Kind() == reflect.Struct {
				continue
			}
			values[i] = val.Field(i).Addr().Interface()
		}
		return
	}
	// Should be a pointer
	for i, fi := range indexes {
		if val.Field(i).Kind() == reflect.Struct { // or anything else
			continue
		}
		values[i] = val.Field(fi).Addr().Interface()
	}

}

/*func GetFields(s Schemer) []interface{} {
	return s.Schema().fields()
}*/
