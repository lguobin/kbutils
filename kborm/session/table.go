package session

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lguobin/kbutils/kborm/log"
	"github.com/lguobin/kbutils/kborm/schema"
)

// Model assigns refTable
func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable returns a Schema instance that contains all parsed fields
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable create a table in database with a model
func (s *Session) CreateTable() error {
	var _SQL string
	var PK string
	var Cons []string
	var PRIMARY_KEY string
	var columns []string
	table := s.RefTable()

	for _, field := range table.Fields {
		if strings.ToUpper(field.Tag) == "primary key" {
			field.Tag = ""
			PK = "PK_" + table.Name
			Cons = append(columns, fmt.Sprintf("%s ", field.Name))
		}
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	if PK != "" {
		Constraint := strings.Join(Cons, ",")
		PRIMARY_KEY = fmt.Sprintf("CONSTRAINT `%s` PRIMARY KEY (%s)", PK, Constraint)
		_SQL = fmt.Sprintf("CREATE TABLE `%s` (%s , %s);", table.Name, desc, PRIMARY_KEY)
	} else {
		_SQL = fmt.Sprintf("CREATE TABLE `%s` (%s);", table.Name, desc)
	}
	_, err := s.Raw(_SQL).Exec()
	return err
}

// DropTable drops a table with the name of model
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// HasTable returns true of the table exists
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
