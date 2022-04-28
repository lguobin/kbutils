package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

// 自定义类型
type CustomJSON struct{}
type CustomText struct{}

var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}

// Get Data Type for mysql Dialect
func (s *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "int(11)"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "float(6,2)"
	case reflect.String:
		return "varchar(255)"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	case reflect.Interface:
		return "json"
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL returns SQL that judge whether the table exists in database
func (s *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "select group_concat(table_name) from information_schema.tables where table_name = ?", args
}
