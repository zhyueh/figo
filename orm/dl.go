package orm

import (
	"database/sql"
	"errors"
	"fmt"
)

type DbRow map[string]interface{}

type dlInterface interface {
	Close()
	Exec(string, ...interface{}) (sql.Result, error)
	Insert(string, ...interface{}) (int64, error)
	Update(string, ...interface{}) (int64, error)
	Delete(string, ...interface{}) (int64, error)
	One(string, ...interface{}) (DbRow, error)
	All(string, ...interface{}) ([]DbRow, error)

	TBegin() error
	TCommit() error
	TRollback() error
	TExec(string, ...interface{}) (sql.Result, error)
	TQuery(string, ...interface{}) ([]DbRow, error)
	TInsert(string, ...interface{}) (int64, error)
	TUpdate(string, ...interface{}) (int64, error)
}

func NewDL(dbType, host, user, password, db string, port int) (dlInterface, error) {
	if dbType == "mysql" {
		connectString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local", user, password, host, port, db)
		return NewMysqlLayer(dbType, connectString)
	}
	return nil, errors.New("un support db layer type")
}
