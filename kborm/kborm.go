package kborm

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lguobin/kbutils/kborm/dialect"
	"github.com/lguobin/kbutils/kborm/log"
	"github.com/lguobin/kbutils/kborm/session"
	_ "github.com/lguobin/kbutils/kborm/sqldriver/mysql"
)

// Engine is the main struct of geeorm, manages all db sessions and transactions.
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine create a instance of Engine
// connect database and ping it to test whether it's alive
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

// Close database connection
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
		return
	}
	log.Info("Close database success")
}

// NewSession creates a new session for next operations
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

// TxFunc will be called between tx.Begin() and tx.Commit()
// https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
type TxFunc func(*session.Session) (interface{}, error)

// Transaction executes sql wrapped in a transaction, then automatically commit if no error occurs
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
			// 如果 err 不为空, 则会回滚
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
		}
	}()
	return f(s)
}

// difference returns a - b
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(valueList ...interface{}) error {
	var err error
	for _, value := range valueList {
		_, err = engine.Transaction(func(s *session.Session) (result interface{}, err error) {
			if !s.Model(value).HasTable() {
				log.Infof("table %s doesn't exist", s.RefTable().Name)
				return nil, s.CreateTable()
			}
			table := s.RefTable()
			rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
			columns, _ := rows.Columns()

			// tx 在执行Query()操作后，rows会维护这个数据库连接
			// 当 tx 想再次调用当前连接进行数据库操作的时候
			// 因为连接还没有断开, 没有调用 rows.Close()
			// tx 无法再从连接池里获取当前连接，所以会提示 busy buffer
			// 必须主动 close()
			rows.Close()

			addCols := difference(table.FieldNames, columns)
			delCols := difference(columns, table.FieldNames)
			log.Infof("added cols %v, deleted cols %v", addCols, delCols)

			for _, col := range addCols {
				f := table.GetField(col)
				sqlStr := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
				if _, err = s.Raw(sqlStr).Exec(); err != nil {
					return
				}
			}

			if len(delCols) == 0 {
				return
			}
			tmp := "tmp_" + table.Name
			fieldStr := strings.Join(table.FieldNames, ", ")
			s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from `%s`;", tmp, fieldStr, table.Name))
			s.Raw(fmt.Sprintf("DROP TABLE `%s`;", table.Name))
			s.Raw(fmt.Sprintf("ALTER TABLE `%s` RENAME TO `%s`;", tmp, table.Name))
			_, err = s.Exec()
			return
		})
		if err != nil {
			break
		}
	}
	return err
}
