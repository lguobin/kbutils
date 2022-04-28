package schema

import (
	"go/ast"
	"reflect"

	"github.com/lguobin/kbutils/kborm/dialect"
)

const TagName = "kborm"

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

// GetField returns field by name
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// Values return the values of dest's member variables
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

type ITableName interface {
	TableName() string
}

// Parse a struct to a Schema instance
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}
	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	// recursion
	Field_alls := search_Struct(modelType, d)
	for _, field := range Field_alls {
		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, field.Name)
		schema.fieldMap[field.Name] = field
	}
	return schema
}

func search_Struct(_struct reflect.Type, d dialect.Dialect) []*Field {
	Fields := make([]*Field, 0)
	for i := 0; i < _struct.NumField(); i++ {
		p := _struct.Field(i)
		if p.Type.Kind() == reflect.Struct {
			temp := search_Struct(p.Type, d)
			for _, v := range temp {
				Fields = append(Fields, v)
			}
		} else {
			if !p.Anonymous && ast.IsExported(p.Name) {
				field_one := &Field{
					Name: p.Name,
					Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
				}
				if v, ok := p.Tag.Lookup(TagName); ok {
					field_one.Tag = v
				}
				Fields = append(Fields, field_one)
			}
		}
	}
	return Fields
}
